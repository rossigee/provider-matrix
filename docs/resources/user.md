# User

**API Version**: `user.matrix.m.crossplane.io/v1beta1`

Matrix users.

## Spec

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `forProvider.localpart` | string | no | Local username |
| `forProvider.displayName` | string | no | Display name |
| `forProvider.passwordSecretRef` | ref | no | Password secret selector (preferred over plaintext) |
| `forProvider.admin` | bool | no | Server admin flag |
| `forProvider.deactivated` | bool | no | Deactivate account |

## Example

See `examples/user/` for working manifests.

## Behavior

- **Create/Update/Delete**: Full managed lifecycle against the Matrix Client-Server API; observe reconciles room/user state.
