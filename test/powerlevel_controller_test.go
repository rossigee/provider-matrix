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

	"github.com/crossplane-contrib/provider-matrix/apis/powerlevel/v1alpha1"
	"github.com/crossplane-contrib/provider-matrix/internal/clients"
)

// MockPowerLevelClient implements clients.Client for PowerLevel testing
type MockPowerLevelClient struct {
	setPowerLevelsFn func(ctx context.Context, roomID string, powerLevels *clients.PowerLevelSpec) error
	getPowerLevelsFn func(ctx context.Context, roomID string) (*clients.PowerLevelContent, error)
}

func (m *MockPowerLevelClient) SetPowerLevels(ctx context.Context, roomID string, powerLevels *clients.PowerLevelSpec) error {
	if m.setPowerLevelsFn != nil {
		return m.setPowerLevelsFn(ctx, roomID, powerLevels)
	}
	return nil
}

func (m *MockPowerLevelClient) GetPowerLevels(ctx context.Context, roomID string) (*clients.PowerLevelContent, error) {
	if m.getPowerLevelsFn != nil {
		return m.getPowerLevelsFn(ctx, roomID)
	}
	return &clients.PowerLevelContent{
		UsersDefault:  intPtr(0),
		EventsDefault: intPtr(50),
		StateDefault:  intPtr(100),
	}, nil
}

// Stub implementations for unused interface methods
func (m *MockPowerLevelClient) CreateUser(ctx context.Context, user *clients.UserSpec) (*clients.User, error) {
	return nil, nil
}

func (m *MockPowerLevelClient) GetUser(ctx context.Context, userID string) (*clients.User, error) {
	return nil, nil
}

func (m *MockPowerLevelClient) UpdateUser(ctx context.Context, userID string, user *clients.UserSpec) (*clients.User, error) {
	return nil, nil
}

func (m *MockPowerLevelClient) DeactivateUser(ctx context.Context, userID string) error {
	return nil
}

func (m *MockPowerLevelClient) CreateRoom(ctx context.Context, room *clients.RoomSpec) (*clients.Room, error) {
	return nil, nil
}

func (m *MockPowerLevelClient) GetRoom(ctx context.Context, roomID string) (*clients.Room, error) {
	return nil, nil
}

func (m *MockPowerLevelClient) UpdateRoom(ctx context.Context, roomID string, room *clients.RoomSpec) (*clients.Room, error) {
	return nil, nil
}

func (m *MockPowerLevelClient) DeleteRoom(ctx context.Context, roomID string) error {
	return nil
}

func (m *MockPowerLevelClient) CreateRoomAlias(ctx context.Context, alias string, roomID string) error {
	return nil
}

func (m *MockPowerLevelClient) GetRoomAlias(ctx context.Context, alias string) (*clients.RoomAlias, error) {
	return nil, nil
}

func (m *MockPowerLevelClient) DeleteRoomAlias(ctx context.Context, alias string) error {
	return nil
}

func (m *MockPowerLevelClient) ListUsers(ctx context.Context, from string, limit int) (*clients.ListUsersResponse, error) {
	return nil, nil
}

func (m *MockPowerLevelClient) ListRooms(ctx context.Context, from string, limit int) (*clients.ListRoomsResponse, error) {
	return nil, nil
}

func (m *MockPowerLevelClient) MakeRoomAdmin(ctx context.Context, roomID, userID string) error {
	return nil
}

func (m *MockPowerLevelClient) BlockRoom(ctx context.Context, roomID string, block bool) error {
	return nil
}

// TestSetPowerLevels tests power level configuration
func TestSetPowerLevels(t *testing.T) {
	mockClient := &MockPowerLevelClient{
		setPowerLevelsFn: func(ctx context.Context, roomID string, powerLevels *clients.PowerLevelSpec) error {
			if powerLevels == nil {
				return errors.New("power levels required")
			}
			return nil
		},
	}

	spec := &clients.PowerLevelSpec{
		RoomID: "!test:example.com",
		PowerLevels: &clients.PowerLevelContent{
			UsersDefault:  intPtr(0),
			EventsDefault: intPtr(50),
			StateDefault:  intPtr(100),
		},
	}

	err := mockClient.SetPowerLevels(context.Background(), spec.RoomID, spec)
	require.NoError(t, err)
}

// TestSetPowerLevelsError tests power level configuration error handling
func TestSetPowerLevelsError(t *testing.T) {
	mockClient := &MockPowerLevelClient{
		setPowerLevelsFn: func(ctx context.Context, roomID string, powerLevels *clients.PowerLevelSpec) error {
			return errors.New("cannot set power levels")
		},
	}

	spec := &clients.PowerLevelSpec{
		RoomID: "!test:example.com",
	}

	err := mockClient.SetPowerLevels(context.Background(), spec.RoomID, spec)
	assert.Error(t, err)
}

// TestGetPowerLevels tests retrieving power level configuration
func TestGetPowerLevels(t *testing.T) {
	mockClient := &MockPowerLevelClient{
		getPowerLevelsFn: func(ctx context.Context, roomID string) (*clients.PowerLevelContent, error) {
			return &clients.PowerLevelContent{
				UsersDefault:  intPtr(0),
				EventsDefault: intPtr(50),
				StateDefault:  intPtr(100),
			}, nil
		},
	}

	content, err := mockClient.GetPowerLevels(context.Background(), "!test:example.com")
	require.NoError(t, err)
	assert.NotNil(t, content)
	assert.Equal(t, 0, *content.UsersDefault)
	assert.Equal(t, 50, *content.EventsDefault)
	assert.Equal(t, 100, *content.StateDefault)
}

// TestPowerLevelResource tests PowerLevel CR creation and manipulation
func TestPowerLevelResource(t *testing.T) {
	pl := &v1alpha1.PowerLevel{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-powerlevel",
		},
		Spec: v1alpha1.PowerLevelSpec{
			ForProvider: v1alpha1.PowerLevelParameters{
				RoomID: "!test:example.com",
				Users: map[string]int{
					"@user:example.com": 50,
				},
			},
		},
	}

	assert.Equal(t, "test-powerlevel", pl.Name)
	assert.Equal(t, "!test:example.com", pl.Spec.ForProvider.RoomID)
	assert.NotNil(t, pl.Spec.ForProvider.Users)
}

// TestPowerLevelProviderConfigReference tests provider config reference handling
func TestPowerLevelProviderConfigReference(t *testing.T) {
	pl := &v1alpha1.PowerLevel{}

	pcRef := &xpcoreapi.ProviderConfigReference{Name: "test-pc"}
	pl.SetProviderConfigReference(pcRef)
	assert.Equal(t, "test-pc", pl.GetProviderConfigReference().Name)
}

// TestPowerLevelConditions tests condition management
func TestPowerLevelConditions(t *testing.T) {
	pl := &v1alpha1.PowerLevel{}

	cond := xpcoreapi.Available()
	pl.SetConditions(cond)

	retrieved := pl.GetCondition(xpcoreapi.TypeReady)
	assert.Equal(t, xpcoreapi.TypeReady, retrieved.Type)
}

// TestPowerLevelManagementPolicies tests management policy handling
func TestPowerLevelManagementPolicies(t *testing.T) {
	pl := &v1alpha1.PowerLevel{}

	policies := xpcoreapi.ManagementPolicies{"*"}
	pl.SetManagementPolicies(policies)

	retrieved := pl.GetManagementPolicies()
	assert.Equal(t, policies, retrieved)
}

// Helper function
func intPtr(i int) *int {
	return &i
}
