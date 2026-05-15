## ADDED Requirements

### Requirement: Open spec scrolled to a specific requirement
When `ModeViewingSpec` is opened from an `indexKindRequirement` item, the viewport SHALL be positioned at the line of the corresponding requirement rather than at the start of the document.

#### Scenario: Open spec from a requirement item
- **WHEN** the index cursor is on a requirement item and the user presses `Enter`
- **THEN** the TUI enters `ModeViewingSpec` and the viewport scrolls to the section of that requirement

#### Scenario: Open spec from the spec item (without requirement target)
- **WHEN** the index cursor is on a spec item (not a requirement) and the user presses `Enter`
- **THEN** the TUI enters `ModeViewingSpec` showing the start of the document (existing behavior)

#### Scenario: Esc from spec opened via requirement returns to the index
- **WHEN** `ModeViewingSpec` was opened from a requirement item and the user presses `Esc`
- **THEN** the TUI returns to `ModeIndex` with the cursor on the requirement item from which it was opened
