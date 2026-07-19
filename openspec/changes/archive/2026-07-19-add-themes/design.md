## Context

Dossier renders content through three independent styling systems:

| System | What it styles | Current behavior | Hardcoded at |
|--------|---------------|------------------|--------------|
| Lipgloss (UI) | Header, tabs, borders, tasks, git, help bar | 22 `var` declarations with ANSI 0-15 colors | `internal/ui/styles.go` |
| Glamour | Markdown artifacts (proposal, design, specs, tasks, config) | `WithStandardStyle("dark")` | `internal/ui/viewport.go:37` |
| Chroma | Code syntax highlighting in git diffs | `styles.Get("monokai")` | `internal/ui/gitdiff.go:45` |

A `Theme` struct exists on `Model` (`internal/ui/model.go:127-128`) but carries only `ViewBg color.Color`, used solely for `tea.NewView` background in `View()`. No mechanism exists to select or switch themes.

The project has specs for `view-background` (how background fills the viewport) and `tui-viewer` (references `tea.View.BackgroundColor`), and PROPUESTAS.md lists themes (#3.1) as a high-priority feature.

## Goals / Non-Goals

**Goals:**
- Select a named theme via `--theme <name>` CLI flag
- Define three built-in themes (`dark`, `light`, `dracula`) mapping glamour style + chroma style + view background color
- Thread the selected glamour style through `ensureRenderer()`
- Thread the selected chroma style through `renderDiff()`/`highlightLine()`
- Provide per-theme background color for `tea.NewView`
- Use Go `flag` package for argument parsing (replaces manual `os.Args` inspection)

**Non-Goals:**
- Migrate the 22 lipgloss styles in `styles.go` to theme-driven colors
- Load themes from external files (YAML, JSON)
- User-defined or custom themes
- Runtime theme switching (requires restart)

## Decisions

### 1. Theme definition: Go structs, not external config

Themes are defined as Go constants (a `map[string]Theme{}` in `internal/ui/themes.go`). Each `Theme` struct maps a user-facing name to the concrete glamour style name and chroma style name.

**Rationale:** Phase 1 only needs 3 themes. Go constants avoid parsing complexity, file I/O, and error handling for malformed config files. External config (YAML) is a natural phase 2 extension.

**Alternatives considered:**
- YAML config file: more flexible but premature for 3 themes. Parsing overhead, error surface, and discovery complexity not worth it yet.
- JSON: same tradeoffs as YAML.

### 2. Theme struct fields

```go
type Theme struct {
    Name         string // user-facing identifier (matches map key)
    GlamourStyle string // passed to glamour.WithStandardStyle()
    ChromaStyle  string // passed to chroma styles.Get()
    ViewBg       color.Color // passed to tea.NewView.BackgroundColor
}
```

`Name` duplicates the map key but serves as a canonical identifier on the struct itself — useful if a Theme ever travels without its map context (e.g., logging, error messages).

### 3. Built-in themes

| Theme | Glamour | Chroma | ViewBg (ANSI) | Rationale |
|-------|---------|--------|---------------|-----------|
| `dark` (default) | `dark` | `monokai` | `#1a1a1a` (234) | Current behavior preserved; monokai is the existing chroma choice |
| `light` | `light` | `github` | `#ffffff` (15) | Natural light counterpart; github is the most popular light chroma style |
| `dracula` | `dracula` | `dracula` | `#282a36` | Only concordant pair between glamour/chroma; popular dark theme |

View background colors are chosen to match each theme's natural background (dark → near-black, light → white, dracula → the dracula spec background).

**Why these three?** One dark (existing default), one light, one concordant. Covers the matrix. Additional themes (`tokyo-night`, `nord`, `catppuccin-*`) can be added as one-liners to the map.

**Alternatives considered:**
- More themes: easy to add later, no architectural change needed. Three validates the system.
- `dracula` as default: breaking change for existing users. `dark` preserves current behavior.

### 4. Chroma style threading: parameter chain, not global mutable state

Current package-level globals (`chromaStyle`, `chromaOnce` in `gitdiff.go:37-41`) are replaced with a cache keyed by style name:

```go
var chromaStyleCache = map[string]*chroma.Style{}

func getChromaStyle(name string) *chroma.Style {
    if s, ok := chromaStyleCache[name]; ok {
        return s
    }
    s := styles.Get(name)
    if s == nil {
        s = styles.Fallback
    }
    chromaStyleCache[name] = s
    return s
}
```

Functions in the call chain gain a `chromaStyleName string` parameter:

```
refreshGitViewport(*Model)         → reads m.theme.ChromaStyle
  → renderDiffContent(*Model)      → passes chromaStyleName through
    → renderDiff(..., styleName)   → passes to highlightLine
      → highlightLine(..., styleName) → calls getChromaStyle(name)
```

**Rationale:** No global mutable state, no sync primitives needed (cache is only written during render, which is single-threaded in Bubble Tea's update loop). Each function explicitly receives the style it needs.

**Alternatives considered:**
- Global var mutated on theme change: simpler call chain, but introduces mutable global state that's fragile to concurrent access and testing.
- Method on Model: couples syntax highlighting to Model, making `highlightLine` untestable without a full Model.

### 5. Glamour style threading: trivial change

`ensureRenderer(width)` already reads from `m.theme` (for cache width tracking). The change is a one-line replacement: `"dark"` → `m.theme.GlamourStyle`. The renderer cache (`glamourRenderer`, `lastRenderWidth`) already gets invalidated on width change; it should also invalidate on theme change. Since the theme is set at startup and never changes at runtime, no extra invalidation logic is needed.

### 6. CLI flag parsing: `flag` package

Replace the current manual `os.Args` inspection in `cmd/dossier/main.go` with Go's standard `flag` package:

```
dossier [--theme <name>] [--version] [--help] [path]
```

`--theme` defaults to `"dark"`. Invalid theme names cause an immediate error with a list of valid names.

**Rationale:** Incidental cleanup. PROPUESTAS.md #1.9 flags the manual parsing as fragile (e.g., `dossier ruta --help` fails). Using `flag` fixes this and gives us `--theme` parsing for free. Also makes `--version` and `--help` standard `flag` output.

### 7. Invalid theme name behavior: fail early

If `--theme` receives an unrecognized name, the program prints available themes and exits with code 1 before starting the TUI. No silent fallback.

**Rationale:** A typo should be visible, not silently ignored. The user chose a theme explicitly; if it doesn't exist, they should know.

### 8. No runtime theme switching

The theme is set at startup and immutable for the session lifetime. No hot-reload, no keybinding to cycle themes.

**Rationale:** Minimizes scope. Recreating glamour/chroma renderers (which have internal caches) and repopulating render caches mid-session adds complexity without proportional value. Users who want a different theme restart the program.

## Risks / Trade-offs

- **[Low] glamour/chroma style name mismatch in future versions**: If glamour or chroma deprecates a style name used by a built-in theme, the style silently falls back (glamour: error creating renderer → raw text; chroma: `styles.Fallback`). **Mitigation:** The chroma layer already has a `Fallback` mechanism. Glamour's `NewTermRenderer` can fail; this change should handle that gracefully (see PROPUESTAS.md #1.3).

- **[Low] ViewBg color may not match glamour/chroma background**: The hardcoded `ViewBg` values are chosen to approximate each theme's look, but glamour and chroma don't expose their background colors programmatically. There's no guarantee of pixel-perfect match. **Mitigation:** Acceptable for phase 1. Colors can be tuned. A future phase could extract background from glamour/chroma configs.

- **[None] Breaking change for existing users**: `--theme dark` preserves the current behavior (glamour "dark" + chroma "monokai"). No migration needed. The only visible difference is the view background color going from terminal default to ANSI 234 (near-black), which is imperceptible on most dark terminals.

## Open Questions

- Should `--theme` be case-insensitive? (e.g., `--theme Dark` → `dark`). Leaning yes, as terminal users expect case-insensitive flags.
- Should `Theme.Colors` (the future lipgloss palette) be a flat struct of `lipgloss.Color` fields, or a map keyed by semantic role names? Flat struct gives compile-time safety but is verbose; map is more concise but runtime-error-prone. Defer to phase 2.
