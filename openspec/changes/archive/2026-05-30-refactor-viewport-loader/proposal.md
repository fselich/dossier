## Why

`loadViewport()` in viewport.go has grown to 155 lines with 5 distinct mode-specific code paths intertwined in a single function. This makes the function hard to read, test, and modify. Splitting it into mode-specific methods improves maintainability without changing behavior.

## What Changes

- Extract 5 mode-specific methods from `loadViewport()`:
  - `loadViewportForIndex()` — handles `ModeIndex` (sync refresh, no glamour)
  - `loadViewportForConfig()` — handles `ModeViewingConfig` (`configToMarkdown` + glamour)
  - `loadViewportForSpec()` — handles `ModeViewingSpec` (requirement extraction + glamour)
  - `loadViewportForTasks()` — handles `TabTasks` in `ModeNormal` (sync task rendering)
  - `loadViewportForArtifact()` — handles proposal/design/specs tabs (cache check + glamour)
- Update `loadViewport()` to delegate to the appropriate method
- Keep exact same logic and behavior — pure extraction refactor

## Capabilities

### New Capabilities

*(none)*

### Modified Capabilities

- `mouse-navigation`: internal restructuring only — `loadViewport()` is called from mouse handlers. No spec-level requirement changes; behavior is identical.

## Impact

- **Affected code**: `internal/ui/viewport.go` (primary), `internal/ui/update.go` (calls `loadViewport()`), `internal/ui/mouse.go` (calls `loadViewport()`)
- **Tests**: existing tests in `view_test.go` should pass without modification
- **No API changes**, no new dependencies
