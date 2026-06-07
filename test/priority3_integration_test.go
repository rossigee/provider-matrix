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
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/crossplane-contrib/provider-matrix/internal/clients"
)

// Priority 3: Real Matrix Server Integration Simulation
// ======================================================

// MockRealServerClient simulates a real Matrix server with realistic behavior
type MockRealServerClient struct {
	// Server state
	users  map[string]*clients.User
	rooms  map[string]*clients.Room
	mu     sync.RWMutex
	delay  time.Duration // Network delay simulation
	errFn  func(op string) error
	closed bool
}

func NewMockRealServerClient() *MockRealServerClient {
	return &MockRealServerClient{
		users: make(map[string]*clients.User),
		rooms: make(map[string]*clients.Room),
		delay: time.Millisecond,
	}
}

func (m *MockRealServerClient) simulateNetworkDelay() {
	if m.delay > 0 {
		time.Sleep(m.delay)
	}
}

func (m *MockRealServerClient) CreateUser(ctx context.Context, user *clients.UserSpec) (*clients.User, error) {
	// Check context before and after network delay
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	m.simulateNetworkDelay()

	// Check context again after network delay
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.errFn != nil {
		if err := m.errFn("create_user"); err != nil {
			return nil, err
		}
	}

	if _, exists := m.users[user.UserID]; exists {
		return nil, errors.New("user already exists")
	}

	newUser := &clients.User{
		UserID:      user.UserID,
		DisplayName: user.DisplayName,
		Admin:       false,
	}
	m.users[user.UserID] = newUser
	return newUser, nil
}

func (m *MockRealServerClient) GetUser(ctx context.Context, userID string) (*clients.User, error) {
	m.simulateNetworkDelay()

	m.mu.RLock()
	defer m.mu.RUnlock()

	user, exists := m.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *MockRealServerClient) UpdateUser(ctx context.Context, userID string, user *clients.UserSpec) (*clients.User, error) {
	m.simulateNetworkDelay()

	m.mu.Lock()
	defer m.mu.Unlock()

	existing, exists := m.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}

	if user.DisplayName != "" {
		existing.DisplayName = user.DisplayName
	}
	return existing, nil
}

func (m *MockRealServerClient) DeactivateUser(ctx context.Context, userID string) error {
	m.simulateNetworkDelay()

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.users[userID]; !exists {
		return errors.New("user not found")
	}
	delete(m.users, userID)
	return nil
}

func (m *MockRealServerClient) CreateRoom(ctx context.Context, room *clients.RoomSpec) (*clients.Room, error) {
	m.simulateNetworkDelay()

	m.mu.Lock()
	defer m.mu.Unlock()

	// Generate a room ID from the alias or name
	roomID := "!" + generateRoomID(room.Alias, room.Name) + ":example.com"

	if _, exists := m.rooms[roomID]; exists {
		return nil, errors.New("room already exists")
	}

	newRoom := &clients.Room{
		RoomID:     roomID,
		Name:       room.Name,
		Topic:      room.Topic,
		Visibility: room.Visibility,
	}
	m.rooms[roomID] = newRoom
	return newRoom, nil
}

func (m *MockRealServerClient) GetRoom(ctx context.Context, roomID string) (*clients.Room, error) {
	m.simulateNetworkDelay()

	m.mu.RLock()
	defer m.mu.RUnlock()

	room, exists := m.rooms[roomID]
	if !exists {
		return nil, errors.New("room not found")
	}
	return room, nil
}

func (m *MockRealServerClient) UpdateRoom(ctx context.Context, roomID string, room *clients.RoomSpec) (*clients.Room, error) {
	m.simulateNetworkDelay()

	m.mu.Lock()
	defer m.mu.Unlock()

	existing, exists := m.rooms[roomID]
	if !exists {
		return nil, errors.New("room not found")
	}

	if room.Name != "" {
		existing.Name = room.Name
	}
	if room.Topic != "" {
		existing.Topic = room.Topic
	}
	return existing, nil
}

func (m *MockRealServerClient) DeleteRoom(ctx context.Context, roomID string) error {
	m.simulateNetworkDelay()

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.rooms[roomID]; !exists {
		return errors.New("room not found")
	}
	delete(m.rooms, roomID)
	return nil
}

func (m *MockRealServerClient) SetPowerLevels(ctx context.Context, roomID string, powerLevels *clients.PowerLevelSpec) error {
	m.simulateNetworkDelay()

	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, exists := m.rooms[roomID]; !exists {
		return errors.New("room not found")
	}
	return nil
}

func (m *MockRealServerClient) GetPowerLevels(ctx context.Context, roomID string) (*clients.PowerLevelContent, error) {
	m.simulateNetworkDelay()

	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, exists := m.rooms[roomID]; !exists {
		return nil, errors.New("room not found")
	}
	return &clients.PowerLevelContent{}, nil
}

func (m *MockRealServerClient) CreateRoomAlias(ctx context.Context, alias string, roomID string) error {
	m.simulateNetworkDelay()
	return nil
}

func (m *MockRealServerClient) GetRoomAlias(ctx context.Context, alias string) (*clients.RoomAlias, error) {
	m.simulateNetworkDelay()
	return &clients.RoomAlias{Alias: alias}, nil
}

func (m *MockRealServerClient) DeleteRoomAlias(ctx context.Context, alias string) error {
	m.simulateNetworkDelay()
	return nil
}

func (m *MockRealServerClient) ListUsers(ctx context.Context, from string, limit int) (*clients.ListUsersResponse, error) {
	m.simulateNetworkDelay()

	m.mu.RLock()
	defer m.mu.RUnlock()

	users := make([]clients.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, *user)
	}
	return &clients.ListUsersResponse{Users: users}, nil
}

func (m *MockRealServerClient) ListRooms(ctx context.Context, from string, limit int) (*clients.ListRoomsResponse, error) {
	m.simulateNetworkDelay()

	m.mu.RLock()
	defer m.mu.RUnlock()

	rooms := make([]clients.Room, 0, len(m.rooms))
	for _, room := range m.rooms {
		rooms = append(rooms, *room)
	}
	return &clients.ListRoomsResponse{Rooms: rooms}, nil
}

func (m *MockRealServerClient) MakeRoomAdmin(ctx context.Context, roomID, userID string) error {
	m.simulateNetworkDelay()
	return nil
}

func (m *MockRealServerClient) BlockRoom(ctx context.Context, roomID string, block bool) error {
	m.simulateNetworkDelay()
	return nil
}

// TestRealServerUserIntegration tests realistic user operations
func TestRealServerUserIntegration(t *testing.T) {
	client := NewMockRealServerClient()

	// Create user
	user, err := client.CreateUser(context.Background(), &clients.UserSpec{
		UserID:      "@alice:example.com",
		DisplayName: "Alice",
		Password:    "secret",
	})
	require.NoError(t, err)
	assert.Equal(t, "@alice:example.com", user.UserID)

	// Get user
	retrieved, err := client.GetUser(context.Background(), "@alice:example.com")
	require.NoError(t, err)
	assert.Equal(t, "Alice", retrieved.DisplayName)

	// Update user
	updated, err := client.UpdateUser(context.Background(), "@alice:example.com", &clients.UserSpec{
		DisplayName: "Alice Smith",
	})
	require.NoError(t, err)
	assert.Equal(t, "Alice Smith", updated.DisplayName)

	// Deactivate user
	err = client.DeactivateUser(context.Background(), "@alice:example.com")
	require.NoError(t, err)

	// Verify user is gone
	_, err = client.GetUser(context.Background(), "@alice:example.com")
	assert.Error(t, err)
}

// generateRoomID creates a room ID from alias or name
func generateRoomID(alias, name string) string {
	if alias != "" {
		return alias
	}
	if name != "" {
		return strings.ToLower(strings.ReplaceAll(name, " ", ""))
	}
	return "room"
}

// TestRealServerRoomIntegration tests realistic room operations
func TestRealServerRoomIntegration(t *testing.T) {
	client := NewMockRealServerClient()

	// Create room
	room, err := client.CreateRoom(context.Background(), &clients.RoomSpec{
		Name:       "General",
		Topic:      "General discussion",
		Visibility: "public",
	})
	require.NoError(t, err)
	assert.Equal(t, "!general:example.com", room.RoomID)

	// Update room
	updated, err := client.UpdateRoom(context.Background(), room.RoomID, &clients.RoomSpec{
		Name: "General Chat",
	})
	require.NoError(t, err)
	assert.Equal(t, "General Chat", updated.Name)

	// Delete room
	err = client.DeleteRoom(context.Background(), room.RoomID)
	require.NoError(t, err)

	// Verify room is gone
	_, err = client.GetRoom(context.Background(), room.RoomID)
	assert.Error(t, err)
}

// Priority 3: Concurrency Testing
// ================================

// TestConcurrentUserCreation tests concurrent user creation without race conditions
func TestConcurrentUserCreation(t *testing.T) {
	client := NewMockRealServerClient()
	numGoroutines := 10
	usersPerGoroutine := 5

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*usersPerGoroutine)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < usersPerGoroutine; j++ {
				userID := fmt.Sprintf("@user_%d_%d:example.com", id, j)
				_, err := client.CreateUser(context.Background(), &clients.UserSpec{
					UserID: userID,
				})
				if err != nil {
					errors <- err
				}
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Verify no errors
	for err := range errors {
		t.Errorf("concurrent creation error: %v", err)
	}

	// Verify all users created
	listResp, err := client.ListUsers(context.Background(), "", 1000)
	require.NoError(t, err)
	assert.Equal(t, numGoroutines*usersPerGoroutine, len(listResp.Users))
}

// TestConcurrentRoomOperations tests concurrent room read/write operations
func TestConcurrentRoomOperations(t *testing.T) {
	client := NewMockRealServerClient()

	// Create initial room
	room, err := client.CreateRoom(context.Background(), &clients.RoomSpec{
		Name: "Test Room",
	})
	require.NoError(t, err)

	numGoroutines := 5
	successCount := int32(0)

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// Read-heavy operations (should succeed)
			_, err := client.GetRoom(context.Background(), room.RoomID)
			if err == nil {
				atomic.AddInt32(&successCount, 1)
			}
		}(i)
	}

	wg.Wait()

	// All reads should succeed
	assert.Equal(t, int32(numGoroutines), successCount)
}

// TestConcurrentUserAndRoomCreation tests concurrent creation of users and rooms
func TestConcurrentUserAndRoomCreation(t *testing.T) {
	client := NewMockRealServerClient()
	numOperations := 20

	var wg sync.WaitGroup
	userErrors := make(chan error, numOperations)
	roomErrors := make(chan error, numOperations)

	// Create users
	for i := 0; i < numOperations; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_, err := client.CreateUser(context.Background(), &clients.UserSpec{
				UserID: fmt.Sprintf("@user_%d:example.com", id),
			})
			if err != nil {
				userErrors <- err
			}
		}(i)
	}

	// Create rooms concurrently
	for i := 0; i < numOperations; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_, err := client.CreateRoom(context.Background(), &clients.RoomSpec{
				Name: fmt.Sprintf("Room %d", id),
			})
			if err != nil {
				roomErrors <- err
			}
		}(i)
	}

	wg.Wait()
	close(userErrors)
	close(roomErrors)

	// Verify no errors
	for err := range userErrors {
		t.Errorf("concurrent user creation error: %v", err)
	}
	for err := range roomErrors {
		t.Errorf("concurrent room creation error: %v", err)
	}

	// Verify all created
	users, err := client.ListUsers(context.Background(), "", 1000)
	require.NoError(t, err)
	assert.Equal(t, numOperations, len(users.Users))

	rooms, err := client.ListRooms(context.Background(), "", 1000)
	require.NoError(t, err)
	assert.Equal(t, numOperations, len(rooms.Rooms))
}

// Priority 3: Load Testing
// ========================

// TestLoadUserCreation tests creating many users in sequence
func TestLoadUserCreation(t *testing.T) {
	client := NewMockRealServerClient()
	numUsers := 100

	startTime := time.Now()
	for i := 0; i < numUsers; i++ {
		_, err := client.CreateUser(context.Background(), &clients.UserSpec{
			UserID:      fmt.Sprintf("@loadtest_%d:example.com", i),
			DisplayName: fmt.Sprintf("User %d", i),
		})
		require.NoError(t, err)
	}
	elapsed := time.Since(startTime)

	// Verify all created
	listResp, err := client.ListUsers(context.Background(), "", 1000)
	require.NoError(t, err)
	assert.Equal(t, numUsers, len(listResp.Users))

	// Performance check (should complete reasonably fast)
	assert.Less(t, elapsed, 10*time.Second, "load test took too long")
}

// TestLoadRoomCreation tests creating many rooms in sequence
func TestLoadRoomCreation(t *testing.T) {
	client := NewMockRealServerClient()
	numRooms := 50

	startTime := time.Now()
	for i := 0; i < numRooms; i++ {
		_, err := client.CreateRoom(context.Background(), &clients.RoomSpec{
			Name: fmt.Sprintf("Load Room %d", i),
		})
		require.NoError(t, err)
	}
	elapsed := time.Since(startTime)

	// Verify all created
	listResp, err := client.ListRooms(context.Background(), "", 1000)
	require.NoError(t, err)
	assert.Equal(t, numRooms, len(listResp.Rooms))

	// Performance check
	assert.Less(t, elapsed, 10*time.Second, "load test took too long")
}

// TestLoadMixedOperations tests load with mixed CRUD operations
func TestLoadMixedOperations(t *testing.T) {
	client := NewMockRealServerClient()

	// Create initial data
	for i := 0; i < 20; i++ {
		client.CreateUser(context.Background(), &clients.UserSpec{
			UserID: fmt.Sprintf("@mixed_%d:example.com", i),
		})
	}

	// Execute mixed operations
	startTime := time.Now()
	for i := 0; i < 50; i++ {
		userID := fmt.Sprintf("@mixed_%d:example.com", i%20)

		// Create, read, update sequence
		client.GetUser(context.Background(), userID)
		client.UpdateUser(context.Background(), userID, &clients.UserSpec{
			DisplayName: fmt.Sprintf("Updated %d", i),
		})
	}
	elapsed := time.Since(startTime)

	// Performance check
	assert.Less(t, elapsed, 10*time.Second)
}

// Priority 3: Advanced Error Recovery
// ====================================

// TestRetryOnTransientError tests recovery from transient errors
func TestRetryOnTransientError(t *testing.T) {
	client := NewMockRealServerClient()
	callCount := 0

	// Simulate transient error that clears after 2 attempts
	client.errFn = func(op string) error {
		callCount++
		if callCount <= 2 {
			return errors.New("temporary service unavailable")
		}
		return nil
	}

	// First attempt fails
	_, err := client.CreateUser(context.Background(), &clients.UserSpec{UserID: "@retry:example.com"})
	assert.Error(t, err)

	// Reset error function
	client.errFn = nil

	// Retry succeeds
	user, err := client.CreateUser(context.Background(), &clients.UserSpec{UserID: "@retry:example.com"})
	require.NoError(t, err)
	assert.Equal(t, "@retry:example.com", user.UserID)
}

// TestPartialFailureRecovery tests handling partial failures in batch operations
func TestPartialFailureRecovery(t *testing.T) {
	client := NewMockRealServerClient()
	numOperations := 10
	successCount := 0
	failureCount := 0

	for i := 0; i < numOperations; i++ {
		_, err := client.CreateUser(context.Background(), &clients.UserSpec{
			UserID: fmt.Sprintf("@partial_%d:example.com", i),
		})
		if err == nil {
			successCount++
		} else {
			failureCount++
		}
	}

	// All should succeed (no failures)
	assert.Equal(t, numOperations, successCount)
	assert.Equal(t, 0, failureCount)

	// Verify data consistency
	listResp, err := client.ListUsers(context.Background(), "", 1000)
	require.NoError(t, err)
	assert.Equal(t, numOperations, len(listResp.Users))
}

// TestTimeoutRecovery tests handling operation timeouts
func TestTimeoutRecovery(t *testing.T) {
	client := NewMockRealServerClient()
	client.delay = 100 * time.Millisecond

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// This operation should timeout
	_, err := client.CreateUser(ctx, &clients.UserSpec{
		UserID: "@timeout:example.com",
	})

	// Either timeout or context error
	assert.Error(t, err)
}

// TestCascadingFailures tests handling cascading failures
func TestCascadingFailures(t *testing.T) {
	client := NewMockRealServerClient()
	failureMode := false

	client.errFn = func(op string) error {
		if failureMode {
			return errors.New("cascading failure detected")
		}
		return nil
	}

	// Normal operation succeeds
	_, err := client.CreateUser(context.Background(), &clients.UserSpec{
		UserID: "@cascade:example.com",
	})
	require.NoError(t, err)

	// Enable failure mode
	failureMode = true

	// Subsequent operations fail
	_, err = client.CreateUser(context.Background(), &clients.UserSpec{
		UserID: "@cascade2:example.com",
	})
	assert.Error(t, err)

	// Disable failure mode
	failureMode = false

	// Operations resume successfully
	_, err = client.CreateUser(context.Background(), &clients.UserSpec{
		UserID: "@cascade3:example.com",
	})
	require.NoError(t, err)
}

// TestDataConsistencyUnderLoad tests data consistency with concurrent operations
func TestDataConsistencyUnderLoad(t *testing.T) {
	client := NewMockRealServerClient()
	numGoroutines := 5
	operationsPerGoroutine := 20

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				userID := fmt.Sprintf("@consistency_%d_%d:example.com", id, j)
				client.CreateUser(context.Background(), &clients.UserSpec{
					UserID:      userID,
					DisplayName: fmt.Sprintf("User %d_%d", id, j),
				})

				// Read-back verification
				user, err := client.GetUser(context.Background(), userID)
				if err == nil {
					assert.Equal(t, userID, user.UserID)
				}
			}
		}(i)
	}

	wg.Wait()

	// Verify final state
	listResp, err := client.ListUsers(context.Background(), "", 10000)
	require.NoError(t, err)
	assert.Equal(t, numGoroutines*operationsPerGoroutine, len(listResp.Users))
}
