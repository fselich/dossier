## Why

The TUI is read-only: the user can view artifacts but cannot edit them without leaving the application. Opening the system editor directly from the TUI eliminates the context switch and makes the editing flow smoother.

## What Changes

- When pressing `e` on any tab with an available artifact, the TUI suspends its event loop, launches `$EDITOR` with the path of the active artifact file and resumes the TUI when the editor is closed
- If `$EDITOR` is not defined, `vi` is attempted as a fallback
- The tab bar updates its content automatically on return (the existing 500 ms polling detects it, or an immediate reload is forced after the editor returns)
- The help bar shows the shortcut `e: edit` when an artifact is available

## Capabilities

### New Capabilities

- `editor-launch`: Ability to suspend the TUI, open a file in `$EDITOR` and resume

### Modified Capabilities

- `tui-viewer`: Adds the `e` keybinding and updates the help bar

## Impact

- Only affects `internal/ui/model.go`
- No new dependencies (uses `os/exec` from the stdlib)
- Requires `tea.ExecProcess` from BubbleTea to correctly hand the terminal over to the editor
