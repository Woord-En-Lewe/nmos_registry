# Task 8: Verification & Sociable Unit Testing

## Summary
Develop a suite of sociable unit tests to verify the integrity of the Business Module. Use an `InMemoryRepository` fake to test the interaction between registration logic, heartbeating, and subscriptions without needing a live HTTP server or SQLite instance.

## Priority
High

## Referenced Specifications
- [NMOS IS-04 Specification](https://specs.amwa.tv/is-04)

## Referenced Files
- `internal/registry/registry_test.go`
- `internal/registry/fakes/fake_repository.go`
