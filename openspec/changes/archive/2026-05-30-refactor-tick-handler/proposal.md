## Why

`handleTick()` in `index.go` is a 143-line god function that mixes three distinct responsibilities: disk polling for the index view, change-list detection in normal mode, and per-artifact content reloading. This makes the code hard to reason about and modify. Extracting methods clarifies the control flow without any behavior change.

## What Changes

- Extract `pollIndexMode()` from `handleTick()` (lines 18-74): handles index-mode disk polling
- Extract `pollNormalModeChanges()` from `handleTick()` (lines 76-104): handles change-list detection in normal mode
- Extract `pollNormalModeContent()` from `handleTick()` (lines 106-154): handles reloading individual change content
- `handleTick()` becomes a 3-line dispatcher delegating to the three new methods

## Capabilities

### New Capabilities

*(none — pure internal refactor)*

### Modified Capabilities

- `change-index`: Internal restructuring of the tick handler; no requirement-level behavior changes

## Non-goals

- Changing any polling logic, timing, or behavior
- Adding or removing features from any mode
- Extracting functions beyond the three specified methods

## Impact

- `internal/ui/index.go`: `handleTick()` split into 4 methods (1 dispatcher + 3 extracted)
- No API changes, no new dependencies
