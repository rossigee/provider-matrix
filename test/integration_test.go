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

package test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"

	"github.com/crossplane-contrib/provider-matrix/apis"
	userv1alpha1 "github.com/crossplane-contrib/provider-matrix/apis/user/v1alpha1"
	roomv1alpha1 "github.com/crossplane-contrib/provider-matrix/apis/room/v1alpha1"
	"github.com/crossplane-contrib/provider-matrix/apis/v1beta1"
)

// IntegrationTestSuite provides integration testing for the Matrix provider
type IntegrationTestSuite struct {
	suite.Suite
	client client.Client
	scheme *runtime.Scheme
}

func TestIntegrationSuite(t *testing.T) {
	// Skip integration tests if no Matrix credentials are provided
	if os.Getenv("MATRIX_ACCESS_TOKEN") == "" {
		t.Skip("Skipping integration tests - MATRIX_ACCESS_TOKEN not set")
	}
	
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) SetupSuite() {
	// Initialize scheme with all APIs
	scheme := runtime.NewScheme()
	err := apis.AddToScheme(scheme)
	require.NoError(suite.T(), err)
	
	// Create fake client for testing
	suite.client = fake.NewClientBuilder().
		WithScheme(scheme).
		Build()
	suite.scheme = scheme
}

func (suite *IntegrationTestSuite) TestProviderConfig() {
	ctx := context.Background()
	
	// Create a ProviderConfig
	pc := &v1beta1.ProviderConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-config",
		},
		Spec: v1beta1.ProviderConfigSpec{
			HomeserverURL: "https://matrix.example.com",
			AdminMode:    boolPtr(true),
			Credentials: v1beta1.ProviderCredentials{
				Source: "Secret",
				CommonCredentialSelectors: v1beta1.CommonCredentialSelectors{
					SecretRef: &v1beta1.SecretKeySelector{
						SecretReference: v1beta1.SecretReference{
							Name:      "matrix-creds",
							Namespace: "default",
						},
						Key: "credentials",
					},
				},
			},
		},
	}
	
	err := suite.client.Create(ctx, pc)
	assert.NoError(suite.T(), err)
	
	// Verify it was created
	retrieved := &v1beta1.ProviderConfig{}
	err = suite.client.Get(ctx, client.ObjectKey{Name: "test-config"}, retrieved)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "https://matrix.example.com", retrieved.Spec.HomeserverURL)
}

func (suite *IntegrationTestSuite) TestUserResource() {
	ctx := context.Background()
	
	// Create a User resource
	user := &userv1alpha1.User{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-user",
		},
		Spec: userv1alpha1.UserSpec{
			ResourceSpec: xpv1.ResourceSpec{
				ProviderConfigReference: &xpv1.Reference{
					Name: "test-config",
				},
			},
			ForProvider: userv1alpha1.UserParameters{
				UserID:      stringPtr("@testuser:example.com"),
				DisplayName: stringPtr("Test User"),
				Admin:       boolPtr(false),
			},
		},
	}
	
	err := suite.client.Create(ctx, user)
	assert.NoError(suite.T(), err)
	
	// Verify it was created
	retrieved := &userv1alpha1.User{}
	err = suite.client.Get(ctx, client.ObjectKey{Name: "test-user"}, retrieved)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "@testuser:example.com", *retrieved.Spec.ForProvider.UserID)
	assert.Equal(suite.T(), "Test User", *retrieved.Spec.ForProvider.DisplayName)
}

func (suite *IntegrationTestSuite) TestRoomResource() {
	ctx := context.Background()
	
	// Create a Room resource
	room := &roomv1alpha1.Room{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-room",
		},
		Spec: roomv1alpha1.RoomSpec{
			ResourceSpec: xpv1.ResourceSpec{
				ProviderConfigReference: &xpv1.Reference{
					Name: "test-config",
				},
			},
			ForProvider: roomv1alpha1.RoomParameters{
				Name:              stringPtr("Test Room"),
				Topic:             stringPtr("A test room for integration testing"),
				Preset:            stringPtr("private_chat"),
				Visibility:        stringPtr("private"),
				EncryptionEnabled: boolPtr(true),
				Invite: []string{
					"@testuser:example.com",
				},
			},
		},
	}
	
	err := suite.client.Create(ctx, room)
	assert.NoError(suite.T(), err)
	
	// Verify it was created
	retrieved := &roomv1alpha1.Room{}
	err = suite.client.Get(ctx, client.ObjectKey{Name: "test-room"}, retrieved)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Test Room", *retrieved.Spec.ForProvider.Name)
	assert.Equal(suite.T(), "private_chat", *retrieved.Spec.ForProvider.Preset)
	assert.True(suite.T(), *retrieved.Spec.ForProvider.EncryptionEnabled)
}

func (suite *IntegrationTestSuite) TestResourceLifecycle() {
	ctx := context.Background()
	
	// Test creating, updating, and deleting resources
	user := &userv1alpha1.User{
		ObjectMeta: metav1.ObjectMeta{
			Name: "lifecycle-user",
		},
		Spec: userv1alpha1.UserSpec{
			ResourceSpec: xpv1.ResourceSpec{
				ProviderConfigReference: &xpv1.Reference{
					Name: "test-config",
				},
			},
			ForProvider: userv1alpha1.UserParameters{
				UserID:      stringPtr("@lifecycle:example.com"),
				DisplayName: stringPtr("Lifecycle User"),
			},
		},
	}
	
	// Create
	err := suite.client.Create(ctx, user)
	assert.NoError(suite.T(), err)
	
	// Update
	user.Spec.ForProvider.DisplayName = stringPtr("Updated Lifecycle User")
	err = suite.client.Update(ctx, user)
	assert.NoError(suite.T(), err)
	
	// Verify update
	retrieved := &userv1alpha1.User{}
	err = suite.client.Get(ctx, client.ObjectKey{Name: "lifecycle-user"}, retrieved)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated Lifecycle User", *retrieved.Spec.ForProvider.DisplayName)
	
	// Delete
	err = suite.client.Delete(ctx, user)
	assert.NoError(suite.T(), err)
}

// Benchmark tests for performance evaluation
func BenchmarkUserCreation(b *testing.B) {
	scheme := runtime.NewScheme()
	err := apis.AddToScheme(scheme)
	require.NoError(b, err)
	
	client := fake.NewClientBuilder().
		WithScheme(scheme).
		Build()
	
	ctx := context.Background()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := &userv1alpha1.User{
			ObjectMeta: metav1.ObjectMeta{
				Name: "bench-user-" + string(rune(i)),
			},
			Spec: userv1alpha1.UserSpec{
				ResourceSpec: xpv1.ResourceSpec{
					ProviderConfigReference: &xpv1.Reference{
						Name: "test-config",
					},
				},
				ForProvider: userv1alpha1.UserParameters{
					UserID:      stringPtr("@bench" + string(rune(i)) + ":example.com"),
					DisplayName: stringPtr("Benchmark User"),
				},
			},
		}
		
		err := client.Create(ctx, user)
		require.NoError(b, err)
	}
}

// Test utilities
func boolPtr(b bool) *bool {
	return &b
}

func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *metav1.Time {
	return &metav1.Time{Time: t}
}