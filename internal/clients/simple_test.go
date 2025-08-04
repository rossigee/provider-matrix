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
	"testing"
	"time"
)

// Test the types without requiring full Kubernetes machinery
func TestUserType(t *testing.T) {
	now := time.Now()
	user := &User{
		UserID:       "@test:example.com",
		DisplayName:  "Test User",
		Admin:        false,
		Deactivated:  false,
		CreationTime: &now,
	}

	if user.UserID != "@test:example.com" {
		t.Errorf("Expected UserID to be @test:example.com, got %s", user.UserID)
	}

	if user.DisplayName != "Test User" {
		t.Errorf("Expected DisplayName to be Test User, got %s", user.DisplayName)
	}

	if user.Admin != false {
		t.Errorf("Expected Admin to be false, got %v", user.Admin)
	}
}

func TestRoomType(t *testing.T) {
	room := &Room{
		RoomID: "!test:example.com",
		Name:   "Test Room",
		Topic:  "A test room",
	}

	if room.RoomID != "!test:example.com" {
		t.Errorf("Expected RoomID to be !test:example.com, got %s", room.RoomID)
	}

	if room.Name != "Test Room" {
		t.Errorf("Expected Name to be Test Room, got %s", room.Name)
	}
}

func TestPowerLevelContent(t *testing.T) {
	defaultLevel := 50
	plc := &PowerLevelContent{
		Users: map[string]int{
			"@admin:example.com": 100,
			"@mod:example.com":   50,
		},
		EventsDefault: &defaultLevel,
	}

	if plc.Users["@admin:example.com"] != 100 {
		t.Errorf("Expected admin power level to be 100, got %d", plc.Users["@admin:example.com"])
	}

	if *plc.EventsDefault != 50 {
		t.Errorf("Expected EventsDefault to be 50, got %d", *plc.EventsDefault)
	}
}
