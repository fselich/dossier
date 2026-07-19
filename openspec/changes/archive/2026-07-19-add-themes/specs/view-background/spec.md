## MODIFIED Requirements

### Requirement: Background is configurable per view
The system SHALL support a per-view background color configuration via a Theme struct on the Model. The background color SHALL be a `color.Color` value. Setting it to `nil` or `lipgloss.NoColor` SHALL disable the background (terminal default fallback). The background color SHALL be populated from the active built-in theme selected via `--theme` flag.

#### Scenario: Different views have different backgrounds
- **WHEN** a future theme configures the index view with one background color and the normal view with another
- **THEN** switching between modes displays the corresponding background color for each view

#### Scenario: Built-in theme provides background color
- **WHEN** the user selects a built-in theme (e.g., `--theme light`)
- **THEN** the view renders with the theme's defined background color

### Requirement: Terminal default fallback
When the view background color is nil or unset (either because the theme defines no background or no theme is loaded), the system SHALL render exactly as it does today: all areas not covered by styled segments display the terminal emulator's default background. No post-processing or wrapping SHALL occur.

#### Scenario: No background configured
- **WHEN** the view background color is nil or unset
- **THEN** the rendering pipeline passes through the view string unchanged, producing the same visual output as the current implementation
