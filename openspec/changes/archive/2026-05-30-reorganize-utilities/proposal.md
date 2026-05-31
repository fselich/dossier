## Why

`extractRequirement` in `ui/viewport.go` and `configToMarkdown` in `ui/view.go` both process openspec-structured content — they belong in `openspec/loader.go` alongside the other content processing logic. Moving them clarifies ownership and reduces the UI package's responsibility.

## What Changes

- Move `extractRequirement` from `ui/viewport.go` to `openspec/loader.go`
- Move `configToMarkdown` from `ui/view.go` to `openspec/loader.go`
- Update imports in `ui/viewport.go` and `ui/view.go` to call the relocated functions

## Capabilities

### Modified Capabilities
- **openspec-loader**: Added requirement extraction (`extractRequirement`) and config formatting (`configToMarkdown`) — functions now live in the loader package
- **requirement-focus-view**: Imports updated to call `extractRequirement` from the openspec package instead of the local ui package

## Impact

- `internal/ui/viewport.go`: `extractRequirement` removed (callers updated to use openspec package)
- `internal/ui/view.go`: `configToMarkdown` removed (callers updated to use openspec package)
- `internal/openspec/loader.go`: two functions added
- No user-facing behavior change
