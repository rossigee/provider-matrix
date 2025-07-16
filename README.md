# Provider Matrix

[![CI](https://github.com/crossplane-contrib/provider-matrix/actions/workflows/ci.yml/badge.svg)](https://github.com/crossplane-contrib/provider-matrix/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/crossplane-contrib/provider-matrix)](https://goreportcard.com/report/github.com/crossplane-contrib/provider-matrix)

## Overview

`provider-matrix` is a [Crossplane](https://crossplane.io/) provider that enables infrastructure and application management for Matrix homeservers and related resources. It allows you to manage Matrix users, rooms, spaces, power levels, and room aliases declaratively using Kubernetes Custom Resources.

## Features

This provider supports the following Matrix resources:

- **User** (`user.matrix.crossplane.io`) - Manage Matrix users with profiles, administrative privileges, devices, and lifecycle
- **Room** (`room.matrix.crossplane.io`) - Create and manage Matrix rooms with custom settings, encryption, and access controls  
- **Space** (`space.matrix.crossplane.io`) - Organize rooms into hierarchical spaces for better organization
- **PowerLevel** (`powerlevel.matrix.crossplane.io`) - Configure granular permissions and power levels within rooms
- **RoomAlias** (`roomalias.matrix.crossplane.io`) - Create human-readable aliases for Matrix rooms

## Quick Start

### Install the Provider

Install the provider by using the following command:

```bash
# Install the Matrix provider
kubectl apply -f examples/provider/config.yaml
```

Notice that the provider will be installed in the `crossplane-system` namespace.

### Create a ProviderConfig

Before creating any Matrix resources, you need to create a `ProviderConfig` that contains the credentials for connecting to your Matrix homeserver:

```bash
kubectl apply -f examples/provider/providerconfig.yaml
```

**Note:** You need to update the secret with your actual Matrix access token.

### Create Matrix Resources

Once the provider is installed and configured, you can create Matrix resources:

```bash
# Create a Matrix user
kubectl apply -f examples/user/user.yaml

# Create a Matrix room  
kubectl apply -f examples/room/room.yaml

# Create a Matrix space
kubectl apply -f examples/space/space.yaml

# Configure room power levels
kubectl apply -f examples/powerlevel/powerlevel.yaml

# Create a room alias
kubectl apply -f examples/roomalias/roomalias.yaml
```

## Configuration

### Provider Configuration

The provider requires access to a Matrix homeserver with appropriate administrative privileges. Configure your credentials using a Kubernetes secret:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: matrix-creds
  namespace: crossplane-system
type: Opaque
data:
  credentials: <base64-encoded-access-token>
```

The ProviderConfig should reference this secret and specify your homeserver URL:

```yaml
apiVersion: matrix.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  homeserverURL: "https://matrix.example.com"
  adminMode: true
  credentials:
    source: Secret
    secretRef:
      namespace: crossplane-system
      name: matrix-creds
      key: credentials
```

### Configuration Options

- `homeserverURL` (required): The URL of your Matrix homeserver
- `adminAPIURL` (optional): The admin API URL (defaults to homeserverURL)
- `userID` (optional): User ID for the Matrix client
- `deviceID` (optional): Device ID for the Matrix client  
- `serverType` (optional): Server type hint (auto, synapse, dendrite, conduit)
- `adminMode` (optional): Enable admin mode for administrative operations

### Access Token

You need a Matrix access token with appropriate permissions:

1. **For Synapse**: Use an admin user's access token or server admin token
2. **For other servers**: Use an access token with sufficient privileges

To obtain an access token:
```bash
# Login and get access token
curl -XPOST -d '{"type":"m.login.password", "user":"admin", "password":"password"}' "https://matrix.example.com/_matrix/client/r0/login"
```

## Resource Examples

### User Management

```yaml
apiVersion: user.matrix.crossplane.io/v1alpha1
kind: User
metadata:
  name: alice
spec:
  forProvider:
    userID: "@alice:example.com"
    displayName: "Alice Wonderland"
    password: "secure_password"
    admin: false
    externalIDs:
      - medium: "email"
        address: "alice@example.com"
        validated: true
  providerConfigRef:
    name: default
```

### Room Creation

```yaml
apiVersion: room.matrix.crossplane.io/v1alpha1
kind: Room
metadata:
  name: team-room
spec:
  forProvider:
    name: "Team Discussion"
    topic: "Private team discussion room"
    preset: "private_chat"
    encryptionEnabled: true
    invite:
      - "@alice:example.com"
      - "@bob:example.com"
  providerConfigRef:
    name: default
```

### Space Organization

```yaml
apiVersion: space.matrix.crossplane.io/v1alpha1
kind: Space
metadata:
  name: company-space
spec:
  forProvider:
    name: "Company Space"
    topic: "Organization-wide space"
    children:
      - roomID: "!team-room:example.com"
        suggested: true
        order: "01"
  providerConfigRef:
    name: default
```

## Architecture

This provider is built using:

- **[mautrix-go](https://github.com/mautrix/go)**: Leading Go Matrix client library (MPL-2.0 licensed)
- **Crossplane Runtime**: For managed resource lifecycle and controller patterns
- **Matrix Client-Server API**: Standard Matrix protocol for basic operations
- **Matrix Admin API**: Server-specific admin endpoints (primarily Synapse)

The provider implements a dual-client architecture:
- **Standard Matrix Client**: For regular Matrix operations available to all users
- **Admin Client**: For administrative operations requiring elevated privileges

## Supported Matrix Servers

This provider is designed to work with any Matrix-compliant homeserver:

- ✅ **Synapse** (full support including admin API)
- ✅ **Dendrite** (basic support via standard Matrix API)
- ✅ **Conduit** (basic support via standard Matrix API)
- ✅ **matrix.org** (limited to standard API operations)

## Limitations

- Admin operations (user creation, room deletion) require Synapse or compatible admin API
- Some features may be server-specific (power level granularity, room settings)
- Federation and media operations are not currently supported
- No support for encrypted message history management

## Development

### Prerequisites

- Go 1.21+
- Docker
- kubectl
- A running Kubernetes cluster
- Access to a Matrix homeserver for testing

### Building

```bash
# Build the provider binary
make build

# Build the provider image
make docker-build

# Build the Crossplane package
make xpkg.build
```

### Testing

```bash
# Run unit tests
make test

# Run linting
make lint

# Run all checks
make ci
```

### Local Development

```bash
# Run the provider locally
make run
```

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Reporting Issues

Please report issues on GitHub: https://github.com/crossplane-contrib/provider-matrix/issues

### Development Status

This provider is in active development. Current status:

- ✅ Core CRD definitions and API types
- ✅ Matrix client implementation using mautrix-go
- ✅ Basic CRUD controllers for all resource types
- ✅ Example manifests and documentation
- 🚧 Unit and integration tests
- 🚧 CI/CD pipeline setup
- 🚧 End-to-end testing with real Matrix servers

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.

## Security

For security concerns, please see our [Security Policy](SECURITY.md).

## Support

- 📖 **Documentation**: https://docs.crossplane.io/
- 💬 **Community**: [Crossplane Slack](https://slack.crossplane.io/)
- 🐛 **Issues**: [GitHub Issues](https://github.com/crossplane-contrib/provider-matrix/issues)
- 🗣️ **Discussions**: [GitHub Discussions](https://github.com/crossplane-contrib/provider-matrix/discussions)