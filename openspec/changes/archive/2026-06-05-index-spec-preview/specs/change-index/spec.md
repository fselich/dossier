## ADDED Requirements

### Requirement: Vista de índice con spec preview bar
The TUI SHALL add a persistent 1-line spec preview bar to the index chrome. The bar SHALL appear between the content area (viewport) and the helpbar, separated from both by horizontal separator lines. The bar SHALL always occupy 1 line in the chrome, regardless of whether the cursor is on a spec item or not. The viewport content height (`contentHeight`) SHALL be reduced by 1 to accommodate the bar.

#### Scenario: Preview bar visible in index
- **WHEN** the mode is `ModeIndex`
- **THEN** the screen shows a 1-line bar between the content area separator and the helpbar separator

#### Scenario: Preview bar not present in other modes
- **WHEN** the mode is `ModeNormal`, `ModeViewingArchive`, `ModeViewingSpec`, or `ModeViewingConfig`
- **THEN** the spec preview bar is not rendered

### Requirement: Preview bar content
When the cursor is on an `indexKindSpec` or `indexKindRequirement` item in the index, the bar SHALL show the spec name followed by ` ┊ ` and the purpose text. The purpose text SHALL be the plain text extracted from between `## Purpose` and the next `##` heading in the spec's `Content`, with markdown syntax stripped. The combined text SHALL be truncated with `…` if it exceeds the available width (`m.width - 2`). When the cursor is on any other item type (active change, archived change), the bar SHALL be empty.

#### Scenario: Cursor on spec shows name and purpose
- **WHEN** the cursor is on a spec item
- **THEN** the bar shows `<spec-name> ┊ <purpose text>` truncated to the terminal width

#### Scenario: Cursor on requirement shows same spec info
- **WHEN** the cursor is on a requirement item within a spec
- **THEN** the bar shows the same content as if the parent spec were selected

#### Scenario: Cursor on active or archived item shows empty bar
- **WHEN** the cursor is on an active change or an archived change
- **THEN** the bar shows no content (empty line)
