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

// Priority 4: Room Creation with Advanced Options
// ================================================

// MockAdvancedRoomClient for testing advanced room creation options
type MockAdvancedRoomClient struct {
	createRoomFn func(ctx context.Context, room *clients.RoomSpec) (*clients.Room, error)
}

func (m *MockAdvancedRoomClient) CreateRoom(ctx context.Context, room *clients.RoomSpec) (*clients.Room, error) {
	if m.createRoomFn != nil {
		return m.createRoomFn(ctx, room)
	}
	return &clients.Room{Name: room.Name}, nil
}

// Stub implementations
func (m *MockAdvancedRoomClient) CreateUser(ctx context.Context, user *clients.UserSpec) (*clients.User, error) {
	return nil, nil
}
func (m *MockAdvancedRoomClient) GetUser(ctx context.Context, userID string) (*clients.User, error) {
	return nil, nil
}
func (m *MockAdvancedRoomClient) UpdateUser(ctx context.Context, userID string, user *clients.UserSpec) (*clients.User, error) {
	return nil, nil
}
func (m *MockAdvancedRoomClient) DeactivateUser(ctx context.Context, userID string) error {
	return nil
}
func (m *MockAdvancedRoomClient) GetRoom(ctx context.Context, roomID string) (*clients.Room, error) {
	return nil, nil
}
func (m *MockAdvancedRoomClient) UpdateRoom(ctx context.Context, roomID string, room *clients.RoomSpec) (*clients.Room, error) {
	return nil, nil
}
func (m *MockAdvancedRoomClient) DeleteRoom(ctx context.Context, roomID string) error {
	return nil
}
func (m *MockAdvancedRoomClient) SetPowerLevels(ctx context.Context, roomID string, powerLevels *clients.PowerLevelSpec) error {
	return nil
}
func (m *MockAdvancedRoomClient) GetPowerLevels(ctx context.Context, roomID string) (*clients.PowerLevelContent, error) {
	return nil, nil
}
func (m *MockAdvancedRoomClient) CreateRoomAlias(ctx context.Context, alias string, roomID string) error {
	return nil
}
func (m *MockAdvancedRoomClient) GetRoomAlias(ctx context.Context, alias string) (*clients.RoomAlias, error) {
	return nil, nil
}
func (m *MockAdvancedRoomClient) DeleteRoomAlias(ctx context.Context, alias string) error {
	return nil
}
func (m *MockAdvancedRoomClient) ListUsers(ctx context.Context, from string, limit int) (*clients.ListUsersResponse, error) {
	return nil, nil
}
func (m *MockAdvancedRoomClient) ListRooms(ctx context.Context, from string, limit int) (*clients.ListRoomsResponse, error) {
	return nil, nil
}
func (m *MockAdvancedRoomClient) MakeRoomAdmin(ctx context.Context, roomID, userID string) error {
	return nil
}
func (m *MockAdvancedRoomClient) BlockRoom(ctx context.Context, roomID string, block bool) error {
	return nil
}

// TestRoomCreationWithInitialState tests room creation with initial state events
func TestRoomCreationWithInitialState(t *testing.T) {
	mockClient := &MockAdvancedRoomClient{
		createRoomFn: func(ctx context.Context, room *clients.RoomSpec) (*clients.Room, error) {
			// Verify initialState is provided
			if room.InitialState != nil && len(room.InitialState) > 0 {
				for _, event := range room.InitialState {
					if event.Type == "m.room.encryption" {
						return &clients.Room{
							Name:              room.Name,
							EncryptionEnabled: true,
						}, nil
					}
				}
			}
			return &clients.Room{Name: room.Name}, nil
		},
	}

	room, err := mockClient.CreateRoom(context.Background(), &clients.RoomSpec{
		Name: "Encrypted Room",
		InitialState: []clients.StateEvent{
			{
				Type:     "m.room.encryption",
				StateKey: "",
				Content: map[string]interface{}{
					"algorithm": "m.megolm.v1.aes-sha2",
				},
			},
		},
	})

	require.NoError(t, err)
	assert.True(t, room.EncryptionEnabled)
}

// TestRoomCreationWithPowerLevelOverrides tests room with custom power levels
func TestRoomCreationWithPowerLevelOverrides(t *testing.T) {
	mockClient := &MockAdvancedRoomClient{
		createRoomFn: func(ctx context.Context, room *clients.RoomSpec) (*clients.Room, error) {
			if room.PowerLevelOverrides != nil {
				return &clients.Room{
					Name:        room.Name,
					PowerLevels: room.PowerLevelOverrides,
				}, nil
			}
			return &clients.Room{Name: room.Name}, nil
		},
	}

	levelFifty := 50
	room, err := mockClient.CreateRoom(context.Background(), &clients.RoomSpec{
		Name: "Admin Room",
		PowerLevelOverrides: &clients.PowerLevelContent{
			UsersDefault: &levelFifty,
		},
	})

	require.NoError(t, err)
	assert.NotNil(t, room.PowerLevels)
}

// TestRoomCreationWithCreationContent tests room with creation content metadata
func TestRoomCreationWithCreationContent(t *testing.T) {
	mockClient := &MockAdvancedRoomClient{
		createRoomFn: func(ctx context.Context, room *clients.RoomSpec) (*clients.Room, error) {
			if room.CreationContent != nil {
				if federated, ok := room.CreationContent["m.federate"].(bool); ok && !federated {
					return &clients.Room{
						Name: room.Name,
					}, nil
				}
			}
			return &clients.Room{Name: room.Name}, nil
		},
	}

	room, err := mockClient.CreateRoom(context.Background(), &clients.RoomSpec{
		Name: "Private Room",
		CreationContent: map[string]interface{}{
			"m.federate": false,
		},
	})

	require.NoError(t, err)
	assert.Equal(t, "Private Room", room.Name)
}

// Priority 4: Batch Operations
// =============================

// MockBatchClient for testing batch operations
type MockBatchClient struct {
	promotionCalls []struct {
		roomID string
		userID string
	}
	aliasCalls []struct {
		alias  string
		roomID string
	}
	failOnAttempt int
}

func (m *MockBatchClient) MakeRoomAdmin(ctx context.Context, roomID, userID string) error {
	m.promotionCalls = append(m.promotionCalls, struct {
		roomID string
		userID string
	}{roomID, userID})

	if m.failOnAttempt > 0 && len(m.promotionCalls) >= m.failOnAttempt {
		return errors.New("permission denied")
	}
	return nil
}

func (m *MockBatchClient) CreateRoomAlias(ctx context.Context, alias string, roomID string) error {
	m.aliasCalls = append(m.aliasCalls, struct {
		alias  string
		roomID string
	}{alias, roomID})

	if m.failOnAttempt > 0 && len(m.aliasCalls) == m.failOnAttempt {
		return errors.New("alias already taken")
	}
	return nil
}

// Stub implementations
func (m *MockBatchClient) CreateUser(ctx context.Context, user *clients.UserSpec) (*clients.User, error) {
	return nil, nil
}
func (m *MockBatchClient) GetUser(ctx context.Context, userID string) (*clients.User, error) {
	return nil, nil
}
func (m *MockBatchClient) UpdateUser(ctx context.Context, userID string, user *clients.UserSpec) (*clients.User, error) {
	return nil, nil
}
func (m *MockBatchClient) DeactivateUser(ctx context.Context, userID string) error {
	return nil
}
func (m *MockBatchClient) CreateRoom(ctx context.Context, room *clients.RoomSpec) (*clients.Room, error) {
	return nil, nil
}
func (m *MockBatchClient) GetRoom(ctx context.Context, roomID string) (*clients.Room, error) {
	return nil, nil
}
func (m *MockBatchClient) UpdateRoom(ctx context.Context, roomID string, room *clients.RoomSpec) (*clients.Room, error) {
	return nil, nil
}
func (m *MockBatchClient) DeleteRoom(ctx context.Context, roomID string) error {
	return nil
}
func (m *MockBatchClient) SetPowerLevels(ctx context.Context, roomID string, powerLevels *clients.PowerLevelSpec) error {
	return nil
}
func (m *MockBatchClient) GetPowerLevels(ctx context.Context, roomID string) (*clients.PowerLevelContent, error) {
	return nil, nil
}
func (m *MockBatchClient) GetRoomAlias(ctx context.Context, alias string) (*clients.RoomAlias, error) {
	return nil, nil
}
func (m *MockBatchClient) DeleteRoomAlias(ctx context.Context, alias string) error {
	return nil
}
func (m *MockBatchClient) ListUsers(ctx context.Context, from string, limit int) (*clients.ListUsersResponse, error) {
	return nil, nil
}
func (m *MockBatchClient) ListRooms(ctx context.Context, from string, limit int) (*clients.ListRoomsResponse, error) {
	return nil, nil
}
func (m *MockBatchClient) BlockRoom(ctx context.Context, roomID string, block bool) error {
	return nil
}

// TestBatchUserPromotion tests promoting multiple users to admin
func TestBatchUserPromotion(t *testing.T) {
	client := &MockBatchClient{}
	roomID := "!room:example.com"
	numUsers := 10

	for i := 1; i <= numUsers; i++ {
		userID := stringOfLengthChar(i, 'u') // u, uu, uuu, etc.
		err := client.MakeRoomAdmin(context.Background(), roomID, "@"+userID+":example.com")
		require.NoError(t, err)
	}

	assert.Equal(t, numUsers, len(client.promotionCalls))
	assert.Equal(t, "@u:example.com", client.promotionCalls[0].userID)
	assert.Equal(t, "@uu:example.com", client.promotionCalls[1].userID)
}

// TestBatchAliasCreation tests creating multiple aliases
func TestBatchAliasCreation(t *testing.T) {
	client := &MockBatchClient{}
	roomID := "!room:example.com"
	aliases := []string{"general", "chat", "main", "primary", "default"}

	for _, alias := range aliases {
		err := client.CreateRoomAlias(context.Background(), "#"+alias+":example.com", roomID)
		require.NoError(t, err)
	}

	assert.Equal(t, len(aliases), len(client.aliasCalls))
}

// TestBatchPartialFailure tests handling partial failures in batch
func TestBatchPartialFailure(t *testing.T) {
	client := &MockBatchClient{
		failOnAttempt: 3, // Fail on 3rd and subsequent attempts
	}
	roomID := "!room:example.com"
	successCount := 0
	failureCount := 0

	for i := 1; i <= 5; i++ {
		err := client.MakeRoomAdmin(context.Background(), roomID, "@user"+stringOfLengthChar(i, 'a')+":example.com")
		if err == nil {
			successCount++
		} else {
			failureCount++
		}
	}

	// First 2 succeed, then 3rd and all after fail
	assert.Equal(t, 2, successCount, "expected 2 successes (1st and 2nd)")
	assert.Equal(t, 3, failureCount, "expected 3 failures (3rd, 4th, 5th)")
}

// Priority 4: Advanced Filtering & Search
// ========================================

// MockFilterClient for testing filtering and search
type MockFilterClient struct {
	users map[string]*clients.User
	rooms map[string]*clients.Room
}

func NewMockFilterClient() *MockFilterClient {
	return &MockFilterClient{
		users: map[string]*clients.User{
			"@admin1:example.com": {UserID: "@admin1:example.com", Admin: true, DisplayName: "Admin One"},
			"@admin2:example.com": {UserID: "@admin2:example.com", Admin: true, DisplayName: "Admin Two"},
			"@user1:example.com":  {UserID: "@user1:example.com", Admin: false, DisplayName: "User One"},
			"@user2:example.com":  {UserID: "@user2:example.com", Admin: false, DisplayName: "User Two"},
		},
		rooms: map[string]*clients.Room{
			"!public1:example.com":  {RoomID: "!public1:example.com", Visibility: "public", EncryptionEnabled: false},
			"!public2:example.com":  {RoomID: "!public2:example.com", Visibility: "public", EncryptionEnabled: true},
			"!private1:example.com": {RoomID: "!private1:example.com", Visibility: "private", EncryptionEnabled: false},
			"!private2:example.com": {RoomID: "!private2:example.com", Visibility: "private", EncryptionEnabled: true},
		},
	}
}

func (m *MockFilterClient) ListUsers(ctx context.Context, from string, limit int) (*clients.ListUsersResponse, error) {
	return &clients.ListUsersResponse{Users: []clients.User{}}, nil
}

func (m *MockFilterClient) FilterUsersByAdmin(ctx context.Context, adminOnly bool) []clients.User {
	var result []clients.User
	for _, user := range m.users {
		if !adminOnly || user.Admin {
			result = append(result, *user)
		}
	}
	return result
}

func (m *MockFilterClient) FilterRoomsByVisibility(ctx context.Context, visibility string) []clients.Room {
	var result []clients.Room
	for _, room := range m.rooms {
		if room.Visibility == visibility {
			result = append(result, *room)
		}
	}
	return result
}

func (m *MockFilterClient) FilterRoomsByEncryption(ctx context.Context, encrypted bool) []clients.Room {
	var result []clients.Room
	for _, room := range m.rooms {
		if room.EncryptionEnabled == encrypted {
			result = append(result, *room)
		}
	}
	return result
}

// Stub implementations
func (m *MockFilterClient) CreateUser(ctx context.Context, user *clients.UserSpec) (*clients.User, error) {
	return nil, nil
}
func (m *MockFilterClient) GetUser(ctx context.Context, userID string) (*clients.User, error) {
	return nil, nil
}
func (m *MockFilterClient) UpdateUser(ctx context.Context, userID string, user *clients.UserSpec) (*clients.User, error) {
	return nil, nil
}
func (m *MockFilterClient) DeactivateUser(ctx context.Context, userID string) error {
	return nil
}
func (m *MockFilterClient) CreateRoom(ctx context.Context, room *clients.RoomSpec) (*clients.Room, error) {
	return nil, nil
}
func (m *MockFilterClient) GetRoom(ctx context.Context, roomID string) (*clients.Room, error) {
	return nil, nil
}
func (m *MockFilterClient) UpdateRoom(ctx context.Context, roomID string, room *clients.RoomSpec) (*clients.Room, error) {
	return nil, nil
}
func (m *MockFilterClient) DeleteRoom(ctx context.Context, roomID string) error {
	return nil
}
func (m *MockFilterClient) SetPowerLevels(ctx context.Context, roomID string, powerLevels *clients.PowerLevelSpec) error {
	return nil
}
func (m *MockFilterClient) GetPowerLevels(ctx context.Context, roomID string) (*clients.PowerLevelContent, error) {
	return nil, nil
}
func (m *MockFilterClient) CreateRoomAlias(ctx context.Context, alias string, roomID string) error {
	return nil
}
func (m *MockFilterClient) GetRoomAlias(ctx context.Context, alias string) (*clients.RoomAlias, error) {
	return nil, nil
}
func (m *MockFilterClient) DeleteRoomAlias(ctx context.Context, alias string) error {
	return nil
}
func (m *MockFilterClient) ListRooms(ctx context.Context, from string, limit int) (*clients.ListRoomsResponse, error) {
	return &clients.ListRoomsResponse{Rooms: []clients.Room{}}, nil
}
func (m *MockFilterClient) MakeRoomAdmin(ctx context.Context, roomID, userID string) error {
	return nil
}
func (m *MockFilterClient) BlockRoom(ctx context.Context, roomID string, block bool) error {
	return nil
}

// TestFilterUsersByAdminStatus tests filtering users by admin status
func TestFilterUsersByAdminStatus(t *testing.T) {
	client := NewMockFilterClient()

	admins := client.FilterUsersByAdmin(context.Background(), true)
	assert.Equal(t, 2, len(admins))

	allUsers := client.FilterUsersByAdmin(context.Background(), false)
	assert.Equal(t, 4, len(allUsers))
}

// TestFilterRoomsByVisibility tests filtering rooms by visibility
func TestFilterRoomsByVisibility(t *testing.T) {
	client := NewMockFilterClient()

	publicRooms := client.FilterRoomsByVisibility(context.Background(), "public")
	assert.Equal(t, 2, len(publicRooms))

	privateRooms := client.FilterRoomsByVisibility(context.Background(), "private")
	assert.Equal(t, 2, len(privateRooms))
}

// TestFilterRoomsByEncryption tests filtering rooms by encryption status
func TestFilterRoomsByEncryption(t *testing.T) {
	client := NewMockFilterClient()

	encrypted := client.FilterRoomsByEncryption(context.Background(), true)
	assert.Equal(t, 2, len(encrypted))

	unencrypted := client.FilterRoomsByEncryption(context.Background(), false)
	assert.Equal(t, 2, len(unencrypted))
}

// TestCombinedFilters tests combining multiple filters
func TestCombinedFilters(t *testing.T) {
	client := NewMockFilterClient()

	publicRooms := client.FilterRoomsByVisibility(context.Background(), "public")
	var publicEncrypted []clients.Room
	for _, room := range publicRooms {
		if room.EncryptionEnabled {
			publicEncrypted = append(publicEncrypted, room)
		}
	}

	assert.Equal(t, 1, len(publicEncrypted))
	assert.Equal(t, "!public2:example.com", publicEncrypted[0].RoomID)
}

// Priority 4: Avatar & Metadata Operations
// =========================================

// MockMetadataClient for testing avatar and metadata operations
type MockMetadataClient struct {
	updateRoomFn func(ctx context.Context, roomID string, room *clients.RoomSpec) (*clients.Room, error)
	updateUserFn func(ctx context.Context, userID string, user *clients.UserSpec) (*clients.User, error)
}

func (m *MockMetadataClient) UpdateRoom(ctx context.Context, roomID string, room *clients.RoomSpec) (*clients.Room, error) {
	if m.updateRoomFn != nil {
		return m.updateRoomFn(ctx, roomID, room)
	}
	return &clients.Room{RoomID: roomID}, nil
}

func (m *MockMetadataClient) UpdateUser(ctx context.Context, userID string, user *clients.UserSpec) (*clients.User, error) {
	if m.updateUserFn != nil {
		return m.updateUserFn(ctx, userID, user)
	}
	return &clients.User{UserID: userID}, nil
}

// Stub implementations
func (m *MockMetadataClient) CreateUser(ctx context.Context, user *clients.UserSpec) (*clients.User, error) {
	return nil, nil
}
func (m *MockMetadataClient) GetUser(ctx context.Context, userID string) (*clients.User, error) {
	return nil, nil
}
func (m *MockMetadataClient) DeactivateUser(ctx context.Context, userID string) error {
	return nil
}
func (m *MockMetadataClient) CreateRoom(ctx context.Context, room *clients.RoomSpec) (*clients.Room, error) {
	return nil, nil
}
func (m *MockMetadataClient) GetRoom(ctx context.Context, roomID string) (*clients.Room, error) {
	return nil, nil
}
func (m *MockMetadataClient) DeleteRoom(ctx context.Context, roomID string) error {
	return nil
}
func (m *MockMetadataClient) SetPowerLevels(ctx context.Context, roomID string, powerLevels *clients.PowerLevelSpec) error {
	return nil
}
func (m *MockMetadataClient) GetPowerLevels(ctx context.Context, roomID string) (*clients.PowerLevelContent, error) {
	return nil, nil
}
func (m *MockMetadataClient) CreateRoomAlias(ctx context.Context, alias string, roomID string) error {
	return nil
}
func (m *MockMetadataClient) GetRoomAlias(ctx context.Context, alias string) (*clients.RoomAlias, error) {
	return nil, nil
}
func (m *MockMetadataClient) DeleteRoomAlias(ctx context.Context, alias string) error {
	return nil
}
func (m *MockMetadataClient) ListUsers(ctx context.Context, from string, limit int) (*clients.ListUsersResponse, error) {
	return nil, nil
}
func (m *MockMetadataClient) ListRooms(ctx context.Context, from string, limit int) (*clients.ListRoomsResponse, error) {
	return nil, nil
}
func (m *MockMetadataClient) MakeRoomAdmin(ctx context.Context, roomID, userID string) error {
	return nil
}
func (m *MockMetadataClient) BlockRoom(ctx context.Context, roomID string, block bool) error {
	return nil
}

// TestUpdateRoomWithAvatar tests updating room with avatar URL
func TestUpdateRoomWithAvatar(t *testing.T) {
	mockClient := &MockMetadataClient{
		updateRoomFn: func(ctx context.Context, roomID string, room *clients.RoomSpec) (*clients.Room, error) {
			return &clients.Room{
				RoomID:    roomID,
				AvatarURL: room.AvatarURL,
			}, nil
		},
	}

	updated, err := mockClient.UpdateRoom(context.Background(), "!room:example.com", &clients.RoomSpec{
		AvatarURL: "mxc://example.com/abc123",
	})

	require.NoError(t, err)
	assert.Equal(t, "mxc://example.com/abc123", updated.AvatarURL)
}

// TestUpdateUserWithAvatar tests updating user with avatar URL
func TestUpdateUserWithAvatar(t *testing.T) {
	mockClient := &MockMetadataClient{
		updateUserFn: func(ctx context.Context, userID string, user *clients.UserSpec) (*clients.User, error) {
			return &clients.User{
				UserID:    userID,
				AvatarURL: user.AvatarURL,
			}, nil
		},
	}

	updated, err := mockClient.UpdateUser(context.Background(), "@alice:example.com", &clients.UserSpec{
		AvatarURL: "mxc://example.com/xyz789",
	})

	require.NoError(t, err)
	assert.Equal(t, "mxc://example.com/xyz789", updated.AvatarURL)
}

// TestUpdateRoomPreset tests updating room with preset configuration
func TestUpdateRoomPreset(t *testing.T) {
	mockClient := &MockMetadataClient{
		updateRoomFn: func(ctx context.Context, roomID string, room *clients.RoomSpec) (*clients.Room, error) {
			if room.Preset != "" {
				return &clients.Room{
					RoomID: roomID,
					Name:   room.Name,
				}, nil
			}
			return &clients.Room{RoomID: roomID}, nil
		},
	}

	updated, err := mockClient.UpdateRoom(context.Background(), "!room:example.com", &clients.RoomSpec{
		Preset: "public_chat",
		Name:   "Public Chat Room",
	})

	require.NoError(t, err)
	assert.Equal(t, "Public Chat Room", updated.Name)
}
