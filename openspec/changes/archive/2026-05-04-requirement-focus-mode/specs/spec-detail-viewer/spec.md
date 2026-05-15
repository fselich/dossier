## MODIFIED Requirements

### Requirement: Open spec scrolled to a specific requirement
When `ModeViewingSpec` is opened from an `indexKindRequirement` item, the TUI SHALL activate focus mode and render only that requirement's block in the viewport, instead of showing the full spec scrolled to that requirement.

#### Scenario: Open spec from a requirement item
- **WHEN** the index cursor is on a requirement item and the user presses `Enter`
- **THEN** the TUI enters `ModeViewingSpec` in focus mode and the viewport shows only the content of that requirement

#### Scenario: Open spec from the spec item (no requirement target)
- **WHEN** the index cursor is on a spec item (not a requirement) and the user presses `Enter`
- **THEN** the TUI enters `ModeViewingSpec` showing the full spec from the beginning (existing behaviour)

#### Scenario: Esc from a spec opened via a requirement returns to the index
- **WHEN** `ModeViewingSpec` was opened from a requirement item and the user presses `Esc`
- **THEN** the TUI returns to `ModeIndex` with the cursor on the requirement item from which it was opened
