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

package powerlevel

import (
	"context"
	"time"

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

	"github.com/crossplane-contrib/provider-matrix/apis/powerlevel/v1alpha1"
	apisv1beta1 "github.com/crossplane-contrib/provider-matrix/apis/v1beta1"
	"github.com/crossplane-contrib/provider-matrix/internal/clients"
	"github.com/crossplane-contrib/provider-matrix/internal/features"
)

const (
	errNotPowerLevel  = "managed resource is not a PowerLevel custom resource"
	errTrackPCUsage   = "cannot track ProviderConfig usage"
	errGetPC          = "cannot get ProviderConfig"
	errGetCreds       = "cannot get credentials"
	errNewClient      = "cannot create new Matrix client"
	errSetPowerLevels = "cannot set Matrix power levels"
	errGetPowerLevels = "cannot get Matrix power levels"
)

// Setup adds a controller that reconciles PowerLevel managed resources.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.PowerLevelKind)

	cps := []managed.ConnectionPublisher{managed.NewAPISecretPublisher(mgr.GetClient(), mgr.GetScheme())}
	if o.Features.Enabled(features.EnableAlphaExternalSecretStores) {
		cps = append(cps, connection.NewDetailsManager(mgr.GetClient(), v1alpha1.PowerLevelGroupVersionKind))
	}

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.PowerLevelGroupVersionKind),
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
		For(&v1alpha1.PowerLevel{}).
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
	cr, ok := mg.(*v1alpha1.PowerLevel)
	if !ok {
		return nil, errors.New(errNotPowerLevel)
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
	cr, ok := mg.(*v1alpha1.PowerLevel)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotPowerLevel)
	}

	roomID := cr.Spec.ForProvider.RoomID
	powerLevels, err := c.service.GetPowerLevels(ctx, roomID)
	if err != nil {
		if clients.IsNotFound(err) {
			return managed.ExternalObservation{
				ResourceExists: false,
			}, nil
		}
		return managed.ExternalObservation{}, errors.Wrap(err, errGetPowerLevels)
	}

	cr.Status.AtProvider = generatePowerLevelObservation(roomID, powerLevels)
	cr.Status.SetConditions(xpv1.Available())

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: isPowerLevelUpToDate(cr, powerLevels),
	}, nil
}

func (c *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.PowerLevel)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotPowerLevel)
	}

	powerLevelSpec := generatePowerLevelSpec(cr)
	err := c.service.SetPowerLevels(ctx, cr.Spec.ForProvider.RoomID, powerLevelSpec)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errSetPowerLevels)
	}

	// Use room ID as external name since power levels are bound to a room
	meta.SetExternalName(cr, cr.Spec.ForProvider.RoomID)

	return managed.ExternalCreation{}, nil
}

func (c *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.PowerLevel)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotPowerLevel)
	}

	powerLevelSpec := generatePowerLevelSpec(cr)
	err := c.service.SetPowerLevels(ctx, cr.Spec.ForProvider.RoomID, powerLevelSpec)
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errSetPowerLevels)
	}

	return managed.ExternalUpdate{}, nil
}

func (c *external) Delete(ctx context.Context, mg resource.Managed) error {
	// Power levels cannot be deleted, only reset to defaults
	// For now, we'll just mark the resource as successfully deleted
	return nil
}

// Helper functions

func generatePowerLevelSpec(cr *v1alpha1.PowerLevel) *clients.PowerLevelSpec {
	spec := &clients.PowerLevelSpec{
		RoomID: cr.Spec.ForProvider.RoomID,
		PowerLevels: &clients.PowerLevelContent{
			Users:  cr.Spec.ForProvider.Users,
			Events: cr.Spec.ForProvider.Events,
		},
	}

	if cr.Spec.ForProvider.EventsDefault != nil {
		spec.PowerLevels.EventsDefault = cr.Spec.ForProvider.EventsDefault
	}
	if cr.Spec.ForProvider.StateDefault != nil {
		spec.PowerLevels.StateDefault = cr.Spec.ForProvider.StateDefault
	}
	if cr.Spec.ForProvider.UsersDefault != nil {
		spec.PowerLevels.UsersDefault = cr.Spec.ForProvider.UsersDefault
	}
	if cr.Spec.ForProvider.Ban != nil {
		spec.PowerLevels.Ban = cr.Spec.ForProvider.Ban
	}
	if cr.Spec.ForProvider.Kick != nil {
		spec.PowerLevels.Kick = cr.Spec.ForProvider.Kick
	}
	if cr.Spec.ForProvider.Redact != nil {
		spec.PowerLevels.Redact = cr.Spec.ForProvider.Redact
	}
	if cr.Spec.ForProvider.Invite != nil {
		spec.PowerLevels.Invite = cr.Spec.ForProvider.Invite
	}

	return spec
}

func generatePowerLevelObservation(roomID string, powerLevels *clients.PowerLevelContent) v1alpha1.PowerLevelObservation {
	obs := v1alpha1.PowerLevelObservation{
		RoomID:       roomID,
		Users:        powerLevels.Users,
		Events:       powerLevels.Events,
		LastModified: &metav1.Time{Time: time.Now()},
	}

	if powerLevels.EventsDefault != nil {
		obs.EventsDefault = *powerLevels.EventsDefault
	}
	if powerLevels.StateDefault != nil {
		obs.StateDefault = *powerLevels.StateDefault
	}
	if powerLevels.UsersDefault != nil {
		obs.UsersDefault = *powerLevels.UsersDefault
	}
	if powerLevels.Ban != nil {
		obs.Ban = *powerLevels.Ban
	}
	if powerLevels.Kick != nil {
		obs.Kick = *powerLevels.Kick
	}
	if powerLevels.Redact != nil {
		obs.Redact = *powerLevels.Redact
	}
	if powerLevels.Invite != nil {
		obs.Invite = *powerLevels.Invite
	}

	return obs
}

func isPowerLevelUpToDate(cr *v1alpha1.PowerLevel, powerLevels *clients.PowerLevelContent) bool {
	// Check user power levels
	if len(cr.Spec.ForProvider.Users) != len(powerLevels.Users) {
		return false
	}
	for userID, level := range cr.Spec.ForProvider.Users {
		if actualLevel, exists := powerLevels.Users[userID]; !exists || actualLevel != level {
			return false
		}
	}

	// Check event power levels
	if len(cr.Spec.ForProvider.Events) != len(powerLevels.Events) {
		return false
	}
	for eventType, level := range cr.Spec.ForProvider.Events {
		if actualLevel, exists := powerLevels.Events[eventType]; !exists || actualLevel != level {
			return false
		}
	}

	// Check default levels
	if cr.Spec.ForProvider.EventsDefault != nil && powerLevels.EventsDefault != nil {
		if *cr.Spec.ForProvider.EventsDefault != *powerLevels.EventsDefault {
			return false
		}
	}
	if cr.Spec.ForProvider.StateDefault != nil && powerLevels.StateDefault != nil {
		if *cr.Spec.ForProvider.StateDefault != *powerLevels.StateDefault {
			return false
		}
	}
	if cr.Spec.ForProvider.UsersDefault != nil && powerLevels.UsersDefault != nil {
		if *cr.Spec.ForProvider.UsersDefault != *powerLevels.UsersDefault {
			return false
		}
	}
	if cr.Spec.ForProvider.Ban != nil && powerLevels.Ban != nil {
		if *cr.Spec.ForProvider.Ban != *powerLevels.Ban {
			return false
		}
	}
	if cr.Spec.ForProvider.Kick != nil && powerLevels.Kick != nil {
		if *cr.Spec.ForProvider.Kick != *powerLevels.Kick {
			return false
		}
	}
	if cr.Spec.ForProvider.Redact != nil && powerLevels.Redact != nil {
		if *cr.Spec.ForProvider.Redact != *powerLevels.Redact {
			return false
		}
	}
	if cr.Spec.ForProvider.Invite != nil && powerLevels.Invite != nil {
		if *cr.Spec.ForProvider.Invite != *powerLevels.Invite {
			return false
		}
	}

	return true
}
