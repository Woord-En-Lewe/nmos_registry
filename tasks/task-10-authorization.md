# Task 10: IS-10 Authorization

**Status:** Future

## Overview

Implement IS-10 NMOS Authorization Specification to secure the Query API and Subscription endpoints. Controllers must authenticate via OAuth 2.0 and present valid JWT tokens with appropriate scopes.

## Why This Matters

- **Node Registration**: Unauthenticated (machine-to-machine via DNS-SD discovery)
- **Controller Access**: Requires IS-10 authentication (Query API, Subscriptions, WebSocket)

## Requirements

### 1. JWT Validation Middleware
- Validate JWT Bearer tokens on incoming requests
- Verify signature using JWKS endpoint
- Check token expiration (`exp` claim)
- Validate audience (`aud` claim)

### 2. Scope Definitions (per BCP-003-02)

| Endpoint | Required Scope | Claims |
|----------|---------------|--------|
| `GET /nodes`, `GET /nodes/{id}` | `x-nmos-query` | `read`: `["nodes/*"]` |
| `GET /devices`, `GET /devices/{id}` | `x-nmos-query` | `read`: `["devices/*"]` |
| `GET /sources/*` | `x-nmos-query` | `read`: `["sources/*"]` |
| `GET /flows/*` | `x-nmos-query` | `read`: `["flows/*"]` |
| `GET /senders/*` | `x-nmos-query` | `read`: `["senders/*"]` |
| `GET /receivers/*` | `x-nmos-query` | `read`: `["receivers/*"]` |
| `POST /subscriptions` | `x-nmos-query` | `write`: `["subscriptions"]` |
| `WS /subscriptions/{id}/ws` | `x-nmos-query` | `read`: matching `resource_path` |

### 3. JWT Claims Structure

JWTs MUST contain namespaced claims under `x-nmos-query`:

```json
{
  "iss": "https://auth.example.com",
  "sub": "client-456",
  "aud": "nmos-registry",
  "exp": 1234567890,
  "client_id": "client-456",
  "x-nmos-query": {
    "read": ["nodes/*", "devices/*", "senders/*", "receivers/*"],
    "write": ["subscriptions"]
  }
}
```

### 4. Endpoints to Validate

**Registration API** (Node → Registry):
- Not authenticated (nodes register without tokens)
- OR use `registration` scope with Client Credentials grant

**Query API** (Controller → Registry):
- Require `Authorization: Bearer <token>`
- Validate `x-nmos-query.read` claims for resource paths

**Subscriptions** (Controller → Registry):
- `POST /subscriptions` requires `x-nmos-query.write: ["subscriptions"]`
- WebSocket upgrade requires `x-nmos-query.read` matching `resource_path`

### 5. Error Responses

Invalid/missing token:
```
401 Unauthorized
WWW-Authenticate: Bearer realm="nmos", error="invalid_token"
```

Insufficient scope:
```
403 Forbidden
WWW-Authenticate: Bearer realm="nmos", error="insufficient_scope"
```

## Implementation Notes

### Files to Create

- `internal/infrastructure/auth/jwt_validator.go` - JWT validation logic
- `internal/infrastructure/auth/middleware.go` - Auth middleware for chi
- `internal/infrastructure/auth/jwks_cache.go` - JWKS key caching

### Files to Modify

- `internal/infrastructure/transport/http/router.go` - Add auth middleware to Query/Subscription routes
- `internal/infrastructure/transport/http/query_handlers.go` - Scope validation
- `internal/infrastructure/transport/http/registration_handlers.go` - Optional registration auth
- `internal/infrastructure/transport/websocket/manager.go` - WebSocket auth validation

### Dependencies

- JWT library (e.g., `github.com/golang-jwt/jwt/v5`)
- JWKS fetching/caching

### TLS Requirement

BCP-003-01 mandates HTTPS. IS-10 forbids plain HTTP in authorized mode.

## References

- [IS-10 Specification](https://specs.amwa.tv/is-10)
- [BCP-003-02 Authorization](https://specs.amwa.tv/bcp-003-02)
- [BCP-003-01 Secure Communication](https://specs.amwa.tv/bcp-003-01)
