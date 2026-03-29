# Task 4: Business Module - Heartbeat & Garbage Collection

## Summary
Implement the `HeartbeatEngine` within the BM to manage Node health timestamps. This includes a background task for garbage collection that prunes resources from the registry when a Node fails to send a heartbeat within the required interval.

## Priority
Medium

## Referenced Specifications
- [IS-04 § Heartbeating](https://specs.amwa.tv/is-04/v1.3/docs/Heartbeating.html)

## Referenced Files
- `internal/registry/heartbeat_engine.go`
- `internal/registry/ports.go`
