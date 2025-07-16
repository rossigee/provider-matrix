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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/crossplane-contrib/provider-matrix/apis/user/v1alpha1"
	"github.com/crossplane-contrib/provider-matrix/internal/clients"
)

func TestGenerateUserSpec(t *testing.T) {
	tests := []struct {
		name string
		cr   *v1alpha1.User
		want *clients.UserSpec
	}{
		{
			name: "basic user spec",
			cr: &v1alpha1.User{
				Spec: v1alpha1.UserSpec{
					ForProvider: v1alpha1.UserParameters{
						UserID:      stringPtr("@alice:example.com"),
						DisplayName: stringPtr("Alice Wonderland"),
						Admin:       boolPtr(false),
					},
				},
			},
			want: &clients.UserSpec{
				UserID:      "@alice:example.com",
				DisplayName: "Alice Wonderland",
				Admin:       false,
			},
		},
		{
			name: "admin user with external IDs",
			cr: &v1alpha1.User{
				Spec: v1alpha1.UserSpec{
					ForProvider: v1alpha1.UserParameters{
						UserID:      stringPtr("@admin:example.com"),
						DisplayName: stringPtr("Admin User"),
						Admin:       boolPtr(true),
						UserType:    stringPtr("admin"),
						ExternalIDs: []v1alpha1.ExternalID{
							{
								Medium:    "email",
								Address:   "admin@example.com",
								Validated: true,
							},
						},
					},
				},
			},
			want: &clients.UserSpec{
				UserID:      "@admin:example.com",
				DisplayName: "Admin User",
				Admin:       true,
				UserType:    "admin",
				ExternalIDs: []clients.ExternalID{
					{
						Medium:    "email",
						Address:   "admin@example.com",
						Validated: true,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateUserSpec(tt.cr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGenerateUserObservation(t *testing.T) {
	now := time.Now()
	user := &clients.User{
		UserID:       "@alice:example.com",
		DisplayName:  "Alice Wonderland",
		Admin:        false,
		Deactivated:  false,
		CreationTime: &now,
		UserType:     "regular",
		ExternalIDs: []clients.ExternalID{
			{
				Medium:    "email",
				Address:   "alice@example.com",
				Validated: true,
			},
		},
		Devices: []clients.Device{
			{
				DeviceID:     "DEVICE123",
				DisplayName:  "Phone",
				LastSeenTime: &now,
			},
		},
	}

	obs := generateUserObservation(user)

	assert.Equal(t, "@alice:example.com", obs.UserID)
	assert.Equal(t, "Alice Wonderland", obs.DisplayName)
	assert.Equal(t, false, obs.Admin)
	assert.Equal(t, false, obs.Deactivated)
	assert.Equal(t, "regular", obs.UserType)
	assert.NotNil(t, obs.CreationTime)
	assert.Len(t, obs.ExternalIDs, 1)
	assert.Equal(t, "email", obs.ExternalIDs[0].Medium)
	assert.Len(t, obs.Devices, 1)
	assert.Equal(t, "DEVICE123", obs.Devices[0].DeviceID)
}

func TestIsUserUpToDate(t *testing.T) {
	tests := []struct {
		name string
		cr   *v1alpha1.User
		user *clients.User
		want bool
	}{
		{
			name: "user is up to date",
			cr: &v1alpha1.User{
				Spec: v1alpha1.UserSpec{
					ForProvider: v1alpha1.UserParameters{
						DisplayName: stringPtr("Alice"),
						Admin:       boolPtr(false),
					},
				},
			},
			user: &clients.User{
				DisplayName: "Alice",
				Admin:       false,
			},
			want: true,
		},
		{
			name: "display name differs",
			cr: &v1alpha1.User{
				Spec: v1alpha1.UserSpec{
					ForProvider: v1alpha1.UserParameters{
						DisplayName: stringPtr("Alice Updated"),
						Admin:       boolPtr(false),
					},
				},
			},
			user: &clients.User{
				DisplayName: "Alice",
				Admin:       false,
			},
			want: false,
		},
		{
			name: "admin status differs",
			cr: &v1alpha1.User{
				Spec: v1alpha1.UserSpec{
					ForProvider: v1alpha1.UserParameters{
						DisplayName: stringPtr("Alice"),
						Admin:       boolPtr(true),
					},
				},
			},
			user: &clients.User{
				DisplayName: "Alice",
				Admin:       false,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isUserUpToDate(tt.cr, tt.user)
			assert.Equal(t, tt.want, got)
		})
	}
}

// Helper functions for creating pointers
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func timePtr(t time.Time) *metav1.Time {
	return &metav1.Time{Time: t}
}