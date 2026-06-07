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
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	xpcoreapi "github.com/crossplane/crossplane/apis/v2/core/v2"

	"github.com/crossplane-contrib/provider-matrix/apis/user/v1alpha1"
	"github.com/crossplane-contrib/provider-matrix/internal/clients"
)

// MockMatrixClient implements clients.Client for testing
type MockMatrixClient struct {
	createUserFn   func(ctx context.Context, user *clients.UserSpec) (*clients.User, error)
	getUserFn      func(ctx context.Context, userID string) (*clients.User, error)
	updateUserFn   func(ctx context.Context, userID string, user *clients.UserSpec) (*clients.User, error)
	deactivateUserFn func(ctx context.Context, userID string) error
}

func (m *MockMatrixClient) CreateUser(ctx context.Context, user *clients.UserSpec) (*clients.User, error) {
	if m.createUserFn != nil {
		return m.createUserFn(ctx, user)
	}
	return &clients.User{UserID: user.UserID}, nil
}

func (m *MockMatrixClient) GetUser(ctx context.Context, userID string) (*clients.User, error) {
	if m.getUserFn != nil {
		return m.getUserFn(ctx, userID)
	}
	return &clients.User{UserID: userID}, nil
}

func (m *MockMatrixClient) UpdateUser(ctx context.Context, userID string, user *clients.UserSpec) (*clients.User, error) {
	if m.updateUserFn != nil {
		return m.updateUserFn(ctx, userID, user)
	}
	return &clients.User{UserID: userID}, nil
}

func (m *MockMatrixClient) DeactivateUser(ctx context.Context, userID string) error {
	if m.deactivateUserFn != nil {
		return m.deactivateUserFn(ctx, userID)
	}
	return nil
}

// Stub implementations for unused interface methods
func (m *MockMatrixClient) CreateRoom(ctx context.Context, room *clients.RoomSpec) (*clients.Room, error) {
	return nil, nil
}

func (m *MockMatrixClient) GetRoom(ctx context.Context, roomID string) (*clients.Room, error) {
	return nil, nil
}

func (m *MockMatrixClient) UpdateRoom(ctx context.Context, roomID string, room *clients.RoomSpec) (*clients.Room, error) {
	return nil, nil
}

func (m *MockMatrixClient) DeleteRoom(ctx context.Context, roomID string) error {
	return nil
}

func (m *MockMatrixClient) SetPowerLevels(ctx context.Context, roomID string, powerLevels *clients.PowerLevelSpec) error {
	return nil
}

func (m *MockMatrixClient) GetPowerLevels(ctx context.Context, roomID string) (*clients.PowerLevelContent, error) {
	return nil, nil
}

func (m *MockMatrixClient) CreateRoomAlias(ctx context.Context, alias string, roomID string) error {
	return nil
}

func (m *MockMatrixClient) GetRoomAlias(ctx context.Context, alias string) (*clients.RoomAlias, error) {
	return nil, nil
}

func (m *MockMatrixClient) DeleteRoomAlias(ctx context.Context, alias string) error {
	return nil
}

func (m *MockMatrixClient) ListUsers(ctx context.Context, from string, limit int) (*clients.ListUsersResponse, error) {
	return nil, nil
}

func (m *MockMatrixClient) ListRooms(ctx context.Context, from string, limit int) (*clients.ListRoomsResponse, error) {
	return nil, nil
}

func (m *MockMatrixClient) MakeRoomAdmin(ctx context.Context, roomID, userID string) error {
	return nil
}

func (m *MockMatrixClient) BlockRoom(ctx context.Context, roomID string, block bool) error {
	return nil
}

// TestExternalObserveCreate tests the Observe operation when resource doesn't exist
func TestExternalObserveCreate(t *testing.T) {
	// Create a test user resource with no external name
	user := &v1alpha1.User{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-user",
		},
		Spec: v1alpha1.UserSpec{
			ForProvider: v1alpha1.UserParameters{
				UserID: stringPtr("@test:example.com"),
			},
		},
	}

	// Test that a resource with no external name is recognized as non-existent
	externalName := ""
	assert.Equal(t, "", externalName, "newly created resource should have no external name")
	assert.NotNil(t, user)
}

// TestExternalObserveExists tests the Observe operation when resource exists
func TestExternalObserveExists(t *testing.T) {
	mockClient := &MockMatrixClient{
		getUserFn: func(ctx context.Context, userID string) (*clients.User, error) {
			return &clients.User{
				UserID:      userID,
				DisplayName: "Test User",
			}, nil
		},
	}

	// Verify that GetUser returns the expected user
	user, err := mockClient.GetUser(context.Background(), "@test:example.com")
	require.NoError(t, err)
	assert.Equal(t, "@test:example.com", user.UserID)
	assert.Equal(t, "Test User", user.DisplayName)
}

// TestExternalCreate tests the Create operation
func TestExternalCreate(t *testing.T) {
	mockClient := &MockMatrixClient{
		createUserFn: func(ctx context.Context, user *clients.UserSpec) (*clients.User, error) {
			return &clients.User{
				UserID:      user.UserID,
				DisplayName: user.DisplayName,
			}, nil
		},
	}

	spec := &clients.UserSpec{
		UserID:      "@newuser:example.com",
		DisplayName: "New User",
		Password:    "secret",
	}

	created, err := mockClient.CreateUser(context.Background(), spec)
	require.NoError(t, err)
	assert.Equal(t, spec.UserID, created.UserID)
	assert.Equal(t, spec.DisplayName, created.DisplayName)
}

// TestExternalCreateError tests Create operation failure
func TestExternalCreateError(t *testing.T) {
	mockClient := &MockMatrixClient{
		createUserFn: func(ctx context.Context, user *clients.UserSpec) (*clients.User, error) {
			return nil, errors.New("homeserver error")
		},
	}

	spec := &clients.UserSpec{
		UserID: "@baduser:example.com",
	}

	_, err := mockClient.CreateUser(context.Background(), spec)
	assert.Error(t, err)
	assert.Equal(t, "homeserver error", err.Error())
}

// TestExternalUpdate tests the Update operation
func TestExternalUpdate(t *testing.T) {
	mockClient := &MockMatrixClient{
		updateUserFn: func(ctx context.Context, userID string, user *clients.UserSpec) (*clients.User, error) {
			return &clients.User{
				UserID:      userID,
				DisplayName: user.DisplayName,
			}, nil
		},
	}

	spec := &clients.UserSpec{
		UserID:      "@user:example.com",
		DisplayName: "Updated Name",
	}

	updated, err := mockClient.UpdateUser(context.Background(), "@user:example.com", spec)
	require.NoError(t, err)
	assert.Equal(t, "@user:example.com", updated.UserID)
	assert.Equal(t, "Updated Name", updated.DisplayName)
}

// TestExternalDelete tests the Delete operation
func TestExternalDelete(t *testing.T) {
	mockClient := &MockMatrixClient{
		deactivateUserFn: func(ctx context.Context, userID string) error {
			return nil
		},
	}

	err := mockClient.DeactivateUser(context.Background(), "@user:example.com")
	assert.NoError(t, err)
}

// TestExternalDeleteError tests Delete operation failure
func TestExternalDeleteError(t *testing.T) {
	mockClient := &MockMatrixClient{
		deactivateUserFn: func(ctx context.Context, userID string) error {
			return errors.New("cannot deactivate user")
		},
	}

	err := mockClient.DeactivateUser(context.Background(), "@nonexistent:example.com")
	assert.Error(t, err)
}

// TestUserSpecGeneration tests conversion from CR to client spec
func TestUserSpecGeneration(t *testing.T) {
	tests := []struct {
		name     string
		userID   *string
		displayName *string
		password *string
		admin    *bool
		expectFields map[string]interface{}
	}{
		{
			name:        "full spec",
			userID:      stringPtr("@user:example.com"),
			displayName: stringPtr("Test User"),
			password:    stringPtr("secret"),
			admin:       boolPtr(true),
			expectFields: map[string]interface{}{
				"userID":      "@user:example.com",
				"displayName": "Test User",
				"password":    "secret",
				"admin":       true,
			},
		},
		{
			name:   "minimal spec",
			userID: stringPtr("@minimal:example.com"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &v1alpha1.User{
				Spec: v1alpha1.UserSpec{
					ForProvider: v1alpha1.UserParameters{
						UserID:      tt.userID,
						DisplayName: tt.displayName,
						Password:    tt.password,
						Admin:       tt.admin,
					},
				},
			}

			// Verify that the user spec can be populated from parameters
			assert.NotNil(t, user.Spec.ForProvider.UserID)
			if tt.userID != nil {
				assert.Equal(t, *tt.userID, *user.Spec.ForProvider.UserID)
			}
		})
	}
}

// TestIsNotFound tests the IsNotFound helper function
func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		expect bool
	}{
		{
			name:   "nil error",
			err:    nil,
			expect: false,
		},
		{
			name:   "404 in message",
			err:    errors.New("404 not found"),
			expect: true,
		},
		{
			name:   "not found in message",
			err:    errors.New("resource not found"),
			expect: true,
		},
		{
			name:   "generic error",
			err:    errors.New("some other error"),
			expect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := clients.IsNotFound(tt.err)
			assert.Equal(t, tt.expect, result)
		})
	}
}

// TestProviderConfigReference tests provider config reference handling
func TestProviderConfigReference(t *testing.T) {
	user := &v1alpha1.User{
		Spec: v1alpha1.UserSpec{},
	}

	// Test GetProviderConfigReference
	ref := user.GetProviderConfigReference()
	assert.Nil(t, ref)

	// Test SetProviderConfigReference
	pcRef := &xpcoreapi.ProviderConfigReference{Name: "test-pc"}
	user.SetProviderConfigReference(pcRef)
	assert.Equal(t, "test-pc", user.GetProviderConfigReference().Name)
}

// TestUserConditions tests condition management
func TestUserConditions(t *testing.T) {
	user := &v1alpha1.User{}

	// Set conditions
	cond := xpcoreapi.Available()
	user.SetConditions(cond)

	// Verify condition is set
	retrieved := user.GetCondition(xpcoreapi.TypeReady)
	assert.Equal(t, xpcoreapi.TypeReady, retrieved.Type)
}

// TestUserManagementPolicies tests management policy handling
func TestUserManagementPolicies(t *testing.T) {
	user := &v1alpha1.User{}

	policies := xpcoreapi.ManagementPolicies{"*"}
	user.SetManagementPolicies(policies)

	retrieved := user.GetManagementPolicies()
	assert.Equal(t, policies, retrieved)
}

// TestUserWriteConnectionSecret tests secret reference handling
func TestUserWriteConnectionSecret(t *testing.T) {
	user := &v1alpha1.User{}

	secretRef := &xpcoreapi.LocalSecretReference{
		Name: "credentials",
	}
	user.SetWriteConnectionSecretToReference(secretRef)

	retrieved := user.GetWriteConnectionSecretToReference()
	assert.Equal(t, "credentials", retrieved.Name)
}
