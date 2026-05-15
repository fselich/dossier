## Why

The TUI loads artifacts once at startup and does not update them. When an external agent (Claude Code, `opsx:apply`) modifies `tasks.md` while the TUI is open, the user sees a stale state until they restart the tool.

## What Changes

- The TUI polls the active change's files every 2 seconds
- `tasks.md`: if the content has changed on disk, it is re-parsed and the view is refreshed; the cursor is restored by task text, not by index
- `proposal.md`, `design.md`, `specs/*/spec.md`: if the content has changed on disk, its entry in the render cache is invalidated; the next access to that tab re-renders with the updated content
- The list of changes (which changes exist) is not reloaded in real time

## Capabilities

### New Capabilities

### Modified Capabilities

- `openspec-loader`: add function to re-read the artifact content of an already-loaded change
- `tui-viewer`: add polling cycle with `tea.Tick` and handling of the change-detected message
- `tasks-toggle`: add cursor restoration by task text after reload

## Impact

- No new dependencies
- Only affects `internal/openspec/loader.go` and `internal/ui/model.go`
- No changes to file format or openspec structure
