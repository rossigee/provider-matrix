# PowerLevel

**API Version**: `powerlevel.matrix.m.crossplane.io/v1beta1`

Room power-level (permission) configuration.

## Spec

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `forProvider.roomID` | string | yes | Target room ID |
| `forProvider.users` | map | no | Per-user levels |
| `forProvider.events` | map | no | Per-event levels |
| `forProvider.ban` | int | no | Ban level |
| `forProvider.kick` | int | no | Kick level |

## Example

See `examples/powerlevel/` for working manifests.

## Behavior

- **Create/Update/Delete**: Full managed lifecycle against the Matrix Client-Server API; observe reconciles room/user state.
