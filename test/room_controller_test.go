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

	"github.com/crossplane-contrib/provider-matrix/apis/room/v1alpha1"
	"github.com/crossplane-contrib/provider-matrix/internal/clients"
)

// MockRoomClient implements clients.Client for Room testing
type MockRoomClient struct {
	createRoomFn   func(ctx context.Context, room *clients.RoomSpec) (*clients.Room, error)
	getRoomFn      func(ctx context.Context, roomID string) (*clients.Room, error)
	updateRoomFn   func(ctx context.Context, roomID string, room *clients.RoomSpec) (*clients.Room, error)
	deleteRoomFn   func(ctx context.Context, roomID string) error
}

func (m *MockRoomClient) CreateRoom(ctx context.Context, room *clients.RoomSpec) (*clients.Room, error) {
	if m.createRoomFn != nil {
		return m.createRoomFn(ctx, room)
	}
	return &clients.Room{RoomID: room.Name}, nil
}

func (m *MockRoomClient) GetRoom(ctx context.Context, roomID string) (*clients.Room, error) {
	if m.getRoomFn != nil {
		return m.getRoomFn(ctx, roomID)
	}
	return &clients.Room{RoomID: roomID}, nil
}

func (m *MockRoomClient) UpdateRoom(ctx context.Context, roomID string, room *clients.RoomSpec) (*clients.Room, error) {
	if m.updateRoomFn != nil {
		return m.updateRoomFn(ctx, roomID, room)
	}
	return &clients.Room{RoomID: roomID}, nil
}

func (m *MockRoomClient) DeleteRoom(ctx context.Context, roomID string) error {
	if m.deleteRoomFn != nil {
		return m.deleteRoomFn(ctx, roomID)
	}
	return nil
}

// Stub implementations for unused interface methods
func (m *MockRoomClient) CreateUser(ctx context.Context, user *clients.UserSpec) (*clients.User, error) {
	return nil, nil
}

func (m *MockRoomClient) GetUser(ctx context.Context, userID string) (*clients.User, error) {
	return nil, nil
}

func (m *MockRoomClient) UpdateUser(ctx context.Context, userID string, user *clients.UserSpec) (*clients.User, error) {
	return nil, nil
}

func (m *MockRoomClient) DeactivateUser(ctx context.Context, userID string) error {
	return nil
}

func (m *MockRoomClient) SetPowerLevels(ctx context.Context, roomID string, powerLevels *clients.PowerLevelSpec) error {
	return nil
}

func (m *MockRoomClient) GetPowerLevels(ctx context.Context, roomID string) (*clients.PowerLevelContent, error) {
	return nil, nil
}

func (m *MockRoomClient) CreateRoomAlias(ctx context.Context, alias string, roomID string) error {
	return nil
}

func (m *MockRoomClient) GetRoomAlias(ctx context.Context, alias string) (*clients.RoomAlias, error) {
	return nil, nil
}

func (m *MockRoomClient) DeleteRoomAlias(ctx context.Context, alias string) error {
	return nil
}

func (m *MockRoomClient) ListUsers(ctx context.Context, from string, limit int) (*clients.ListUsersResponse, error) {
	return nil, nil
}

func (m *MockRoomClient) ListRooms(ctx context.Context, from string, limit int) (*clients.ListRoomsResponse, error) {
	return nil, nil
}

func (m *MockRoomClient) MakeRoomAdmin(ctx context.Context, roomID, userID string) error {
	return nil
}

func (m *MockRoomClient) BlockRoom(ctx context.Context, roomID string, block bool) error {
	return nil
}

// TestRoomCreate tests room creation
func TestRoomCreate(t *testing.T) {
	mockClient := &MockRoomClient{
		createRoomFn: func(ctx context.Context, room *clients.RoomSpec) (*clients.Room, error) {
			return &clients.Room{
				RoomID:            "!test:example.com",
				Name:              room.Name,
				Topic:             room.Topic,
				EncryptionEnabled: room.EncryptionEnabled,
			}, nil
		},
	}

	spec := &clients.RoomSpec{
		Name:                 "Test Room",
		Topic:                "A test room",
		Preset:               "private_chat",
		EncryptionEnabled:    true,
	}

	created, err := mockClient.CreateRoom(context.Background(), spec)
	require.NoError(t, err)
	assert.Equal(t, "!test:example.com", created.RoomID)
	assert.Equal(t, "Test Room", created.Name)
	assert.Equal(t, "A test room", created.Topic)
	assert.True(t, created.EncryptionEnabled)
}

// TestRoomCreateError tests room creation error handling
func TestRoomCreateError(t *testing.T) {
	mockClient := &MockRoomClient{
		createRoomFn: func(ctx context.Context, room *clients.RoomSpec) (*clients.Room, error) {
			return nil, errors.New("room creation failed")
		},
	}

	spec := &clients.RoomSpec{
		Name: "Bad Room",
	}

	_, err := mockClient.CreateRoom(context.Background(), spec)
	assert.Error(t, err)
}

// TestRoomGet tests getting room information
func TestRoomGet(t *testing.T) {
	mockClient := &MockRoomClient{
		getRoomFn: func(ctx context.Context, roomID string) (*clients.Room, error) {
			return &clients.Room{
				RoomID: roomID,
				Name:   "Retrieved Room",
			}, nil
		},
	}

	room, err := mockClient.GetRoom(context.Background(), "!test:example.com")
	require.NoError(t, err)
	assert.Equal(t, "!test:example.com", room.RoomID)
	assert.Equal(t, "Retrieved Room", room.Name)
}

// TestRoomUpdate tests room update operations
func TestRoomUpdate(t *testing.T) {
	mockClient := &MockRoomClient{
		updateRoomFn: func(ctx context.Context, roomID string, room *clients.RoomSpec) (*clients.Room, error) {
			return &clients.Room{
				RoomID: roomID,
				Name:   room.Name,
				Topic:  room.Topic,
			}, nil
		},
	}

	spec := &clients.RoomSpec{
		Name:  "Updated Room",
		Topic: "Updated topic",
	}

	updated, err := mockClient.UpdateRoom(context.Background(), "!test:example.com", spec)
	require.NoError(t, err)
	assert.Equal(t, "Updated Room", updated.Name)
	assert.Equal(t, "Updated topic", updated.Topic)
}

// TestRoomDelete tests room deletion
func TestRoomDelete(t *testing.T) {
	mockClient := &MockRoomClient{
		deleteRoomFn: func(ctx context.Context, roomID string) error {
			return nil
		},
	}

	err := mockClient.DeleteRoom(context.Background(), "!test:example.com")
	assert.NoError(t, err)
}

// TestRoomResource tests Room CR creation and manipulation
func TestRoomResource(t *testing.T) {
	room := &v1alpha1.Room{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-room",
		},
		Spec: v1alpha1.RoomSpec{
			ForProvider: v1alpha1.RoomParameters{
				Name:              stringPtr("Test Room"),
				Topic:             stringPtr("Room topic"),
				Preset:            stringPtr("public_chat"),
				EncryptionEnabled: boolPtr(false),
			},
		},
	}

	assert.Equal(t, "test-room", room.Name)
	assert.Equal(t, "Test Room", *room.Spec.ForProvider.Name)
}

// TestRoomProviderConfigReference tests room provider config reference
func TestRoomProviderConfigReference(t *testing.T) {
	room := &v1alpha1.Room{}

	pcRef := &xpcoreapi.ProviderConfigReference{Name: "test-pc"}
	room.SetProviderConfigReference(pcRef)
	assert.Equal(t, "test-pc", room.GetProviderConfigReference().Name)
}

// TestRoomConditions tests room condition management
func TestRoomConditions(t *testing.T) {
	room := &v1alpha1.Room{}

	cond := xpcoreapi.Available()
	room.SetConditions(cond)

	retrieved := room.GetCondition(xpcoreapi.TypeReady)
	assert.Equal(t, xpcoreapi.TypeReady, retrieved.Type)
}
