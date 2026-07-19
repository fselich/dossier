# theme-system Specification

## Purpose
Defines the centralized theme system that allows selecting a named visual theme via `--theme` CLI flag, mapping to glamour markdown style, chroma syntax-highlighting style, view background color, and diff background colors.

## Requirements

### Requirement: Built-in themes define styling for all renderers
The system SHALL provide a set of named built-in themes. Each theme SHALL specify a glamour markdown style, a chroma syntax-highlighting style, a view background color, and diff add/remove background colors. Theme names SHALL be case-insensitive when specified via the CLI flag.

The built-in themes SHALL be:

| Theme | Glamour style | Chroma style | Background | Diff Add | Diff Remove |
|-------|--------------|--------------|------------|----------|-------------|
| `none` (default) | `dark` | `monokai` | Terminal default (nil) | `#1a3a1a` | `#3a1a1a` |
| `dark` | `dark` | `monokai` | Dark grey (#1a1a1a) | `#1a3a1a` | `#3a1a1a` |
| `light` | `light` | `github` | White (#ffffff) | `#e6ffed` | `#ffeef0` |
| `dracula` | `dracula` | `dracula` | Dracula (#282a36) | `#1f3425` | `#3d1f26` |

#### Scenario: Theme struct carries all renderer configurations
- **WHEN** a built-in theme is defined
- **THEN** it includes a `Name` string, a `GlamourStyle` string, a `ChromaStyle` string, a `ViewBg` color, a `DiffAddBg` string, and a `DiffRemoveBg` string

#### Scenario: Theme name is case-insensitive
- **WHEN** the user specifies `--theme DARK` or `--theme Dark`
- **THEN** the system resolves it to the built-in `dark` theme

### Requirement: Theme selection via `--theme` CLI flag
The system SHALL accept a `--theme <name>` CLI flag that selects which built-in theme to use. If the flag is omitted, the `none` theme SHALL be used (dark UI colors, terminal default background). If the flag value does not match any built-in theme name (case-insensitive), the system SHALL print the list of valid theme names to stderr and exit with code 1.

The system SHALL use Go's `flag` package for argument parsing, accepting `--theme`, `--version`, `--help`, and an optional positional `[path]` argument.

#### Scenario: No `--theme` flag selects default none theme
- **WHEN** the user runs `dossier` without `--theme`
- **THEN** the `none` theme is used (dark UI colors, terminal default background)

#### Scenario: Valid `--theme` flag selects matching theme
- **WHEN** the user runs `dossier --theme dracula`
- **THEN** the `dracula` theme is used (glamour "dracula", chroma "dracula", dracula background)

#### Scenario: Invalid `--theme` flag exits with error
- **WHEN** the user runs `dossier --theme nonexistent`
- **THEN** the system prints an error listing valid theme names and exits with code 1

#### Scenario: `--theme` with a path argument
- **WHEN** the user runs `dossier --theme light openspec/changes/my-feature`
- **THEN** the TUI opens in single-change mode with the `light` theme applied

### Requirement: Glamour renderer uses theme's glamour style
When creating the glamour `TermRenderer` for markdown rendering, the system SHALL pass the active theme's `GlamourStyle` to `glamour.WithStandardStyle()` instead of a hardcoded value. If glamour fails to create a renderer for the given style, the system SHALL fall back to displaying raw markdown text without crashing. If `GlamourStyle` is empty, the system SHALL fall back to `"dark"`.

#### Scenario: Dark theme uses glamour "dark" style
- **WHEN** the active theme is `dark`
- **THEN** glamour renders markdown using the "dark" style configuration

#### Scenario: Light theme uses glamour "light" style
- **WHEN** the active theme is `light`
- **THEN** glamour renders markdown using the "light" style configuration

#### Scenario: Glamour style resolution failure falls back to raw text
- **WHEN** the theme specifies a glamour style name that glamour cannot resolve
- **THEN** the system displays the raw markdown content without ANSI styling, and no panic or crash occurs

### Requirement: Chroma syntax highlighting uses theme's chroma style
When highlighting code in git diffs, the system SHALL use the active theme's `ChromaStyle` to look up the chroma style. If chroma cannot resolve the style name, it SHALL fall back to chroma's built-in `Fallback` style. The chroma style SHALL be cached by style name to avoid repeated lookups. If `ChromaStyle` is empty, the system SHALL fall back to `"monokai"`.

#### Scenario: Dark theme uses chroma "monokai"
- **WHEN** the active theme is `dark` and a git diff is displayed
- **THEN** code syntax is highlighted using chroma's "monokai" style

#### Scenario: Dracula theme uses chroma "dracula"
- **WHEN** the active theme is `dracula` and a git diff is displayed
- **THEN** code syntax is highlighted using chroma's "dracula" style

#### Scenario: Chroma style not found uses fallback
- **WHEN** the theme specifies a chroma style name that chroma cannot resolve
- **THEN** code syntax is highlighted using chroma's `Fallback` style without error

### Requirement: View background color is set from theme
The system SHALL set `tea.View.BackgroundColor` to the active theme's `ViewBg` color when creating the view. If `ViewBg` is nil or unset, the terminal's default background SHALL be used.

#### Scenario: Dark theme sets dark grey background
- **WHEN** the active theme is `dark`
- **THEN** `tea.View.BackgroundColor` is set, filling the entire viewport with dark grey

#### Scenario: Light theme sets white background
- **WHEN** the active theme is `light`
- **THEN** `tea.View.BackgroundColor` is set, filling the entire viewport with white

#### Scenario: Nil background leaves terminal default
- **WHEN** a theme has `ViewBg` set to nil
- **THEN** the terminal emulator's default background is visible

### Requirement: Diff background colors are theme-configurable
The system SHALL use the active theme's `DiffAddBg` and `DiffRemoveBg` colors when rendering git diff backgrounds for added and removed lines. If either color is empty, the system SHALL fall back to dark defaults (`#1a3a1a` for added, `#3a1a1a` for removed).

#### Scenario: Light theme uses light diff backgrounds
- **WHEN** the active theme is `light` and a git diff is displayed
- **THEN** added lines have background `#e6ffed` and removed lines have background `#ffeef0`

#### Scenario: Dark theme uses dark diff backgrounds
- **WHEN** the active theme is `dark` and a git diff is displayed
- **THEN** added lines have background `#1a3a1a` and removed lines have background `#3a1a1a`

### Requirement: Theme defines UI chrome colors
Each theme SHALL define a `ThemeColors` struct with semantic color roles used by all lipgloss UI styles. The struct SHALL contain fields for `PrimaryFg`, `MutedFg`, `MidFg`, `AccentBlue`, `AccentYellow`, `AccentCyan`, `AccentGreen`, `AccentRed`, `AccentMagenta`, `ActiveBg`, and `ActiveFg`. Each field SHALL be a `lipgloss.Color`.

#### Scenario: Dark theme uses current hardcoded values
- **WHEN** the active theme is `dark`
- **THEN** `PrimaryFg` is `"15"`, `MutedFg` is `"8"`, `AccentBlue` is `"12"`, matching the pre-theme behavior

#### Scenario: Light theme uses dark-adapted colors
- **WHEN** the active theme is `light`
- **THEN** `PrimaryFg` is `"0"` (black, readable on white background)
- **AND** `AccentYellow` is `"3"` (darker yellow, readable on white)
- **AND** `AccentRed` is `"1"` (red, not bright red `"9"`)
- **AND** `MidFg` is `"8"` (dark gray, pending tasks readable) and `MutedFg` is `"7"` (light gray, done tasks dimmed)

#### Scenario: none theme inherits dark theme colors
- **WHEN** the active theme is `none`
- **THEN** `ThemeColors` is identical to the `dark` theme
- **AND** only `ViewBg` differs (nil vs dark grey)

### Requirement: ThemeStyles are pre-built from ThemeColors
The system SHALL provide a `BuildStyles(ThemeColors) ThemeStyles` function that constructs all lipgloss styles from a `ThemeColors` value. The resulting `ThemeStyles` SHALL be stored on the `Theme` struct. Each style SHALL be constructed once at startup.

#### Scenario: BuildStyles produces all styles from colors
- **WHEN** `BuildStyles(c)` is called with a `ThemeColors` value
- **THEN** it returns a `ThemeStyles` with all ~24 lipgloss styles populated
- **AND** each style uses the corresponding color from `c`

#### Scenario: Header style uses AccentBlue
- **WHEN** the theme's `AccentBlue` is `"4"`
- **THEN** `ThemeStyles.Header` is `lipgloss.NewStyle().Bold(true).Foreground(Color("4"))`

### Requirement: UI rendering uses theme styles
All UI rendering functions SHALL use `m.theme.Styles.Xxx` instead of package-level `var` style declarations. No rendering code SHALL reference hardcoded ANSI color values outside of `ThemeColors` definitions.

#### Scenario: Tab bar uses theme styles
- **WHEN** rendering the tab bar
- **THEN** the active tab uses `m.theme.Styles.TabActive`
- **AND** inactive tabs use `m.theme.Styles.TabInactive`
- **AND** disabled tabs use `m.theme.Styles.TabDisabled`

#### Scenario: Git status file list uses theme styles
- **WHEN** rendering the git status file list
- **THEN** modified files use `m.theme.Styles.GitModified`
- **AND** added files use `m.theme.Styles.GitAdded`
- **AND** cursor marks use `m.theme.Styles.TaskCursorMark`

#### Scenario: Task inline markdown uses theme styles
- **WHEN** rendering inline markdown in task items
- **THEN** code spans use `m.theme.Styles.TaskCodeCyan` (pending) or `m.theme.Styles.TaskCodeDone` (done)

### Requirement: Dracula theme has adjusted AccentYellow
The `dracula` theme SHALL use `AccentYellow` = `"3"` (instead of `"11"`) because bright yellow clashes with the dracula background (`#282a36`). All other UI colors SHALL match the `dark` theme.

#### Scenario: Dracula sections use darker yellow
- **WHEN** the active theme is `dracula`
- **THEN** section headers and git modified files render with `AccentYellow` = `"3"`

### Requirement: BaseText style applies PrimaryFg to unstyled text
The system SHALL provide a `BaseText` style in `ThemeStyles` that applies only `PrimaryFg` (no bold, no background). Non-selected item names in the index view SHALL use `BaseText` to ensure they are readable regardless of the active theme's background color.

#### Scenario: BaseText style uses PrimaryFg
- **WHEN** `BuildStyles(c)` is called with a `ThemeColors` value
- **THEN** `ThemeStyles.BaseText` is `lipgloss.NewStyle().Foreground(c.PrimaryFg)`

#### Scenario: Active change names use BaseText when not selected
- **WHEN** rendering an active change in the index view and the item is not the cursor
- **THEN** the change name is rendered with `m.theme.Styles.BaseText`

#### Scenario: Light theme non-selected names are readable
- **WHEN** the active theme is `light` (background `#ffffff`) and an item is not the cursor
- **THEN** the item name uses `PrimaryFg` = `"0"` (black), providing contrast against the white background
