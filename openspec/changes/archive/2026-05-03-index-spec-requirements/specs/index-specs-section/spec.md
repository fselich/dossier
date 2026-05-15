## MODIFIED Requirements

### Requirement: Project specs loading
The loader SHALL expose a `LoadProjectSpecs()` function that reads `openspec/specs/` in the current working directory and returns a list of `ProjectSpec`, each with its `Name`, the `RequirementCount` obtained by counting `### Requirement: ` lines in the corresponding `spec.md`, and `RequirementNames []string` with the names of each requirement in the order they appear. Specs SHALL be sorted alphabetically by name.

#### Scenario: Specs available
- **WHEN** `openspec/specs/` contains two or more subdirectories with `spec.md`
- **THEN** `LoadProjectSpecs()` returns one entry per spec with the correct requirement count and the list of requirement names in order of appearance

#### Scenario: Specs directory absent
- **WHEN** `openspec/specs/` does not exist or is empty
- **THEN** `LoadProjectSpecs()` returns an empty list without error

#### Scenario: spec.md with no requirements
- **WHEN** a `spec.md` contains no `### Requirement:` line
- **THEN** the corresponding `ProjectSpec` has `RequirementCount` 0 and `RequirementNames` empty, and still appears in the list

### Requirement: Specs section in the index
In `ModeIndex`, the TUI SHALL display a "Specifications" section below the "Archived Changes" section. The section SHALL list each spec on a line with its name on the left and `N requirements` in secondary style on the right, aligned in two columns. The width of the name column SHALL adjust to the longest name. If there are no specs, the section SHALL show a message indicating that no specs are available.

#### Scenario: Specs present
- **WHEN** the mode is `ModeIndex` and `LoadProjectSpecs()` returns at least one spec
- **THEN** the screen shows a "Specifications" section with each spec in the format `name  N requirements`, with the count column aligned

#### Scenario: Column alignment with different name lengths
- **WHEN** there are specs with names of different lengths
- **THEN** all `N requirements` counts appear aligned in the same column

#### Scenario: No specs available
- **WHEN** the mode is `ModeIndex` and `LoadProjectSpecs()` returns an empty list
- **THEN** the "Specifications" section shows the message "No specifications available"

#### Scenario: Specs loaded when entering the index
- **WHEN** the user enters `ModeIndex`
- **THEN** the spec list is loaded from disk at that moment

## ADDED Requirements

### Requirement: Expand spec to see its requirements
In `ModeIndex`, pressing `Space` on a spec item SHALL toggle its expanded/collapsed state. When expanded, the names of its requirements SHALL appear indented below the spec, one per line, as navigable items. Pressing `Space` again SHALL collapse the list. The expanded state of each spec SHALL be independent of the cursor and SHALL reset when leaving and re-entering `ModeIndex`.

#### Scenario: Expand a spec
- **WHEN** the cursor is on a spec item and the user presses `Space`
- **THEN** the requirement names of the spec appear indented below it as navigable items

#### Scenario: Collapse an expanded spec
- **WHEN** the cursor is on an expanded spec item and the user presses `Space`
- **THEN** the requirement items disappear and the spec is shown again as a single line

#### Scenario: Multiple specs expanded simultaneously
- **WHEN** the user expands two different specs
- **THEN** both show their requirements simultaneously and independently

#### Scenario: Expanded state does not affect another spec's cursor
- **WHEN** spec A is expanded and the user moves the cursor to spec B and presses `Space`
- **THEN** spec B expands without collapsing spec A

#### Scenario: Expand state resets on re-entering the index
- **WHEN** the user has expanded specs, leaves `ModeIndex`, and re-enters
- **THEN** all specs appear collapsed

### Requirement: Requirements as navigable items in the index
Requirements of an expanded spec SHALL be navigable items in `indexItems` with kind `indexKindRequirement`. The cursor SHALL be able to move through them with `j`/`k` continuously alongside the rest of the index items. Requirement items SHALL NOT be selectable with `Space`.

#### Scenario: j/k traverses requirements of an expanded spec
- **WHEN** spec A is expanded and the cursor is on the last requirement of A
- **THEN** pressing `j` moves the cursor to the next item (spec B or the first requirement of B if expanded)

#### Scenario: Space on a requirement has no effect
- **WHEN** the cursor is on a requirement item and the user presses `Space`
- **THEN** no change occurs in the index

### Requirement: Cursor snap when collapsing a spec with cursor inside
When the user collapses a spec whose cursor was positioned on one of its requirements, the cursor SHALL jump to the collapsed spec item.

#### Scenario: Collapse spec with cursor on a requirement
- **WHEN** the cursor is on a requirement of spec A and the user moves the cursor to spec A and presses `Space` to collapse
- **THEN** the cursor is positioned on the spec A item

#### Scenario: Collapsing spec without cursor inside does not move the cursor
- **WHEN** the cursor is on an item other than spec A and its requirements, and the user presses `Space` on spec A
- **THEN** the cursor remains on the item where it was
