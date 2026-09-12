# RoomAlias

**API Version**: `roomalias.matrix.m.crossplane.io/v1beta1`

Human-readable room aliases.

## Spec

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `forProvider.alias` | string | yes | Alias (e.g. #general:server) |
| `forProvider.roomID` | string | yes | Target room ID |
| `forProvider.setAsCanonical` | bool | no | Set as canonical alias |
| `forProvider.altAliases` | array | no | Alternative aliases |

## Example

See `examples/roomalias/` for working manifests.

## Behavior

- **Create/Update/Delete**: Full managed lifecycle against the Matrix Client-Server API; observe reconciles room/user state.
