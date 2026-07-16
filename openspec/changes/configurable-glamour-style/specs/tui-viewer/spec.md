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
