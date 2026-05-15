# specs-subnav Specification

## Purpose
Implements the sub-navigation chip row visible in the `specs` tab, which displays the available specs of a change and allows cycling through them with the `3` key.

## Requirements


### Requirement: Sub-navegación de specs
When the active tab is `specs` and at least one spec is available, the TUI SHALL display a chip row as the first line inside the content block, immediately after the horizontal separator (`├───┤`), with the name of each spec. The chip of the currently visible spec SHALL be shown with active style (same as an active tab). The other chips SHALL be shown with inactive style. The row SHALL be static (it is not part of the scrollable viewport area).

#### Scenario: Un solo spec
- **WHEN** the change has a single spec and the active tab is `specs`
- **THEN** the chip row is shown as the first line of the content block, with one chip representing that spec marked as active

#### Scenario: Múltiples specs
- **WHEN** the change has two or more specs and the active tab is `specs`
- **THEN** the chip row is shown as the first line of the content block; the visible spec's chip is active and the others are inactive

#### Scenario: Fila ausente en otras tabs
- **WHEN** the active tab is not `specs`
- **THEN** no specs chip row is shown in the content block

#### Scenario: La fila no desaparece al hacer scroll
- **WHEN** the active tab is `specs` and the user scrolls down
- **THEN** the chip row remains visible as the first line of the content block

#### Scenario: Separación visual entre tab bar y chips de specs
- **WHEN** the active tab is `specs`
- **THEN** the horizontal separator `├───┤` appears between the tab bar and the specs chip row

### Requirement: Ciclo de specs con tecla 3
The `3` key SHALL have dual behavior: if the active tab is not `specs`, it switches to it showing the last selected spec (or the first if there is no previous selection). If the active tab is already `specs`, it SHALL advance to the next spec in the list, wrapping back to the first when the last is reached.

#### Scenario: Entrar a specs desde otra tab
- **WHEN** the active tab is `proposal` and the user presses `3`
- **THEN** the active tab becomes `specs` and the previously selected spec is shown (or the first one)

#### Scenario: Ciclar al siguiente spec
- **WHEN** the active tab is `specs`, there are 3 specs and the active spec is the second, and the user presses `3`
- **THEN** the active spec becomes the third and the viewport shows its content

#### Scenario: Ciclar desde el último spec vuelve al primero
- **WHEN** the active tab is `specs`, the active spec is the last one, and the user presses `3`
- **THEN** the active spec becomes the first one

#### Scenario: Un solo spec, pulsar 3 no cambia nada visible
- **WHEN** the active tab is `specs` and there is only one spec, and the user presses `3`
- **THEN** the active spec remains the same and the content does not change

### Requirement: Ajuste de altura de contenido con sub-nav visible
When the specs sub-nav is visible, the content area SHALL reduce its height by 1 line to accommodate the extra row, preventing the viewport from overflowing outside the box.

#### Scenario: Altura reducida en tab specs
- **WHEN** the active tab is `specs` and specs are available
- **THEN** the viewport has 1 fewer line of height than in the other tabs

#### Scenario: Altura normal en otras tabs
- **WHEN** the active tab is `proposal`, `design`, or `tasks`
- **THEN** the viewport has the standard height
