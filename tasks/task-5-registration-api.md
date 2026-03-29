# Task 5: Infrastructure - Registration API Handlers

## Summary
Implement the HTTP handlers for the Registration API using `go-chi/chi`. These handlers will translate incoming HTTP requests (POST, DELETE) into calls to the `ResourceManager` and `HeartbeatEngine` in the BM.

## Priority
High

## Referenced Specifications
- [IS-04 § Registration API](https://specs.amwa.tv/is-04/v1.3/docs/Registration_API.html)
- [IS-04 § Health Endpoints](https://specs.amwa.tv/is-04/v1.3/docs/Health_Endpoints.html)

## Referenced Files
- `internal/infrastructure/transport/http/registration_handlers.go`
- `internal/infrastructure/transport/http/router.go`
