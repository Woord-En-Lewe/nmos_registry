# IS-04 Compliant NMOS Registry Implementation Plan

## Overview
This plan ensures the NMOS Registry fully complies with AMWA IS-04 v1.3 specification for Discovery and Registration.

---

## Phase 1: Error Handling Fix (Immediate Bug)

### Problem
Currently, validation errors (e.g., "parent node not found") return HTTP 500 (Internal Server Error). Per IS-04 spec, these should return 400 (Bad Request).

### IS-04 Requirement (Error Conditions section)
> A registry MAY issue a `400` (Bad Request) code for:
> - The request body does not meet the JSON schema for that resource type
> - The `id` included in the request has already been used by another resource type held in the registry
> - The `version` included in the request is earlier than the matching resource already held in the registry
> - A parent resource ID has been modified (for example the `node_id` in a Device registration is modified during an update)
> - The parent resource referred to either doesn't exist in the registry or the ID matches the wrong type of resource

### Files to Create/Modify

#### 1. `internal/registry/errors.go` (NEW)
Define error types for proper HTTP status code handling:
- `ValidationError` - for 400 Bad Request cases
- `ConflictError` - for 409 Conflict cases

#### 2. `internal/registry/resource_manager.go` (MODIFY)
Wrap parent-not-found errors with `ValidationError`:
- `RegisterDevice` line 40 - parent node not found
- `RegisterSource` line 78 - parent device not found
- `RegisterFlow` lines 98, 103 - parent source/device not found
- `RegisterSender` lines 117, 123 - parent device/flow not found
- `RegisterReceiver` line 138 - parent device not found

#### 3. `internal/infrastructure/transport/http/registration_handlers.go` (MODIFY)
Update error handling at line ~131 to check error types:
- Return 400 for `ValidationError`
- Return 409 for `ConflictError`
- Return 500 only for unhandled errors

#### 4. `internal/infrastructure/transport/http/registration_handlers_test.go` (MODIFY)
Add tests for:
- 400 returned when parent doesn't exist
- 400 returned for invalid resource data

---

## Phase 2: IS-04 Required Validations

### 2.1 Version Comparison
IS-04 version format is `<seconds>:<nanoseconds>` (TAI timestamp).

**Files:**
- `internal/registry/version.go` (NEW) - `CompareVersions(v1, v2 string) int`

### 2.2 API Version Tracking Per Resource
IS-04 requires tracking which API version was used when registering a resource, and returning 409 when there's a mismatch.

**Files:**
- `internal/infrastructure/persistence/db/schema.sql` - add `api_version` column
- `internal/infrastructure/persistence/db/queries.sql` - update queries
- Regenerate `internal/infrastructure/persistence/db/models.go` and `queries.sql.go`
- `internal/registry/models.go` - add `APIVersion` field to all resource structs
- `internal/registry/resource_manager.go` - check API version on registration
- `internal/infrastructure/persistence/repository.go` - handle APIVersion in CRUD

**409 Conflict** when registering at same resource ID but different API version than existing.

### 2.3 ID Collision Detection
Before registering, check if ID exists as a different resource type.

**Files:**
- `internal/registry/resource_manager.go` - add ID collision check in all Register methods

### 2.4 Parent ID Modification Detection
On update, reject if parent ID has changed.

**Files:**
- `internal/registry/resource_manager.go` - add parent ID change check in update paths

### 2.5 JSON Schema Validation
IS-04 strongly recommends JSON schema validation.

**Files:**
- `internal/infrastructure/transport/http/schemas/*.json` (NEW) - IS-04 JSON schemas
- `internal/infrastructure/transport/http/registration_handlers.go` - validate against schema before registration

---

## Phase 3: Database Migration

### SQL Schema Updates
Add `api_version` column to all resource tables:
- nodes
- devices
- sources
- flows
- senders
- receivers

---

## Phase 4: Testing

Add comprehensive tests for:
- Version comparison edge cases (equal, earlier, later, invalid format)
- API version conflict (409)
- ID collision across resource types (400)
- Parent ID modification rejection (400)
- JSON schema validation failures (400)
- Parent resource not found (400)
- Successful registration (201/200)

---

## Summary of File Changes

| File | Action |
|------|--------|
| `internal/registry/errors.go` | Create |
| `internal/registry/version.go` | Create |
| `internal/registry/resource_manager.go` | Modify |
| `internal/registry/models.go` | Modify |
| `internal/infrastructure/transport/http/registration_handlers.go` | Modify |
| `internal/infrastructure/transport/http/schemas/*.json` | Create |
| `internal/infrastructure/persistence/db/schema.sql` | Modify |
| `internal/infrastructure/persistence/db/queries.sql` | Modify |
| `internal/infrastructure/persistence/db/models.go` | Regenerate |
| `internal/infrastructure/persistence/db/queries.sql.go` | Regenerate |
| `internal/infrastructure/persistence/repository.go` | Modify |
| `internal/infrastructure/transport/http/registration_handlers_test.go` | Modify |
| `code_map.md` | Update |

---

## Status

- [x] Phase 1: Error Handling Fix
- [ ] Phase 2: IS-04 Required Validations
- [ ] Phase 3: Database Migration
- [ ] Phase 4: Testing

---

## Phase 1 Completed

### Changes Made

1. **Created `internal/registry/errors.go`**:
   - Added `ErrResourceNotFound` sentinel error
   - Added `ValidationError` struct for 400 Bad Request cases
   - Added `ConflictError` struct for 409 Conflict cases
   - Added constructor functions: `NewValidationError`, `NewValidationErrorf`, `NewConflictError`, `NewConflictErrorf`
   - Added wrapper functions: `WrapValidationError`, `WrapConflictError`

2. **Modified `internal/registry/resource_manager.go`**:
   - All `Register*` methods now wrap `ErrResourceNotFound` with `ValidationError`
   - Parent validation errors now return `ValidationError` instead of `fmt.Errorf`

3. **Modified `internal/infrastructure/transport/http/registration_handlers.go`**:
   - Added `errors` import
   - Changed error handling to check for `ValidationError` and return 400
   - Only returns 500 for non-validation errors

4. **Modified `internal/registry/fake_repository_test.go`**:
   - Updated all `Get*` methods to return `ErrResourceNotFound` instead of custom fmt errors

### Result
- Parent not found errors now return HTTP 400 Bad Request (was 500)
- All tests pass