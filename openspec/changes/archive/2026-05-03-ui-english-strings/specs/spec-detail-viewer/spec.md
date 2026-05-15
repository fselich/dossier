## MODIFIED Requirements

### Requirement: Spec viewing mode
The TUI SHALL implement a `ModeViewingSpec` mode that displays the content of `openspec/specs/<name>/spec.md` rendered as markdown in the read-only viewport. The mode SHALL be activated by pressing `Enter` on a spec in `ModeIndex`. `Esc` SHALL return to `ModeIndex`. No editing, tabs, or subnav SHALL exist in this mode.

#### Scenario: Open spec from the index
- **WHEN** the index cursor is on a spec and the user presses `Enter`
- **THEN** the TUI enters `ModeViewingSpec` and displays the content of the selected spec's `spec.md` rendered as markdown

#### Scenario: Content scroll
- **WHEN** the mode is `ModeViewingSpec` and the user presses `j` or `k`
- **THEN** the viewport scrolls down or up respectively

#### Scenario: Return to index
- **WHEN** the mode is `ModeViewingSpec` and the user presses `Esc`
- **THEN** the TUI returns to `ModeIndex` with the cursor on the spec that was being viewed

#### Scenario: Header in spec viewing mode
- **WHEN** the mode is `ModeViewingSpec`
- **THEN** the header shows `<project>  ·  <spec-name>  [spec]`

#### Scenario: HelpBar in spec viewing mode
- **WHEN** the mode is `ModeViewingSpec`
- **THEN** the help bar shows `j/k: scroll  Esc: index  q: quit`
