# Space

**API Version**: `space.matrix.m.crossplane.io/v1beta1`

Matrix spaces (room hierarchies).

## Spec

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `forProvider.name` | string | no | Space name |
| `forProvider.topic` | string | no | Topic |
| `forProvider.alias` | string | no | Canonical alias |
| `forProvider.invite` | array | no | User IDs to invite |

## Example

See `examples/space/` for working manifests.

## Behavior

- **Create/Update/Delete**: Full managed lifecycle against the Matrix Client-Server API; observe reconciles room/user state.
