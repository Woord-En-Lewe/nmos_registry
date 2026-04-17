# Task 9: IS-09 System Parameters API

**Status:** Pending

## Overview

Implement IS-09 NMOS System Parameters Specification to provide global configuration parameters that Nodes can fetch at startup. This allows configuring heartbeat intervals and other system-wide parameters without recompiling.

## Requirements

### 1. System API Endpoint
- Add `/x-nmos/system/v1.0/global` GET endpoint
- Return global configuration resource with system parameters

### 2. Global Configuration Resource
The `/global` response must include:

```json
{
  "id": "<system-unique-id>",
  "version": "<timestamp>",
  "is04": {
    "heartbeat_interval": 5000
  },
  "ptp": {
    "announce_receipt_timeout": <int>,
    "domain_number": <int>
  },
  "syslogv2": {
    "hostname": "<string>",
    "port": <int>
  },
  "syslog": {
    "hostname": "<string>",
    "port": <int>
  }
}
```

### 3. DNS-SD Advertisement
- Add `_nmos-system._tcp` service type for System API discovery
- Include required TXT records: `api_proto`, `api_ver`, `pri`

### 4. Integration with HeartbeatEngine
- HeartbeatEngine should accept configurable heartbeat interval
- Fetch from System API or use IS-04 default (5000ms) if unavailable

### 5. Error Handling
- Return proper IS-04/IS-09 error format: `{code, error, debug}`
- Handle DNS-SD discovery failures gracefully

## Missing Features (Non-Blocking)

### DNS-SD TXT Records Missing `api_proto`
Currently the mDNS announcer advertises Registration and Query APIs but is missing the required `api_proto` TXT record:

**Current:**
```go
TXTRecords: map[string]string{
    "api_ver": RegistryAPIVersion,
    "pri":     formatPriority(c.Priority),
}
```

**Required (IS-04):**
```go
TXTRecords: map[string]string{
    "api_proto": "http",  // or "https" if TLS configured
    "api_ver": RegistryAPIVersion,
    "pri":     formatPriority(c.Priority),
}
```

### CORS Implementation
IS-09 requires proper CORS headers. The current middleware has `Access-Control-Allow-Methods` missing `PUT, PATCH, HEAD`:

**Current:**
```go
w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
```

**Required (IS-09):**
```go
w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, HEAD, OPTIONS, DELETE")
```

## Files to Modify/Create

- `internal/infrastructure/transport/http/system_handlers.go` - New System API handlers
- `internal/infrastructure/transport/http/router.go` - Add System API routes
- `internal/infrastructure/mdns/announcer.go` - Add `_nmos-system._tcp` advertisement
- `internal/infrastructure/mdns/constants.go` - Add System service type constant
- `internal/infrastructure/mdns/models.go` - Add System service config
- `cmd/server/main.go` - Wire up System API and announcer
- `internal/registry/heartbeat_engine.go` - Accept configurable intervals

## References

- [IS-09 Specification](https://specs.amwa.tv/is-09)
- IS-04 Section on DNS-SD discovery
