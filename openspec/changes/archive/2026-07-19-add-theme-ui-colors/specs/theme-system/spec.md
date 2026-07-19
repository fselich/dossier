## ADDED Requirements

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
