# change-index Specification (Delta)

## ADDED Requirements

### Requirement: Filtrar el índice con /

While the mode is `ModeIndex`, pressing `/` SHALL enter a filter typing state where the help bar is replaced by a prompt showing `/` followed by the current query text. While in this typing state, every printable character typed SHALL be appended to the query, every `Backspace` SHALL remove the last character, and the index items SHALL be filtered in real-time using case-insensitive substring matching against:
- The change name for active changes
- The change name for archived changes
- The spec name for specification items
- The requirement name for requirement items

Pressing `Enter` while in the filter typing state SHALL confirm the query: the filter remains applied, the help bar returns to its normal state but with a `[/query]` indicator appended, and the user can navigate the filtered list with the usual keys. Pressing `Esc` while in the filter typing state SHALL cancel the editing and restore the filter state that was active before `/` was pressed (either the previously applied filter or no filter).

While a filter is applied but the user is not typing, `Esc` SHALL clear the filter and show all items. Pressing `/` again SHALL re-enter the filter typing state with the current query pre-filled for editing.

#### Scenario: Pressing / enters filter typing state
- **WHEN** the mode is `ModeIndex` and the user presses `/`
- **THEN** the help bar shows `/` with a visible input cursor and no filtering occurs yet

#### Scenario: Typing filters items in real-time
- **WHEN** the mode is `ModeIndex`, the user presses `/`, then types "foo"
- **THEN** within the same frame, only items whose name contains "foo" (case-insensitive) remain visible, and items without "foo" are hidden

#### Scenario: Type the query, press Enter, filter stays applied
- **WHEN** the mode is `ModeIndex`, the user types "bar" after `/` and presses `Enter`
- **THEN** the filter remains active, the help bar returns with normal bindings plus a `[/bar]` indicator, and only items matching "bar" are shown

#### Scenario: Esc during typing cancels input without changing filter
- **WHEN** the mode is `ModeIndex`, a filter "foo" is active, the user presses `/` to edit, types "bar", then presses `Esc`
- **THEN** the typing is cancelled, the filter reverts to "foo", and items matching "foo" are shown

#### Scenario: Esc with filter active but not typing clears the filter
- **WHEN** the mode is `ModeIndex`, a filter "foo" is active, and the user presses `Esc`
- **THEN** the filter is cleared and all items are shown

#### Scenario: Esc with no filter active quits the application
- **WHEN** the mode is `ModeIndex`, no filter is active, and the user presses `Esc`
- **THEN** the application quits

#### Scenario: Backspace during typing removes last character
- **WHEN** the mode is `ModeIndex`, the user types "foobar" after `/`, then presses `Backspace` twice
- **THEN** the query becomes "foob" and the filter updates accordingly

#### Scenario: / reopens editing with pre-filled query
- **WHEN** the mode is `ModeIndex`, a filter "foo" is active, and the user presses `/`
- **THEN** the filter typing state opens with the query pre-filled as "foo" and the user can edit it

#### Scenario: Case-insensitive matching
- **WHEN** the mode is `ModeIndex`, a change named "MyFeature" exists, and the user types `/myfeature`
- **THEN** "MyFeature" is shown as a matching item

### Requirement: Secciones sin coincidencias muestran mensaje

When a filter is active and a section (Active Changes, Specifications, or Archived Changes) has no items matching the filter, that section SHALL display a message "No items match '<query>'" in help style instead of listing the items of that section. Sections that already have no items without a filter SHALL keep their existing "No active changes" / "No specifications available" / "No archived changes" messages.

#### Scenario: Sección activa sin match muestra mensaje
- **WHEN** the mode is `ModeIndex`, a filter is active, and no active change matches it
- **THEN** the Active Changes section shows "No items match '<query>'" instead of the list of changes

#### Scenario: Otras secciones con matches aún se muestran
- **WHEN** the mode is `ModeIndex`, a filter matches some specs but no active changes and no archived changes
- **THEN** the Active Changes section shows the no-match message, the Specifications section shows matching specs, and the Archived Changes section shows the no-match message

#### Scenario: Sin filtro se mantienen los mensajes originales
- **WHEN** the mode is `ModeIndex`, no filter is active, and there are no active changes
- **THEN** the Active Changes section shows "No active changes" (unchanged)

### Requirement: Cursor preservado al filtrar

When a filter is applied or changed, the cursor SHALL be preserved on the same logical item if it still matches the filter. If the current item no longer matches, the cursor SHALL move to the first item in the filtered list. If no items match the filter, the cursor SHALL be set to 0 (no items selectable).

#### Scenario: Cursor stays on same item when it still matches
- **WHEN** the cursor is on a change named "data-export", the user types `/data`, and the cursor is at position 2
- **THEN** the cursor remains at position 2 if "data-export" is the third item in the filtered list (moves with the indirection)

#### Scenario: Cursor moves to first item when current item is filtered out
- **WHEN** the cursor is on a change named "auth" and the user types `/data`
- **THEN** the cursor moves to the first item in the filtered list (or to 0 if no items match)
