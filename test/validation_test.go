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
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestValidateMatrixID tests Matrix ID validation logic
func TestValidateMatrixID(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		shouldErr bool
		reason    string
	}{
		// Valid User IDs
		{
			name:      "valid user ID basic",
			id:        "@alice:example.com",
			shouldErr: false,
		},
		{
			name:      "valid user ID with underscore",
			id:        "@alice_bob:example.com",
			shouldErr: false,
		},
		{
			name:      "valid user ID with number",
			id:        "@alice123:example.com",
			shouldErr: false,
		},
		{
			name:      "valid user ID with dot",
			id:        "@alice.bob:example.com",
			shouldErr: false,
		},
		// Valid Room IDs
		{
			name:      "valid room ID",
			id:        "!abc123:example.com",
			shouldErr: false,
		},
		{
			name:      "valid room ID with long hash",
			id:        "!abcdef1234567890:example.com",
			shouldErr: false,
		},
		// Valid Alias IDs
		{
			name:      "valid alias ID",
			id:        "#general:example.com",
			shouldErr: false,
		},
		{
			name:      "valid alias with underscore",
			id:        "#general_chat:example.com",
			shouldErr: false,
		},
		// Valid with alternative servers
		{
			name:      "valid with matrix.org",
			id:        "@alice:matrix.org",
			shouldErr: false,
		},
		{
			name:      "valid with localhost",
			id:        "@alice:localhost",
			shouldErr: false,
		},
		{
			name:      "valid with port",
			id:        "@alice:example.com:8008",
			shouldErr: false,
		},
		// Invalid IDs
		{
			name:      "missing prefix",
			id:        "alice:example.com",
			shouldErr: true,
			reason:    "no @ prefix for user",
		},
		{
			name:      "missing colon",
			id:        "@aliceexample.com",
			shouldErr: true,
			reason:    "missing domain separator",
		},
		{
			name:      "missing host",
			id:        "@alice:",
			shouldErr: true,
			reason:    "no domain after colon",
		},
		{
			name:      "invalid prefix",
			id:        "$alice:example.com",
			shouldErr: true,
			reason:    "invalid prefix character",
		},
		{
			name:      "empty string",
			id:        "",
			shouldErr: true,
			reason:    "empty ID",
		},
		{
			name:      "only prefix",
			id:        "@",
			shouldErr: true,
			reason:    "incomplete ID",
		},
		{
			name:      "whitespace in ID",
			id:        "@alice bob:example.com",
			shouldErr: true,
			reason:    "spaces not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate Matrix ID validation logic
			valid := validateTestMatrixID(tt.id)
			if tt.shouldErr {
				assert.False(t, valid, "expected ID to be invalid: "+tt.reason)
			} else {
				assert.True(t, valid, "expected ID to be valid")
			}
		})
	}
}

// TestExtractDomain tests domain extraction from Matrix IDs
func TestExtractDomain(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		expected  string
		shouldErr bool
	}{
		{
			name:      "extract from user ID",
			id:        "@alice:example.com",
			expected:  "example.com",
			shouldErr: false,
		},
		{
			name:      "extract from room ID",
			id:        "!abc123:matrix.org",
			expected:  "matrix.org",
			shouldErr: false,
		},
		{
			name:      "extract from alias",
			id:        "#general:localhost",
			expected:  "localhost",
			shouldErr: false,
		},
		{
			name:      "extract with port",
			id:        "@alice:example.com:8008",
			expected:  "example.com:8008",
			shouldErr: false,
		},
		{
			name:      "extract IPv4",
			id:        "@alice:192.168.1.1",
			expected:  "192.168.1.1",
			shouldErr: false,
		},
		{
			name:      "invalid - no colon",
			id:        "@aliceexample.com",
			expected:  "",
			shouldErr: true,
		},
		{
			name:      "invalid - empty domain",
			id:        "@alice:",
			expected:  "",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			domain, err := extractTestDomain(tt.id)
			if tt.shouldErr {
				assert.Error(t, err, "expected extraction to fail")
			} else {
				assert.NoError(t, err, "expected extraction to succeed")
				assert.Equal(t, tt.expected, domain)
			}
		})
	}
}

// TestMatrixIDEdgeCases tests edge cases and boundary conditions
func TestMatrixIDEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		shouldErr bool
	}{
		// Length boundaries
		{
			name:      "very long local part",
			id:        "@" + stringOfLength(255) + ":example.com",
			shouldErr: false,
		},
		{
			name:      "very long domain",
			id:        "@alice:" + stringOfLength(255),
			shouldErr: false,
		},
		// Special cases
		{
			name:      "underscore in local",
			id:        "@alice_bob_charlie:example.com",
			shouldErr: false,
		},
		{
			name:      "numbers in domain",
			id:        "@alice:example123.com",
			shouldErr: false,
		},
		{
			name:      "hyphen in domain",
			id:        "@alice:example-test.com",
			shouldErr: false,
		},
		// Unicode (should fail in most Matrix implementations)
		{
			name:      "unicode in local",
			id:        "@αλίς:example.com",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := validateTestMatrixID(tt.id)
			if tt.shouldErr {
				assert.False(t, valid, "expected validation to fail for: "+tt.id)
			} else {
				assert.True(t, valid, "expected validation to pass for: "+tt.id)
			}
		})
	}
}

// Helper functions for testing

func validateTestMatrixID(id string) bool {
	if id == "" {
		return false
	}

	if len(id) < 3 {
		return false
	}

	prefix := id[0]
	if prefix != '@' && prefix != '!' && prefix != '#' {
		return false
	}

	colonIndex := -1
	for i, ch := range id {
		if ch == ':' {
			colonIndex = i
			break
		}
	}

	if colonIndex == -1 || colonIndex == 1 {
		return false
	}

	localPart := id[1:colonIndex]
	domain := id[colonIndex+1:]

	if domain == "" {
		return false
	}

	// Validate characters - only alphanumeric, underscore, hyphen, dot, colon (for port)
	for _, ch := range localPart {
		if !((ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '_' || ch == '.' || ch == '-') {
			return false
		}
	}

	// Check for whitespace
	for _, ch := range id {
		if ch == ' ' || ch == '\t' || ch == '\n' {
			return false
		}
	}

	return true
}

func extractTestDomain(id string) (string, error) {
	colonIndex := -1
	for i, ch := range id {
		if ch == ':' {
			colonIndex = i
			break
		}
	}

	if colonIndex == -1 || colonIndex >= len(id)-1 {
		return "", stringErr("invalid Matrix ID format")
	}

	domain := id[colonIndex+1:]
	if domain == "" {
		return "", stringErr("empty domain")
	}

	return domain, nil
}

func stringOfLength(length int) string {
	result := ""
	for i := 0; i < length; i++ {
		result += "a"
	}
	return result
}

func stringErr(msg string) error {
	return stringError{msg: msg}
}

type stringError struct {
	msg string
}

func (e stringError) Error() string {
	return e.msg
}
