## ADDED Requirements

### Requirement: BaseText style applies PrimaryFg to unstyled text
The system SHALL provide a `BaseText` style in `ThemeStyles` that applies only `PrimaryFg` (no bold, no background). Non-selected item names in the index view SHALL use `BaseText` to ensure they are readable regardless of the active theme's background color.

#### Scenario: BaseText style uses PrimaryFg
- **WHEN** `BuildStyles(c)` is called with a `ThemeColors` value
- **THEN** `ThemeStyles.BaseText` is `lipgloss.NewStyle().Foreground(c.PrimaryFg)`

#### Scenario: Active change names use BaseText when not selected
- **WHEN** rendering an active change in the index view and the item is not the cursor
- **THEN** the change name is rendered with `m.theme.Styles.BaseText`

#### Scenario: Archived change names use BaseText when not selected
- **WHEN** rendering an archived change in the index view and the item is not the cursor
- **THEN** the change name is rendered with `m.theme.Styles.BaseText`

#### Scenario: Spec names use BaseText when not selected
- **WHEN** rendering a specification name in the index view and the item is not the cursor
- **THEN** the specification name is rendered with `m.theme.Styles.BaseText`

#### Scenario: Light theme non-selected names are readable
- **WHEN** the active theme is `light` (background `#ffffff`) and an item is not the cursor
- **THEN** the item name uses `PrimaryFg` = `"0"` (black), providing contrast against the white background
