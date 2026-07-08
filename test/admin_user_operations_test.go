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
	"github.com/crossplane-contrib/provider-matrix/internal/clients"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// MockAdminUserClient for testing admin-specific user operations
type MockAdminUserClient struct {
	createUserAdminFn     func(ctx context.Context, user *clients.UserSpec) (*clients.User, error)
	getUserAdminFn        func(ctx context.Context, userID string) (*clients.User, error)
	updateUserAdminFn     func(ctx context.Context, userID string, user *clients.UserSpec) (*clients.User, error)
	deactivateUserAdminFn func(ctx context.Context, userID string) error
}

func (m *MockAdminUserClient) CreateUser(ctx context.Context, user *clients.UserSpec) (*clients.User, error) {
	if m.createUserAdminFn != nil {
		return m.createUserAdminFn(ctx, user)
	}
	return &clients.User{UserID: user.UserID}, nil
}

func (m *MockAdminUserClient) GetUser(ctx context.Context, userID string) (*clients.User, error) {
	if m.getUserAdminFn != nil {
		return m.getUserAdminFn(ctx, userID)
	}
	return &clients.User{UserID: userID}, nil
}

func (m *MockAdminUserClient) UpdateUser(ctx context.Context, userID string, user *clients.UserSpec) (*clients.User, error) {
	if m.updateUserAdminFn != nil {
		return m.updateUserAdminFn(ctx, userID, user)
	}
	return &clients.User{UserID: userID}, nil
}

func (m *MockAdminUserClient) DeactivateUser(ctx context.Context, userID string) error {
	if m.deactivateUserAdminFn != nil {
		return m.deactivateUserAdminFn(ctx, userID)
	}
	return nil
}

// Stub implementations for unused interface methods
func (m *MockAdminUserClient) CreateRoom(ctx context.Context, room *clients.RoomSpec) (*clients.Room, error) {
	return nil, nil
}

func (m *MockAdminUserClient) GetRoom(ctx context.Context, roomID string) (*clients.Room, error) {
	return nil, nil
}

func (m *MockAdminUserClient) UpdateRoom(ctx context.Context, roomID string, room *clients.RoomSpec) (*clients.Room, error) {
	return nil, nil
}

func (m *MockAdminUserClient) DeleteRoom(ctx context.Context, roomID string) error {
	return nil
}

func (m *MockAdminUserClient) SetPowerLevels(ctx context.Context, roomID string, powerLevels *clients.PowerLevelSpec) error {
	return nil
}

func (m *MockAdminUserClient) GetPowerLevels(ctx context.Context, roomID string) (*clients.PowerLevelContent, error) {
	return nil, nil
}

func (m *MockAdminUserClient) CreateRoomAlias(ctx context.Context, alias string, roomID string) error {
	return nil
}

func (m *MockAdminUserClient) GetRoomAlias(ctx context.Context, alias string) (*clients.RoomAlias, error) {
	return nil, nil
}

func (m *MockAdminUserClient) DeleteRoomAlias(ctx context.Context, alias string) error {
	return nil
}

func (m *MockAdminUserClient) ListUsers(ctx context.Context, from string, limit int) (*clients.ListUsersResponse, error) {
	return nil, nil
}

func (m *MockAdminUserClient) ListRooms(ctx context.Context, from string, limit int) (*clients.ListRoomsResponse, error) {
	return nil, nil
}

func (m *MockAdminUserClient) MakeRoomAdmin(ctx context.Context, roomID, userID string) error {
	return nil
}

func (m *MockAdminUserClient) BlockRoom(ctx context.Context, roomID string, block bool) error {
	return nil
}

// TestAdminCreateUser tests admin user creation
func TestAdminCreateUser(t *testing.T) {
	mockClient := &MockAdminUserClient{
		createUserAdminFn: func(ctx context.Context, user *clients.UserSpec) (*clients.User, error) {
			return &clients.User{
				UserID:      user.UserID,
				DisplayName: user.DisplayName,
				Admin:       true,
			}, nil
		},
	}

	spec := &clients.UserSpec{
		UserID:      "@admin:example.com",
		DisplayName: "Admin User",
		Password:    "secret",
	}

	created, err := mockClient.CreateUser(context.Background(), spec)
	require.NoError(t, err)
	assert.Equal(t, "@admin:example.com", created.UserID)
	assert.True(t, created.Admin)
}

// TestAdminGetUser tests admin user retrieval
func TestAdminGetUser(t *testing.T) {
	mockClient := &MockAdminUserClient{
		getUserAdminFn: func(ctx context.Context, userID string) (*clients.User, error) {
			return &clients.User{
				UserID:      userID,
				DisplayName: "Admin User",
				Admin:       true,
			}, nil
		},
	}

	user, err := mockClient.GetUser(context.Background(), "@admin:example.com")
	require.NoError(t, err)
	assert.True(t, user.Admin)
}

// TestAdminUpdateUser tests admin user updates
func TestAdminUpdateUser(t *testing.T) {
	mockClient := &MockAdminUserClient{
		updateUserAdminFn: func(ctx context.Context, userID string, user *clients.UserSpec) (*clients.User, error) {
			return &clients.User{
				UserID:      userID,
				DisplayName: user.DisplayName,
				Admin:       true,
			}, nil
		},
	}

	spec := &clients.UserSpec{
		DisplayName: "Updated Admin",
	}

	updated, err := mockClient.UpdateUser(context.Background(), "@admin:example.com", spec)
	require.NoError(t, err)
	assert.Equal(t, "Updated Admin", updated.DisplayName)
	assert.True(t, updated.Admin)
}

// TestAdminDeactivateUser tests admin user deactivation
func TestAdminDeactivateUser(t *testing.T) {
	mockClient := &MockAdminUserClient{
		deactivateUserAdminFn: func(ctx context.Context, userID string) error {
			if userID == "" {
				return errors.New("userID required")
			}
			return nil
		},
	}

	err := mockClient.DeactivateUser(context.Background(), "@user:example.com")
	require.NoError(t, err)
}

// TestAdminDeactivateUserError tests deactivation error handling
func TestAdminDeactivateUserError(t *testing.T) {
	mockClient := &MockAdminUserClient{
		deactivateUserAdminFn: func(ctx context.Context, userID string) error {
			return errors.New("cannot deactivate system user")
		},
	}

	err := mockClient.DeactivateUser(context.Background(), "@system:example.com")
	assert.Error(t, err)
}

// TestAdminUserWorkflow tests complete admin user workflow
func TestAdminUserWorkflow(t *testing.T) {
	callLog := []string{}

	mockClient := &MockAdminUserClient{
		createUserAdminFn: func(ctx context.Context, user *clients.UserSpec) (*clients.User, error) {
			callLog = append(callLog, "create")
			return &clients.User{UserID: user.UserID, Admin: true}, nil
		},
		getUserAdminFn: func(ctx context.Context, userID string) (*clients.User, error) {
			callLog = append(callLog, "get")
			return &clients.User{UserID: userID, Admin: true}, nil
		},
		updateUserAdminFn: func(ctx context.Context, userID string, user *clients.UserSpec) (*clients.User, error) {
			callLog = append(callLog, "update")
			return &clients.User{UserID: userID, DisplayName: user.DisplayName, Admin: true}, nil
		},
	}

	// Create admin user
	_, err := mockClient.CreateUser(context.Background(), &clients.UserSpec{UserID: "@newadmin:example.com"})
	require.NoError(t, err)

	// Get user to verify
	user, err := mockClient.GetUser(context.Background(), "@newadmin:example.com")
	require.NoError(t, err)
	assert.True(t, user.Admin)

	// Update user details
	_, err = mockClient.UpdateUser(context.Background(), "@newadmin:example.com", &clients.UserSpec{DisplayName: "New Admin"})
	require.NoError(t, err)

	// Verify workflow
	assert.Equal(t, []string{"create", "get", "update"}, callLog)
}

// TestAdminBulkUserOperations tests multiple user operations
func TestAdminBulkUserOperations(t *testing.T) {
	callCount := 0
	mockClient := &MockAdminUserClient{
		createUserAdminFn: func(ctx context.Context, user *clients.UserSpec) (*clients.User, error) {
			callCount++
			return &clients.User{UserID: user.UserID}, nil
		},
	}

	// Create multiple users
	for i := 1; i <= 5; i++ {
		userID := "@user" + stringOfLengthChar(i, 'a') + ":example.com"
		_, err := mockClient.CreateUser(context.Background(), &clients.UserSpec{
			UserID: userID,
		})
		require.NoError(t, err)
	}

	assert.Equal(t, 5, callCount, "expected 5 user creation calls")
}

// Helper function
func stringOfLengthChar(count int, ch rune) string {
	result := ""
	for i := 0; i < count; i++ {
		result += string(ch)
	}
	return result
}
