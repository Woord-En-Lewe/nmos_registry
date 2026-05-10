# NMOS Registry

Go 1.26.1 single-module project implementing AMWA IS-04 v1.3 (NMOS device discovery/registration).

## Build & Run

```bash
go build -o nmos_registry .
./nmos_registry                     # port 8080
PORT=9000 ./nmos_registry           # custom port
```

## Test

```bash
go test ./...
```

Single package: no sub-modules, no package boundaries to worry about.

## Architecture

- Entry point: `cmd/server/main.go`
- Business logic: `internal/registry/`
- HTTP/WS transport: `internal/infrastructure/transport/http/`
- Persistence: `internal/infrastructure/persistence/` (SQLite in-memory)
- mDNS: `internal/infrastructure/mdns/`
- JSON schema validator: `internal/infrastructure/schema/validator.go`

## Schema Files

JSON schemas are loaded from `./schemas/` (relative to working directory, not relative to the binary). Schema compilation errors are silently skipped at startup — if a schema file is missing or invalid, it simply won't be registered for validation.

## Database

In-memory SQLite with shared cache: `file::memory:?cache=shared&_pragma=foreign_keys(1)`. Every test run starts fresh.

## Dependencies

- `github.com/go-chi/chi/v5` — HTTP routing
- `github.com/mattn/go-sqlite3` — SQLite driver
- `github.com/gorilla/websocket` — WebSocket support
- `github.com/grandcat/zeroconf` — mDNS/DNS-SD
- `github.com/santhosh-tekuri/jsonschema/v6` — JSON schema validation

## NMOS Specification Lookup

When working with NMOS specs (IS-04, IS-05, etc.), always verify against authoritative sources in this order:

1. **Repo documents** — Check `./schemas/` and any spec-related files in the repo first
2. **MCP servers** — Use available MCP servers to fetch specs if configured
3. **Web search** — Fall back to web search only if the above are unavailable

Do not assume general training data is current or accurate for NMOS specifications.