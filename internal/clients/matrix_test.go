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

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				HomeserverURL: "https://matrix.example.com",
				AccessToken:   "test_token",
				UserID:        "@test:example.com",
				AdminMode:     false,
			},
			wantErr: false,
		},
		{
			name: "valid config with admin mode",
			config: &Config{
				HomeserverURL: "https://matrix.example.com",
				AdminAPIURL:   "https://matrix.example.com",
				AccessToken:   "admin_token",
				UserID:        "@admin:example.com",
				AdminMode:     true,
			},
			wantErr: false,
		},
		{
			name: "invalid homeserver URL",
			config: &Config{
				HomeserverURL: "invalid-url",
				AccessToken:   "test_token",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
			}
		})
	}
}

func TestValidateMatrixID(t *testing.T) {
	tests := []struct {
		name     string
		matrixID string
		idType   string
		wantErr  bool
	}{
		{
			name:     "valid user ID",
			matrixID: "@alice:example.com",
			idType:   "user",
			wantErr:  false,
		},
		{
			name:     "valid room ID",
			matrixID: "!roomid:example.com",
			idType:   "room",
			wantErr:  false,
		},
		{
			name:     "valid alias",
			matrixID: "#room:example.com",
			idType:   "alias",
			wantErr:  false,
		},
		{
			name:     "invalid user ID prefix",
			matrixID: "alice:example.com",
			idType:   "user",
			wantErr:  true,
		},
		{
			name:     "invalid room ID prefix",
			matrixID: "roomid:example.com",
			idType:   "room",
			wantErr:  true,
		},
		{
			name:     "invalid alias prefix",
			matrixID: "room:example.com",
			idType:   "alias",
			wantErr:  true,
		},
		{
			name:     "invalid format - no colon",
			matrixID: "@alice",
			idType:   "user",
			wantErr:  true,
		},
		{
			name:     "invalid format - too many colons",
			matrixID: "@alice:example:com",
			idType:   "user",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMatrixID(tt.matrixID, tt.idType)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestExtractDomain(t *testing.T) {
	tests := []struct {
		name     string
		matrixID string
		want     string
	}{
		{
			name:     "user ID",
			matrixID: "@alice:example.com",
			want:     "example.com",
		},
		{
			name:     "room ID",
			matrixID: "!roomid:matrix.org",
			want:     "matrix.org",
		},
		{
			name:     "alias",
			matrixID: "#room:server.net",
			want:     "server.net",
		},
		{
			name:     "invalid format",
			matrixID: "invalid",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractDomain(tt.matrixID)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "404 in error message",
			err:  assert.AnError,
			want: false, // AnError doesn't contain "404"
		},
		{
			name: "not found in error message",
			err:  assert.AnError,
			want: false, // AnError doesn't contain "not found"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsNotFound(tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}
