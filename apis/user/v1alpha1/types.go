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

// UserParameters define the desired state of a Matrix User
type UserParameters struct {
	// UserID is the Matrix user ID (e.g., @alice:example.com)
	// If not provided, will be generated from localpart and homeserver domain
	// +kubebuilder:validation:Pattern="^@[a-zA-Z0-9._=/-]+:[a-zA-Z0-9.-]+$"
	UserID *string `json:"userID,omitempty"`

	// Localpart is the local part of the Matrix user ID (before the @)
	// Required if UserID is not provided
	// +kubebuilder:validation:Pattern="^[a-zA-Z0-9._=/-]+$"
	Localpart *string `json:"localpart,omitempty"`

	// Password for the user account. Will be auto-generated if not provided.
	// Note: Use passwordSecretRef for secure password management
	Password *string `json:"password,omitempty"`

	// PasswordSecretRef references a Secret containing the user password
	PasswordSecretRef *xpv1.SecretKeySelector `json:"passwordSecretRef,omitempty"`

	// DisplayName is the user's display name
	DisplayName *string `json:"displayName,omitempty"`

	// AvatarURL is the user's avatar URL (mxc:// URL)
	// +kubebuilder:validation:Pattern="^mxc://.*"
	AvatarURL *string `json:"avatarURL,omitempty"`

	// Admin indicates if the user should have server admin privileges
	// +kubebuilder:default=false
	Admin *bool `json:"admin,omitempty"`

	// Deactivated indicates if the user account should be deactivated
	// +kubebuilder:default=false
	Deactivated *bool `json:"deactivated,omitempty"`

	// ExternalIDs are third-party identifiers (3PIDs) associated with the user
	ExternalIDs []ExternalID `json:"externalIDs,omitempty"`

	// UserType specifies the type of user account
	// +kubebuilder:validation:Enum=regular;guest;support
	// +kubebuilder:default="regular"
	UserType *string `json:"userType,omitempty"`

	// ExpireTime is when the user account expires (for guest users)
	ExpireTime *metav1.Time `json:"expireTime,omitempty"`
}

// ExternalID represents a third-party identifier associated with a user
type ExternalID struct {
	// Medium is the type of identifier (email, msisdn)
	// +kubebuilder:validation:Enum=email;msisdn
	Medium string `json:"medium"`

	// Address is the actual identifier value
	// +kubebuilder:validation:Required
	Address string `json:"address"`

	// Validated indicates if the identifier has been validated
	// +kubebuilder:default=false
	Validated *bool `json:"validated,omitempty"`
}

// UserObservation reflects the observed state of a Matrix User
type UserObservation struct {
	// UserID is the full Matrix user ID
	UserID string `json:"userID,omitempty"`

	// DisplayName is the current display name
	DisplayName string `json:"displayName,omitempty"`

	// AvatarURL is the current avatar URL
	AvatarURL string `json:"avatarURL,omitempty"`

	// Admin indicates if the user has admin privileges
	Admin bool `json:"admin,omitempty"`

	// Deactivated indicates if the user is deactivated
	Deactivated bool `json:"deactivated,omitempty"`

	// CreationTime is when the user was created
	CreationTime *metav1.Time `json:"creationTime,omitempty"`

	// LastSeenTime is when the user was last seen
	LastSeenTime *metav1.Time `json:"lastSeenTime,omitempty"`

	// Devices is a list of devices associated with the user
	Devices []Device `json:"devices,omitempty"`

	// ExternalIDs are the validated external identifiers
	ExternalIDs []ExternalID `json:"externalIDs,omitempty"`

	// UserType is the type of user account
	UserType string `json:"userType,omitempty"`

	// ShadowBanned indicates if the user is shadow banned
	ShadowBanned bool `json:"shadowBanned,omitempty"`
}

// Device represents a Matrix device
type Device struct {
	// DeviceID is the unique device identifier
	DeviceID string `json:"deviceID,omitempty"`

	// DisplayName is the device display name
	DisplayName string `json:"displayName,omitempty"`

	// LastSeenIP is the last IP address the device was seen from
	LastSeenIP string `json:"lastSeenIP,omitempty"`

	// LastSeenTime is when the device was last seen
	LastSeenTime *metav1.Time `json:"lastSeenTime,omitempty"`
}

// A UserSpec defines the desired state of a User.
type UserSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       UserParameters `json:"forProvider"`
}

// A UserStatus represents the observed state of a User.
type UserStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          UserObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A User is a managed resource that represents a Matrix User
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,matrix}
type User struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UserSpec   `json:"spec"`
	Status UserStatus `json:"status,omitempty"`
}

// GetProviderConfigReference returns the provider config reference.
func (u *User) GetProviderConfigReference() *xpv1.Reference {
	return u.Spec.ProviderConfigReference
}

// SetProviderConfigReference sets the provider config reference.
func (u *User) SetProviderConfigReference(ref *xpv1.Reference) {
	u.Spec.ProviderConfigReference = ref
}

// GetCondition returns the condition with the given type.
func (u *User) GetCondition(ct xpv1.ConditionType) xpv1.Condition {
	return u.Status.GetCondition(ct)
}

// SetConditions sets the conditions.
func (u *User) SetConditions(c ...xpv1.Condition) {
	u.Status.SetConditions(c...)
}

// GetDeletionPolicy returns the deletion policy.
func (u *User) GetDeletionPolicy() xpv1.DeletionPolicy {
	return u.Spec.DeletionPolicy
}

// SetDeletionPolicy sets the deletion policy.
func (u *User) SetDeletionPolicy(p xpv1.DeletionPolicy) {
	u.Spec.DeletionPolicy = p
}

// GetManagementPolicies returns the management policies.
func (u *User) GetManagementPolicies() xpv1.ManagementPolicies {
	return u.Spec.ManagementPolicies
}

// SetManagementPolicies sets the management policies.
func (u *User) SetManagementPolicies(p xpv1.ManagementPolicies) {
	u.Spec.ManagementPolicies = p
}

// GetPublishConnectionDetailsTo returns the publish connection details to configuration.
func (u *User) GetPublishConnectionDetailsTo() *xpv1.PublishConnectionDetailsTo {
	return u.Spec.PublishConnectionDetailsTo
}

// SetPublishConnectionDetailsTo sets the publish connection details to configuration.
func (u *User) SetPublishConnectionDetailsTo(p *xpv1.PublishConnectionDetailsTo) {
	u.Spec.PublishConnectionDetailsTo = p
}

// GetWriteConnectionSecretToReference returns the write connection secret to reference.
func (u *User) GetWriteConnectionSecretToReference() *xpv1.SecretReference {
	return u.Spec.WriteConnectionSecretToReference
}

// SetWriteConnectionSecretToReference sets the write connection secret to reference.
func (u *User) SetWriteConnectionSecretToReference(r *xpv1.SecretReference) {
	u.Spec.WriteConnectionSecretToReference = r
}

// +kubebuilder:object:root=true

// UserList contains a list of User
type UserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []User `json:"items"`
}
