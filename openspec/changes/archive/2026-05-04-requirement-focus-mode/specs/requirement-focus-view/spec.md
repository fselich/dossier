## ADDED Requirements

### Requirement: Focused single-requirement rendering
When `ModeViewingSpec` is activated from an `indexKindRequirement` item, the TUI SHALL render only that requirement's block (from `### Requirement: <name>` to the next `### Requirement:` or end of file) instead of the full spec.

#### Scenario: Opening a requirement from the index shows only that requirement
- **WHEN** the index cursor is on a requirement item and the user presses `Enter`
- **THEN** the viewport shows only the content of that requirement, without the rest of the spec

#### Scenario: Requirement not found shows error message
- **WHEN** the requirement name does not exist in the spec content
- **THEN** the viewport shows `(spec not available)`

### Requirement: Navigation between requirements with h/l in focus mode
In focus mode, `h` and `l` SHALL navigate to the previous and next requirement within the same spec, respectively, updating the viewport with the new requirement.

#### Scenario: Navigate to the next requirement with l
- **WHEN** focus mode is active and the user presses `l`
- **THEN** the viewport shows the next requirement in `RequirementNames`; if it is the last one, it wraps to the first

#### Scenario: Navigate to the previous requirement with h
- **WHEN** focus mode is active and the user presses `h`
- **THEN** the viewport shows the previous requirement in `RequirementNames`; if it is the first one, it wraps to the last

#### Scenario: Spec with a single requirement
- **WHEN** focus mode is active, the spec has only one requirement, and the user presses `h` or `l`
- **THEN** the viewport continues to show the same requirement (wraps to itself)

### Requirement: Req N/M counter in the header during focus mode
In focus mode, the header SHALL show `<project>  ·  <spec-name>  ·  Req N/M` where `N` is the 1-based position of the current requirement and `M` is the total number of requirements in the spec.

#### Scenario: Header in focus mode
- **WHEN** focus mode is active
- **THEN** the header shows `<project>  ·  <spec-name>  ·  Req N/M`

#### Scenario: Header in ModeViewingSpec without focus mode
- **WHEN** `ModeViewingSpec` was opened from a spec item (not a requirement item)
- **THEN** the header shows `<project>  ·  <spec-name>  [spec]` (existing behaviour)

### Requirement: HelpBar in focus mode
In focus mode the helpbar SHALL show the relevant controls: navigation between requirements, return to index, and quit.

#### Scenario: HelpBar in focus mode
- **WHEN** focus mode is active
- **THEN** the helpbar shows `h/l: prev/next req  Esc: index  q: quit`
