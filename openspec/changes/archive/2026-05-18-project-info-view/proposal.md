## Why

The TUI has no way to surface the project's `openspec/config.yaml` — the context block (tech stack, domain description) and writing rules are invisible once you're inside the tool. Users have to leave the TUI and open the file manually when they want to remind themselves of the project setup.

## What Changes

- Add a new `ModeViewingConfig` UI mode that renders the project config in a full-screen viewport
- Parse `openspec/config.yaml` into a `ProjectConfig` struct (context + rules; schema field ignored)
- Build a markdown representation of the config and render it via Glamour, consistent with how specs and archive artifacts are displayed
- Wire the `i` keybinding from `ModeIndex` and `ModeNormal` to enter the new mode
- Add `Esc`/`q` to exit back to the previous mode
- Update the help bar in `ModeViewingConfig` with navigation hints

## Capabilities

### New Capabilities

- `project-config-view`: Read-only view of `openspec/config.yaml` content (context + rules) rendered as markdown inside the existing TUI viewport pattern.

### Modified Capabilities

<!-- No existing spec requirements are changing. -->

## Impact

- `internal/openspec/loader.go`: new `LoadConfig()` function and `ProjectConfig` struct
- `internal/ui/model.go`: new `ModeViewingConfig` constant; `Model` gains a `projectConfig` field and a `prevMode` field to restore the correct mode on exit
- `internal/ui/view.go`: new `viewConfig()` rendering path; `renderHeader()` and `renderHelpBar()` extended for the new mode
- `internal/ui/update.go`: `i` keybinding wired in `ModeIndex` and `ModeNormal`; `Esc`/`q` exit logic for `ModeViewingConfig`
- No new dependencies; Glamour and Lipgloss already present

## Non-goals

- Editing `config.yaml` from within the TUI
- Showing the raw YAML or the `schema` field
- Hot-reloading the config when the file changes on disk
