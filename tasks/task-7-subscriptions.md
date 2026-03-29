# Task 7: Subscription Logic & WebSocket Transport

## Summary
Implement the `SubscriptionEngine` in the BM and the corresponding WebSocket transport in the IM. This enables real-time notifications for resource changes as specified in IS-04.

## Priority
Low

## Referenced Specifications
- [IS-04 § Query API - Subscriptions](https://specs.amwa.tv/is-04/v1.3/docs/Query_API_-_Subscriptions.html)
- [IS-04 § Query API - WebSockets](https://specs.amwa.tv/is-04/v1.3/docs/Query_API_-_WebSockets.html)

## Referenced Files
- `internal/registry/subscription_engine.go`
- `internal/infrastructure/transport/websocket/manager.go`
- `internal/registry/ports.go`
