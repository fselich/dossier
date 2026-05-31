# change-index Specification (Delta)

## MODIFIED Requirements

### Requirement: Vista índice de pantalla completa

The TUI SHALL implement a `ModeIndex` mode that occupies the full screen with the same TUI chrome (borders, header, helpbar). The index SHALL show three sections: "Active Changes" with active changes, "Specifications" with the project specs in `openspec/specs/`, and "Archived Changes" with changes in `openspec/changes/archive/`; the three separated by a section line. When a view background color is configured, the entire index view SHALL render with that background color filling the full terminal viewport, including all whitespace areas between elements and the empty area below the box frame.

*Implementation note: Index state fields (`indexItems`, `indexCursor`, `expandedSpecs`, `specSortBySuffix`, `specOrder`, `archiveChanges`, `archiveCursor`) are now grouped under `Model.index IndexState`. Scenarios unchanged.*

### Requirement: Navegación en el índice

The cursor SHALL be able to move through all items (active and archived) with `j` (down) and `k` (up). Section separators are not selectable items. The cursor SHALL NOT go past the first or last item.

*Implementation note: All index cursor logic now accesses `m.index.Cursor` (formerly `m.indexCursor`).*

### Requirement: Seleccionar un change con Enter

Pressing `Enter` on an item or left-clicking an already-selected item SHALL close the index and open the selected change. If it is an active change, the mode switches to `ModeNormal` with that active change. If it is an archived change, the mode switches to `ModeViewingArchive` with that archived change.

*Implementation note: Archive changes and cursor now accessed via `m.index.ArchiveChanges` and `m.index.ArchiveCursor`.*

### Requirement: Helpbar del índice

The helpbar in `ModeIndex` SHALL show navigation hints and SHALL reflect the current sort mode via the `s` binding label:
- When sort mode is **name**: `j/k: navigate  Enter: open  Space: expand  click: select  s: sort by suffix  Esc: quit`
- When sort mode is **suffix**: `j/k: navigate  Enter: open  Space: expand  click: select  s: sort by name  Esc: quit`

*Implementation note: Sort state now accessed via `m.index.SpecSortBySuffix`.*

### Requirement: Actualización en tiempo real del índice

While the mode is `ModeIndex`, the TUI SHALL detect on each tick (≤ 500 ms) whether the list of active changes, the list of archived changes, or the list of project specs has changed on disk. If any structural change is detected, the index SHALL reload all three lists, rebuild the navigable items, and refresh the viewport. Additionally, when no structural change is detected, the TUI SHALL reload the task content of each active change from disk and, if any task content has changed, SHALL rebuild the index items and refresh the viewport so that progress bars reflect the latest task completion state. The cursor SHALL be preserved if the resulting index has at least as many items as the current position; otherwise it SHALL move to the last available item.

*Implementation note: All index-related fields now accessed via `m.index`. Behavior unchanged.*

#### Scenario: Nuevo spec aparece en disco mientras el índice está abierto

- **WHEN** the mode is `ModeIndex` and a new directory is created in `openspec/specs/`
- **THEN** within a maximum of 500 ms the index shows the new spec in the "Specifications" section without user intervention

#### Scenario: Tareas actualizadas en disco mientras el índice está abierto

- **WHEN** the mode is `ModeIndex` and the `tasks.md` file of an active change is externally modified
- **THEN** within a maximum of 500 ms the progress bar for that change in the index reflects the updated `done/total` count without user intervention
