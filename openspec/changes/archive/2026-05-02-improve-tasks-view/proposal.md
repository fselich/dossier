## Why

The TUI tasks view has three visual deficiencies: items that are not under the cursor do not word-wrap (long texts are cut off), inline markdown in tasks (bold, code) is not rendered, and there is no global progress bar showing the overall progress of the change.

## What Changes

- All task items word-wrap to the terminal width, not just the selected item
- Each task's text goes through a mini inline-renderer that converts `**bold**` and `` `code` `` to ANSI styles
- A global progress bar is added at the top of the tasks view, showing the total completed tasks over the change total

## Capabilities

### New Capabilities

### Modified Capabilities

- `tui-viewer`: word wrap on all task items; global progress bar at the top of the view
- `tasks-toggle`: inline markdown rendering (bold, code) in each task's text

## Impact

- Only affects `internal/ui/model.go`
- No new dependencies
- No changes to the `tasks.md` format or toggle logic
