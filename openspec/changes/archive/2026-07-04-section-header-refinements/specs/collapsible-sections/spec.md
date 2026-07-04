## MODIFIED Requirements

### Requirement: Collapse state is visually indicated

The system SHALL show a visual indicator on collapsed section headers.

#### Scenario: No indicator when expanded
- **WHEN** a section is expanded
- **THEN** the section header SHALL show no collapse indicator — just the section name and count

#### Scenario: Ellipsis when collapsed
- **WHEN** a section is collapsed
- **THEN** the section header SHALL display a unicode ellipsis `…` after the count, rendered in muted help text style

## ADDED Requirements

### Requirement: Enter is a no-op on section headers

The system SHALL NOT navigate anywhere when Enter is pressed on a section header.

#### Scenario: Enter on section does nothing
- **WHEN** the cursor is on a section header and the user presses `Enter`
- **THEN** nothing happens — the cursor stays on the section header and no view changes

#### Scenario: Enter on other items still works
- **WHEN** the cursor is on a change, spec, requirement, or archived change and the user presses `Enter`
- **THEN** the item SHALL open normally (change viewer, spec viewer, archive viewer respectively)
