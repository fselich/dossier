## ADDED Requirements

### Requirement: Sort toggle for Specifications section
In `ModeIndex`, pressing `s` SHALL toggle the sort order of the Specifications section between two modes: **name** (alphabetical by full spec name, the default) and **suffix** (alphabetical by the last `-`-delimited segment of the name). The toggle SHALL be session-only; on startup the sort is always **name**.

#### Scenario: Toggle from name to suffix
- **WHEN** the sort mode is **name** and the user presses `s` in `ModeIndex`
- **THEN** the Specifications section is re-rendered with specs ordered by their last `-`-delimited segment

#### Scenario: Toggle from suffix to name
- **WHEN** the sort mode is **suffix** and the user presses `s` in `ModeIndex`
- **THEN** the Specifications section is re-rendered with specs in full-name alphabetical order

#### Scenario: Default sort on startup
- **WHEN** the TUI starts and enters `ModeIndex`
- **THEN** specs are listed in alphabetical order by full name (mode is **name**)

### Requirement: Suffix extraction
The suffix of a spec name SHALL be the substring after the last `-` character. If the name contains no `-`, the entire name is the suffix.

#### Scenario: Multi-segment name
- **WHEN** the spec name is `archive-viewer`
- **THEN** its sort suffix is `viewer`

#### Scenario: Single-segment name
- **WHEN** the spec name is `loader`
- **THEN** its sort suffix is `loader`

#### Scenario: Tie-breaking within suffix sort
- **WHEN** two specs share the same suffix (e.g. `archive-viewer` and `tui-viewer`)
- **THEN** their relative order within the tie SHALL follow their full-name alphabetical order (stable sort)

### Requirement: Cursor preservation across sort toggle
When the sort order is toggled, the cursor SHALL remain on the same spec (or requirement child) it was on before the toggle.

#### Scenario: Cursor on spec item when toggling
- **WHEN** the cursor is on a spec item and the user presses `s`
- **THEN** after the sort toggle the cursor is still on the same spec, now at its new position in the list

#### Scenario: Cursor on active or archived item when toggling
- **WHEN** the cursor is on an active change or archived change item and the user presses `s`
- **THEN** the cursor position is unchanged (active and archived sections do not reorder)

### Requirement: Live reload preserves sort mode
When the tick detects a filesystem change while in `ModeIndex` and reloads the spec list, the current sort mode SHALL be preserved.

#### Scenario: New spec appears while in suffix sort mode
- **WHEN** the sort mode is **suffix** and a new spec directory appears on disk
- **THEN** the new spec is inserted at the correct position in the suffix-sorted list within the next poll cycle
