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
	"time"
)

// User represents a Matrix user
type User struct {
	UserID       string            `json:"user_id"`
	DisplayName  string            `json:"displayname,omitempty"`
	AvatarURL    string            `json:"avatar_url,omitempty"`
	Admin        bool              `json:"admin"`
	Deactivated  bool              `json:"deactivated"`
	CreationTime *time.Time        `json:"creation_ts,omitempty"`
	LastSeenTime *time.Time        `json:"last_seen_ts,omitempty"`
	UserType     string            `json:"user_type,omitempty"`
	ExternalIDs  []ExternalID      `json:"external_ids,omitempty"`
	Devices      []Device          `json:"devices,omitempty"`
}

// UserSpec represents the parameters for creating/updating a user
type UserSpec struct {
	UserID      string        `json:"user_id,omitempty"`
	Localpart   string        `json:"localpart,omitempty"`
	Password    string        `json:"password,omitempty"`
	DisplayName string        `json:"displayname,omitempty"`
	AvatarURL   string        `json:"avatar_url,omitempty"`
	Admin       bool          `json:"admin"`
	Deactivated bool          `json:"deactivated"`
	UserType    string        `json:"user_type,omitempty"`
	ExternalIDs []ExternalID  `json:"external_ids,omitempty"`
	ExpireTime  *time.Time    `json:"expire_time,omitempty"`
}

// ExternalID represents a third-party identifier
type ExternalID struct {
	Medium    string `json:"medium"`
	Address   string `json:"address"`
	Validated bool   `json:"validated"`
}

// Device represents a Matrix device
type Device struct {
	DeviceID     string     `json:"device_id"`
	DisplayName  string     `json:"display_name,omitempty"`
	LastSeenIP   string     `json:"last_seen_ip,omitempty"`
	LastSeenTime *time.Time `json:"last_seen_ts,omitempty"`
}

// Room represents a Matrix room
type Room struct {
	RoomID            string              `json:"room_id"`
	Name              string              `json:"name,omitempty"`
	Topic             string              `json:"topic,omitempty"`
	Alias             string              `json:"canonical_alias,omitempty"`
	AvatarURL         string              `json:"avatar,omitempty"`
	Creator           string              `json:"creator,omitempty"`
	CreationTime      *time.Time          `json:"creation_ts,omitempty"`
	RoomVersion       string              `json:"room_version,omitempty"`
	JoinedMembers     int                 `json:"joined_members"`
	InvitedMembers    int                 `json:"invited_members"`
	Visibility        string              `json:"visibility,omitempty"`
	GuestAccess       string              `json:"guest_access,omitempty"`
	HistoryVisibility string              `json:"history_visibility,omitempty"`
	JoinRules         string              `json:"join_rules,omitempty"`
	EncryptionEnabled bool                `json:"encryption,omitempty"`
	PowerLevels       *PowerLevelContent  `json:"power_levels,omitempty"`
	State             []StateEvent        `json:"state,omitempty"`
}

// RoomSpec represents the parameters for creating/updating a room
type RoomSpec struct {
	Name                string                 `json:"name,omitempty"`
	Topic               string                 `json:"topic,omitempty"`
	Alias               string                 `json:"room_alias_name,omitempty"`
	Preset              string                 `json:"preset,omitempty"`
	Visibility          string                 `json:"visibility,omitempty"`
	RoomVersion         string                 `json:"room_version,omitempty"`
	CreationContent     map[string]interface{} `json:"creation_content,omitempty"`
	InitialState        []StateEvent           `json:"initial_state,omitempty"`
	Invite              []string               `json:"invite,omitempty"`
	PowerLevelOverrides *PowerLevelContent     `json:"power_level_content_override,omitempty"`
	GuestAccess         string                 `json:"guest_access,omitempty"`
	HistoryVisibility   string                 `json:"history_visibility,omitempty"`
	JoinRules           string                 `json:"join_rules,omitempty"`
	EncryptionEnabled   bool                   `json:"encryption,omitempty"`
	AvatarURL           string                 `json:"avatar_url,omitempty"`
}

// StateEvent represents a Matrix state event
type StateEvent struct {
	Type     string                 `json:"type"`
	StateKey string                 `json:"state_key"`
	Content  map[string]interface{} `json:"content"`
}

// PowerLevelContent defines power levels for room events and users
type PowerLevelContent struct {
	Users         map[string]int `json:"users,omitempty"`
	Events        map[string]int `json:"events,omitempty"`
	EventsDefault *int           `json:"events_default,omitempty"`
	StateDefault  *int           `json:"state_default,omitempty"`
	UsersDefault  *int           `json:"users_default,omitempty"`
	Ban           *int           `json:"ban,omitempty"`
	Kick          *int           `json:"kick,omitempty"`
	Redact        *int           `json:"redact,omitempty"`
	Invite        *int           `json:"invite,omitempty"`
}

// PowerLevelSpec represents the parameters for setting power levels
type PowerLevelSpec struct {
	RoomID       string             `json:"room_id"`
	PowerLevels  *PowerLevelContent `json:"power_levels"`
}

// RoomAlias represents a Matrix room alias
type RoomAlias struct {
	Alias  string `json:"alias"`
	RoomID string `json:"room_id"`
}

// Space represents a Matrix space (special type of room)
type Space struct {
	Room          // Embedded room fields
	SpaceType     string      `json:"type"`          // Should be "m.space"
	Children      []SpaceChild `json:"children,omitempty"`
}

// SpaceChild represents a child room or space within a space
type SpaceChild struct {
	RoomID      string   `json:"room_id"`
	Via         []string `json:"via,omitempty"`
	Order       string   `json:"order,omitempty"`
	Suggested   bool     `json:"suggested,omitempty"`
}

// SpaceSpec represents the parameters for creating/updating a space
type SpaceSpec struct {
	RoomSpec          // Embedded room spec fields
	Children []SpaceChild `json:"children,omitempty"`
}

// AdminResponse represents a generic admin API response
type AdminResponse struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
}

// ListUsersResponse represents the response from listing users
type ListUsersResponse struct {
	Users      []User `json:"users"`
	Total      int    `json:"total"`
	NextToken  string `json:"next_token,omitempty"`
	PrevToken  string `json:"prev_token,omitempty"`
}

// ListRoomsResponse represents the response from listing rooms
type ListRoomsResponse struct {
	Rooms     []Room `json:"rooms"`
	Total     int    `json:"total"`
	NextToken string `json:"next_token,omitempty"`
	PrevToken string `json:"prev_token,omitempty"`
}