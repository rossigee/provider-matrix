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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// A ProviderConfigSpec defines the desired state of a ProviderConfig.
type ProviderConfigSpec struct {
	// Credentials required to authenticate to this provider.
	Credentials ProviderCredentials `json:"credentials"`

	// HomeserverURL is the base URL for the Matrix homeserver.
	// Examples: https://matrix.example.com, https://dendrite.example.com
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern="^https?://.*"
	HomeserverURL string `json:"homeserverURL"`

	// AdminAPIURL is the base URL for admin API operations (if different from homeserver URL).
	// Used for Synapse admin API operations. If not set, will use HomeserverURL.
	// +kubebuilder:validation:Pattern="^https?://.*"
	AdminAPIURL *string `json:"adminAPIURL,omitempty"`

	// UserID is the Matrix user ID for the provider.
	// Format: @localpart:domain
	// +kubebuilder:validation:Pattern="^@[a-zA-Z0-9._=/-]+:[a-zA-Z0-9.-]+$"
	UserID *string `json:"userID,omitempty"`

	// DeviceID is the device ID to use for authentication.
	DeviceID *string `json:"deviceID,omitempty"`

	// ServerType indicates the type of Matrix server (for API compatibility).
	// +kubebuilder:validation:Enum=synapse;dendrite;conduit;auto
	// +kubebuilder:default="auto"
	ServerType *string `json:"serverType,omitempty"`

	// AdminMode enables administrative operations when supported.
	// +kubebuilder:default=false
	AdminMode *bool `json:"adminMode,omitempty"`
}

// ProviderCredentials required to authenticate.
type ProviderCredentials struct {
	// Source of the provider credentials.
	// +kubebuilder:validation:Enum=Secret;InjectedIdentity;Environment;Filesystem
	Source xpv1.CredentialsSource `json:"source"`

	xpv1.CommonCredentialSelectors `json:",inline"`
}

// A ProviderConfigStatus reflects the observed state of a ProviderConfig.
type ProviderConfigStatus struct {
	xpv1.ProviderConfigStatus `json:",inline"`
}

// +kubebuilder:object:root=true

// A ProviderConfig configures a Matrix provider.
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="SECRET-NAME",type="string",JSONPath=".spec.credentials.secretRef.name",priority=1
type ProviderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProviderConfigSpec   `json:"spec"`
	Status ProviderConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProviderConfigList contains a list of ProviderConfig.
type ProviderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProviderConfig `json:"items"`
}

// +kubebuilder:object:root=true

// A ProviderConfigUsage indicates that a resource is using a ProviderConfig.
type ProviderConfigUsage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	xpv1.ProviderConfigUsage `json:",inline"`
}

// +kubebuilder:object:root=true

// ProviderConfigUsageList contains a list of ProviderConfigUsage
type ProviderConfigUsageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProviderConfigUsage `json:"items"`
}