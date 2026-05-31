## MODIFIED Requirements

### Requirement: Focused single-requirement rendering
When `ModeViewingSpec` is activated from an `indexKindRequirement` item, the TUI SHALL render only the block of that requirement (from `### Requirement: <name>` to the next `### Requirement:` or end of file) using `openspec.ExtractRequirement` instead of a local `extractRequirement` function.

#### Scenario: Abrir requirement desde el índice muestra solo ese requirement
- **WHEN** the index cursor is on a requirement item and the user presses `Enter`
- **THEN** the viewport shows only the content of that requirement, extracted via `openspec.ExtractRequirement`

#### Scenario: Requirement no encontrado muestra mensaje de error
- **WHEN** the requirement name does not exist in the spec content
- **THEN** the viewport shows `(spec not available)`
