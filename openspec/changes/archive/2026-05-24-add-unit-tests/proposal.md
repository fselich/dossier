## Why

Only 2 test files (~174 lines) exist for ~1300 lines of production code across `internal/openspec` and `internal/ui`. Refactoring without test coverage is high-risk — any change can silently break functionality. The P0 items (especially `inject-root-path`) require tests as a safety net before structural changes.

## What Changes

- Add `internal/openspec/loader_test.go`: Tests for all public loader functions (`Load`, `LoadFromPath`, `LoadProjectSpecs`, `ParseTasks`, `ToggleTask`, `FindCursorByText`, `ListChangeNames`, `ListArchiveChanges`, `ReloadChange`). Uses `t.TempDir()` for filesystem isolation.
- Add `internal/ui/view_test.go`: Tests for rendering helpers (`extractRequirement`, `renderTasksContent`, `renderIndexContent`, `buildIndexItems`, `firstAvailableTab`).
- Achieve coverage > 60% on `internal/openspec` and > 40% on `internal/ui`.

## Capabilities

### New Capabilities

_None. This change adds test coverage only; it does not introduce new functional capabilities._

### Modified Capabilities

_None. Existing behavior is not changed; tests verify current behavior._

## Impact

- New files: `internal/openspec/loader_test.go`, `internal/ui/view_test.go`
- No changes to production code
- `go test ./...` becomes the primary verification step before any refactoring
- Testable functions can be validated with `go test -race -count=1 ./...`
