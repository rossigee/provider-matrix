# Room

**API Version**: `room.matrix.m.crossplane.io/v1beta1`

Matrix chat rooms.

## Spec

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `forProvider.name` | string | no | Room name |
| `forProvider.topic` | string | no | Room topic |
| `forProvider.alias` | string | no | Canonical alias |
| `forProvider.preset` | string | no | Join preset |
| `forProvider.visibility` | string | no | public/private |
| `forProvider.roomVersion` | string | no | Room version |
| `forProvider.initialState` | array | no | Initial state events |

## Example

See `examples/room/` for working manifests.

## Behavior

- **Create/Update/Delete**: Full managed lifecycle against the Matrix Client-Server API; observe reconciles room/user state.
