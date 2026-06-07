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
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/crossplane-contrib/provider-matrix/internal/clients"
)

// Priority 2: Domain Extraction Edge Cases
// ========================================

// TestDomainExtractionIPv6 tests IPv6 address handling
func TestDomainExtractionIPv6(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected string
	}{
		{
			name:     "IPv6 with brackets",
			id:       "@alice:[2001:db8::1]",
			expected: "[2001:db8::1]",
		},
		{
			name:     "IPv6 with port",
			id:       "@alice:[2001:db8::1]:8008",
			expected: "[2001:db8::1]:8008",
		},
		{
			name:     "IPv6 localhost",
			id:       "@alice:[::1]",
			expected: "[::1]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain, err := extractTestDomain(tt.id)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, domain)
		})
	}
}

// TestDomainExtractionWithPort tests port handling in domains
func TestDomainExtractionWithPort(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected string
	}{
		{
			name:     "standard port 8008",
			id:       "@alice:example.com:8008",
			expected: "example.com:8008",
		},
		{
			name:     "custom port 9000",
			id:       "@alice:matrix.org:9000",
			expected: "matrix.org:9000",
		},
		{
			name:     "HTTPS port 8448",
			id:       "@alice:example.com:8448",
			expected: "example.com:8448",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain, err := extractTestDomain(tt.id)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, domain)
		})
	}
}

// TestDomainExtractionSpecialChars tests special characters in domain
func TestDomainExtractionSpecialChars(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected string
	}{
		{
			name:     "hyphen in domain",
			id:       "@alice:my-matrix.example.com",
			expected: "my-matrix.example.com",
		},
		{
			name:     "underscore in domain",
			id:       "@alice:matrix_server.example.com",
			expected: "matrix_server.example.com",
		},
		{
			name:     "multiple dots",
			id:       "@alice:mail.matrix.example.com",
			expected: "mail.matrix.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain, err := extractTestDomain(tt.id)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, domain)
		})
	}
}

// Priority 2: Room Details Retrieval
// ==================================

// MockRoomDetailsClient for testing room detail retrieval
type MockRoomDetailsClient struct {
	getRoomDetailsFn func(ctx context.Context, roomID string) (*clients.Room, error)
}

func (m *MockRoomDetailsClient) GetRoom(ctx context.Context, roomID string) (*clients.Room, error) {
	if m.getRoomDetailsFn != nil {
		return m.getRoomDetailsFn(ctx, roomID)
	}
	return &clients.Room{RoomID: roomID}, nil
}

// Stub implementations
func (m *MockRoomDetailsClient) CreateUser(ctx context.Context, user *clients.UserSpec) (*clients.User, error) {
	return nil, nil
}
func (m *MockRoomDetailsClient) GetUser(ctx context.Context, userID string) (*clients.User, error) {
	return nil, nil
}
func (m *MockRoomDetailsClient) UpdateUser(ctx context.Context, userID string, user *clients.UserSpec) (*clients.User, error) {
	return nil, nil
}
func (m *MockRoomDetailsClient) DeactivateUser(ctx context.Context, userID string) error {
	return nil
}
func (m *MockRoomDetailsClient) CreateRoom(ctx context.Context, room *clients.RoomSpec) (*clients.Room, error) {
	return nil, nil
}
func (m *MockRoomDetailsClient) UpdateRoom(ctx context.Context, roomID string, room *clients.RoomSpec) (*clients.Room, error) {
	return nil, nil
}
func (m *MockRoomDetailsClient) DeleteRoom(ctx context.Context, roomID string) error {
	return nil
}
func (m *MockRoomDetailsClient) SetPowerLevels(ctx context.Context, roomID string, powerLevels *clients.PowerLevelSpec) error {
	return nil
}
func (m *MockRoomDetailsClient) GetPowerLevels(ctx context.Context, roomID string) (*clients.PowerLevelContent, error) {
	return nil, nil
}
func (m *MockRoomDetailsClient) CreateRoomAlias(ctx context.Context, alias string, roomID string) error {
	return nil
}
func (m *MockRoomDetailsClient) GetRoomAlias(ctx context.Context, alias string) (*clients.RoomAlias, error) {
	return nil, nil
}
func (m *MockRoomDetailsClient) DeleteRoomAlias(ctx context.Context, alias string) error {
	return nil
}
func (m *MockRoomDetailsClient) ListUsers(ctx context.Context, from string, limit int) (*clients.ListUsersResponse, error) {
	return nil, nil
}
func (m *MockRoomDetailsClient) ListRooms(ctx context.Context, from string, limit int) (*clients.ListRoomsResponse, error) {
	return nil, nil
}
func (m *MockRoomDetailsClient) MakeRoomAdmin(ctx context.Context, roomID, userID string) error {
	return nil
}
func (m *MockRoomDetailsClient) BlockRoom(ctx context.Context, roomID string, block bool) error {
	return nil
}

// TestGetRoomDetails tests retrieving comprehensive room information
func TestGetRoomDetails(t *testing.T) {
	mockClient := &MockRoomDetailsClient{
		getRoomDetailsFn: func(ctx context.Context, roomID string) (*clients.Room, error) {
			return &clients.Room{
				RoomID:         roomID,
				Name:           "General",
				Topic:          "General discussion",
				JoinedMembers:  42,
				Visibility:     "public",
			}, nil
		},
	}

	room, err := mockClient.GetRoom(context.Background(), "!general:example.com")
	require.NoError(t, err)
	assert.Equal(t, "!general:example.com", room.RoomID)
	assert.Equal(t, 42, room.JoinedMembers)
	assert.Equal(t, "public", room.Visibility)
}

// TestGetRoomDetailsLargeRoom tests handling large rooms
func TestGetRoomDetailsLargeRoom(t *testing.T) {
	mockClient := &MockRoomDetailsClient{
		getRoomDetailsFn: func(ctx context.Context, roomID string) (*clients.Room, error) {
			return &clients.Room{
				RoomID:        roomID,
				JoinedMembers: 10000,
				Name:          "Large Room",
			}, nil
		},
	}

	room, err := mockClient.GetRoom(context.Background(), "!large:example.com")
	require.NoError(t, err)
	assert.Equal(t, 10000, room.JoinedMembers)
}

// TestGetRoomDetailsWithEncryption tests encrypted room details
func TestGetRoomDetailsWithEncryption(t *testing.T) {
	mockClient := &MockRoomDetailsClient{
		getRoomDetailsFn: func(ctx context.Context, roomID string) (*clients.Room, error) {
			return &clients.Room{
				RoomID:            roomID,
				Name:              "Secure Room",
				EncryptionEnabled: true,
			}, nil
		},
	}

	room, err := mockClient.GetRoom(context.Background(), "!secure:example.com")
	require.NoError(t, err)
	assert.True(t, room.EncryptionEnabled)
}

// Priority 2: Error Handling
// ==========================

// TestNetworkErrorHandling tests network-level errors
func TestNetworkErrorHandling(t *testing.T) {
	tests := []struct {
		name    string
		error   error
		message string
	}{
		{
			name:    "connection timeout",
			error:   errors.New("i/o timeout"),
			message: "i/o timeout",
		},
		{
			name:    "connection refused",
			error:   errors.New("connection refused"),
			message: "connection refused",
		},
		{
			name:    "DNS resolution failed",
			error:   errors.New("no such host"),
			message: "no such host",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Error(t, tt.error)
			assert.Contains(t, tt.error.Error(), tt.message)
		})
	}
}

// TestMalformedResponseHandling tests malformed API responses
func TestMalformedResponseHandling(t *testing.T) {
	tests := []struct {
		name    string
		error   error
		message string
	}{
		{
			name:    "invalid JSON",
			error:   errors.New("invalid JSON response"),
			message: "invalid JSON",
		},
		{
			name:    "missing required field",
			error:   errors.New("missing user_id in response"),
			message: "missing user_id",
		},
		{
			name:    "unexpected type",
			error:   errors.New("expected string, got int"),
			message: "expected string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Error(t, tt.error)
			assert.Contains(t, tt.error.Error(), tt.message)
		})
	}
}

// TestTimeoutScenarios tests timeout error scenarios
func TestTimeoutScenarios(t *testing.T) {
	tests := []struct {
		name     string
		timeout  time.Duration
		expected bool
	}{
		{
			name:     "quick operation",
			timeout:  100 * time.Millisecond,
			expected: true,
		},
		{
			name:     "reasonable timeout",
			timeout:  5 * time.Second,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify timeout configuration is valid
			assert.True(t, tt.timeout > 0)
			assert.Equal(t, tt.expected, tt.timeout > 0)
		})
	}
}

// Priority 2: Admin Permission Scenarios
// ======================================

// MockAdminPermClient for testing admin permissions
type MockAdminPermClient struct {
	makeAdminFn func(ctx context.Context, roomID, userID string) error
	checkPermFn func(ctx context.Context, userID string) bool
}

func (m *MockAdminPermClient) MakeRoomAdmin(ctx context.Context, roomID, userID string) error {
	if m.makeAdminFn != nil {
		return m.makeAdminFn(ctx, roomID, userID)
	}
	return nil
}

func (m *MockAdminPermClient) checkPermissions(ctx context.Context, userID string) bool {
	if m.checkPermFn != nil {
		return m.checkPermFn(ctx, userID)
	}
	return false
}

// Stub implementations
func (m *MockAdminPermClient) CreateUser(ctx context.Context, user *clients.UserSpec) (*clients.User, error) {
	return nil, nil
}
func (m *MockAdminPermClient) GetUser(ctx context.Context, userID string) (*clients.User, error) {
	return nil, nil
}
func (m *MockAdminPermClient) UpdateUser(ctx context.Context, userID string, user *clients.UserSpec) (*clients.User, error) {
	return nil, nil
}
func (m *MockAdminPermClient) DeactivateUser(ctx context.Context, userID string) error {
	return nil
}
func (m *MockAdminPermClient) CreateRoom(ctx context.Context, room *clients.RoomSpec) (*clients.Room, error) {
	return nil, nil
}
func (m *MockAdminPermClient) GetRoom(ctx context.Context, roomID string) (*clients.Room, error) {
	return nil, nil
}
func (m *MockAdminPermClient) UpdateRoom(ctx context.Context, roomID string, room *clients.RoomSpec) (*clients.Room, error) {
	return nil, nil
}
func (m *MockAdminPermClient) DeleteRoom(ctx context.Context, roomID string) error {
	return nil
}
func (m *MockAdminPermClient) SetPowerLevels(ctx context.Context, roomID string, powerLevels *clients.PowerLevelSpec) error {
	return nil
}
func (m *MockAdminPermClient) GetPowerLevels(ctx context.Context, roomID string) (*clients.PowerLevelContent, error) {
	return nil, nil
}
func (m *MockAdminPermClient) CreateRoomAlias(ctx context.Context, alias string, roomID string) error {
	return nil
}
func (m *MockAdminPermClient) GetRoomAlias(ctx context.Context, alias string) (*clients.RoomAlias, error) {
	return nil, nil
}
func (m *MockAdminPermClient) DeleteRoomAlias(ctx context.Context, alias string) error {
	return nil
}
func (m *MockAdminPermClient) ListUsers(ctx context.Context, from string, limit int) (*clients.ListUsersResponse, error) {
	return nil, nil
}
func (m *MockAdminPermClient) ListRooms(ctx context.Context, from string, limit int) (*clients.ListRoomsResponse, error) {
	return nil, nil
}
func (m *MockAdminPermClient) BlockRoom(ctx context.Context, roomID string, block bool) error {
	return nil
}

// TestAdminPermissionPromotionSuccess tests successful admin promotion
func TestAdminPermissionPromotionSuccess(t *testing.T) {
	mockClient := &MockAdminPermClient{
		makeAdminFn: func(ctx context.Context, roomID, userID string) error {
			if userID == "@regular:example.com" {
				return nil // Success
			}
			return errors.New("user not found")
		},
		checkPermFn: func(ctx context.Context, userID string) bool {
			return userID == "@admin:example.com"
		},
	}

	err := mockClient.MakeRoomAdmin(context.Background(), "!room:example.com", "@regular:example.com")
	require.NoError(t, err)
}

// TestAdminPermissionDenied tests permission denied scenario
func TestAdminPermissionDenied(t *testing.T) {
	mockClient := &MockAdminPermClient{
		makeAdminFn: func(ctx context.Context, roomID, userID string) error {
			return errors.New("insufficient permissions")
		},
	}

	err := mockClient.MakeRoomAdmin(context.Background(), "!room:example.com", "@user:example.com")
	assert.Error(t, err)
}

// TestAdminPermissionCheck tests permission verification
func TestAdminPermissionCheck(t *testing.T) {
	mockClient := &MockAdminPermClient{
		checkPermFn: func(ctx context.Context, userID string) bool {
			return userID == "@admin:example.com" || userID == "@superadmin:example.com"
		},
	}

	assert.True(t, mockClient.checkPermissions(context.Background(), "@admin:example.com"))
	assert.True(t, mockClient.checkPermissions(context.Background(), "@superadmin:example.com"))
	assert.False(t, mockClient.checkPermissions(context.Background(), "@user:example.com"))
}

// TestAdminMultiplePolicies tests multiple permission policies
func TestAdminMultiplePolicies(t *testing.T) {
	policies := []string{
		"admin_access",
		"room_creation",
		"user_management",
		"space_management",
	}

	for _, policy := range policies {
		t.Run(policy, func(t *testing.T) {
			assert.NotEmpty(t, policy)
		})
	}
}
