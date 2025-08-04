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
)

func TestValidateMatrixIDUtil(t *testing.T) {
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

			if (err != nil) != tt.wantErr {
				t.Errorf("validateMatrixID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExtractDomainUtil(t *testing.T) {
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
		{
			name:     "no domain part",
			matrixID: "@user:",
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
