## Why

The current TUI has no visual separation between interface zones: the header, tab bar, content and help bar read as a continuous block. Adding borders or separator lines gives visual structure and makes it easier to perceive the zones at a glance.

## What Changes

- A horizontal separator line is added between the tab bar and the content area
- A horizontal separator line is added between the content area (or the global progress bar) and the help bar
- The header and tab bar are enclosed in an ASCII box (top border + minimal side border), and the help bar forms the bottom zone of the box

## Capabilities

### New Capabilities

### Modified Capabilities

- `tui-viewer`: adds visual separators between the TUI layout zones

## Impact

- Only affects `internal/ui/model.go` and `internal/ui/styles.go`
- No new dependencies
- No changes to business logic or artifact parsing
