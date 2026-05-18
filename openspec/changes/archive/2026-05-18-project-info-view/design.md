## Context

The TUI is structured around a `Mode` int and a single `Model` struct (Bubble Tea Elm architecture). Modes already in use: `ModeNormal`, `ModeIndex`, `ModeViewingArchive`, `ModeViewingSpec`. Each mode has its own rendering path in `view.go` and key handling in `update.go`. Config loading currently only reads `schema` and `created` fields from `.openspec.yaml` change metadata — `openspec/config.yaml` is not parsed at all.

The `ModeViewingSpec` pattern is the closest analogue: it enters a full-screen viewport, renders markdown via Glamour asynchronously, and exits with `Esc`/`q`. This change mirrors that pattern exactly.

## Goals / Non-Goals

**Goals:**
- Render `openspec/config.yaml` context and rules as readable markdown in a full-screen viewport
- Enter from both `ModeIndex` (keybinding `i`) and `ModeNormal` (keybinding `i`)
- Exit back to the mode the user came from (`Esc` or `q`)
- Reuse the existing async Glamour render pipeline (no new rendering machinery)

**Non-Goals:**
- Editing the config from the TUI
- Hot-reloading when `config.yaml` changes on disk (config is loaded once at startup)
- Displaying the raw `schema` field
- Parsing config.yaml rules beyond what's needed for display

## Decisions

### 1. Load config once at startup, not on demand

Config is loaded in `main.go` alongside `openspec.Load()`, stored in `Model.projectConfig`. No polling tick needed — the file rarely changes and there's no precedent for reloading it live.

**Alternative considered**: reload on every `i` keypress. Rejected — adds complexity with no real benefit given how static this file is.

### 2. Add `prevMode` field to restore the caller's mode on exit

When the user presses `i` from `ModeNormal`, they expect `Esc` to return to `ModeNormal`, not `ModeIndex`. A single `prevMode Mode` field on the model captures where we came from.

**Alternative considered**: always return to `ModeIndex`. Rejected — breaks the mental model if the user pressed `i` while browsing a change.

### 3. Convert config to markdown at render time, not at load time

`LoadConfig()` returns a `ProjectConfig` struct with typed fields. The markdown string is assembled in `view.go` when entering the mode, then passed to the existing `glamourRender` goroutine via `renderedConfigMsg`.

**Alternative considered**: store pre-built markdown in `ProjectConfig`. Rejected — rendering logic belongs in the UI layer.

### 4. Reuse existing `specRenderedMsg`-style async pattern

A new `renderedConfigMsg` carries the Glamour output back to the event loop, identical to how `specRenderedMsg` works. This keeps the render pipeline uniform.

### 5. Rules rendered as a markdown list per key

YAML rules are `map[string][]string` (keyed by artifact type). Each key becomes a `### <key>` heading with a bullet list. This produces clean, readable output without exposing YAML structure to the user.

## Risks / Trade-offs

- **config.yaml missing or malformed** → `LoadConfig()` returns a zero-value struct with empty fields; the view shows an empty viewport rather than crashing. Acceptable for a read-only display feature.
- **Very long context block** → viewport scroll handles it; no truncation needed.
- **`prevMode` adds state** → minimal, one field. Risk of stale state if a future refactor changes mode transitions, but the pattern is already used implicitly elsewhere.

## Open Questions

- None. Scope is fully defined.
