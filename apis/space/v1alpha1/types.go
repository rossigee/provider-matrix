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

// SpaceParameters define the desired state of a Matrix Space
type SpaceParameters struct {
	// Name is the human-readable name for the space
	Name *string `json:"name,omitempty"`

	// Topic is the topic/description for the space
	Topic *string `json:"topic,omitempty"`

	// Alias is the space alias (e.g., #space:matrix.org)
	// +kubebuilder:validation:Pattern="^#[a-zA-Z0-9._=/-]+:[a-zA-Z0-9.-]+$"
	Alias *string `json:"alias,omitempty"`

	// Visibility controls space visibility in the directory
	// +kubebuilder:validation:Enum=public;private
	// +kubebuilder:default="private"
	Visibility *string `json:"visibility,omitempty"`

	// RoomVersion specifies the Matrix room version to use for the space
	// +kubebuilder:validation:Pattern="^[0-9]+$|^[0-9]+\.[0-9]+$"
	RoomVersion *string `json:"roomVersion,omitempty"`

	// CreationContent is additional content for the m.room.create event
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Type=object
	CreationContent *runtime.RawExtension `json:"creationContent,omitempty"`

	// InitialState is a list of state events to set in the new space
	InitialState []StateEvent `json:"initialState,omitempty"`

	// Invite is a list of user IDs to invite to the space
	Invite []string `json:"invite,omitempty"`

	// PowerLevelOverrides allows customizing power levels for the space
	PowerLevelOverrides *PowerLevelContent `json:"powerLevelOverrides,omitempty"`

	// GuestAccess controls whether guests can join the space
	// +kubebuilder:validation:Enum=can_join;forbidden
	// +kubebuilder:default="forbidden"
	GuestAccess *string `json:"guestAccess,omitempty"`

	// HistoryVisibility controls message history visibility
	// +kubebuilder:validation:Enum=invited;joined;shared;world_readable
	// +kubebuilder:default="shared"
	HistoryVisibility *string `json:"historyVisibility,omitempty"`

	// JoinRules controls who can join the space
	// +kubebuilder:validation:Enum=public;invite;restricted
	// +kubebuilder:default="invite"
	JoinRules *string `json:"joinRules,omitempty"`

	// AvatarURL is the space's avatar image URL (mxc:// URL)
	// +kubebuilder:validation:Pattern="^mxc://.*"
	AvatarURL *string `json:"avatarURL,omitempty"`

	// Children defines the child rooms and spaces within this space
	Children []SpaceChild `json:"children,omitempty"`
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

// PowerLevelContent defines power levels for space events and users
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

// SpaceChild represents a child room or space within a space
type SpaceChild struct {
	// RoomID is the Matrix room or space ID to include as a child
	// +kubebuilder:validation:Pattern="^![a-zA-Z0-9]+:[a-zA-Z0-9.-]+$"
	// +kubebuilder:validation:Required
	RoomID string `json:"roomID"`

	// Via is a list of servers that can be used to join the child
	Via []string `json:"via,omitempty"`

	// Order is used to sort children in the space
	Order *string `json:"order,omitempty"`

	// Suggested indicates if this child is a suggested room
	// +kubebuilder:default=false
	Suggested *bool `json:"suggested,omitempty"`
}

// SpaceObservation reflects the observed state of a Matrix Space
type SpaceObservation struct {
	// SpaceID is the Matrix space ID (also a room ID)
	SpaceID string `json:"spaceID,omitempty"`

	// Name is the current space name
	Name string `json:"name,omitempty"`

	// Topic is the current space topic
	Topic string `json:"topic,omitempty"`

	// Alias is the canonical space alias
	Alias string `json:"alias,omitempty"`

	// AvatarURL is the current space avatar URL
	AvatarURL string `json:"avatarURL,omitempty"`

	// Creator is the user ID of the space creator
	Creator string `json:"creator,omitempty"`

	// CreationTime is when the space was created
	CreationTime *metav1.Time `json:"creationTime,omitempty"`

	// RoomVersion is the room version used for the space
	RoomVersion string `json:"roomVersion,omitempty"`

	// JoinedMembers is the number of joined members
	JoinedMembers int `json:"joinedMembers,omitempty"`

	// InvitedMembers is the number of invited members
	InvitedMembers int `json:"invitedMembers,omitempty"`

	// Visibility is the current space visibility
	Visibility string `json:"visibility,omitempty"`

	// GuestAccess is the current guest access setting
	GuestAccess string `json:"guestAccess,omitempty"`

	// HistoryVisibility is the current history visibility setting
	HistoryVisibility string `json:"historyVisibility,omitempty"`

	// JoinRules is the current join rules setting
	JoinRules string `json:"joinRules,omitempty"`

	// Children contains the current child rooms and spaces
	Children []SpaceChild `json:"children,omitempty"`

	// State contains current space state events
	State []StateEvent `json:"state,omitempty"`

	// PowerLevels contains current power level settings
	PowerLevels *PowerLevelContent `json:"powerLevels,omitempty"`
}

// A SpaceSpec defines the desired state of a Space.
type SpaceSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       SpaceParameters `json:"forProvider"`
}

// A SpaceStatus represents the observed state of a Space.
type SpaceStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          SpaceObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Space is a managed resource that represents a Matrix Space
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,matrix}
type Space struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SpaceSpec   `json:"spec"`
	Status SpaceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SpaceList contains a list of Space
type SpaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Space `json:"items"`
}
