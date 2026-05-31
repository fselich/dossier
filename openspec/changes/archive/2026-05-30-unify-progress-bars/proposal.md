## Why

Three separate progress bar rendering functions (`renderTabBar`, `renderActiveItem`, `progressBar`) all implement the same logic — building a `[███···]` style bar from `done`/`total` with color switch at completion. Two view functions (`viewIndexContent`/`viewConfigContent`) are also duplicates. Consolidating eliminates copy-paste drift and reduces maintenance surface.

## What Changes

- Extract a single `renderProgressBar(done, total, width int) string` function in `tasks.go`
- Replace `renderTabBar()` in `view.go` to call `renderProgressBar`
- Replace `renderActiveItem()` in `index.go` to call `renderProgressBar`
- Replace existing `progressBar()` in `tasks.go` to call `renderProgressBar`
- Merge `viewIndexContent` and `viewConfigContent` duplicates in `view.go` into one function

## Capabilities

### Modified Capabilities
- **progress-bar-complete-style**: Unified function signature replaces three ad-hoc implementations with a single parameterized `renderProgressBar(done, total, width int) string`

## Impact

- `internal/ui/view.go`: `renderTabBar`, `viewIndexContent`, `viewConfigContent` modified/removed
- `internal/ui/index.go`: `renderActiveItem` updated
- `internal/ui/tasks.go`: `progressBar` replaced, new `renderProgressBar` added
- No user-facing behavior change — rendering output stays identical
