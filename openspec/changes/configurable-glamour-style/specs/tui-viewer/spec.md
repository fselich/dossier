## MODIFIED Requirements

### Requirement: Render de markdown con glamour
The TUI SHALL render `proposal`, `design`, and `specs` artifacts using glamour with the width of the content area. The content area SHALL be scrollable with `j`/`k` or the arrow keys. The Glamour standard style SHALL default to `dark`. If `DOSSIER_GLAMOUR_STYLE` is set, the TUI SHALL use its value as the Glamour standard style when creating the markdown renderer.

#### Scenario: Scroll en contenido largo
- **WHEN** the artifact has more content than the screen height and the user presses `j`
- **THEN** the content scrolls down one line

#### Scenario: Wrap de glamour ajustado al ancho
- **WHEN** the terminal is 80 columns wide
- **THEN** glamour renders the markdown without exceeding those 80 columns

#### Scenario: Default Glamour style
- **WHEN** `DOSSIER_GLAMOUR_STYLE` is unset
- **THEN** the markdown renderer uses Glamour's `dark` standard style

#### Scenario: Environment-selected Glamour style
- **WHEN** `DOSSIER_GLAMOUR_STYLE` is set to `light`
- **THEN** the markdown renderer uses Glamour's `light` standard style

### Requirement: Tabs de artifact
The TUI SHALL show a tab bar with tabs `proposal`, `design`, `tasks`, `specs`, and `code` when the git tab is available. Tabs for absent artifacts SHALL be shown visually disabled and not selectable. Available inactive tabs SHALL remain readable on both light and dark terminal palettes. When the tab bar includes a task progress indicator, that indicator SHALL be compact and SHALL NOT consume all remaining horizontal space on wide terminals.

#### Scenario: Available inactive tabs remain readable
- **WHEN** the terminal uses a light color palette and an available tab is inactive
- **THEN** the tab label is rendered with a visible foreground color

#### Scenario: Top progress indicator is compact
- **WHEN** tasks exist and the terminal is wide
- **THEN** the tab-bar progress indicator uses a bounded width instead of filling all remaining columns
