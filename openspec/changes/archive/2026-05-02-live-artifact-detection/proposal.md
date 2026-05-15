## Why

The TUI starts with the disk state at that instant and does not update tab availability or the change list as new files are created. If `/opsx:propose` generates `proposal.md` while the TUI is open, the tab stays disabled until restart. The progress counter in the tab bar also does not reflect the TUI's own toggles until the next polling cycle (up to 2 s of visible lag).

## What Changes

- The polling interval changes from 2 s to 500 ms to detect changes in less than half a second
- The tick handler detects artifact presence changes (absent → present and vice versa), not just content changes
- If the TUI started with no active changes, the tick attempts to reload the change list from disk
- Task toggles within the TUI update in-memory state immediately, without waiting for the next tick

## Capabilities

### New Capabilities

### Modified Capabilities

- `tui-viewer`: faster polling, artifact presence detection, change list reload from empty state

## Impact

- Only affects `internal/ui/model.go` and `internal/openspec/loader.go`
- No new dependencies
- The biggest change is the interval: from `2*time.Second` to `500*time.Millisecond`
