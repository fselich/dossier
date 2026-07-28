## ADDED Requirements

### Requirement: Tasks cursor supports page-jump navigation
The system SHALL move the tasks tab cursor by a screen's worth of rendered lines (matching the current viewport height) when the user presses `<pgdown>` or `<pgup>`, rather than treating those keys as a no-op.

#### Scenario: Page-down jumps the cursor forward by one screen
- **WHEN** the user presses `<pgdown>` while the tasks tab is in the live (non-archive) view
- **THEN** the cursor SHALL move to the item whose rendered line is approximately one viewport height below the cursor's current rendered line

#### Scenario: Page-up jumps the cursor backward by one screen
- **WHEN** the user presses `<pgup>` while the tasks tab is in the live (non-archive) view
- **THEN** the cursor SHALL move to the item whose rendered line is approximately one viewport height above the cursor's current rendered line

#### Scenario: Page-down near the end of the list clamps to the last item
- **WHEN** the user presses `<pgdown>` and fewer than one screen's worth of lines remain below the cursor
- **THEN** the cursor SHALL move to the last item in the list

#### Scenario: Page-up near the start of the list clamps to the first item
- **WHEN** the user presses `<pgup>` and fewer than one screen's worth of lines remain above the cursor
- **THEN** the cursor SHALL move to the first item in the list

### Requirement: Archived tasks support page scrolling
The system SHALL scroll the tasks viewport by a full page when the user presses `<pgdown>` or `<pgup>` while viewing an archived (read-only) change's tasks tab, matching the page-scroll behavior of every other document tab.

#### Scenario: Page-down scrolls archived tasks
- **WHEN** the user presses `<pgdown>` while viewing an archived change's tasks tab
- **THEN** the viewport SHALL scroll down by one page

#### Scenario: Page-up scrolls archived tasks
- **WHEN** the user presses `<pgup>` while viewing an archived change's tasks tab
- **THEN** the viewport SHALL scroll up by one page
</content>
