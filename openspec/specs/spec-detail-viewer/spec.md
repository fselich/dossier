# spec-detail-viewer Specification

## Purpose
Allows viewing the content of a project spec (`openspec/specs/<name>/spec.md`) rendered as markdown directly from the index, in read-only mode with scrolling.

## Requirements

### Requirement: Modo de visualización de spec
The TUI SHALL implement a `ModeViewingSpec` mode that displays the content of `openspec/specs/<name>/spec.md` rendered as markdown in a read-only viewport. The mode SHALL be activated by pressing `Enter` on a spec in `ModeIndex`. `Esc` SHALL return to `ModeIndex`. There SHALL be no editing, tabs, or subnav in this mode.

#### Scenario: Abrir spec desde el índice
- **WHEN** the index cursor is on a spec and the user presses `Enter`
- **THEN** the TUI enters `ModeViewingSpec` and shows the content of the selected spec's `spec.md` rendered as markdown

#### Scenario: Scroll del contenido
- **WHEN** the mode is `ModeViewingSpec` and the user presses `j` or `k`
- **THEN** the viewport scrolls down or up respectively

#### Scenario: Volver al índice
- **WHEN** the mode is `ModeViewingSpec` and the user presses `Esc`
- **THEN** the TUI returns to `ModeIndex` with the cursor on the spec that was being viewed

#### Scenario: Header en modo visualización de spec
- **WHEN** the mode is `ModeViewingSpec`
- **THEN** the header shows `<project>  ·  <spec-name>  [spec]`

#### Scenario: HelpBar en modo visualización de spec
- **WHEN** the mode is `ModeViewingSpec`
- **THEN** the help bar shows `j/k: scroll  Esc: index  q: quit`

### Requirement: Campo Content en ProjectSpec
`openspec.ProjectSpec` SHALL include a `Content string` field with the raw content of `spec.md`. `LoadProjectSpecs()` SHALL read and store this content when loading specs.

#### Scenario: Content poblado
- **WHEN** `LoadProjectSpecs()` processes a spec with `spec.md` present
- **THEN** the returned `ProjectSpec` has `Content` with the full text of the file

#### Scenario: Content vacío si spec.md ausente
- **WHEN** a spec directory does not contain `spec.md`
- **THEN** `Content` is an empty string and the spec still appears in the list

### Requirement: Abrir spec con scroll a un requirement específico
When `ModeViewingSpec` is opened from an `indexKindRequirement` item, the TUI SHALL activate focus mode and render only the block of that requirement in the viewport, instead of showing the full spec scrolled to that requirement.

#### Scenario: Abrir spec desde un requirement item
- **WHEN** the index cursor is on a requirement item and the user presses `Enter`
- **THEN** the TUI enters `ModeViewingSpec` in focus mode and the viewport shows only the content of that requirement

#### Scenario: Abrir spec desde el item del spec (sin requirement target)
- **WHEN** the index cursor is on a spec item (not a requirement) and the user presses `Enter`
- **THEN** the TUI enters `ModeViewingSpec` showing the full spec from the beginning (existing behavior)

#### Scenario: Esc desde spec abierto vía requirement vuelve al índice
- **WHEN** `ModeViewingSpec` was opened from a requirement item and the user presses `Esc`
- **THEN** the TUI returns to `ModeIndex` with the cursor on the requirement item from which it was opened
