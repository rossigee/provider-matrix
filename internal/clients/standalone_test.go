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

package clients_test

import (
	"fmt"
	"strings"
	"testing"
)

// Copy of validateMatrixID for testing
func validateMatrixID(matrixID, idType string) error {
	switch idType {
	case "user":
		if !strings.HasPrefix(matrixID, "@") {
			return fmt.Errorf("user ID must start with @")
		}
	case "room":
		if !strings.HasPrefix(matrixID, "!") {
			return fmt.Errorf("room ID must start with !")
		}
	case "alias":
		if !strings.HasPrefix(matrixID, "#") {
			return fmt.Errorf("room alias must start with #")
		}
	}

	parts := strings.Split(matrixID[1:], ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid Matrix ID format: %s", matrixID)
	}

	return nil
}

// Copy of extractDomain for testing
func extractDomain(matrixID string) string {
	parts := strings.Split(matrixID, ":")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
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
			name:     "invalid format - no colon",
			matrixID: "@alice",
			idType:   "user",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMatrixID(tt.matrixID, tt.idType)

			if (err != nil) != tt.wantErr {
				t.Errorf("validateMatrixID() error = %v, wantErr %v", err, tt.wantErr)
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
			if got != tt.want {
				t.Errorf("extractDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
