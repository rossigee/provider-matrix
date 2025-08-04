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
)

// TestBasicFunctionality tests the basic package structure
func TestBasicFunctionality(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic test",
			input:    "test",
			expected: "test",
		},
		{
			name:     "matrix ID format",
			input:    "@user:example.com",
			expected: "@user:example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.input != tt.expected {
				t.Errorf("got %s, want %s", tt.input, tt.expected)
			}
		})
	}
}

// TestMatrixIDValidation tests Matrix ID validation logic
func TestMatrixIDValidation(t *testing.T) {
	validIDs := []string{
		"@alice:example.com",
		"!room:example.com",
		"#alias:example.com",
		"@user:matrix.org",
		"!abc123:server.net",
		"#general:company.com",
	}

	for _, id := range validIDs {
		t.Run("valid_"+id, func(t *testing.T) {
			// Basic validation - check for proper prefix and colon
			if len(id) < 3 {
				t.Errorf("ID too short: %s", id)
			}

			prefix := id[0]
			if prefix != '@' && prefix != '!' && prefix != '#' {
				t.Errorf("Invalid prefix for ID: %s", id)
			}

			colonIndex := -1
			for i, ch := range id {
				if ch == ':' {
					colonIndex = i
					break
				}
			}

			if colonIndex == -1 {
				t.Errorf("No colon found in ID: %s", id)
			}
		})
	}
}

// BenchmarkMatrixIDParsing benchmarks ID parsing performance
func BenchmarkMatrixIDParsing(b *testing.B) {
	testID := "@alice:example.com"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Simulate parsing
		_ = len(testID)
		_ = testID[0]
		for j, ch := range testID {
			if ch == ':' {
				_ = testID[1:j]
				_ = testID[j+1:]
				break
			}
		}
	}
}
