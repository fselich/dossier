## Why

The current archive selector is a modal overlay that is poorly integrated visually into the TUI. It also does not allow the user to navigate fluidly between active and archived changes from a single point. It is replaced by a full-screen index view that unifies both worlds.

## What Changes

- **BREAKING**: `ModeArchivePicker` and its modal overlay are removed
- `ModeIndex` is introduced: a full-screen view with two sections (active and archived)
- Active changes are shown with a progress bar `[█░] N/M`
- Archived changes are shown as a list with a date
- Navigation is simplified: `a`/`Esc` from any change view goes to the index; `Esc` from the index quits the app; `q` quits from anywhere

## Capabilities

### New Capabilities

- `change-index`: Full-screen index view that lists active changes (with progress) and archived changes (with date), with `j`/`k` navigation and `Enter` to select

### Modified Capabilities

- `tui-viewer`: Changes the mode model and navigation transitions — `ModeArchivePicker` is removed, `ModeIndex` is added, and the `a` and `Esc` shortcuts are redefined
- `archive-viewer`: Access to archived changes now always goes through `ModeIndex` instead of the modal

## Impact

- `internal/ui/model.go`: removal of `ModeArchivePicker`, addition of `ModeIndex`, redefinition of state transitions
- `internal/ui/styles.go`: modal styles can be removed; new styles for the index view
- No changes to `internal/openspec/` — `ListArchiveChanges()` already exists and serves the purpose
