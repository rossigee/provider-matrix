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

package clients

import (
	"context"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
	"github.com/pkg/errors"
)

// getIntValue returns the value of an int pointer or a default value
func getIntValue(ptr *int, defaultValue int) int {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
}

// User operations

// CreateUser creates a new Matrix user
func (c *matrixClient) CreateUser(ctx context.Context, userSpec *UserSpec) (*User, error) {
	// Use admin API if available and enabled
	if c.adminClient != nil {
		return c.adminClient.createUser(ctx, userSpec)
	}

	// Fallback to standard user registration (limited functionality)
	return nil, errors.New("user creation requires admin API access")
}

// GetUser retrieves user information
func (c *matrixClient) GetUser(ctx context.Context, userID string) (*User, error) {
	// Validate user ID format
	if err := validateMatrixID(userID, "user"); err != nil {
		return nil, errors.Wrap(err, "invalid user ID")
	}

	// Try admin API first if available
	if c.adminClient != nil {
		return c.adminClient.getUser(ctx, userID)
	}

	// Fallback to profile API for basic info
	profile, err := c.client.GetProfile(ctx, id.UserID(userID))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user profile")
	}

	user := &User{
		UserID:      userID,
		DisplayName: profile.DisplayName,
		AvatarURL:   profile.AvatarURL.String(),
	}

	return user, nil
}

// UpdateUser updates user information
func (c *matrixClient) UpdateUser(ctx context.Context, userID string, userSpec *UserSpec) (*User, error) {
	if err := validateMatrixID(userID, "user"); err != nil {
		return nil, errors.Wrap(err, "invalid user ID")
	}

	// Use admin API if available
	if c.adminClient != nil {
		return c.adminClient.updateUser(ctx, userID, userSpec)
	}

	// Fallback to basic profile updates
	if userSpec.DisplayName != "" {
		err := c.client.SetDisplayName(ctx, userSpec.DisplayName)
		if err != nil {
			return nil, errors.Wrap(err, "failed to set display name")
		}
	}

	if userSpec.AvatarURL != "" {
		avatarURL, err := id.ParseContentURI(userSpec.AvatarURL)
		if err != nil {
			return nil, errors.Wrap(err, "invalid avatar URL")
		}
		err = c.client.SetAvatarURL(ctx, avatarURL)
		if err != nil {
			return nil, errors.Wrap(err, "failed to set avatar URL")
		}
	}

	return c.GetUser(ctx, userID)
}

// DeactivateUser deactivates a user account
func (c *matrixClient) DeactivateUser(ctx context.Context, userID string) error {
	if c.adminClient == nil {
		return errors.New("user deactivation requires admin API access")
	}

	if err := validateMatrixID(userID, "user"); err != nil {
		return errors.Wrap(err, "invalid user ID")
	}

	return c.adminClient.deactivateUser(ctx, userID)
}

// Room operations

// CreateRoom creates a new Matrix room
func (c *matrixClient) CreateRoom(ctx context.Context, roomSpec *RoomSpec) (*Room, error) {
	// Build mautrix room creation request
	req := &mautrix.ReqCreateRoom{
		Name:         roomSpec.Name,
		Topic:        roomSpec.Topic,
		RoomAliasName: roomSpec.Alias,
		Preset:       roomSpec.Preset,
		Visibility:   roomSpec.Visibility,
		RoomVersion:  roomSpec.RoomVersion,
		CreationContent: roomSpec.CreationContent,
		Invite:       make([]id.UserID, len(roomSpec.Invite)),
	}

	// Convert invite list
	for i, userID := range roomSpec.Invite {
		req.Invite[i] = id.UserID(userID)
	}

	// Convert initial state
	for _, state := range roomSpec.InitialState {
		req.InitialState = append(req.InitialState, &event.Event{
			Type:     event.Type{Type: state.Type},
			StateKey: &state.StateKey,
			Content:  event.Content{Parsed: state.Content},
		})
	}

	// Set power level overrides if provided
	if roomSpec.PowerLevelOverrides != nil {
		// Convert user IDs in power levels
		userLevels := make(map[id.UserID]int)
		for userID, level := range roomSpec.PowerLevelOverrides.Users {
			userLevels[id.UserID(userID)] = level
		}
		
		req.PowerLevelOverride = &event.PowerLevelsEventContent{
			Users:           userLevels,
			Events:          roomSpec.PowerLevelOverrides.Events,
			EventsDefault:   getIntValue(roomSpec.PowerLevelOverrides.EventsDefault, 0),
			StateDefaultPtr: roomSpec.PowerLevelOverrides.StateDefault,
			UsersDefault:    getIntValue(roomSpec.PowerLevelOverrides.UsersDefault, 0),
			BanPtr:          roomSpec.PowerLevelOverrides.Ban,
			KickPtr:         roomSpec.PowerLevelOverrides.Kick,
			RedactPtr:       roomSpec.PowerLevelOverrides.Redact,
			InvitePtr:       roomSpec.PowerLevelOverrides.Invite,
		}
	}

	// Create the room
	resp, err := c.client.CreateRoom(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create room")
	}

	// Set additional room state if needed
	roomID := resp.RoomID.String()
	
	if roomSpec.GuestAccess != "" {
		_, err = c.client.SendStateEvent(ctx, resp.RoomID, event.StateGuestAccess, "", &event.GuestAccessEventContent{
			GuestAccess: event.GuestAccess(roomSpec.GuestAccess),
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to set guest access")
		}
	}

	if roomSpec.HistoryVisibility != "" {
		_, err = c.client.SendStateEvent(ctx, resp.RoomID, event.StateHistoryVisibility, "", &event.HistoryVisibilityEventContent{
			HistoryVisibility: event.HistoryVisibility(roomSpec.HistoryVisibility),
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to set history visibility")
		}
	}

	if roomSpec.JoinRules != "" {
		_, err = c.client.SendStateEvent(ctx, resp.RoomID, event.StateJoinRules, "", &event.JoinRulesEventContent{
			JoinRule: event.JoinRule(roomSpec.JoinRules),
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to set join rules")
		}
	}

	if roomSpec.EncryptionEnabled {
		_, err = c.client.SendStateEvent(ctx, resp.RoomID, event.StateEncryption, "", &event.EncryptionEventContent{
			Algorithm: id.AlgorithmMegolmV1,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to enable encryption")
		}
	}

	if roomSpec.AvatarURL != "" {
		avatarURL, err := id.ParseContentURI(roomSpec.AvatarURL)
		if err != nil {
			return nil, errors.Wrap(err, "invalid avatar URL")
		}
		_, err = c.client.SendStateEvent(ctx, resp.RoomID, event.StateRoomAvatar, "", &event.RoomAvatarEventContent{
			URL: avatarURL,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to set avatar")
		}
	}

	return c.GetRoom(ctx, roomID)
}

// GetRoom retrieves room information
func (c *matrixClient) GetRoom(ctx context.Context, roomID string) (*Room, error) {
	if err := validateMatrixID(roomID, "room"); err != nil {
		return nil, errors.Wrap(err, "invalid room ID")
	}

	roomIDObj := id.RoomID(roomID)

	// Try admin API first for comprehensive info
	if c.adminClient != nil {
		room, err := c.adminClient.getRoomDetails(ctx, roomID)
		if err == nil {
			return room, nil
		}
		// Fall back to standard API if admin fails
	}

	// Get basic room state using standard API
	room := &Room{
		RoomID: roomID,
	}

	// Get room name
	var nameContent event.RoomNameEventContent
	err := c.client.StateEvent(ctx, roomIDObj, event.StateRoomName, "", &nameContent)
	if err == nil {
		room.Name = nameContent.Name
	}

	// Get room topic
	var topicContent event.TopicEventContent
	err = c.client.StateEvent(ctx, roomIDObj, event.StateTopic, "", &topicContent)
	if err == nil {
		room.Topic = topicContent.Topic
	}

	// Get canonical alias
	var aliasContent event.CanonicalAliasEventContent
	err = c.client.StateEvent(ctx, roomIDObj, event.StateCanonicalAlias, "", &aliasContent)
	if err == nil && aliasContent.Alias != "" {
		room.Alias = aliasContent.Alias.String()
	}

	// Get avatar
	var avatarContent event.RoomAvatarEventContent
	err = c.client.StateEvent(ctx, roomIDObj, event.StateRoomAvatar, "", &avatarContent)
	if err == nil {
		room.AvatarURL = avatarContent.URL.String()
	}

	// Get power levels
	var powerContent event.PowerLevelsEventContent
	err = c.client.StateEvent(ctx, roomIDObj, event.StatePowerLevels, "", &powerContent)
	if err == nil {
		// Convert user IDs from mautrix format to our format
		users := make(map[string]int)
		for userID, level := range powerContent.Users {
			users[string(userID)] = level
		}
		
		room.PowerLevels = &PowerLevelContent{
			Users:         users,
			Events:        powerContent.Events,
			EventsDefault: &powerContent.EventsDefault,
			StateDefault:  powerContent.StateDefaultPtr,
			UsersDefault:  &powerContent.UsersDefault,
			Ban:           powerContent.BanPtr,
			Kick:          powerContent.KickPtr,
			Redact:        powerContent.RedactPtr,
			Invite:        powerContent.InvitePtr,
		}
	}

	return room, nil
}

// UpdateRoom updates room information
func (c *matrixClient) UpdateRoom(ctx context.Context, roomID string, roomSpec *RoomSpec) (*Room, error) {
	if err := validateMatrixID(roomID, "room"); err != nil {
		return nil, errors.Wrap(err, "invalid room ID")
	}

	roomIDObj := id.RoomID(roomID)

	// Update room name
	if roomSpec.Name != "" {
		_, err := c.client.SendStateEvent(ctx, roomIDObj, event.StateRoomName, "", &event.RoomNameEventContent{
			Name: roomSpec.Name,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to update room name")
		}
	}

	// Update room topic
	if roomSpec.Topic != "" {
		_, err := c.client.SendStateEvent(ctx, roomIDObj, event.StateTopic, "", &event.TopicEventContent{
			Topic: roomSpec.Topic,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to update room topic")
		}
	}

	// Update other room settings as needed...
	// (Similar pattern for other state events)

	return c.GetRoom(ctx, roomID)
}

// DeleteRoom deletes a room
func (c *matrixClient) DeleteRoom(ctx context.Context, roomID string) error {
	if c.adminClient == nil {
		return errors.New("room deletion requires admin API access")
	}

	if err := validateMatrixID(roomID, "room"); err != nil {
		return errors.Wrap(err, "invalid room ID")
	}

	options := map[string]interface{}{
		"block": false,
		"purge": true,
	}

	return c.adminClient.deleteRoom(ctx, roomID, options)
}

// Power level operations

// SetPowerLevels sets power levels in a room
func (c *matrixClient) SetPowerLevels(ctx context.Context, roomID string, powerLevels *PowerLevelSpec) error {
	if err := validateMatrixID(roomID, "room"); err != nil {
		return errors.Wrap(err, "invalid room ID")
	}

	roomIDObj := id.RoomID(roomID)
	
	// Convert user IDs to mautrix format
	users := make(map[id.UserID]int)
	for userID, level := range powerLevels.PowerLevels.Users {
		users[id.UserID(userID)] = level
	}
	
	content := &event.PowerLevelsEventContent{
		Users:           users,
		Events:          powerLevels.PowerLevels.Events,
		EventsDefault:   getIntValue(powerLevels.PowerLevels.EventsDefault, 0),
		StateDefaultPtr: powerLevels.PowerLevels.StateDefault,
		UsersDefault:    getIntValue(powerLevels.PowerLevels.UsersDefault, 0),
		BanPtr:          powerLevels.PowerLevels.Ban,
		KickPtr:         powerLevels.PowerLevels.Kick,
		RedactPtr:       powerLevels.PowerLevels.Redact,
		InvitePtr:       powerLevels.PowerLevels.Invite,
	}

	_, err := c.client.SendStateEvent(ctx, roomIDObj, event.StatePowerLevels, "", content)
	if err != nil {
		return errors.Wrap(err, "failed to set power levels")
	}

	return nil
}

// GetPowerLevels retrieves power levels from a room
func (c *matrixClient) GetPowerLevels(ctx context.Context, roomID string) (*PowerLevelContent, error) {
	if err := validateMatrixID(roomID, "room"); err != nil {
		return nil, errors.Wrap(err, "invalid room ID")
	}

	roomIDObj := id.RoomID(roomID)
	var powerContent event.PowerLevelsEventContent
	err := c.client.StateEvent(ctx, roomIDObj, event.StatePowerLevels, "", &powerContent)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get power levels")
	}

	// Convert user IDs from mautrix format to our format
	users := make(map[string]int)
	for userID, level := range powerContent.Users {
		users[string(userID)] = level
	}

	return &PowerLevelContent{
		Users:         users,
		Events:        powerContent.Events,
		EventsDefault: &powerContent.EventsDefault,
		StateDefault:  powerContent.StateDefaultPtr,
		UsersDefault:  &powerContent.UsersDefault,
		Ban:           powerContent.BanPtr,
		Kick:          powerContent.KickPtr,
		Redact:        powerContent.RedactPtr,
		Invite:        powerContent.InvitePtr,
	}, nil
}

// Room alias operations

// CreateRoomAlias creates a room alias
func (c *matrixClient) CreateRoomAlias(ctx context.Context, alias string, roomID string) error {
	if err := validateMatrixID(alias, "alias"); err != nil {
		return errors.Wrap(err, "invalid alias")
	}
	if err := validateMatrixID(roomID, "room"); err != nil {
		return errors.Wrap(err, "invalid room ID")
	}

	aliasID := id.RoomAlias(alias)
	roomIDObj := id.RoomID(roomID)

	_, err := c.client.CreateAlias(ctx, aliasID, roomIDObj)
	if err != nil {
		return errors.Wrap(err, "failed to create room alias")
	}

	return nil
}

// GetRoomAlias retrieves room alias information
func (c *matrixClient) GetRoomAlias(ctx context.Context, alias string) (*RoomAlias, error) {
	if err := validateMatrixID(alias, "alias"); err != nil {
		return nil, errors.Wrap(err, "invalid alias")
	}

	aliasID := id.RoomAlias(alias)
	resp, err := c.client.ResolveAlias(ctx, aliasID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve alias")
	}

	return &RoomAlias{
		Alias:  alias,
		RoomID: resp.RoomID.String(),
	}, nil
}

// DeleteRoomAlias deletes a room alias
func (c *matrixClient) DeleteRoomAlias(ctx context.Context, alias string) error {
	if err := validateMatrixID(alias, "alias"); err != nil {
		return errors.Wrap(err, "invalid alias")
	}

	aliasID := id.RoomAlias(alias)
	_, err := c.client.DeleteAlias(ctx, aliasID)
	if err != nil {
		return errors.Wrap(err, "failed to delete room alias")
	}

	return nil
}