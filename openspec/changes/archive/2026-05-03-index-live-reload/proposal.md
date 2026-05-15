## Why

The index view (`ModeIndex`) loads its data once on entry and never refreshes it. If a spec is created or archived while the index is open, the change is invisible until the user exits and re-enters the view. The polling tick already runs every 500 ms but unconditionally skips `ModeIndex`, so the fix is targeted and low-risk.

## What Changes

- `handleTick()` detects changes to active changes, archived changes, and project specs while in `ModeIndex`, and refreshes the index in-place when any list has changed
- Two new cheap loader functions: `ListArchiveNames()` and `ListSpecNames()` — directory listing only, no file reads — to make tick-level detection affordable
- The index cursor is preserved if the item under it still exists after a reload; reset to 0 only if it would go out of bounds

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `change-index`: new requirement for real-time refresh of the index while `ModeIndex` is active
- `openspec-loader`: two new functions (`ListArchiveNames`, `ListSpecNames`) for cheap directory-level change detection

## Impact

- `internal/openspec/loader.go`: add `ListArchiveNames()` and `ListSpecNames()`
- `internal/ui/model.go`: extend `handleTick()` to poll while in `ModeIndex`
