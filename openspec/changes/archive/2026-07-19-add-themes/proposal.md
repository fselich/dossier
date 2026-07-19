## Why

Dossier has three independent color/styling systems — UI widgets (22 hardcoded ANSI color vars in `styles.go`), markdown rendering (glamour, hardcoded to `"dark"`), and code syntax highlighting (chroma, hardcoded to `"monokai"`) — with no way to switch between them or select a coherent theme. A `Theme` struct exists but only carries a single `ViewBg` field. Users need at minimum a light/dark toggle, and the groundwork to eventually theme all three systems from a single name.

## What Changes

- Add `--theme <name>` CLI flag that selects a named theme (default: `dark`)
- Define `Theme` struct with `Name`, `GlamourStyle`, `ChromaStyle`, and a future `Colors` field for UI palette
- Define three built-in themes as a `map[string]Theme`: `dark`, `light`, `dracula`
- Thread the selected theme's glamour style into `ensureRenderer()` replacing the hardcoded `"dark"`
- Thread the selected theme's chroma style into `highlightLine()` replacing the hardcoded `"monokai"`
- Set `tea.View.BackgroundColor` from the theme (already wired, needs a sensible default per theme)
- Use Go's `flag` package for argument parsing (incidental cleanup: replaces manual `os.Args` inspection)

## Capabilities

### New Capabilities
- `theme-system`: centralized theme selection via `--theme` flag, with built-in themes mapping glamour styles, chroma styles, and view background colors; extensible for future UI palette migration

### Modified Capabilities
- `view-background`: the `ViewBg` field on `Theme` gains concrete values per built-in theme instead of being always nil
- `tui-viewer`: the `View()` method already reads `m.theme.ViewBg`; this change populates it meaningfully for each theme

## Impact

- `cmd/dossier/main.go`: argument parsing (flag package), theme lookup, pass Theme to `ui.New`
- `internal/ui/model.go`: `Theme` struct gains `Name`, `GlamourStyle`, `ChromaStyle` fields; `New()` accepts `Theme`
- `internal/ui/viewport.go`: `ensureRenderer()` reads glamour style from `m.theme` instead of hardcoded string
- `internal/ui/gitdiff.go`: `highlightLine()` and `renderDiff()` accept chroma style name parameter; `initChromaStyle()` is replaced by a theme-aware cache
- `internal/ui/styles.go`: no changes in this phase (left for future UI palette migration)

## Non-goals

- Migrating the 22 hardcoded styles in `styles.go` to `Theme.Colors` (phase 2)
- Loading themes from external YAML/JSON files
- User-defined custom themes
- `tokyo-night` or additional themes beyond the initial three
