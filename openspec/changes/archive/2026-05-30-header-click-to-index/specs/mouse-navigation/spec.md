# mouse-navigation Specification (Delta)

## ADDED Requirements

### Requirement: Header click navigates to index

In `ModeNormal` and `ModeViewingArchive`, the TUI SHALL enter `ModeIndex` when the user left-clicks on the header row (screen Y=1). In all other modes, clicking the header row SHALL be ignored. The behavior SHALL be the same as pressing the `a` or `Esc` key.

#### Scenario: Click header in normal mode enters index

- **WHEN** the mode is `ModeNormal` and the user left-clicks at Y=1
- **THEN** the TUI enters `ModeIndex`

#### Scenario: Click header in archive view enters index

- **WHEN** the mode is `ModeViewingArchive` and the user left-clicks at Y=1
- **THEN** the TUI enters `ModeIndex`

#### Scenario: Click header in index mode does nothing

- **WHEN** the mode is `ModeIndex` and the user left-clicks at Y=1
- **THEN** no mode change occurs (already in index)

#### Scenario: Click header in spec viewer does nothing

- **WHEN** the mode is `ModeViewingSpec` and the user left-clicks at Y=1
- **THEN** no mode change occurs
