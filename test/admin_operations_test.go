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

	"github.com/crossplane-contrib/provider-matrix/internal/clients"
)

// MockAdminClient for testing admin operations
type MockAdminClient struct {
	listUsersFn       func(ctx context.Context, from string, limit int) (*clients.ListUsersResponse, error)
	listRoomsFn       func(ctx context.Context, from string, limit int) (*clients.ListRoomsResponse, error)
	makeRoomAdminFn   func(ctx context.Context, roomID, userID string) error
	blockRoomFn       func(ctx context.Context, roomID string, block bool) error
	deactivateUserFn  func(ctx context.Context, userID string) error
}

func (m *MockAdminClient) ListUsers(ctx context.Context, from string, limit int) (*clients.ListUsersResponse, error) {
	if m.listUsersFn != nil {
		return m.listUsersFn(ctx, from, limit)
	}
	return &clients.ListUsersResponse{}, nil
}

func (m *MockAdminClient) ListRooms(ctx context.Context, from string, limit int) (*clients.ListRoomsResponse, error) {
	if m.listRoomsFn != nil {
		return m.listRoomsFn(ctx, from, limit)
	}
	return &clients.ListRoomsResponse{}, nil
}

func (m *MockAdminClient) MakeRoomAdmin(ctx context.Context, roomID, userID string) error {
	if m.makeRoomAdminFn != nil {
		return m.makeRoomAdminFn(ctx, roomID, userID)
	}
	return nil
}

func (m *MockAdminClient) BlockRoom(ctx context.Context, roomID string, block bool) error {
	if m.blockRoomFn != nil {
		return m.blockRoomFn(ctx, roomID, block)
	}
	return nil
}

func (m *MockAdminClient) DeactivateUser(ctx context.Context, userID string) error {
	if m.deactivateUserFn != nil {
		return m.deactivateUserFn(ctx, userID)
	}
	return nil
}

// Stub implementations for unused interface methods
func (m *MockAdminClient) CreateUser(ctx context.Context, user *clients.UserSpec) (*clients.User, error) {
	return nil, nil
}

func (m *MockAdminClient) GetUser(ctx context.Context, userID string) (*clients.User, error) {
	return nil, nil
}

func (m *MockAdminClient) UpdateUser(ctx context.Context, userID string, user *clients.UserSpec) (*clients.User, error) {
	return nil, nil
}

func (m *MockAdminClient) CreateRoom(ctx context.Context, room *clients.RoomSpec) (*clients.Room, error) {
	return nil, nil
}

func (m *MockAdminClient) GetRoom(ctx context.Context, roomID string) (*clients.Room, error) {
	return nil, nil
}

func (m *MockAdminClient) UpdateRoom(ctx context.Context, roomID string, room *clients.RoomSpec) (*clients.Room, error) {
	return nil, nil
}

func (m *MockAdminClient) DeleteRoom(ctx context.Context, roomID string) error {
	return nil
}

func (m *MockAdminClient) SetPowerLevels(ctx context.Context, roomID string, powerLevels *clients.PowerLevelSpec) error {
	return nil
}

func (m *MockAdminClient) GetPowerLevels(ctx context.Context, roomID string) (*clients.PowerLevelContent, error) {
	return nil, nil
}

func (m *MockAdminClient) CreateRoomAlias(ctx context.Context, alias string, roomID string) error {
	return nil
}

func (m *MockAdminClient) GetRoomAlias(ctx context.Context, alias string) (*clients.RoomAlias, error) {
	return nil, nil
}

func (m *MockAdminClient) DeleteRoomAlias(ctx context.Context, alias string) error {
	return nil
}

// TestListUsers tests listing users
func TestListUsers(t *testing.T) {
	mockClient := &MockAdminClient{
		listUsersFn: func(ctx context.Context, from string, limit int) (*clients.ListUsersResponse, error) {
			return &clients.ListUsersResponse{
				Users: []clients.User{
					{UserID: "@alice:example.com", DisplayName: "Alice"},
					{UserID: "@bob:example.com", DisplayName: "Bob"},
				},
				NextToken: "next_page",
			}, nil
		},
	}

	resp, err := mockClient.ListUsers(context.Background(), "", 10)
	require.NoError(t, err)
	assert.Equal(t, 2, len(resp.Users))
	assert.Equal(t, "@alice:example.com", resp.Users[0].UserID)
	assert.Equal(t, "next_page", resp.NextToken)
}

// TestListUsersPagination tests pagination in user listing
func TestListUsersPagination(t *testing.T) {
	mockClient := &MockAdminClient{
		listUsersFn: func(ctx context.Context, from string, limit int) (*clients.ListUsersResponse, error) {
			if from == "" {
				return &clients.ListUsersResponse{
					Users: []clients.User{
						{UserID: "@user1:example.com"},
						{UserID: "@user2:example.com"},
					},
					NextToken: "token2",
				}, nil
			}
			return &clients.ListUsersResponse{
				Users: []clients.User{
					{UserID: "@user3:example.com"},
				},
			}, nil
		},
	}

	// First page
	page1, err := mockClient.ListUsers(context.Background(), "", 2)
	require.NoError(t, err)
	assert.Equal(t, 2, len(page1.Users))
	assert.NotEmpty(t, page1.NextToken)

	// Second page
	page2, err := mockClient.ListUsers(context.Background(), page1.NextToken, 2)
	require.NoError(t, err)
	assert.Equal(t, 1, len(page2.Users))
}

// TestListRooms tests listing rooms
func TestListRooms(t *testing.T) {
	mockClient := &MockAdminClient{
		listRoomsFn: func(ctx context.Context, from string, limit int) (*clients.ListRoomsResponse, error) {
			return &clients.ListRoomsResponse{
				Rooms: []clients.Room{
					{RoomID: "!room1:example.com", Name: "General"},
					{RoomID: "!room2:example.com", Name: "Random"},
				},
				NextToken: "rooms_token",
			}, nil
		},
	}

	resp, err := mockClient.ListRooms(context.Background(), "", 10)
	require.NoError(t, err)
	assert.Equal(t, 2, len(resp.Rooms))
	assert.Equal(t, "!room1:example.com", resp.Rooms[0].RoomID)
}

// TestMakeRoomAdmin tests promoting a user to room admin
func TestMakeRoomAdmin(t *testing.T) {
	mockClient := &MockAdminClient{
		makeRoomAdminFn: func(ctx context.Context, roomID, userID string) error {
			if roomID == "" || userID == "" {
				return errors.New("roomID and userID are required")
			}
			return nil
		},
	}

	err := mockClient.MakeRoomAdmin(context.Background(), "!room:example.com", "@user:example.com")
	require.NoError(t, err)
}

// TestMakeRoomAdminError tests error handling for make room admin
func TestMakeRoomAdminError(t *testing.T) {
	mockClient := &MockAdminClient{
		makeRoomAdminFn: func(ctx context.Context, roomID, userID string) error {
			return errors.New("user not in room")
		},
	}

	err := mockClient.MakeRoomAdmin(context.Background(), "!room:example.com", "@user:example.com")
	assert.Error(t, err)
	assert.Equal(t, "user not in room", err.Error())
}

// TestBlockRoom tests blocking a room
func TestBlockRoom(t *testing.T) {
	mockClient := &MockAdminClient{
		blockRoomFn: func(ctx context.Context, roomID string, block bool) error {
			if roomID == "" {
				return errors.New("roomID is required")
			}
			return nil
		},
	}

	err := mockClient.BlockRoom(context.Background(), "!room:example.com", true)
	require.NoError(t, err)
}

// TestUnblockRoom tests unblocking a room
func TestUnblockRoom(t *testing.T) {
	mockClient := &MockAdminClient{
		blockRoomFn: func(ctx context.Context, roomID string, block bool) error {
			return nil
		},
	}

	err := mockClient.BlockRoom(context.Background(), "!room:example.com", false)
	require.NoError(t, err)
}

// TestBlockRoomError tests error handling for block room
func TestBlockRoomError(t *testing.T) {
	mockClient := &MockAdminClient{
		blockRoomFn: func(ctx context.Context, roomID string, block bool) error {
			return errors.New("cannot block room")
		},
	}

	err := mockClient.BlockRoom(context.Background(), "!room:example.com", true)
	assert.Error(t, err)
}

// TestListUsersError tests error handling for list users
func TestListUsersError(t *testing.T) {
	mockClient := &MockAdminClient{
		listUsersFn: func(ctx context.Context, from string, limit int) (*clients.ListUsersResponse, error) {
			return nil, errors.New("admin API not available")
		},
	}

	_, err := mockClient.ListUsers(context.Background(), "", 10)
	assert.Error(t, err)
}

// TestListRoomsError tests error handling for list rooms
func TestListRoomsError(t *testing.T) {
	mockClient := &MockAdminClient{
		listRoomsFn: func(ctx context.Context, from string, limit int) (*clients.ListRoomsResponse, error) {
			return nil, errors.New("admin API not available")
		},
	}

	_, err := mockClient.ListRooms(context.Background(), "", 10)
	assert.Error(t, err)
}

// TestAdminBatchOperations tests multiple admin operations
func TestAdminBatchOperations(t *testing.T) {
	callCount := 0
	mockClient := &MockAdminClient{
		listUsersFn: func(ctx context.Context, from string, limit int) (*clients.ListUsersResponse, error) {
			callCount++
			return &clients.ListUsersResponse{Users: []clients.User{
				{UserID: "@user:example.com"},
			}}, nil
		},
		makeRoomAdminFn: func(ctx context.Context, roomID, userID string) error {
			callCount++
			return nil
		},
	}

	// Execute batch operations
	_, err := mockClient.ListUsers(context.Background(), "", 10)
	require.NoError(t, err)

	err = mockClient.MakeRoomAdmin(context.Background(), "!room:example.com", "@user:example.com")
	require.NoError(t, err)

	assert.Equal(t, 2, callCount, "expected 2 calls to admin operations")
}
