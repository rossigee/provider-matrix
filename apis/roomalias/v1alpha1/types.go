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

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// RoomAliasParameters define the desired state of a Matrix Room Alias
type RoomAliasParameters struct {
	// Alias is the room alias to create (e.g., #example:matrix.org)
	// +kubebuilder:validation:Pattern="^#[a-zA-Z0-9._=/-]+:[a-zA-Z0-9.-]+$"
	// +kubebuilder:validation:Required
	Alias string `json:"alias"`

	// RoomID is the Matrix room ID that this alias should point to
	// +kubebuilder:validation:Pattern="^![a-zA-Z0-9]+:[a-zA-Z0-9.-]+$"
	// +kubebuilder:validation:Required
	RoomID string `json:"roomID"`

	// SetAsCanonical determines if this alias should be set as the canonical alias for the room
	// +kubebuilder:default=false
	SetAsCanonical *bool `json:"setAsCanonical,omitempty"`

	// AltAliases is a list of alternative aliases to publish for the room
	AltAliases []string `json:"altAliases,omitempty"`
}

// RoomAliasObservation reflects the observed state of a Matrix Room Alias
type RoomAliasObservation struct {
	// Alias is the room alias
	Alias string `json:"alias,omitempty"`

	// RoomID is the Matrix room ID that this alias points to
	RoomID string `json:"roomID,omitempty"`

	// IsCanonical indicates if this is the canonical alias for the room
	IsCanonical bool `json:"isCanonical,omitempty"`

	// IsPublished indicates if this alias is published in the room directory
	IsPublished bool `json:"isPublished,omitempty"`

	// CreationTime is when the alias was created
	CreationTime *metav1.Time `json:"creationTime,omitempty"`

	// Servers is a list of servers that know about this alias
	Servers []string `json:"servers,omitempty"`
}

// A RoomAliasSpec defines the desired state of a RoomAlias.
type RoomAliasSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       RoomAliasParameters `json:"forProvider"`
}

// A RoomAliasStatus represents the observed state of a RoomAlias.
type RoomAliasStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          RoomAliasObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A RoomAlias is a managed resource that represents a Matrix Room Alias
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="ALIAS",type="string",JSONPath=".spec.forProvider.alias"
// +kubebuilder:printcolumn:name="ROOM-ID",type="string",JSONPath=".spec.forProvider.roomID"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,matrix}
type RoomAlias struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RoomAliasSpec   `json:"spec"`
	Status RoomAliasStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RoomAliasList contains a list of RoomAlias
type RoomAliasList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RoomAlias `json:"items"`
}
