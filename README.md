# NMOS Registry

An implementation of the [AMWA IS-04 specification](https://specs.amwa.tv/is-04/) for discovery and registration of networked media resources.

## Overview

NMOS Registry provides centralized discovery and registration for media devices and resources following the NMOS (Networked Media Open Specifications) standards. It implements both the Registration and Query APIs defined in IS-04 v1.3.

## Features

- **Registration API** - Resource lifecycle management (Nodes, Devices, Sources, Flows, Senders, Receivers)
- **Query API** - Read-only access to registered resources with RQL filtering and pagination
- **Heartbeating** - Automatic cleanup of expired nodes
- **WebSocket Subscriptions** - Real-time resource change notifications
- **mDNS/DNS-SD** - Automatic service advertisement for discovery

## Quick Start

```bash
# Build
go build -o nmos_registry .

# Run (defaults to port 8080)
./nmos_registry

# Run on custom port
PORT=9000 ./nmos_registry
```

## API Endpoints

### Registration API (`/x-nmos/registration/v1.3/`)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/resource` | Register/update a resource |
| DELETE | `/resource/{type}/{id}` | Unregister a resource |
| POST | `/health/nodes/{nodeId}` | Heartbeat for a node |

### Query API (`/x-nmos/query/v1.3/`)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/nodes` | List all nodes |
| GET | `/nodes/{id}` | Get node by ID |
| GET | `/devices` | List all devices |
| GET | `/devices/{id}` | Get device by ID |
| GET | `/sources` | List all sources |
| GET | `/sources/{id}` | Get source by ID |
| GET | `/flows` | List all flows |
| GET | `/flows/{id}` | Get flow by ID |
| GET | `/senders` | List all senders |
| GET | `/senders/{id}` | Get sender by ID |
| GET | `/receivers` | List all receivers |
| GET | `/receivers/{id}` | Get receiver by ID |
| GET | `/subscriptions` | List active subscriptions |
| POST | `/subscriptions` | Create a subscription |
| GET | `/subscriptions/{id}` | Get subscription details |
| DELETE | `/subscriptions/{id}` | Delete a subscription |
| GET | `/subscriptions/{id}/ws` | WebSocket for live updates |

## DNS-SD Records

For wide-area DNS-SD discovery (not mDNS), the following records are required per IS-04 v1.3.

### Service Types

| Service | DNS Service Type | Description |
|---------|------------------|-------------|
| Registration API | `_nmos-register._tcp` | v1.3 Registration API |
| Query API | `_nmos-query._tcp` | Query API |

### DNS Record Hierarchy

```
_services._dns-sd._udp.<domain>    PTR    _nmos-register._tcp.<domain>
_services._dns-sd._udp.<domain>    PTR    _nmos-query._tcp.<domain>

_nmos-register._tcp.<domain>       PTR    <instance>._nmos-register._tcp.<domain>
_nmos-query._tcp.<domain>          PTR    <instance>._nmos-query._tcp.<domain>

<instance>._nmos-register._tcp.<domain>    SRV    0 0 8080 <host>.<domain>
<instance>._nmos-register._tcp.<domain>    TXT    "api_ver=v1.3" "pri=0100"

<instance>._nmos-query._tcp.<domain>    SRV    0 0 8080 <host>.<domain>
<instance>._nmos-query._tcp.<domain>    TXT    "api_ver=v1.3" "pri=0100"
```

### TXT Record Attributes

| Key | Value | Description |
|-----|-------|-------------|
| `api_ver` | `v1.3` | API version |
| `pri` | `0100` | Priority (0000-9999, lower = higher priority) |

### Example BIND9 Zone Configuration

```
; Service discovery enumeration
_services._dns-sd._udp    PTR     @
b._dns-sd._udp           PTR     @
lb._dns-sd._udp          PTR     @

; NMOS service type advertisements
_services._dns-sd._udp    PTR     _nmos-register._tcp
_services._dns-sd._udp    PTR     _nmos-query._tcp

; PTR records for service instances
_nmos-register._tcp    PTR     nmos-reg._nmos-register._tcp.example.com.
_nmos-query._tcp       PTR     nmos-reg._nmos-query._tcp.example.com.

; SRV and TXT records for Registration API
nmos-reg._nmos-register._tcp    SRV    0 0 8080  nmos-reg.example.com.
nmos-reg._nmos-register._tcp    TXT    "api_ver=v1.3" "pri=0100"

; SRV and TXT records for Query API
nmos-reg._nmos-query._tcp       SRV    0 0 8080  nmos-reg.example.com.
nmos-reg._nmos-query._tcp       TXT    "api_ver=v1.3" "pri=0100"

; Host A record
nmos-reg.example.com.    A    192.168.1.100
```

### MikroTik Configuration

MikroTik RouterOS does not support PTR records. For full DNS-SD compliance, use a complete DNS server (BIND, CoreDNS). SRV-only configuration may work with some clients:

```routeros
/ip dns static
add name=nmos-reg.example.com type=A address=192.168.1.100
add name="_nmos-register._tcp.example.com" type=SRV srv-port=8080 target=nmos-reg.example.com
add name="_nmos-query._tcp.example.com" type=SRV srv-port=8080 target=nmos-reg.example.com
add name="_nmos-register._tcp.example.com" type=TXT data="api_ver=v1.3" "pri=0100"
add name="_nmos-query._tcp.example.com" type=TXT data="api_ver=v1.3" "pri=0100"
/ip dns set allow-remote-requests=yes
```

## Architecture

```
internal/
├── registry/                    # Business logic
│   ├── models.go               # Domain models
│   ├── ports.go                # Interface definitions
│   ├── resource_manager.go     # Resource lifecycle
│   ├── heartbeat_engine.go     # Health monitoring
│   └── subscription_engine.go  # Query subscriptions
└── infrastructure/
    ├── persistence/            # SQLite storage
    ├── transport/              # HTTP/WebSocket
    └── mdns/                   # mDNS/DNS-SD discovery
```

## Requirements

- Go 1.25.5+
- SQLite (via modernc.org/sqlite)
