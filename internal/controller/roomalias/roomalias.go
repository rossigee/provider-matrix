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

package roomalias

import (
	"context"
	"time"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/v2/pkg/controller"
	"github.com/crossplane/crossplane-runtime/v2/pkg/event"
	"github.com/crossplane/crossplane-runtime/v2/pkg/meta"
	"github.com/crossplane/crossplane-runtime/v2/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/v2/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/v2/pkg/resource"

	"github.com/crossplane-contrib/provider-matrix/apis/roomalias/v1alpha1"
	apisv1beta1 "github.com/crossplane-contrib/provider-matrix/apis/v1beta1"
	"github.com/crossplane-contrib/provider-matrix/internal/clients"
)

const (
	errNotRoomAlias    = "managed resource is not a RoomAlias custom resource"
	errTrackPCUsage    = "cannot track ProviderConfig usage"
	errGetPC           = "cannot get ProviderConfig"
	errGetCreds        = "cannot get credentials"
	errNewClient       = "cannot create new Matrix client"
	errCreateRoomAlias = "cannot create Matrix room alias"
	errGetRoomAlias    = "cannot get Matrix room alias"
	errDeleteRoomAlias = "cannot delete Matrix room alias"
)

// Setup adds a controller that reconciles RoomAlias managed resources.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.RoomAliasKind)

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.RoomAliasGroupVersionKind),
		managed.WithExternalConnecter(&connector{
			kube:         mgr.GetClient(),
			usage:        clients.NewProviderConfigUsageTracker(mgr.GetClient()),
			newServiceFn: clients.NewClient,
		}),
		managed.WithLogger(o.Logger.WithValues("controller", name)),
		managed.WithPollInterval(o.PollInterval),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		WithEventFilter(resource.DesiredStateChanged()).
		For(&v1alpha1.RoomAlias{}).
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
	cr, ok := mg.(*v1alpha1.RoomAlias)
	if !ok {
		return nil, errors.New(errNotRoomAlias)
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
	cr, ok := mg.(*v1alpha1.RoomAlias)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotRoomAlias)
	}

	alias := cr.Spec.ForProvider.Alias
	roomAlias, err := c.service.GetRoomAlias(ctx, alias)
	if err != nil {
		if clients.IsNotFound(err) {
			return managed.ExternalObservation{
				ResourceExists: false,
			}, nil
		}
		return managed.ExternalObservation{}, errors.Wrap(err, errGetRoomAlias)
	}

	cr.Status.AtProvider = generateRoomAliasObservation(roomAlias)
	cr.Status.SetConditions(xpv1.Available())

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: isRoomAliasUpToDate(cr, roomAlias),
	}, nil
}

func (c *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.RoomAlias)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotRoomAlias)
	}

	alias := cr.Spec.ForProvider.Alias
	roomID := cr.Spec.ForProvider.RoomID

	err := c.service.CreateRoomAlias(ctx, alias, roomID)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreateRoomAlias)
	}

	meta.SetExternalName(cr, alias)

	return managed.ExternalCreation{}, nil
}

func (c *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	// Room aliases cannot be updated, only recreated
	// If the room ID changes, we need to delete and recreate
	cr, ok := mg.(*v1alpha1.RoomAlias)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotRoomAlias)
	}

	alias := cr.Spec.ForProvider.Alias
	roomID := cr.Spec.ForProvider.RoomID

	// Delete existing alias
	err := c.service.DeleteRoomAlias(ctx, alias)
	if err != nil && !clients.IsNotFound(err) {
		return managed.ExternalUpdate{}, errors.Wrap(err, errDeleteRoomAlias)
	}

	// Create with new room ID
	err = c.service.CreateRoomAlias(ctx, alias, roomID)
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errCreateRoomAlias)
	}

	return managed.ExternalUpdate{}, nil
}

func (c *external) Delete(ctx context.Context, mg resource.Managed) (managed.ExternalDelete, error) {
	cr, ok := mg.(*v1alpha1.RoomAlias)
	if !ok {
		return managed.ExternalDelete{}, errors.New(errNotRoomAlias)
	}

	alias := meta.GetExternalName(cr)
	if alias == "" {
		alias = cr.Spec.ForProvider.Alias
	}

	if alias == "" {
		return managed.ExternalDelete{}, nil
	}

	return managed.ExternalDelete{}, errors.Wrap(c.service.DeleteRoomAlias(ctx, alias), errDeleteRoomAlias)
}

// Disconnect closes the external client.
func (c *external) Disconnect(ctx context.Context) error {
	return nil // No special disconnect logic needed
}

// Helper functions

func generateRoomAliasObservation(roomAlias *clients.RoomAlias) v1alpha1.RoomAliasObservation {
	obs := v1alpha1.RoomAliasObservation{
		Alias:        roomAlias.Alias,
		RoomID:       roomAlias.RoomID,
		IsCanonical:  false, // This would need to be determined by checking room state
		IsPublished:  true,  // Assume published if alias exists
		CreationTime: &metav1.Time{Time: time.Now()},
		Servers:      []string{}, // Would need to be extracted from resolve response
	}

	return obs
}

func isRoomAliasUpToDate(cr *v1alpha1.RoomAlias, roomAlias *clients.RoomAlias) bool {
	// Check if the alias points to the correct room
	if cr.Spec.ForProvider.RoomID != roomAlias.RoomID {
		return false
	}

	return true
}
