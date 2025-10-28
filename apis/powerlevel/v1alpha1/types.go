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

	xpv1 "github.com/crossplane/crossplane-runtime/v2/apis/common/v1"
)

// PowerLevelParameters define the desired state of room power levels
type PowerLevelParameters struct {
	// RoomID is the Matrix room ID to manage power levels for
	// +kubebuilder:validation:Pattern="^![a-zA-Z0-9]+:[a-zA-Z0-9.-]+$"
	// +kubebuilder:validation:Required
	RoomID string `json:"roomID"`

	// Users maps user IDs to their power levels in the room
	Users map[string]int `json:"users,omitempty"`

	// Events maps event types to required power levels
	Events map[string]int `json:"events,omitempty"`

	// EventsDefault is the default power level required to send events
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:default=0
	EventsDefault *int `json:"eventsDefault,omitempty"`

	// StateDefault is the default power level required to send state events
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:default=50
	StateDefault *int `json:"stateDefault,omitempty"`

	// UsersDefault is the default power level for users in the room
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:default=0
	UsersDefault *int `json:"usersDefault,omitempty"`

	// Ban is the power level required to ban users
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:default=50
	Ban *int `json:"ban,omitempty"`

	// Kick is the power level required to kick users
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:default=50
	Kick *int `json:"kick,omitempty"`

	// Redact is the power level required to redact events
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:default=50
	Redact *int `json:"redact,omitempty"`

	// Invite is the power level required to invite users
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	// +kubebuilder:default=0
	Invite *int `json:"invite,omitempty"`
}

// PowerLevelObservation reflects the observed state of room power levels
type PowerLevelObservation struct {
	// RoomID is the Matrix room ID
	RoomID string `json:"roomID,omitempty"`

	// Users contains the current user power levels
	Users map[string]int `json:"users,omitempty"`

	// Events contains the current event type power levels
	Events map[string]int `json:"events,omitempty"`

	// EventsDefault is the current default power level for events
	EventsDefault int `json:"eventsDefault,omitempty"`

	// StateDefault is the current default power level for state events
	StateDefault int `json:"stateDefault,omitempty"`

	// UsersDefault is the current default power level for users
	UsersDefault int `json:"usersDefault,omitempty"`

	// Ban is the current power level required to ban users
	Ban int `json:"ban,omitempty"`

	// Kick is the current power level required to kick users
	Kick int `json:"kick,omitempty"`

	// Redact is the current power level required to redact events
	Redact int `json:"redact,omitempty"`

	// Invite is the current power level required to invite users
	Invite int `json:"invite,omitempty"`

	// LastModified is when the power levels were last modified
	LastModified *metav1.Time `json:"lastModified,omitempty"`
}

// A PowerLevelSpec defines the desired state of a PowerLevel.
type PowerLevelSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       PowerLevelParameters `json:"forProvider"`
}

// A PowerLevelStatus represents the observed state of a PowerLevel.
type PowerLevelStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          PowerLevelObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A PowerLevel is a managed resource that represents Matrix room power levels
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="ROOM-ID",type="string",JSONPath=".spec.forProvider.roomID"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,matrix}
type PowerLevel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PowerLevelSpec   `json:"spec"`
	Status PowerLevelStatus `json:"status,omitempty"`
}

// GetProviderConfigReference returns the provider config reference.
func (p *PowerLevel) GetProviderConfigReference() *xpv1.Reference {
	return p.Spec.ProviderConfigReference
}

// SetProviderConfigReference sets the provider config reference.
func (p *PowerLevel) SetProviderConfigReference(ref *xpv1.Reference) {
	p.Spec.ProviderConfigReference = ref
}

// GetCondition returns the condition with the given type.
func (p *PowerLevel) GetCondition(ct xpv1.ConditionType) xpv1.Condition {
	return p.Status.GetCondition(ct)
}

// SetConditions sets the conditions.
func (p *PowerLevel) SetConditions(c ...xpv1.Condition) {
	p.Status.SetConditions(c...)
}

// GetDeletionPolicy returns the deletion policy.
func (p *PowerLevel) GetDeletionPolicy() xpv1.DeletionPolicy {
	return p.Spec.DeletionPolicy
}

// SetDeletionPolicy sets the deletion policy.
func (p *PowerLevel) SetDeletionPolicy(dp xpv1.DeletionPolicy) {
	p.Spec.DeletionPolicy = dp
}

// GetManagementPolicies returns the management policies.
func (p *PowerLevel) GetManagementPolicies() xpv1.ManagementPolicies {
	return p.Spec.ManagementPolicies
}

// SetManagementPolicies sets the management policies.
func (p *PowerLevel) SetManagementPolicies(mp xpv1.ManagementPolicies) {
	p.Spec.ManagementPolicies = mp
}



// GetWriteConnectionSecretToReference returns the write connection secret to reference.
func (p *PowerLevel) GetWriteConnectionSecretToReference() *xpv1.SecretReference {
	return p.Spec.WriteConnectionSecretToReference
}

// SetWriteConnectionSecretToReference sets the write connection secret to reference.
func (p *PowerLevel) SetWriteConnectionSecretToReference(s *xpv1.SecretReference) {
	p.Spec.WriteConnectionSecretToReference = s
}

// +kubebuilder:object:root=true

// PowerLevelList contains a list of PowerLevel
type PowerLevelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PowerLevel `json:"items"`
}
