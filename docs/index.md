# Provider Matrix Documentation

A Crossplane v2 provider for managing Matrix chat resources. All resources are namespaced (`*.matrix.m.crossplane.io/v1beta1`) with full multi-tenancy support.

## Resource Documentation

### Messaging

| Resource | API Group | Description |
|----------|-----------|-------------|
| [Room](resources/room.md) | `room.matrix.m.crossplane.io/v1beta1` | Matrix chat rooms |
| [RoomAlias](resources/roomalias.md) | `roomalias.matrix.m.crossplane.io/v1beta1` | Room aliases |
| [PowerLevel](resources/powerlevel.md) | `powerlevel.matrix.m.crossplane.io/v1beta1` | Room permission levels |

### User Management

| Resource | API Group | Description |
|----------|-----------|-------------|
| [User](resources/user.md) | `user.matrix.m.crossplane.io/v1beta1` | Matrix users |

### Administration

| Resource | API Group | Description |
|----------|-----------|-------------|
| [Space](resources/space.md) | `space.matrix.m.crossplane.io/v1beta1` | Matrix spaces |

## API Coverage Gaps

Matrix Client-Server / Admin API surface not yet modeled: devices and E2E key backup management, server ACLs, room member invites/kicks/bans as resources, server notices, push rules, media retention config, federation allowlists, and Synapse admin endpoints (user media, experimental features, background updates).
