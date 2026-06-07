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

	"github.com/crossplane-contrib/provider-matrix/apis/roomalias/v1alpha1"
	"github.com/crossplane-contrib/provider-matrix/internal/clients"
)

// MockRoomAliasClient implements clients.Client for RoomAlias testing
type MockRoomAliasClient struct {
	createRoomAliasFn func(ctx context.Context, alias string, roomID string) error
	getRoomAliasFn    func(ctx context.Context, alias string) (*clients.RoomAlias, error)
	deleteRoomAliasFn func(ctx context.Context, alias string) error
}

func (m *MockRoomAliasClient) CreateRoomAlias(ctx context.Context, alias string, roomID string) error {
	if m.createRoomAliasFn != nil {
		return m.createRoomAliasFn(ctx, alias, roomID)
	}
	return nil
}

func (m *MockRoomAliasClient) GetRoomAlias(ctx context.Context, alias string) (*clients.RoomAlias, error) {
	if m.getRoomAliasFn != nil {
		return m.getRoomAliasFn(ctx, alias)
	}
	return &clients.RoomAlias{
		Alias:  alias,
		RoomID: "!test:example.com",
	}, nil
}

func (m *MockRoomAliasClient) DeleteRoomAlias(ctx context.Context, alias string) error {
	if m.deleteRoomAliasFn != nil {
		return m.deleteRoomAliasFn(ctx, alias)
	}
	return nil
}

// Stub implementations for unused interface methods
func (m *MockRoomAliasClient) CreateUser(ctx context.Context, user *clients.UserSpec) (*clients.User, error) {
	return nil, nil
}

func (m *MockRoomAliasClient) GetUser(ctx context.Context, userID string) (*clients.User, error) {
	return nil, nil
}

func (m *MockRoomAliasClient) UpdateUser(ctx context.Context, userID string, user *clients.UserSpec) (*clients.User, error) {
	return nil, nil
}

func (m *MockRoomAliasClient) DeactivateUser(ctx context.Context, userID string) error {
	return nil
}

func (m *MockRoomAliasClient) CreateRoom(ctx context.Context, room *clients.RoomSpec) (*clients.Room, error) {
	return nil, nil
}

func (m *MockRoomAliasClient) GetRoom(ctx context.Context, roomID string) (*clients.Room, error) {
	return nil, nil
}

func (m *MockRoomAliasClient) UpdateRoom(ctx context.Context, roomID string, room *clients.RoomSpec) (*clients.Room, error) {
	return nil, nil
}

func (m *MockRoomAliasClient) DeleteRoom(ctx context.Context, roomID string) error {
	return nil
}

func (m *MockRoomAliasClient) SetPowerLevels(ctx context.Context, roomID string, powerLevels *clients.PowerLevelSpec) error {
	return nil
}

func (m *MockRoomAliasClient) GetPowerLevels(ctx context.Context, roomID string) (*clients.PowerLevelContent, error) {
	return nil, nil
}

func (m *MockRoomAliasClient) ListUsers(ctx context.Context, from string, limit int) (*clients.ListUsersResponse, error) {
	return nil, nil
}

func (m *MockRoomAliasClient) ListRooms(ctx context.Context, from string, limit int) (*clients.ListRoomsResponse, error) {
	return nil, nil
}

func (m *MockRoomAliasClient) MakeRoomAdmin(ctx context.Context, roomID, userID string) error {
	return nil
}

func (m *MockRoomAliasClient) BlockRoom(ctx context.Context, roomID string, block bool) error {
	return nil
}

// TestCreateRoomAlias tests room alias creation
func TestCreateRoomAlias(t *testing.T) {
	mockClient := &MockRoomAliasClient{
		createRoomAliasFn: func(ctx context.Context, alias string, roomID string) error {
			if alias == "" || roomID == "" {
				return errors.New("alias and roomID are required")
			}
			return nil
		},
	}

	err := mockClient.CreateRoomAlias(context.Background(), "#general:example.com", "!test:example.com")
	require.NoError(t, err)
}

// TestCreateRoomAliasError tests room alias creation error handling
func TestCreateRoomAliasError(t *testing.T) {
	mockClient := &MockRoomAliasClient{
		createRoomAliasFn: func(ctx context.Context, alias string, roomID string) error {
			return errors.New("alias already in use")
		},
	}

	err := mockClient.CreateRoomAlias(context.Background(), "#taken:example.com", "!test:example.com")
	assert.Error(t, err)
	assert.Equal(t, "alias already in use", err.Error())
}

// TestGetRoomAlias tests retrieving room alias information
func TestGetRoomAlias(t *testing.T) {
	mockClient := &MockRoomAliasClient{
		getRoomAliasFn: func(ctx context.Context, alias string) (*clients.RoomAlias, error) {
			return &clients.RoomAlias{
				Alias:  alias,
				RoomID: "!target:example.com",
			}, nil
		},
	}

	alias, err := mockClient.GetRoomAlias(context.Background(), "#general:example.com")
	require.NoError(t, err)
	assert.Equal(t, "#general:example.com", alias.Alias)
	assert.Equal(t, "!target:example.com", alias.RoomID)
}

// TestGetRoomAliasNotFound tests getting non-existent alias
func TestGetRoomAliasNotFound(t *testing.T) {
	mockClient := &MockRoomAliasClient{
		getRoomAliasFn: func(ctx context.Context, alias string) (*clients.RoomAlias, error) {
			return nil, errors.New("alias not found")
		},
	}

	_, err := mockClient.GetRoomAlias(context.Background(), "#nonexistent:example.com")
	assert.Error(t, err)
}

// TestDeleteRoomAlias tests room alias deletion
func TestDeleteRoomAlias(t *testing.T) {
	mockClient := &MockRoomAliasClient{
		deleteRoomAliasFn: func(ctx context.Context, alias string) error {
			return nil
		},
	}

	err := mockClient.DeleteRoomAlias(context.Background(), "#general:example.com")
	assert.NoError(t, err)
}

// TestDeleteRoomAliasError tests deletion error handling
func TestDeleteRoomAliasError(t *testing.T) {
	mockClient := &MockRoomAliasClient{
		deleteRoomAliasFn: func(ctx context.Context, alias string) error {
			return errors.New("cannot delete alias")
		},
	}

	err := mockClient.DeleteRoomAlias(context.Background(), "#protected:example.com")
	assert.Error(t, err)
}

// TestRoomAliasResource tests RoomAlias CR creation and manipulation
func TestRoomAliasResource(t *testing.T) {
	alias := &v1alpha1.RoomAlias{
		ObjectMeta: metav1.ObjectMeta{
			Name: "general-alias",
		},
		Spec: v1alpha1.RoomAliasSpec{
			ForProvider: v1alpha1.RoomAliasParameters{
				Alias:  "#general:example.com",
				RoomID: "!abc123:example.com",
			},
		},
	}

	assert.Equal(t, "general-alias", alias.Name)
	assert.Equal(t, "#general:example.com", alias.Spec.ForProvider.Alias)
	assert.Equal(t, "!abc123:example.com", alias.Spec.ForProvider.RoomID)
}

// TestRoomAliasProviderConfigReference tests provider config reference handling
func TestRoomAliasProviderConfigReference(t *testing.T) {
	alias := &v1alpha1.RoomAlias{}

	pcRef := &xpcoreapi.ProviderConfigReference{Name: "test-pc"}
	alias.SetProviderConfigReference(pcRef)
	assert.Equal(t, "test-pc", alias.GetProviderConfigReference().Name)
}

// TestRoomAliasConditions tests condition management
func TestRoomAliasConditions(t *testing.T) {
	alias := &v1alpha1.RoomAlias{}

	cond := xpcoreapi.Available()
	alias.SetConditions(cond)

	retrieved := alias.GetCondition(xpcoreapi.TypeReady)
	assert.Equal(t, xpcoreapi.TypeReady, retrieved.Type)
}

// TestRoomAliasManagementPolicies tests management policy handling
func TestRoomAliasManagementPolicies(t *testing.T) {
	alias := &v1alpha1.RoomAlias{}

	policies := xpcoreapi.ManagementPolicies{"*"}
	alias.SetManagementPolicies(policies)

	retrieved := alias.GetManagementPolicies()
	assert.Equal(t, policies, retrieved)
}

// TestRoomAliasWriteConnectionSecret tests secret reference handling
func TestRoomAliasWriteConnectionSecret(t *testing.T) {
	alias := &v1alpha1.RoomAlias{}

	secretRef := &xpcoreapi.LocalSecretReference{
		Name: "alias-config",
	}
	alias.SetWriteConnectionSecretToReference(secretRef)

	retrieved := alias.GetWriteConnectionSecretToReference()
	assert.Equal(t, "alias-config", retrieved.Name)
}
