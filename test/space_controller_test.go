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
	"github.com/crossplane-contrib/provider-matrix/apis/space/v1alpha1"
	"github.com/crossplane-contrib/provider-matrix/internal/clients"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

// MockSpaceClient implements clients.Client for Space testing
type MockSpaceClient struct {
	createRoomFn func(ctx context.Context, room *clients.RoomSpec) (*clients.Room, error)
	getRoomFn    func(ctx context.Context, roomID string) (*clients.Room, error)
	updateRoomFn func(ctx context.Context, roomID string, room *clients.RoomSpec) (*clients.Room, error)
	deleteRoomFn func(ctx context.Context, roomID string) error
}

func (m *MockSpaceClient) CreateRoom(ctx context.Context, room *clients.RoomSpec) (*clients.Room, error) {
	if m.createRoomFn != nil {
		return m.createRoomFn(ctx, room)
	}
	return &clients.Room{RoomID: "!space:example.com"}, nil
}

func (m *MockSpaceClient) GetRoom(ctx context.Context, roomID string) (*clients.Room, error) {
	if m.getRoomFn != nil {
		return m.getRoomFn(ctx, roomID)
	}
	return &clients.Room{RoomID: roomID}, nil
}

func (m *MockSpaceClient) UpdateRoom(ctx context.Context, roomID string, room *clients.RoomSpec) (*clients.Room, error) {
	if m.updateRoomFn != nil {
		return m.updateRoomFn(ctx, roomID, room)
	}
	return &clients.Room{RoomID: roomID}, nil
}

func (m *MockSpaceClient) DeleteRoom(ctx context.Context, roomID string) error {
	if m.deleteRoomFn != nil {
		return m.deleteRoomFn(ctx, roomID)
	}
	return nil
}

// Stub implementations for unused interface methods
func (m *MockSpaceClient) CreateUser(ctx context.Context, user *clients.UserSpec) (*clients.User, error) {
	return nil, nil
}

func (m *MockSpaceClient) GetUser(ctx context.Context, userID string) (*clients.User, error) {
	return nil, nil
}

func (m *MockSpaceClient) UpdateUser(ctx context.Context, userID string, user *clients.UserSpec) (*clients.User, error) {
	return nil, nil
}

func (m *MockSpaceClient) DeactivateUser(ctx context.Context, userID string) error {
	return nil
}

func (m *MockSpaceClient) SetPowerLevels(ctx context.Context, roomID string, powerLevels *clients.PowerLevelSpec) error {
	return nil
}

func (m *MockSpaceClient) GetPowerLevels(ctx context.Context, roomID string) (*clients.PowerLevelContent, error) {
	return nil, nil
}

func (m *MockSpaceClient) CreateRoomAlias(ctx context.Context, alias string, roomID string) error {
	return nil
}

func (m *MockSpaceClient) GetRoomAlias(ctx context.Context, alias string) (*clients.RoomAlias, error) {
	return nil, nil
}

func (m *MockSpaceClient) DeleteRoomAlias(ctx context.Context, alias string) error {
	return nil
}

func (m *MockSpaceClient) ListUsers(ctx context.Context, from string, limit int) (*clients.ListUsersResponse, error) {
	return nil, nil
}

func (m *MockSpaceClient) ListRooms(ctx context.Context, from string, limit int) (*clients.ListRoomsResponse, error) {
	return nil, nil
}

func (m *MockSpaceClient) MakeRoomAdmin(ctx context.Context, roomID, userID string) error {
	return nil
}

func (m *MockSpaceClient) BlockRoom(ctx context.Context, roomID string, block bool) error {
	return nil
}

// TestSpaceCreate tests space creation
func TestSpaceCreate(t *testing.T) {
	mockClient := &MockSpaceClient{
		createRoomFn: func(ctx context.Context, room *clients.RoomSpec) (*clients.Room, error) {
			return &clients.Room{
				RoomID: "!space:example.com",
				Name:   room.Name,
				Topic:  room.Topic,
			}, nil
		},
	}

	spec := &clients.RoomSpec{
		Name:  "Engineering",
		Topic: "Engineering organization space",
	}

	created, err := mockClient.CreateRoom(context.Background(), spec)
	require.NoError(t, err)
	assert.Equal(t, "!space:example.com", created.RoomID)
	assert.Equal(t, "Engineering", created.Name)
}

// TestSpaceGet tests retrieving space information
func TestSpaceGet(t *testing.T) {
	mockClient := &MockSpaceClient{
		getRoomFn: func(ctx context.Context, roomID string) (*clients.Room, error) {
			return &clients.Room{
				RoomID: roomID,
				Name:   "Operations",
				Topic:  "Operations space",
			}, nil
		},
	}

	space, err := mockClient.GetRoom(context.Background(), "!space:example.com")
	require.NoError(t, err)
	assert.Equal(t, "!space:example.com", space.RoomID)
	assert.Equal(t, "Operations", space.Name)
}

// TestSpaceUpdate tests space updates
func TestSpaceUpdate(t *testing.T) {
	mockClient := &MockSpaceClient{
		updateRoomFn: func(ctx context.Context, roomID string, room *clients.RoomSpec) (*clients.Room, error) {
			return &clients.Room{
				RoomID: roomID,
				Name:   room.Name,
				Topic:  room.Topic,
			}, nil
		},
	}

	spec := &clients.RoomSpec{
		Name:  "Updated Space",
		Topic: "Updated topic",
	}

	updated, err := mockClient.UpdateRoom(context.Background(), "!space:example.com", spec)
	require.NoError(t, err)
	assert.Equal(t, "Updated Space", updated.Name)
}

// TestSpaceDelete tests space deletion
func TestSpaceDelete(t *testing.T) {
	mockClient := &MockSpaceClient{
		deleteRoomFn: func(ctx context.Context, roomID string) error {
			return nil
		},
	}

	err := mockClient.DeleteRoom(context.Background(), "!space:example.com")
	assert.NoError(t, err)
}

// TestSpaceResource tests Space CR creation
func TestSpaceResource(t *testing.T) {
	space := &v1alpha1.Space{
		ObjectMeta: metav1.ObjectMeta{
			Name: "eng-space",
		},
		Spec: v1alpha1.SpaceSpec{
			ForProvider: v1alpha1.SpaceParameters{
				Name:  stringPtr("Engineering"),
				Topic: stringPtr("Engineering organization"),
			},
		},
	}

	assert.Equal(t, "eng-space", space.Name)
	assert.Equal(t, "Engineering", *space.Spec.ForProvider.Name)
}

// TestSpaceResourceMetadata tests Space resource metadata
func TestSpaceResourceMetadata(t *testing.T) {
	space := &v1alpha1.Space{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "eng-space",
			Namespace: "default",
		},
	}

	assert.Equal(t, "eng-space", space.Name)
	assert.Equal(t, "default", space.Namespace)
}
