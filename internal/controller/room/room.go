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

package room

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
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/crossplane-contrib/provider-matrix/apis/room/v1alpha1"
	apisv1beta1 "github.com/crossplane-contrib/provider-matrix/apis/v1beta1"
	"github.com/crossplane-contrib/provider-matrix/internal/clients"
	"github.com/crossplane-contrib/provider-matrix/internal/features"
)

const (
	errNotRoom       = "managed resource is not a Room custom resource"
	errTrackPCUsage  = "cannot track ProviderConfig usage"
	errGetPC         = "cannot get ProviderConfig"
	errGetCreds      = "cannot get credentials"
	errNewClient     = "cannot create new Matrix client"
	errCreateRoom    = "cannot create Matrix room"
	errGetRoom       = "cannot get Matrix room"
	errUpdateRoom    = "cannot update Matrix room"
	errDeleteRoom    = "cannot delete Matrix room"
)

// Setup adds a controller that reconciles Room managed resources.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.RoomGroupKind)

	cps := []managed.ConnectionPublisher{managed.NewAPISecretPublisher(mgr.GetClient(), mgr.GetScheme())}
	if o.Features.Enabled(features.EnableAlphaExternalSecretStores) {
		cps = append(cps, connection.NewDetailsManager(mgr.GetClient(), apisv1beta1.StoreConfigGroupVersionKind))
	}

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.RoomGroupVersionKind),
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
		For(&v1alpha1.Room{}).
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
	cr, ok := mg.(*v1alpha1.Room)
	if !ok {
		return nil, errors.New(errNotRoom)
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
	cr, ok := mg.(*v1alpha1.Room)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotRoom)
	}

	roomID := cr.GetAnnotations()[resource.AnnotationKeyExternalName]
	if roomID == "" {
		return managed.ExternalObservation{
			ResourceExists: false,
		}, nil
	}

	room, err := c.service.GetRoom(ctx, roomID)
	if err != nil {
		if clients.IsNotFound(err) {
			return managed.ExternalObservation{
				ResourceExists: false,
			}, nil
		}
		return managed.ExternalObservation{}, errors.Wrap(err, errGetRoom)
	}

	cr.Status.AtProvider = generateRoomObservation(room)
	cr.Status.SetConditions(xpv1.Available())

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: isRoomUpToDate(cr, room),
	}, nil
}

func (c *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Room)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotRoom)
	}

	roomSpec := generateRoomSpec(cr)
	room, err := c.service.CreateRoom(ctx, roomSpec)
	if err != nil {
		return managed.ExternalCreation{}, errors.Wrap(err, errCreateRoom)
	}

	cr.SetAnnotations(map[string]string{
		resource.AnnotationKeyExternalName: room.RoomID,
	})

	return managed.ExternalCreation{
		ExternalNameAssigned: true,
	}, nil
}

func (c *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.Room)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotRoom)
	}

	roomID := cr.GetAnnotations()[resource.AnnotationKeyExternalName]
	roomSpec := generateRoomSpec(cr)
	_, err := c.service.UpdateRoom(ctx, roomID, roomSpec)
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errUpdateRoom)
	}

	return managed.ExternalUpdate{}, nil
}

func (c *external) Delete(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*v1alpha1.Room)
	if !ok {
		return errors.New(errNotRoom)
	}

	roomID := cr.GetAnnotations()[resource.AnnotationKeyExternalName]
	if roomID == "" {
		return nil
	}

	return errors.Wrap(c.service.DeleteRoom(ctx, roomID), errDeleteRoom)
}

// Helper functions

func generateRoomSpec(cr *v1alpha1.Room) *clients.RoomSpec {
	spec := &clients.RoomSpec{}

	if cr.Spec.ForProvider.Name != nil {
		spec.Name = *cr.Spec.ForProvider.Name
	}
	if cr.Spec.ForProvider.Topic != nil {
		spec.Topic = *cr.Spec.ForProvider.Topic
	}
	if cr.Spec.ForProvider.Alias != nil {
		spec.Alias = *cr.Spec.ForProvider.Alias
	}
	if cr.Spec.ForProvider.Preset != nil {
		spec.Preset = *cr.Spec.ForProvider.Preset
	}
	if cr.Spec.ForProvider.Visibility != nil {
		spec.Visibility = *cr.Spec.ForProvider.Visibility
	}
	if cr.Spec.ForProvider.RoomVersion != nil {
		spec.RoomVersion = *cr.Spec.ForProvider.RoomVersion
	}

	spec.CreationContent = cr.Spec.ForProvider.CreationContent
	spec.Invite = cr.Spec.ForProvider.Invite

	// Convert initial state
	for _, state := range cr.Spec.ForProvider.InitialState {
		spec.InitialState = append(spec.InitialState, clients.StateEvent{
			Type:     state.Type,
			StateKey: state.StateKey,
			Content:  state.Content,
		})
	}

	// Convert power level overrides
	if cr.Spec.ForProvider.PowerLevelOverrides != nil {
		spec.PowerLevelOverrides = &clients.PowerLevelContent{
			Users:         cr.Spec.ForProvider.PowerLevelOverrides.Users,
			Events:        cr.Spec.ForProvider.PowerLevelOverrides.Events,
			EventsDefault: cr.Spec.ForProvider.PowerLevelOverrides.EventsDefault,
			StateDefault:  cr.Spec.ForProvider.PowerLevelOverrides.StateDefault,
			UsersDefault:  cr.Spec.ForProvider.PowerLevelOverrides.UsersDefault,
			Ban:           cr.Spec.ForProvider.PowerLevelOverrides.Ban,
			Kick:          cr.Spec.ForProvider.PowerLevelOverrides.Kick,
			Redact:        cr.Spec.ForProvider.PowerLevelOverrides.Redact,
			Invite:        cr.Spec.ForProvider.PowerLevelOverrides.Invite,
		}
	}

	if cr.Spec.ForProvider.GuestAccess != nil {
		spec.GuestAccess = *cr.Spec.ForProvider.GuestAccess
	}
	if cr.Spec.ForProvider.HistoryVisibility != nil {
		spec.HistoryVisibility = *cr.Spec.ForProvider.HistoryVisibility
	}
	if cr.Spec.ForProvider.JoinRules != nil {
		spec.JoinRules = *cr.Spec.ForProvider.JoinRules
	}
	if cr.Spec.ForProvider.EncryptionEnabled != nil {
		spec.EncryptionEnabled = *cr.Spec.ForProvider.EncryptionEnabled
	}
	if cr.Spec.ForProvider.AvatarURL != nil {
		spec.AvatarURL = *cr.Spec.ForProvider.AvatarURL
	}

	return spec
}

func generateRoomObservation(room *clients.Room) v1alpha1.RoomObservation {
	obs := v1alpha1.RoomObservation{
		RoomID:            room.RoomID,
		Name:              room.Name,
		Topic:             room.Topic,
		Alias:             room.Alias,
		AvatarURL:         room.AvatarURL,
		Creator:           room.Creator,
		RoomVersion:       room.RoomVersion,
		JoinedMembers:     room.JoinedMembers,
		InvitedMembers:    room.InvitedMembers,
		Visibility:        room.Visibility,
		GuestAccess:       room.GuestAccess,
		HistoryVisibility: room.HistoryVisibility,
		JoinRules:         room.JoinRules,
		EncryptionEnabled: room.EncryptionEnabled,
	}

	if room.CreationTime != nil {
		obs.CreationTime = &metav1.Time{Time: *room.CreationTime}
	}

	// Convert state events
	for _, state := range room.State {
		obs.State = append(obs.State, v1alpha1.StateEvent{
			Type:     state.Type,
			StateKey: state.StateKey,
			Content:  state.Content,
		})
	}

	// Convert power levels
	if room.PowerLevels != nil {
		obs.PowerLevels = &v1alpha1.PowerLevelContent{
			Users:         room.PowerLevels.Users,
			Events:        room.PowerLevels.Events,
			EventsDefault: room.PowerLevels.EventsDefault,
			StateDefault:  room.PowerLevels.StateDefault,
			UsersDefault:  room.PowerLevels.UsersDefault,
			Ban:           room.PowerLevels.Ban,
			Kick:          room.PowerLevels.Kick,
			Redact:        room.PowerLevels.Redact,
			Invite:        room.PowerLevels.Invite,
		}
	}

	return obs
}

func isRoomUpToDate(cr *v1alpha1.Room, room *clients.Room) bool {
	// Check name
	if cr.Spec.ForProvider.Name != nil && *cr.Spec.ForProvider.Name != room.Name {
		return false
	}

	// Check topic
	if cr.Spec.ForProvider.Topic != nil && *cr.Spec.ForProvider.Topic != room.Topic {
		return false
	}

	// Check alias
	if cr.Spec.ForProvider.Alias != nil && *cr.Spec.ForProvider.Alias != room.Alias {
		return false
	}

	// Check guest access
	if cr.Spec.ForProvider.GuestAccess != nil && *cr.Spec.ForProvider.GuestAccess != room.GuestAccess {
		return false
	}

	// Check history visibility
	if cr.Spec.ForProvider.HistoryVisibility != nil && *cr.Spec.ForProvider.HistoryVisibility != room.HistoryVisibility {
		return false
	}

	// Check join rules
	if cr.Spec.ForProvider.JoinRules != nil && *cr.Spec.ForProvider.JoinRules != room.JoinRules {
		return false
	}

	// Check encryption
	if cr.Spec.ForProvider.EncryptionEnabled != nil && *cr.Spec.ForProvider.EncryptionEnabled != room.EncryptionEnabled {
		return false
	}

	// Check avatar URL
	if cr.Spec.ForProvider.AvatarURL != nil && *cr.Spec.ForProvider.AvatarURL != room.AvatarURL {
		return false
	}

	return true
}