# requirement-focus-view Specification

## Purpose

TBD — Focused viewing mode for a single requirement, activated when entering from a requirement item in the index.

## Requirements

### Requirement: Focused single-requirement rendering
When `ModeViewingSpec` is activated from an `indexKindRequirement` item, the TUI SHALL render only the block of that requirement (from `### Requirement: <name>` to the next `### Requirement:` or end of file) instead of the full spec.

#### Scenario: Abrir requirement desde el índice muestra solo ese requirement
- **WHEN** the index cursor is on a requirement item and the user presses `Enter`
- **THEN** the viewport shows only the content of that requirement, without the rest of the spec

#### Scenario: Requirement no encontrado muestra mensaje de error
- **WHEN** the requirement name does not exist in the spec content
- **THEN** the viewport shows `(spec not available)`

### Requirement: Navegación entre requirements con h/l en focus mode
In focus mode, `h` and `l` SHALL navigate to the previous and next requirement within the same spec, respectively, updating the viewport with the new requirement.

#### Scenario: Navegar al siguiente requirement con l
- **WHEN** focus mode is active and the user presses `l`
- **THEN** the viewport shows the next requirement in `RequirementNames`; if it is the last one, it wraps to the first

#### Scenario: Navegar al requirement anterior con h
- **WHEN** focus mode is active and the user presses `h`
- **THEN** the viewport shows the previous requirement in `RequirementNames`; if it is the first one, it wraps to the last

#### Scenario: Spec con un único requirement
- **WHEN** focus mode is active, the spec has only one requirement, and the user presses `h` or `l`
- **THEN** the viewport continues showing the same requirement (wraps to itself)

### Requirement: Contador Req N/M en el header durante focus mode
In focus mode, the header SHALL show `<project>  ·  <spec-name>  ·  Req N/M` where `N` is the 1-based position of the current requirement and `M` is the total number of requirements in the spec.

#### Scenario: Header en focus mode
- **WHEN** focus mode is active
- **THEN** the header shows `<project>  ·  <spec-name>  ·  Req N/M`

#### Scenario: Header en ModeViewingSpec sin focus mode
- **WHEN** `ModeViewingSpec` was opened from a spec item (not a requirement item)
- **THEN** the header shows `<project>  ·  <spec-name>  [spec]` (existing behavior)

### Requirement: HelpBar en focus mode
In focus mode the helpbar SHALL show the relevant controls: navigation between requirements, return to index, and quit.

#### Scenario: HelpBar en focus mode
- **WHEN** focus mode is active
- **THEN** the helpbar shows `h/l: prev/next req  Esc: index  q: quit`
