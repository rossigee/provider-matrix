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
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/id"

	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane-contrib/provider-matrix/apis/v1beta1"
)

const (
	// DefaultTimeout for Matrix API operations
	defaultTimeout = 30 * time.Second
)

// Client interface for Matrix API operations
type Client interface {
	// User operations
	CreateUser(ctx context.Context, user *UserSpec) (*User, error)
	GetUser(ctx context.Context, userID string) (*User, error)
	UpdateUser(ctx context.Context, userID string, user *UserSpec) (*User, error)
	DeactivateUser(ctx context.Context, userID string) error

	// Room operations
	CreateRoom(ctx context.Context, room *RoomSpec) (*Room, error)
	GetRoom(ctx context.Context, roomID string) (*Room, error)
	UpdateRoom(ctx context.Context, roomID string, room *RoomSpec) (*Room, error)
	DeleteRoom(ctx context.Context, roomID string) error

	// Power level operations
	SetPowerLevels(ctx context.Context, roomID string, powerLevels *PowerLevelSpec) error
	GetPowerLevels(ctx context.Context, roomID string) (*PowerLevelContent, error)

	// Room alias operations
	CreateRoomAlias(ctx context.Context, alias string, roomID string) error
	GetRoomAlias(ctx context.Context, alias string) (*RoomAlias, error)
	DeleteRoomAlias(ctx context.Context, alias string) error
}

// Config holds the configuration for the Matrix client
type Config struct {
	HomeserverURL string
	AdminAPIURL   string
	AccessToken   string
	UserID        string
	DeviceID      string
	ServerType    string
	AdminMode     bool
	HTTPClient    *http.Client
}

// matrixClient implements the Client interface using mautrix-go
type matrixClient struct {
	config      *Config
	client      *mautrix.Client
	adminClient *adminClient
}

// NewClient creates a new Matrix client
func NewClient(config *Config) (Client, error) {
	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{
			Timeout: defaultTimeout,
		}
	}

	// Create mautrix client
	client, err := mautrix.NewClient(config.HomeserverURL, "", "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create mautrix client")
	}

	client.AccessToken = config.AccessToken
	client.UserID = id.UserID(config.UserID)
	client.DeviceID = id.DeviceID(config.DeviceID)
	client.Client = config.HTTPClient

	// Create admin client if admin mode is enabled
	var adminClient *adminClient
	if config.AdminMode {
		adminClient = newAdminClient(config)
	}

	return &matrixClient{
		config:      config,
		client:      client,
		adminClient: adminClient,
	}, nil
}

// GetConfig extracts the configuration from the provider config
func GetConfig(ctx context.Context, c client.Client, mg resource.Managed) (*Config, error) {
	switch {
	case mg.GetProviderConfigReference() != nil:
		return UseProviderConfig(ctx, c, mg)
	default:
		return nil, errors.New("no credentials specified")
	}
}

// UseProviderConfig extracts configuration from a ProviderConfig
func UseProviderConfig(ctx context.Context, c client.Client, mg resource.Managed) (*Config, error) {
	pc := &v1beta1.ProviderConfig{}
	if err := c.Get(ctx, types.NamespacedName{Name: mg.GetProviderConfigReference().Name}, pc); err != nil {
		return nil, errors.Wrap(err, "cannot get referenced ProviderConfig")
	}

	// TODO: Fix ProviderConfigUsage tracking with newer crossplane-runtime
	// t := resource.NewProviderConfigUsageTracker(c, &v1beta1.ProviderConfigUsage{})
	// if err := t.Track(ctx, mg); err != nil {
	// 	return nil, errors.Wrap(err, "cannot track ProviderConfig usage")
	// }

	credBytes, err := resource.CommonCredentialExtractor(ctx, pc.Spec.Credentials.Source, c, pc.Spec.Credentials.CommonCredentialSelectors)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get credentials")
	}

	if len(credBytes) == 0 {
		return nil, errors.New("matrix access token not found in credentials")
	}
	accessToken := string(credBytes)

	adminAPIURL := pc.Spec.HomeserverURL
	if pc.Spec.AdminAPIURL != nil {
		adminAPIURL = *pc.Spec.AdminAPIURL
	}

	serverType := "auto"
	if pc.Spec.ServerType != nil {
		serverType = *pc.Spec.ServerType
	}

	adminMode := false
	if pc.Spec.AdminMode != nil {
		adminMode = *pc.Spec.AdminMode
	}

	userID := ""
	if pc.Spec.UserID != nil {
		userID = *pc.Spec.UserID
	}

	deviceID := ""
	if pc.Spec.DeviceID != nil {
		deviceID = *pc.Spec.DeviceID
	}

	return &Config{
		HomeserverURL: pc.Spec.HomeserverURL,
		AdminAPIURL:   adminAPIURL,
		AccessToken:   accessToken,
		UserID:        userID,
		DeviceID:      deviceID,
		ServerType:    serverType,
		AdminMode:     adminMode,
	}, nil
}

// IsNotFound checks if an error represents a "not found" condition
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}

	// Check for Matrix-specific not found errors
	if mautrixErr, ok := err.(mautrix.HTTPError); ok {
		return mautrixErr.RespError != nil && mautrixErr.RespError.ErrCode == "M_NOT_FOUND"
	}

	// Check for HTTP 404
	if strings.Contains(err.Error(), "404") || strings.Contains(strings.ToLower(err.Error()), "not found") {
		return true
	}

	return false
}

// Helper method to validate Matrix IDs
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

// Helper method to extract domain from Matrix ID
func extractDomain(matrixID string) string {
	parts := strings.Split(matrixID, ":")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}
