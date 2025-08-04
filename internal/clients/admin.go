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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// adminClient handles Matrix admin API operations (primarily for Synapse)
type adminClient struct {
	config     *Config
	httpClient *http.Client
	baseURL    string
}

// newAdminClient creates a new admin API client
func newAdminClient(config *Config) *adminClient {
	baseURL := config.AdminAPIURL
	if baseURL == "" {
		baseURL = config.HomeserverURL
	}

	return &adminClient{
		config:     config,
		httpClient: config.HTTPClient,
		baseURL:    baseURL,
	}
}

// makeRequest makes an HTTP request to the admin API
func (c *adminClient) makeRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal request body")
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	url := fmt.Sprintf("%s%s", strings.TrimSuffix(c.baseURL, "/"), path)
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.AccessToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "crossplane-provider-matrix")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute request")
	}

	return resp, nil
}

// handleResponse processes the HTTP response
func (c *adminClient) handleResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return errors.Errorf("admin API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if target != nil {
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			return errors.Wrap(err, "failed to decode response")
		}
	}

	return nil
}

// User admin operations

// createUser creates a new user via admin API
func (c *adminClient) createUser(ctx context.Context, userSpec *UserSpec) (*User, error) {
	path := fmt.Sprintf("/_synapse/admin/v2/users/%s", url.PathEscape(userSpec.UserID))

	resp, err := c.makeRequest(ctx, "PUT", path, userSpec)
	if err != nil {
		return nil, err
	}

	var user User
	if err := c.handleResponse(resp, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// getUser retrieves user information via admin API
func (c *adminClient) getUser(ctx context.Context, userID string) (*User, error) {
	path := fmt.Sprintf("/_synapse/admin/v2/users/%s", url.PathEscape(userID))

	resp, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var user User
	if err := c.handleResponse(resp, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// updateUser updates user information via admin API
func (c *adminClient) updateUser(ctx context.Context, userID string, userSpec *UserSpec) (*User, error) {
	path := fmt.Sprintf("/_synapse/admin/v2/users/%s", url.PathEscape(userID))

	resp, err := c.makeRequest(ctx, "PUT", path, userSpec)
	if err != nil {
		return nil, err
	}

	var user User
	if err := c.handleResponse(resp, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// deactivateUser deactivates a user via admin API
func (c *adminClient) deactivateUser(ctx context.Context, userID string) error {
	path := fmt.Sprintf("/_synapse/admin/v1/deactivate/%s", url.PathEscape(userID))

	resp, err := c.makeRequest(ctx, "POST", path, map[string]interface{}{
		"erase": false,
	})
	if err != nil {
		return err
	}

	return c.handleResponse(resp, nil)
}

// listUsers lists users via admin API
func (c *adminClient) listUsers(ctx context.Context, from string, limit int) (*ListUsersResponse, error) {
	path := "/_synapse/admin/v2/users"

	params := url.Values{}
	if from != "" {
		params.Set("from", from)
	}
	if limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}

	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	resp, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result ListUsersResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Room admin operations

// deleteRoom deletes a room via admin API
func (c *adminClient) deleteRoom(ctx context.Context, roomID string, options map[string]interface{}) error {
	path := fmt.Sprintf("/_synapse/admin/v1/rooms/%s/delete", url.PathEscape(roomID))

	if options == nil {
		options = make(map[string]interface{})
	}

	resp, err := c.makeRequest(ctx, "POST", path, options)
	if err != nil {
		return err
	}

	return c.handleResponse(resp, nil)
}

// getRoomDetails gets detailed room information via admin API
func (c *adminClient) getRoomDetails(ctx context.Context, roomID string) (*Room, error) {
	path := fmt.Sprintf("/_synapse/admin/v1/rooms/%s", url.PathEscape(roomID))

	resp, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var room Room
	if err := c.handleResponse(resp, &room); err != nil {
		return nil, err
	}

	return &room, nil
}

// listRooms lists rooms via admin API
func (c *adminClient) listRooms(ctx context.Context, from string, limit int) (*ListRoomsResponse, error) {
	path := "/_synapse/admin/v1/rooms"

	params := url.Values{}
	if from != "" {
		params.Set("from", from)
	}
	if limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}

	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	resp, err := c.makeRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result ListRoomsResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// makeRoomAdmin grants admin privileges to a user in a room
func (c *adminClient) makeRoomAdmin(ctx context.Context, roomID, userID string) error {
	path := fmt.Sprintf("/_synapse/admin/v1/rooms/%s/make_room_admin", url.PathEscape(roomID))

	body := map[string]interface{}{
		"user_id": userID,
	}

	resp, err := c.makeRequest(ctx, "POST", path, body)
	if err != nil {
		return err
	}

	return c.handleResponse(resp, nil)
}

// blockRoom blocks a room from being joined
func (c *adminClient) blockRoom(ctx context.Context, roomID string, block bool) error {
	path := fmt.Sprintf("/_synapse/admin/v1/rooms/%s/block", url.PathEscape(roomID))

	body := map[string]interface{}{
		"block": block,
	}

	resp, err := c.makeRequest(ctx, "PUT", path, body)
	if err != nil {
		return err
	}

	return c.handleResponse(resp, nil)
}
