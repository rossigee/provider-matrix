/*
Copyright 2025 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package user

import (
	"context"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/connection"
	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/crossplane-contrib/provider-matrix/apis/user/v1alpha1"
	apisv1beta1 "github.com/crossplane-contrib/provider-matrix/apis/v1beta1"
	"github.com/crossplane-contrib/provider-matrix/internal/clients"
	"github.com/crossplane-contrib/provider-matrix/internal/features"
)

const (
	errNotUser        = "managed resource is not a User custom resource"
	errTrackPCUsage   = "cannot track ProviderConfig usage"
	errGetPC          = "cannot get ProviderConfig"
	errGetCreds       = "cannot get credentials"
	errNewClient      = "cannot create new Matrix client"
	errCreateUser     = "cannot create Matrix user"
	errGetUser        = "cannot get Matrix user"
	errUpdateUser     = "cannot update Matrix user"
	errDeactivateUser = "cannot deactivate Matrix user"
)

// Setup adds a controller that reconciles User managed resources.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.UserKind)

	cps := []managed.ConnectionPublisher{managed.NewAPISecretPublisher(mgr.GetClient(), mgr.GetScheme())}
	if o.Features.Enabled(features.EnableAlphaExternalSecretStores) {
		cps = append(cps, connection.NewDetailsManager(mgr.GetClient(), v1alpha1.UserGroupVersionKind))
	}

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.UserGroupVersionKind),
		managed.WithExternalConnecter(&connector{
			kube:         mgr.GetClient(),
			usage:        resource.NewProviderConfigUsageTracker(mgr.GetClient(), &apisv1beta1.ProviderConfigUsage{}),
			newServiceFn: clients.NewClient,
		}),
		managed.WithLogger(o.Logger.WithValues("controller", name)),
		managed.WithPollInterval(o.PollInterval),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
		managed.WithConnectionPublishers(cps...))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		WithEventFilter(resource.DesiredStateChanged()).
		For(&v1alpha1.User{}).
		Complete(ratelimiter.NewReconciler(name, r, o.GlobalRateLimiter))
}

// A connector is expected to produce an ExternalClient when its Connect method
// is called.
type connector struct {
	kube         client.Client
	usage        resource.Tracker
	newServiceFn func(config *clients.Config) (clients.Client, error)
}

// Connect typically produces an ExternalClient by:
// 1. Tracking that the managed resource is using a ProviderConfig.
// 2. Getting the managed resource's ProviderConfig.
// 3. Getting the credentials specified by the ProviderConfig.
// 4. Using the credentials to form a client.
func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.User)
	if !ok {
		return nil, errors.New(errNotUser)
	}

	if err := c.usage.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, errTrackPCUsage)
	}

	pc := &apisv1beta1.ProviderConfig{}
	if err := c.kube.Get(ctx, types.NamespacedName{Name: cr.GetProviderConfigReference().Name}, pc); err != nil {
		return nil, errors.Wrap(err, errGetPC)
	}

	config, err := clients.GetConfig(ctx, c.kube, mg)
	if err != nil {
		return nil, errors.Wrap(err, errGetCreds)
	}

	service, err := c.newServiceFn(config)
	if err != nil {
		return nil, errors.Wrap(err, errNewClient)
	}

	return &external{service: service}, nil
}

// An ExternalClient observes, then either creates, updates, or deletes an
// external resource to ensure it reflects the managed resource's desired state.
type external struct {
	service clients.Client
}

func (c *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.User)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotUser)
	}

	userID := meta.GetExternalName(cr)
	if userID == "" {
		return managed.ExternalObservation{
			ResourceExists: false,
		}, nil
	}

	user, err := c.service.GetUser(ctx, userID)
	if err != nil {
		if clients.IsNotFound(err) {
			return managed.ExternalObservation{
				ResourceExists: false,
			}, nil
		}
		return managed.ExternalObservation{}, errors.Wrap(err, errGetUser)
	}

	cr.Status.AtProvider = generateUserObservation(user)
	cr.Status.SetConditions(xpv1.Available())

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: isUserUpToDate(cr, user),
	}, nil
}

func (c *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.User)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotUser)
	}

	userSpec := generateUserSpec(cr)
	user, err := c.service.CreateUser(ctx, userSpec)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreateUser)
	}

	meta.SetExternalName(cr, user.UserID)

	return managed.ExternalCreation{}, nil
}

func (c *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.User)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotUser)
	}

	userID := meta.GetExternalName(cr)
	userSpec := generateUserSpec(cr)
	_, err := c.service.UpdateUser(ctx, userID, userSpec)
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errUpdateUser)
	}

	return managed.ExternalUpdate{}, nil
}

func (c *external) Delete(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*v1alpha1.User)
	if !ok {
		return errors.New(errNotUser)
	}

	userID := meta.GetExternalName(cr)
	if userID == "" {
		return nil
	}

	return errors.Wrap(c.service.DeactivateUser(ctx, userID), errDeactivateUser)
}

// Helper functions

func generateUserSpec(cr *v1alpha1.User) *clients.UserSpec {
	spec := &clients.UserSpec{}

	if cr.Spec.ForProvider.UserID != nil {
		spec.UserID = *cr.Spec.ForProvider.UserID
	}
	if cr.Spec.ForProvider.Localpart != nil {
		spec.Localpart = *cr.Spec.ForProvider.Localpart
	}
	if cr.Spec.ForProvider.Password != nil {
		spec.Password = *cr.Spec.ForProvider.Password
	}
	if cr.Spec.ForProvider.DisplayName != nil {
		spec.DisplayName = *cr.Spec.ForProvider.DisplayName
	}
	if cr.Spec.ForProvider.AvatarURL != nil {
		spec.AvatarURL = *cr.Spec.ForProvider.AvatarURL
	}
	if cr.Spec.ForProvider.Admin != nil {
		spec.Admin = *cr.Spec.ForProvider.Admin
	}
	if cr.Spec.ForProvider.Deactivated != nil {
		spec.Deactivated = *cr.Spec.ForProvider.Deactivated
	}
	if cr.Spec.ForProvider.UserType != nil {
		spec.UserType = *cr.Spec.ForProvider.UserType
	}

	// Convert external IDs
	for _, extID := range cr.Spec.ForProvider.ExternalIDs {
		validated := false
		if extID.Validated != nil {
			validated = *extID.Validated
		}
		spec.ExternalIDs = append(spec.ExternalIDs, clients.ExternalID{
			Medium:    extID.Medium,
			Address:   extID.Address,
			Validated: validated,
		})
	}

	if cr.Spec.ForProvider.ExpireTime != nil {
		spec.ExpireTime = &cr.Spec.ForProvider.ExpireTime.Time
	}

	return spec
}

func generateUserObservation(user *clients.User) v1alpha1.UserObservation {
	obs := v1alpha1.UserObservation{
		UserID:      user.UserID,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
		Admin:       user.Admin,
		Deactivated: user.Deactivated,
		UserType:    user.UserType,
	}

	if user.CreationTime != nil {
		obs.CreationTime = &metav1.Time{Time: *user.CreationTime}
	}
	if user.LastSeenTime != nil {
		obs.LastSeenTime = &metav1.Time{Time: *user.LastSeenTime}
	}

	// Convert external IDs
	for _, extID := range user.ExternalIDs {
		validated := &extID.Validated
		obs.ExternalIDs = append(obs.ExternalIDs, v1alpha1.ExternalID{
			Medium:    extID.Medium,
			Address:   extID.Address,
			Validated: validated,
		})
	}

	// Convert devices
	for _, device := range user.Devices {
		deviceObs := v1alpha1.Device{
			DeviceID:    device.DeviceID,
			DisplayName: device.DisplayName,
			LastSeenIP:  device.LastSeenIP,
		}
		if device.LastSeenTime != nil {
			deviceObs.LastSeenTime = &metav1.Time{Time: *device.LastSeenTime}
		}
		obs.Devices = append(obs.Devices, deviceObs)
	}

	return obs
}

func isUserUpToDate(cr *v1alpha1.User, user *clients.User) bool {
	// Check display name
	if cr.Spec.ForProvider.DisplayName != nil && *cr.Spec.ForProvider.DisplayName != user.DisplayName {
		return false
	}

	// Check avatar URL
	if cr.Spec.ForProvider.AvatarURL != nil && *cr.Spec.ForProvider.AvatarURL != user.AvatarURL {
		return false
	}

	// Check admin status
	if cr.Spec.ForProvider.Admin != nil && *cr.Spec.ForProvider.Admin != user.Admin {
		return false
	}

	// Check deactivated status
	if cr.Spec.ForProvider.Deactivated != nil && *cr.Spec.ForProvider.Deactivated != user.Deactivated {
		return false
	}

	// Check user type
	if cr.Spec.ForProvider.UserType != nil && *cr.Spec.ForProvider.UserType != user.UserType {
		return false
	}

	return true
}
