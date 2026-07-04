# Collapsible Sections

## Purpose

The index view sections (Active Changes, Specifications, Archived Changes) can be collapsed and expanded to let users focus on relevant content and reduce scrolling.

## Requirements

### Requirement: Index sections are navigable

The system SHALL include section header items in the index item list so the cursor can land on section headers.

#### Scenario: Cursor lands on section header
- **WHEN** the index view is rendered with cursor at position 0
- **THEN** the cursor SHALL be on the first section header (not skipping past it)

#### Scenario: Arrow navigation moves through sections
- **WHEN** the user presses `j` or `down` repeatedly in the index view
- **THEN** the cursor SHALL move through section headers and items in sequence

### Requirement: Sections can be collapsed and expanded

The system SHALL support collapsing and expanding index sections via the `Space` key.

#### Scenario: Collapse a section
- **WHEN** the cursor is on an expanded section header and the user presses `Space`
- **THEN** the section SHALL collapse and all child items of that section SHALL be hidden

#### Scenario: Expand a section
- **WHEN** the cursor is on a collapsed section header and the user presses `Space`
- **THEN** the section SHALL expand and its child items SHALL become visible

#### Scenario: Space on spec still expands requirements
- **WHEN** the cursor is on a spec item and the user presses `Space`
- **THEN** the spec SHALL toggle its requirement expansion as before

### Requirement: Collapse state is visually indicated

The system SHALL show a visual indicator on collapsed section headers.

#### Scenario: No indicator when expanded
- **WHEN** a section is expanded
- **THEN** the section header SHALL show no collapse indicator — just the section name and count

#### Scenario: Ellipsis when collapsed
- **WHEN** a section is collapsed
- **THEN** the section header SHALL display a unicode ellipsis `…` after the count, rendered in muted help text style

### Requirement: Collapse state persists across rebuilds

The system SHALL preserve section collapse state across index view rebuilds (tick polling, filter changes, mode switches).

#### Scenario: Collapse survives poll
- **WHEN** a section is collapsed and the index view is rebuilt on a tick poll
- **THEN** the section SHALL remain collapsed after the rebuild

### Requirement: Filtering respects collapse

The system SHALL not display items inside a collapsed section when a filter is active.

#### Scenario: Filtered items hidden in collapsed section
- **WHEN** a section is collapsed and a filter is applied
- **THEN** items that would match the filter but belong to the collapsed section SHALL NOT be shown

### Requirement: Help bar shows section toggle action

The system SHALL update the help bar to indicate that `Space` toggles sections.

#### Scenario: Help bar updated
- **WHEN** the index view is displayed
- **THEN** the help bar SHALL show the `Space` action for toggling sections

### Requirement: Enter is a no-op on section headers

The system SHALL NOT navigate anywhere when Enter is pressed on a section header.

#### Scenario: Enter on section does nothing
- **WHEN** the cursor is on a section header and the user presses `Enter`
- **THEN** nothing happens — the cursor stays on the section header and no view changes

#### Scenario: Enter on other items still works
- **WHEN** the cursor is on a change, spec, requirement, or archived change and the user presses `Enter`
- **THEN** the item SHALL open normally (change viewer, spec viewer, archive viewer respectively)
