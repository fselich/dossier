## Why

The empty state (no active changes) shows `q: salir` — the only Spanish string left in the UI. It also gives no hint that the index view is accessible from here, leaving the user with no visible path to navigate to archived changes or specs.

## What Changes

- Replace `q: salir` with `a: index  q: quit` in `emptyView()`

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

_(none — single string change, no spec-level behavior change)_

## Impact

- `internal/ui/model.go`: one string in `emptyView()`
