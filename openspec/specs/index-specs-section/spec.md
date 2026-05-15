# index-specs-section Specification

## Purpose
Defines the "Specs" section of the index: loading project specs from `openspec/specs/`, displaying them with requirement counts in `ModeIndex`, and cursor navigation to open each spec in `ModeViewingSpec`.
## Requirements
### Requirement: Carga de specs del proyecto
The loader SHALL expose a `LoadProjectSpecs()` function that reads `openspec/specs/` in the current working directory and returns a list of `ProjectSpec`, each with its `Name`, the `RequirementCount` obtained by counting `### Requirement: ` lines in the corresponding `spec.md`, and `RequirementNames []string` with the names of each requirement in the order they appear. Specs SHALL be sorted alphabetically by name.

#### Scenario: Specs disponibles
- **WHEN** `openspec/specs/` contains two or more subdirectories with `spec.md`
- **THEN** `LoadProjectSpecs()` returns one entry per spec with the correct requirement count and the list of requirement names in order of appearance

#### Scenario: Directorio specs ausente
- **WHEN** `openspec/specs/` does not exist or is empty
- **THEN** `LoadProjectSpecs()` returns an empty list without error

#### Scenario: spec.md sin requisitos
- **WHEN** a `spec.md` contains no `### Requirement:` lines
- **THEN** the corresponding `ProjectSpec` has `RequirementCount` 0 and `RequirementNames` empty, and still appears in the list

### Requirement: Sección Specs en el índice
In `ModeIndex`, the TUI SHALL show a "Specifications" section between the "Active Changes" section and the "Archived Changes" section. The section SHALL list each spec on one line with its name on the left and `N requirements` in secondary style on the right, aligned in two columns. The name column width SHALL adjust to the longest name. If there are no specs, the section SHALL show a message indicating that no specs are available.

#### Scenario: Specs presentes
- **WHEN** the mode is `ModeIndex` and `LoadProjectSpecs()` returns at least one spec
- **THEN** the screen shows a "Specifications" section with each spec in the format `name  N requirements`, with the count column aligned, positioned between "Active Changes" and "Archived Changes"

#### Scenario: Alineación de columnas con nombres de distinta longitud
- **WHEN** there are specs with names of different lengths
- **THEN** all `N requirements` counts appear aligned in the same column

#### Scenario: Sin specs disponibles
- **WHEN** the mode is `ModeIndex` and `LoadProjectSpecs()` returns an empty list
- **THEN** the "Specifications" section shows the message "No specifications available"

#### Scenario: Specs cargados al entrar al índice
- **WHEN** the user enters `ModeIndex`
- **THEN** the spec list is loaded from disk at that moment (same as archived changes)

#### Scenario: Orden de secciones en pantalla
- **WHEN** the mode is `ModeIndex` and active changes, specs, and archived changes exist
- **THEN** the "Specifications" section appears below "Active Changes" and above "Archived Changes"

### Requirement: Specs no seleccionables en el índice
Specs listed in the "Specs" section of `ModeIndex` SHALL be navigable. They SHALL be part of `indexItems` with kind `indexKindSpec`, the cursor SHALL be positionable on them with `j`/`k`, and `Enter` SHALL open `ModeViewingSpec` for the spec under the cursor.

#### Scenario: Cursor entra en la sección Specs
- **WHEN** the cursor is on the last navigable item before the Specs section (last archived item) and the user presses `j`
- **THEN** the cursor advances to the first spec in the "Specs" section

#### Scenario: Enter sobre un spec activa la visualización
- **WHEN** the index cursor is on a spec and the user presses `Enter`
- **THEN** the TUI enters `ModeViewingSpec` showing the content of that spec

#### Scenario: Cursor no sobrepasa el último spec
- **WHEN** the cursor is on the last spec and the user presses `j`
- **THEN** the cursor does not move (existing boundary behavior)

### Requirement: Expandir spec para ver sus requirements
In `ModeIndex`, pressing `Space` on a spec item SHALL toggle its expanded/collapsed state. When expanded, the names of its requirements SHALL appear indented below the spec, one per line, as navigable items. Pressing `Space` again SHALL collapse the list. The expanded state of each spec SHALL be independent of the cursor and SHALL reset when leaving and re-entering `ModeIndex`.

#### Scenario: Expandir un spec
- **WHEN** the cursor is on a spec item and the user presses `Space`
- **THEN** the requirement names of the spec appear indented below it as navigable items

#### Scenario: Colapsar un spec expandido
- **WHEN** the cursor is on an expanded spec item and the user presses `Space`
- **THEN** the requirement items disappear and the spec is shown again as a single line

#### Scenario: Múltiples specs expandidos simultáneamente
- **WHEN** the user expands two different specs
- **THEN** both show their requirements simultaneously and independently

#### Scenario: Estado expandido no afecta al cursor de otro spec
- **WHEN** spec A is expanded and the user moves the cursor to spec B and presses `Space`
- **THEN** spec B expands without collapsing spec A

#### Scenario: Expand state se resetea al re-entrar al índice
- **WHEN** the user has expanded specs, leaves `ModeIndex`, and re-enters it
- **THEN** all specs appear collapsed

### Requirement: Requirements como ítems navegables en el índice
The requirements of an expanded spec SHALL be navigable items in `indexItems` with kind `indexKindRequirement`. The cursor SHALL be able to move through them with `j`/`k` continuously alongside the rest of the index items. Requirement items SHALL NOT be selectable with `Space`.

#### Scenario: j/k atraviesa requirements de un spec expandido
- **WHEN** spec A is expanded and the cursor is on the last requirement of A
- **THEN** pressing `j` moves the cursor to the next item (spec B or the first requirement of B if expanded)

#### Scenario: Space sobre un requirement no tiene efecto
- **WHEN** the cursor is on a requirement item and the user presses `Space`
- **THEN** no change occurs in the index

### Requirement: Cursor snap al colapsar spec con cursor interior
When the user collapses a spec whose cursor was positioned on one of its requirements, the cursor SHALL jump to the collapsed spec item.

#### Scenario: Colapsar spec con cursor en un requirement
- **WHEN** the cursor is on a requirement of spec A and the user moves the cursor to spec A and presses `Space` to collapse
- **THEN** the cursor is positioned on the spec A item

#### Scenario: Colapsar spec sin cursor interior no mueve el cursor
- **WHEN** the cursor is on an item other than spec A and its requirements, and the user presses `Space` on spec A
- **THEN** the cursor remains on the item where it was
