# per-mode-handlers Specification

## Purpose

Each UI mode has its own `update*()` method in a dedicated file, with `update.go` acting as a thin dispatcher that delegates to the active mode. Keys are defined via `key.Binding` structs instead of raw string comparisons.

## ADDED Requirements

### Requirement: Dispatcher delegates to per-mode update functions
The `Update()` method SHALL handle infrastructure messages (`WindowSizeMsg`, `renderedMsg`, `tickMsg`, etc.) and then delegate key/mouse handling to a mode-specific `update*()` method via a `switch` on `m.mode`. The mode-specific methods SHALL return `(tea.Model, tea.Cmd)`.

#### Scenario: Key pressed in ModeNormal
- **WHEN** a `tea.KeyPressMsg` arrives and mode is `ModeNormal`
- **THEN** `Update()` delegates to `updateViewer(msg)` which returns the updated model and command

#### Scenario: Key pressed in ModeIndex
- **WHEN** a `tea.KeyPressMsg` arrives and mode is `ModeIndex`
- **THEN** `Update()` delegates to `updateIndex(msg)` which returns the updated model and command

#### Scenario: Key pressed in ModeViewingSpec
- **WHEN** a `tea.KeyPressMsg` arrives and mode is `ModeViewingSpec`
- **THEN** `Update()` delegates to `updateSpec(msg)` which returns the updated model and command

#### Scenario: Key pressed in ModeViewingConfig
- **WHEN** a `tea.KeyPressMsg` arrives and mode is `ModeViewingConfig`
- **THEN** `Update()` delegates to `updateConfig(msg)` which returns the updated model and command

### Requirement: Per-mode files for update handlers
Each mode SHALL have its own file containing the mode-specific `update*()` method and its keybinding definitions.

#### Scenario: viewer.go exists
- **WHEN** the codebase is inspected
- **THEN** `internal/ui/viewer.go` contains `updateViewer()` for `ModeNormal` and `ModeViewingArchive`

#### Scenario: index.go has update method
- **WHEN** the codebase is inspected
- **THEN** `internal/ui/index.go` contains `updateIndex()` for `ModeIndex`

#### Scenario: spec.go exists
- **WHEN** the codebase is inspected
- **THEN** `internal/ui/spec.go` contains `updateSpec()` for `ModeViewingSpec`

#### Scenario: config.go exists
- **WHEN** the codebase is inspected
- **THEN** `internal/ui/config.go` contains `updateConfig()` for `ModeViewingConfig`

### Requirement: Key bindings defined with key.Binding structs
Each per-mode file SHALL define a `keymap` struct using `key.Binding` from the Bubbles `key` package. Key matching SHALL use `key.Matches(msg, binding)` instead of raw string comparisons.

#### Scenario: Viewer keymap defined
- **WHEN** the viewer mode is active
- **THEN** keybindings for artifact switching, navigation, editing, and quitting are defined as `key.Binding` values

#### Scenario: Key matches with key.Matches
- **WHEN** a key is pressed and compared against a `key.Binding`
- **THEN** `key.Matches(msg, binding)` returns true if the key matches any bound key in the binding

#### Scenario: Disabled binding does not match
- **WHEN** a `key.Binding` is marked disabled via `.SetEnabled(false)`
- **THEN** `key.Matches(msg, binding)` returns false regardless of the key pressed

### Requirement: Help bar generated from key bindings
The help bar SHALL use the `help` bubble's `ShortHelpView()` method, passing the active mode's relevant `key.Binding` slice. The manual `renderHelpBar()` switch SHALL be replaced by this auto-generated output.

#### Scenario: Tasks tab help bar
- **WHEN** the active tab is tasks and mode is ModeNormal
- **THEN** the help bar shows bindings for navigation, toggle, edit, index, and quit

#### Scenario: Disabled keys hidden from help
- **WHEN** a `key.Binding` is disabled
- **THEN** the help bar omits that key from the displayed shortcuts

### Requirement: User-facing keybindings unchanged
The per-mode refactor SHALL NOT change any user-facing keybinding behavior. All existing keys (`q`, `h`, `l`, `j`, `k`, `1`-`4`, `tab`, `shift+tab`, `enter`, `space`, `a`, `s`, `e`, `i`, `esc`, `ctrl+c`, `up`, `down`) SHALL produce the same results as before.

#### Scenario: All existing keybindings preserved
- **WHEN** any existing test case for key handling is run
- **THEN** the test passes without modification to the expected behavior assertions
