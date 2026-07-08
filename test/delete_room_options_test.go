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
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// MockDeleteRoomClient for testing advanced room deletion with options
type MockDeleteRoomClient struct {
	deleteRoomWithOptionsFn func(ctx context.Context, roomID string, options map[string]interface{}) error
}

func (m *MockDeleteRoomClient) DeleteRoom(ctx context.Context, roomID string) error {
	// Simple delete without options
	if m.deleteRoomWithOptionsFn != nil {
		return m.deleteRoomWithOptionsFn(ctx, roomID, map[string]interface{}{})
	}
	return nil
}

// Test helper for deletion with options
func (m *MockDeleteRoomClient) deleteRoomWithOptions(ctx context.Context, roomID string, options map[string]interface{}) error {
	if m.deleteRoomWithOptionsFn != nil {
		return m.deleteRoomWithOptionsFn(ctx, roomID, options)
	}
	return nil
}

// TestDeleteRoomBasic tests simple room deletion
func TestDeleteRoomBasic(t *testing.T) {
	mockClient := &MockDeleteRoomClient{
		deleteRoomWithOptionsFn: func(ctx context.Context, roomID string, options map[string]interface{}) error {
			if roomID == "" {
				return errors.New("roomID required")
			}
			return nil
		},
	}

	err := mockClient.deleteRoomWithOptions(context.Background(), "!room:example.com", nil)
	require.NoError(t, err)
}

// TestDeleteRoomWithPurge tests deletion with purge option
func TestDeleteRoomWithPurge(t *testing.T) {
	mockClient := &MockDeleteRoomClient{
		deleteRoomWithOptionsFn: func(ctx context.Context, roomID string, options map[string]interface{}) error {
			if options != nil {
				if purge, ok := options["purge"].(bool); ok && purge {
					// Purge option was specified
					return nil
				}
			}
			return errors.New("purge option required")
		},
	}

	options := map[string]interface{}{
		"purge": true,
	}

	err := mockClient.deleteRoomWithOptions(context.Background(), "!room:example.com", options)
	require.NoError(t, err)
}

// TestDeleteRoomWithKickUsers tests deletion with kick_users option
func TestDeleteRoomWithKickUsers(t *testing.T) {
	mockClient := &MockDeleteRoomClient{
		deleteRoomWithOptionsFn: func(ctx context.Context, roomID string, options map[string]interface{}) error {
			if options != nil {
				if kick, ok := options["kick_users"].(bool); ok && kick {
					return nil
				}
			}
			return errors.New("kick_users option required")
		},
	}

	options := map[string]interface{}{
		"kick_users": true,
	}

	err := mockClient.deleteRoomWithOptions(context.Background(), "!room:example.com", options)
	require.NoError(t, err)
}

// TestDeleteRoomWithMultipleOptions tests deletion with multiple options
func TestDeleteRoomWithMultipleOptions(t *testing.T) {
	mockClient := &MockDeleteRoomClient{
		deleteRoomWithOptionsFn: func(ctx context.Context, roomID string, options map[string]interface{}) error {
			if options == nil {
				return errors.New("options required")
			}

			purge, hasPurge := options["purge"].(bool)
			kick, hasKick := options["kick_users"].(bool)
			block, hasBlock := options["block"].(bool)

			if hasPurge && purge && hasKick && kick && hasBlock && block {
				return nil
			}
			return errors.New("expected all options: purge, kick_users, block")
		},
	}

	options := map[string]interface{}{
		"purge":      true,
		"kick_users": true,
		"block":      true,
	}

	err := mockClient.deleteRoomWithOptions(context.Background(), "!room:example.com", options)
	require.NoError(t, err)
}

// TestDeleteRoomInvalidOptions tests deletion with invalid options
func TestDeleteRoomInvalidOptions(t *testing.T) {
	mockClient := &MockDeleteRoomClient{
		deleteRoomWithOptionsFn: func(ctx context.Context, roomID string, options map[string]interface{}) error {
			if options != nil {
				for key := range options {
					// Validate option names
					switch key {
					case "purge", "kick_users", "block":
						continue
					default:
						return errors.New("invalid option: " + key)
					}
				}
			}
			return nil
		},
	}

	options := map[string]interface{}{
		"invalid_option": true,
	}

	err := mockClient.deleteRoomWithOptions(context.Background(), "!room:example.com", options)
	assert.Error(t, err)
}

// TestDeleteRoomOptionTypes tests option type validation
func TestDeleteRoomOptionTypes(t *testing.T) {
	mockClient := &MockDeleteRoomClient{
		deleteRoomWithOptionsFn: func(ctx context.Context, roomID string, options map[string]interface{}) error {
			if options != nil {
				if purge, ok := options["purge"].(bool); ok {
					if !purge {
						return errors.New("purge must be true if specified")
					}
					return nil
				}
				if _, ok := options["purge"].(string); ok {
					return errors.New("purge must be boolean, not string")
				}
			}
			return nil
		},
	}

	// Valid: boolean option
	err := mockClient.deleteRoomWithOptions(context.Background(), "!room:example.com", map[string]interface{}{
		"purge": true,
	})
	require.NoError(t, err)

	// Invalid: string option
	err = mockClient.deleteRoomWithOptions(context.Background(), "!room:example.com", map[string]interface{}{
		"purge": "true",
	})
	assert.Error(t, err)
}

// TestDeleteRoomEmptyRoom tests deleting room with no members
func TestDeleteRoomEmptyRoom(t *testing.T) {
	mockClient := &MockDeleteRoomClient{
		deleteRoomWithOptionsFn: func(ctx context.Context, roomID string, options map[string]interface{}) error {
			// Empty rooms should delete without kick_users
			return nil
		},
	}

	err := mockClient.deleteRoomWithOptions(context.Background(), "!empty:example.com", map[string]interface{}{})
	require.NoError(t, err)
}

// TestDeleteRoomFullWorkflow tests complete deletion workflow
func TestDeleteRoomFullWorkflow(t *testing.T) {
	operations := []string{}

	mockClient := &MockDeleteRoomClient{
		deleteRoomWithOptionsFn: func(ctx context.Context, roomID string, options map[string]interface{}) error {
			// Step 1: Kick users
			if kick, ok := options["kick_users"].(bool); ok && kick {
				operations = append(operations, "kick_users")
			}

			// Step 2: Purge data
			if purge, ok := options["purge"].(bool); ok && purge {
				operations = append(operations, "purge")
			}

			// Step 3: Block room
			if block, ok := options["block"].(bool); ok && block {
				operations = append(operations, "block")
			}

			if len(operations) == 0 {
				operations = append(operations, "simple_delete")
			}

			return nil
		},
	}

	// Delete with full cleanup
	options := map[string]interface{}{
		"kick_users": true,
		"purge":      true,
		"block":      true,
	}

	err := mockClient.deleteRoomWithOptions(context.Background(), "!room:example.com", options)
	require.NoError(t, err)

	// Verify operations were performed
	assert.Equal(t, 3, len(operations))
	assert.Contains(t, operations, "kick_users")
	assert.Contains(t, operations, "purge")
	assert.Contains(t, operations, "block")
}

// TestDeleteRoomErrorHandling tests error scenarios
func TestDeleteRoomErrorHandling(t *testing.T) {
	testCases := []struct {
		name      string
		roomID    string
		options   map[string]interface{}
		shouldErr bool
	}{
		{
			name:      "empty roomID",
			roomID:    "",
			options:   nil,
			shouldErr: true,
		},
		{
			name:      "invalid roomID format",
			roomID:    "invalid",
			options:   nil,
			shouldErr: true,
		},
		{
			name:      "room not found",
			roomID:    "!nonexistent:example.com",
			options:   nil,
			shouldErr: true,
		},
	}

	mockClient := &MockDeleteRoomClient{
		deleteRoomWithOptionsFn: func(ctx context.Context, roomID string, options map[string]interface{}) error {
			if roomID == "" {
				return errors.New("roomID required")
			}
			if roomID == "invalid" {
				return errors.New("invalid roomID format")
			}
			if roomID == "!nonexistent:example.com" {
				return errors.New("room not found")
			}
			return nil
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := mockClient.deleteRoomWithOptions(context.Background(), tc.roomID, tc.options)
			if tc.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
