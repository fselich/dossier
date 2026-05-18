## Purpose

The project config view is a full-screen read-only TUI mode that displays the contents of `openspec/config.yaml`, rendering the project's context and rules as formatted markdown. It is accessible from both `ModeNormal` and `ModeIndex` via the `i` keybinding.

## Requirements

### Requirement: User can open project config view
The TUI SHALL provide a `i` keybinding that opens a full-screen read-only view of the project's `openspec/config.yaml` content, accessible from both `ModeNormal` and `ModeIndex`.

#### Scenario: Open from index
- **WHEN** the user is in `ModeIndex` and presses `i`
- **THEN** the TUI transitions to `ModeViewingConfig` and renders the config content in the viewport

#### Scenario: Open from normal mode
- **WHEN** the user is in `ModeNormal` (browsing a change) and presses `i`
- **THEN** the TUI transitions to `ModeViewingConfig` and renders the config content in the viewport

### Requirement: Config view renders context and rules as markdown
The config view SHALL display the `context` field as a prose section and the `rules` field as grouped bullet lists, rendered via Glamour. The `schema` field SHALL NOT be displayed.

#### Scenario: Context displayed
- **WHEN** the config view is open
- **THEN** the viewport shows a `## Context` heading followed by the context prose

#### Scenario: Rules displayed
- **WHEN** the config view is open and the config contains rules
- **THEN** each rule key appears as a `### <key>` heading with its items as a bullet list

#### Scenario: Empty config
- **WHEN** `openspec/config.yaml` is missing or has no content
- **THEN** the viewport shows an empty view without crashing

### Requirement: User can exit config view
The config view SHALL support `Esc` and `q` to exit, returning the user to the mode they came from.

#### Scenario: Exit to index
- **WHEN** the user opened the config view from `ModeIndex` and presses `Esc` or `q`
- **THEN** the TUI returns to `ModeIndex`

#### Scenario: Exit to normal mode
- **WHEN** the user opened the config view from `ModeNormal` and presses `Esc` or `q`
- **THEN** the TUI returns to `ModeNormal`

### Requirement: Config view header identifies the current screen
The header bar in `ModeViewingConfig` SHALL display the project name and a `[config]` label, consistent with how `[archive]` and `[spec]` labels work in other modes.

#### Scenario: Header label
- **WHEN** the config view is open
- **THEN** the header reads `<project-name>  ·  project config`

### Requirement: Config view help bar shows navigation hints
The help bar in `ModeViewingConfig` SHALL show `j/k: scroll  i/Esc: back  q: quit`. The help bars in `ModeIndex` and `ModeNormal` SHALL include `i: info` as a hint for opening this view.

#### Scenario: Help bar content in config view
- **WHEN** the config view is open
- **THEN** the help bar displays `j/k: scroll  i/Esc: back  q: quit`

#### Scenario: Help bar hint in index and normal mode
- **WHEN** the user is in `ModeIndex` or `ModeNormal`
- **THEN** the help bar includes `i: info`
