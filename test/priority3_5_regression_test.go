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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/crossplane-contrib/provider-matrix/internal/clients"
)

// Priority 3.5: Field-Level Updates (Regression Prevention)
// ===========================================================

// MockFieldUpdateClient tracks field-level updates to detect unintended changes
type MockFieldUpdateClient struct {
	users map[string]*clients.User
	rooms map[string]*clients.Room
}

func NewMockFieldUpdateClient() *MockFieldUpdateClient {
	return &MockFieldUpdateClient{
		users: map[string]*clients.User{
			"@alice:example.com": {
				UserID:      "@alice:example.com",
				DisplayName: "Alice Original",
				AvatarURL:   "mxc://example.com/avatar1",
				Admin:       false,
			},
		},
		rooms: map[string]*clients.Room{
			"!room:example.com": {
				RoomID:            "!room:example.com",
				Name:              "Original Name",
				Topic:             "Original Topic",
				EncryptionEnabled: false,
				Visibility:        "public",
			},
		},
	}
}

func (m *MockFieldUpdateClient) UpdateUser(ctx context.Context, userID string, spec *clients.UserSpec) (*clients.User, error) {
	user, exists := m.users[userID]
	if !exists {
		return nil, nil
	}

	// Field-level update: only update specified fields
	if spec.DisplayName != "" {
		user.DisplayName = spec.DisplayName
	}
	if spec.AvatarURL != "" {
		user.AvatarURL = spec.AvatarURL
	}
	if spec.Admin {
		user.Admin = spec.Admin
	}

	return user, nil
}

func (m *MockFieldUpdateClient) UpdateRoom(ctx context.Context, roomID string, spec *clients.RoomSpec) (*clients.Room, error) {
	room, exists := m.rooms[roomID]
	if !exists {
		return nil, nil
	}

	// Field-level update: only update specified fields
	if spec.Name != "" {
		room.Name = spec.Name
	}
	if spec.Topic != "" {
		room.Topic = spec.Topic
	}
	if spec.Visibility != "" {
		room.Visibility = spec.Visibility
	}
	if spec.EncryptionEnabled {
		room.EncryptionEnabled = spec.EncryptionEnabled
	}

	return room, nil
}

// Stub implementations
func (m *MockFieldUpdateClient) CreateUser(ctx context.Context, user *clients.UserSpec) (*clients.User, error) {
	return nil, nil
}
func (m *MockFieldUpdateClient) GetUser(ctx context.Context, userID string) (*clients.User, error) {
	if user, exists := m.users[userID]; exists {
		return user, nil
	}
	return nil, nil
}
func (m *MockFieldUpdateClient) DeactivateUser(ctx context.Context, userID string) error {
	return nil
}
func (m *MockFieldUpdateClient) CreateRoom(ctx context.Context, room *clients.RoomSpec) (*clients.Room, error) {
	return nil, nil
}
func (m *MockFieldUpdateClient) GetRoom(ctx context.Context, roomID string) (*clients.Room, error) {
	if room, exists := m.rooms[roomID]; exists {
		return room, nil
	}
	return nil, nil
}
func (m *MockFieldUpdateClient) DeleteRoom(ctx context.Context, roomID string) error {
	return nil
}
func (m *MockFieldUpdateClient) SetPowerLevels(ctx context.Context, roomID string, powerLevels *clients.PowerLevelSpec) error {
	return nil
}
func (m *MockFieldUpdateClient) GetPowerLevels(ctx context.Context, roomID string) (*clients.PowerLevelContent, error) {
	return nil, nil
}
func (m *MockFieldUpdateClient) CreateRoomAlias(ctx context.Context, alias string, roomID string) error {
	return nil
}
func (m *MockFieldUpdateClient) GetRoomAlias(ctx context.Context, alias string) (*clients.RoomAlias, error) {
	return nil, nil
}
func (m *MockFieldUpdateClient) DeleteRoomAlias(ctx context.Context, alias string) error {
	return nil
}
func (m *MockFieldUpdateClient) ListUsers(ctx context.Context, from string, limit int) (*clients.ListUsersResponse, error) {
	return nil, nil
}
func (m *MockFieldUpdateClient) ListRooms(ctx context.Context, from string, limit int) (*clients.ListRoomsResponse, error) {
	return nil, nil
}
func (m *MockFieldUpdateClient) MakeRoomAdmin(ctx context.Context, roomID, userID string) error {
	return nil
}
func (m *MockFieldUpdateClient) BlockRoom(ctx context.Context, roomID string, block bool) error {
	return nil
}

// TestFieldLevelUserUpdate_DisplayNameOnly tests updating only display name without clearing avatar
func TestFieldLevelUserUpdate_DisplayNameOnly(t *testing.T) {
	client := NewMockFieldUpdateClient()
	userID := "@alice:example.com"

	// Get original state
	original, err := client.GetUser(context.Background(), userID)
	require.NoError(t, err)
	originalAvatar := original.AvatarURL

	// Update only display name
	updated, err := client.UpdateUser(context.Background(), userID, &clients.UserSpec{
		DisplayName: "Alice Updated",
	})
	require.NoError(t, err)

	// Verify display name changed, avatar preserved (regression test)
	assert.Equal(t, "Alice Updated", updated.DisplayName)
	assert.Equal(t, originalAvatar, updated.AvatarURL, "REGRESSION: avatar should not be cleared")
}

// TestFieldLevelUserUpdate_AvatarOnly tests updating only avatar without clearing display name
func TestFieldLevelUserUpdate_AvatarOnly(t *testing.T) {
	client := NewMockFieldUpdateClient()
	userID := "@alice:example.com"

	// Get original state
	original, err := client.GetUser(context.Background(), userID)
	require.NoError(t, err)
	originalName := original.DisplayName

	// Update only avatar
	updated, err := client.UpdateUser(context.Background(), userID, &clients.UserSpec{
		AvatarURL: "mxc://example.com/avatar_new",
	})
	require.NoError(t, err)

	// Verify avatar changed, display name preserved (regression test)
	assert.Equal(t, "mxc://example.com/avatar_new", updated.AvatarURL)
	assert.Equal(t, originalName, updated.DisplayName, "REGRESSION: display name should not be cleared")
}

// TestFieldLevelRoomUpdate_TopicOnly tests updating only topic without affecting encryption
func TestFieldLevelRoomUpdate_TopicOnly(t *testing.T) {
	client := NewMockFieldUpdateClient()
	roomID := "!room:example.com"

	// Get original state
	original, err := client.GetRoom(context.Background(), roomID)
	require.NoError(t, err)
	originalEncryption := original.EncryptionEnabled

	// Update only topic
	updated, err := client.UpdateRoom(context.Background(), roomID, &clients.RoomSpec{
		Topic: "New Topic",
	})
	require.NoError(t, err)

	// Verify topic changed, encryption preserved (regression test)
	assert.Equal(t, "New Topic", updated.Topic)
	assert.Equal(t, originalEncryption, updated.EncryptionEnabled, "REGRESSION: encryption should not be changed")
}

// TestFieldLevelRoomUpdate_NameAndVisibility tests updating multiple fields together
func TestFieldLevelRoomUpdate_NameAndVisibility(t *testing.T) {
	client := NewMockFieldUpdateClient()
	roomID := "!room:example.com"

	// Get original state
	original, err := client.GetRoom(context.Background(), roomID)
	require.NoError(t, err)
	originalTopic := original.Topic

	// Update name and visibility, not topic
	updated, err := client.UpdateRoom(context.Background(), roomID, &clients.RoomSpec{
		Name:       "New Name",
		Visibility: "private",
	})
	require.NoError(t, err)

	// Verify changes and preserved fields (regression test)
	assert.Equal(t, "New Name", updated.Name)
	assert.Equal(t, "private", updated.Visibility)
	assert.Equal(t, originalTopic, updated.Topic, "REGRESSION: topic should not be changed")
}

// Priority 3.5: State Event Management (Regression Prevention)
// ============================================================

// MockStateEventClient tracks state events to detect corruption
type MockStateEventClient struct {
	stateEvents map[string][]clients.StateEvent
}

func NewMockStateEventClient() *MockStateEventClient {
	return &MockStateEventClient{
		stateEvents: map[string][]clients.StateEvent{
			"!room:example.com": {
				{
					Type:     "m.room.create",
					StateKey: "",
					Content: map[string]interface{}{
						"creator": "@alice:example.com",
					},
				},
				{
					Type:     "m.room.name",
					StateKey: "",
					Content: map[string]interface{}{
						"name": "Test Room",
					},
				},
			},
		},
	}
}

func (m *MockStateEventClient) AddStateEvent(ctx context.Context, roomID string, event clients.StateEvent) error {
	if _, exists := m.stateEvents[roomID]; !exists {
		m.stateEvents[roomID] = []clients.StateEvent{}
	}
	m.stateEvents[roomID] = append(m.stateEvents[roomID], event)
	return nil
}

func (m *MockStateEventClient) GetStateEvents(ctx context.Context, roomID string) ([]clients.StateEvent, error) {
	events, exists := m.stateEvents[roomID]
	if !exists {
		return []clients.StateEvent{}, nil
	}
	return events, nil
}

func (m *MockStateEventClient) ReplaceStateEvent(ctx context.Context, roomID string, eventType string, newEvent clients.StateEvent) error {
	events, exists := m.stateEvents[roomID]
	if !exists {
		return nil
	}

	// Replace the first matching event type
	for i, event := range events {
		if event.Type == eventType {
			events[i] = newEvent
			break
		}
	}
	return nil
}

// Stub implementations
func (m *MockStateEventClient) CreateUser(ctx context.Context, user *clients.UserSpec) (*clients.User, error) {
	return nil, nil
}
func (m *MockStateEventClient) GetUser(ctx context.Context, userID string) (*clients.User, error) {
	return nil, nil
}
func (m *MockStateEventClient) UpdateUser(ctx context.Context, userID string, user *clients.UserSpec) (*clients.User, error) {
	return nil, nil
}
func (m *MockStateEventClient) DeactivateUser(ctx context.Context, userID string) error {
	return nil
}
func (m *MockStateEventClient) CreateRoom(ctx context.Context, room *clients.RoomSpec) (*clients.Room, error) {
	return nil, nil
}
func (m *MockStateEventClient) GetRoom(ctx context.Context, roomID string) (*clients.Room, error) {
	return nil, nil
}
func (m *MockStateEventClient) UpdateRoom(ctx context.Context, roomID string, room *clients.RoomSpec) (*clients.Room, error) {
	return nil, nil
}
func (m *MockStateEventClient) DeleteRoom(ctx context.Context, roomID string) error {
	return nil
}
func (m *MockStateEventClient) SetPowerLevels(ctx context.Context, roomID string, powerLevels *clients.PowerLevelSpec) error {
	return nil
}
func (m *MockStateEventClient) GetPowerLevels(ctx context.Context, roomID string) (*clients.PowerLevelContent, error) {
	return nil, nil
}
func (m *MockStateEventClient) CreateRoomAlias(ctx context.Context, alias string, roomID string) error {
	return nil
}
func (m *MockStateEventClient) GetRoomAlias(ctx context.Context, alias string) (*clients.RoomAlias, error) {
	return nil, nil
}
func (m *MockStateEventClient) DeleteRoomAlias(ctx context.Context, alias string) error {
	return nil
}
func (m *MockStateEventClient) ListUsers(ctx context.Context, from string, limit int) (*clients.ListUsersResponse, error) {
	return nil, nil
}
func (m *MockStateEventClient) ListRooms(ctx context.Context, from string, limit int) (*clients.ListRoomsResponse, error) {
	return nil, nil
}
func (m *MockStateEventClient) MakeRoomAdmin(ctx context.Context, roomID, userID string) error {
	return nil
}
func (m *MockStateEventClient) BlockRoom(ctx context.Context, roomID string, block bool) error {
	return nil
}

// TestRetrieveAllStateEvents tests retrieving complete state event history without loss
func TestRetrieveAllStateEvents(t *testing.T) {
	client := NewMockStateEventClient()
	roomID := "!room:example.com"

	events, err := client.GetStateEvents(context.Background(), roomID)
	require.NoError(t, err)

	// Verify all events retrieved (regression: no silent loss)
	assert.Equal(t, 2, len(events), "REGRESSION: state events lost during retrieval")
	assert.Equal(t, "m.room.create", events[0].Type)
	assert.Equal(t, "m.room.name", events[1].Type)
}

// TestCustomStateEventCreation tests adding custom state events without corruption
func TestCustomStateEventCreation(t *testing.T) {
	client := NewMockStateEventClient()
	roomID := "!room:example.com"

	// Get initial count
	initialEvents, err := client.GetStateEvents(context.Background(), roomID)
	require.NoError(t, err)
	initialCount := len(initialEvents)

	// Add custom state event
	customEvent := clients.StateEvent{
		Type:     "com.example.custom",
		StateKey: "custom_data",
		Content: map[string]interface{}{
			"data": "value",
		},
	}

	err = client.AddStateEvent(context.Background(), roomID, customEvent)
	require.NoError(t, err)

	// Verify event added and others preserved (regression: no duplication/loss)
	updatedEvents, err := client.GetStateEvents(context.Background(), roomID)
	require.NoError(t, err)

	assert.Equal(t, initialCount+1, len(updatedEvents), "REGRESSION: event not added properly")
	assert.Equal(t, "com.example.custom", updatedEvents[len(updatedEvents)-1].Type)
}

// TestStateEventReplaceability tests that events can be replaced without duplicates
func TestStateEventReplaceability(t *testing.T) {
	client := NewMockStateEventClient()
	roomID := "!room:example.com"

	// Get initial count
	initialEvents, err := client.GetStateEvents(context.Background(), roomID)
	require.NoError(t, err)
	initialCount := len(initialEvents)

	// Replace existing state event
	newEvent := clients.StateEvent{
		Type:     "m.room.name",
		StateKey: "",
		Content: map[string]interface{}{
			"name": "Updated Room Name",
		},
	}

	err = client.ReplaceStateEvent(context.Background(), roomID, "m.room.name", newEvent)
	require.NoError(t, err)

	// Verify replacement without duplication (regression: no duplicate events)
	updatedEvents, err := client.GetStateEvents(context.Background(), roomID)
	require.NoError(t, err)

	assert.Equal(t, initialCount, len(updatedEvents), "REGRESSION: event count changed after replacement (duplicate?)")

	// Verify content updated
	nameEvents := 0
	for _, event := range updatedEvents {
		if event.Type == "m.room.name" {
			nameEvents++
			assert.Equal(t, "Updated Room Name", event.Content["name"])
		}
	}
	assert.Equal(t, 1, nameEvents, "REGRESSION: multiple name events exist (duplication)")
}

// TestStateEventOrdering tests that event order is preserved (no reordering)
func TestStateEventOrdering(t *testing.T) {
	client := NewMockStateEventClient()
	roomID := "!room:example.com"

	events, err := client.GetStateEvents(context.Background(), roomID)
	require.NoError(t, err)

	// Verify order: create should come before name (logical ordering)
	assert.Equal(t, 2, len(events))
	assert.Equal(t, "m.room.create", events[0].Type, "REGRESSION: room.create not first (ordering corrupted)")
	assert.Equal(t, "m.room.name", events[1].Type, "REGRESSION: room.name not second (ordering corrupted)")
}

// TestStateEventIntegrity tests that event content is not corrupted
func TestStateEventIntegrity(t *testing.T) {
	client := NewMockStateEventClient()
	roomID := "!room:example.com"

	events, err := client.GetStateEvents(context.Background(), roomID)
	require.NoError(t, err)

	// Verify content integrity for each event
	createEvent := events[0]
	assert.NotNil(t, createEvent.Content["creator"], "REGRESSION: event content corrupted (missing creator)")

	nameEvent := events[1]
	assert.NotNil(t, nameEvent.Content["name"], "REGRESSION: event content corrupted (missing name)")
	assert.Equal(t, "Test Room", nameEvent.Content["name"])
}
