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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
)

// RoomParameters define the desired state of a Matrix Room
type RoomParameters struct {
	// Name is the human-readable name for the room
	Name *string `json:"name,omitempty"`

	// Topic is the topic/description for the room
	Topic *string `json:"topic,omitempty"`

	// Alias is the room alias (e.g., #example:matrix.org)
	// +kubebuilder:validation:Pattern="^#[a-zA-Z0-9._=/-]+:[a-zA-Z0-9.-]+$"
	Alias *string `json:"alias,omitempty"`

	// Preset determines the room's configuration template
	// +kubebuilder:validation:Enum=private_chat;public_chat;trusted_private_chat
	// +kubebuilder:default="private_chat"
	Preset *string `json:"preset,omitempty"`

	// Visibility controls room visibility in the directory
	// +kubebuilder:validation:Enum=public;private
	// +kubebuilder:default="private"
	Visibility *string `json:"visibility,omitempty"`

	// RoomVersion specifies the Matrix room version to use
	// +kubebuilder:validation:Pattern="^[0-9]+$|^[0-9]+\.[0-9]+$"
	RoomVersion *string `json:"roomVersion,omitempty"`

	// CreationContent is additional content for the m.room.create event
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Type=object
	CreationContent *runtime.RawExtension `json:"creationContent,omitempty"`

	// InitialState is a list of state events to set in the new room
	InitialState []StateEvent `json:"initialState,omitempty"`

	// Invite is a list of user IDs to invite to the room
	Invite []string `json:"invite,omitempty"`

	// PowerLevelOverrides allows customizing power levels for the room
	PowerLevelOverrides *PowerLevelContent `json:"powerLevelOverrides,omitempty"`

	// GuestAccess controls whether guests can join the room
	// +kubebuilder:validation:Enum=can_join;forbidden
	// +kubebuilder:default="forbidden"
	GuestAccess *string `json:"guestAccess,omitempty"`

	// HistoryVisibility controls message history visibility
	// +kubebuilder:validation:Enum=invited;joined;shared;world_readable
	// +kubebuilder:default="shared"
	HistoryVisibility *string `json:"historyVisibility,omitempty"`

	// JoinRules controls who can join the room
	// +kubebuilder:validation:Enum=public;invite;restricted;knock
	// +kubebuilder:default="invite"
	JoinRules *string `json:"joinRules,omitempty"`

	// EncryptionEnabled indicates if the room should be encrypted
	// +kubebuilder:default=false
	EncryptionEnabled *bool `json:"encryptionEnabled,omitempty"`

	// AvatarURL is the room's avatar image URL (mxc:// URL)
	// +kubebuilder:validation:Pattern="^mxc://.*"
	AvatarURL *string `json:"avatarURL,omitempty"`
}

// StateEvent represents a Matrix state event
type StateEvent struct {
	// Type is the event type
	// +kubebuilder:validation:Required
	Type string `json:"type"`

	// StateKey is the state key for the event
	StateKey string `json:"stateKey"`

	// Content is the event content
	// +kubebuilder:validation:Required
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Type=object
	Content runtime.RawExtension `json:"content"`
}

// PowerLevelContent defines power levels for room events and users
type PowerLevelContent struct {
	// Users maps user IDs to their power levels
	Users map[string]int `json:"users,omitempty"`

	// Events maps event types to required power levels
	Events map[string]int `json:"events,omitempty"`

	// EventsDefault is the default power level for events
	EventsDefault *int `json:"eventsDefault,omitempty"`

	// StateDefault is the default power level for state events
	StateDefault *int `json:"stateDefault,omitempty"`

	// UsersDefault is the default power level for users
	UsersDefault *int `json:"usersDefault,omitempty"`

	// Ban is the power level required to ban users
	Ban *int `json:"ban,omitempty"`

	// Kick is the power level required to kick users
	Kick *int `json:"kick,omitempty"`

	// Redact is the power level required to redact events
	Redact *int `json:"redact,omitempty"`

	// Invite is the power level required to invite users
	Invite *int `json:"invite,omitempty"`
}

// RoomObservation reflects the observed state of a Matrix Room
type RoomObservation struct {
	// RoomID is the Matrix room ID
	RoomID string `json:"roomID,omitempty"`

	// Name is the current room name
	Name string `json:"name,omitempty"`

	// Topic is the current room topic
	Topic string `json:"topic,omitempty"`

	// Alias is the canonical room alias
	Alias string `json:"alias,omitempty"`

	// AvatarURL is the current room avatar URL
	AvatarURL string `json:"avatarURL,omitempty"`

	// Creator is the user ID of the room creator
	Creator string `json:"creator,omitempty"`

	// CreationTime is when the room was created
	CreationTime *metav1.Time `json:"creationTime,omitempty"`

	// RoomVersion is the room version
	RoomVersion string `json:"roomVersion,omitempty"`

	// JoinedMembers is the number of joined members
	JoinedMembers int `json:"joinedMembers,omitempty"`

	// InvitedMembers is the number of invited members
	InvitedMembers int `json:"invitedMembers,omitempty"`

	// Visibility is the current room visibility
	Visibility string `json:"visibility,omitempty"`

	// GuestAccess is the current guest access setting
	GuestAccess string `json:"guestAccess,omitempty"`

	// HistoryVisibility is the current history visibility setting
	HistoryVisibility string `json:"historyVisibility,omitempty"`

	// JoinRules is the current join rules setting
	JoinRules string `json:"joinRules,omitempty"`

	// EncryptionEnabled indicates if the room is encrypted
	EncryptionEnabled bool `json:"encryptionEnabled,omitempty"`

	// State contains current room state events
	State []StateEvent `json:"state,omitempty"`

	// PowerLevels contains current power level settings
	PowerLevels *PowerLevelContent `json:"powerLevels,omitempty"`
}

// A RoomSpec defines the desired state of a Room.
type RoomSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       RoomParameters `json:"forProvider"`
}

// A RoomStatus represents the observed state of a Room.
type RoomStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          RoomObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Room is a managed resource that represents a Matrix Room
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,matrix}
type Room struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RoomSpec   `json:"spec"`
	Status RoomStatus `json:"status,omitempty"`
}

// GetProviderConfigReference returns the provider config reference.
func (r *Room) GetProviderConfigReference() *xpv1.Reference {
	return r.Spec.ProviderConfigReference
}

// SetProviderConfigReference sets the provider config reference.
func (r *Room) SetProviderConfigReference(ref *xpv1.Reference) {
	r.Spec.ProviderConfigReference = ref
}

// GetCondition returns the condition with the given type.
func (r *Room) GetCondition(ct xpv1.ConditionType) xpv1.Condition {
	return r.Status.GetCondition(ct)
}

// SetConditions sets the conditions.
func (r *Room) SetConditions(c ...xpv1.Condition) {
	r.Status.SetConditions(c...)
}

// GetDeletionPolicy returns the deletion policy.
func (r *Room) GetDeletionPolicy() xpv1.DeletionPolicy {
	return r.Spec.DeletionPolicy
}

// SetDeletionPolicy sets the deletion policy.
func (r *Room) SetDeletionPolicy(p xpv1.DeletionPolicy) {
	r.Spec.DeletionPolicy = p
}

// GetManagementPolicies returns the management policies.
func (r *Room) GetManagementPolicies() xpv1.ManagementPolicies {
	return r.Spec.ManagementPolicies
}

// SetManagementPolicies sets the management policies.
func (r *Room) SetManagementPolicies(p xpv1.ManagementPolicies) {
	r.Spec.ManagementPolicies = p
}



// GetWriteConnectionSecretToReference returns the write connection secret to reference.
func (r *Room) GetWriteConnectionSecretToReference() *xpv1.SecretReference {
	return r.Spec.WriteConnectionSecretToReference
}

// SetWriteConnectionSecretToReference sets the write connection secret to reference.
func (r *Room) SetWriteConnectionSecretToReference(s *xpv1.SecretReference) {
	r.Spec.WriteConnectionSecretToReference = s
}

// +kubebuilder:object:root=true

// RoomList contains a list of Room
type RoomList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Room `json:"items"`
}
