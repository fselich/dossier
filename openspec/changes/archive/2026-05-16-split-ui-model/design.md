## Context

`internal/ui/model.go` is 1527 lines — 75% of the entire project codebase. It mixes type declarations, Elm-architecture lifecycle, viewport management, index view logic, task view logic, and shared rendering in a single file. The file compiles and works correctly; the problem is purely navigational and organizational.

The `ui` package has one other file (`styles.go`, 58 lines). Splitting into multiple files within the same package requires zero interface or API changes.

## Goals / Non-Goals

**Goals:**
- Split `model.go` into six focused files, each covering one semantic area
- No behavior changes, no API changes, no new abstractions
- Files should be independently navigable: a developer looking for task rendering knows to open `tasks.go`

**Non-Goals:**
- Changing the Elm-architecture structure (Model/Update/View)
- Introducing sub-packages or new interfaces
- Refactoring internals of any function
- Changing `styles.go`

## Decisions

### Keep `Update` in a single file

The `Update` function (~340 lines) is a switch over message types and key bindings. It is a routing function — its value is showing all cases in one place. Splitting it by mode (e.g., `update_index.go`, `update_normal.go`) would make it harder to trace what happens on a given keypress.

**Alternatives considered:**
- Split by mode: rejected — forces the reader to jump files to understand a single key binding
- Split by message type (key vs non-key): rejected — artificial boundary with no semantic gain

### Split by semantic area, not by type category

Each file owns a complete concern: all of the index view (data + rendering), all of the task view (data + rendering), etc. This is preferable to splitting by layer (e.g., "all rendering functions" in one file) because changes tend to touch one concern at a time.

### Stay in the same package (`ui`)

No sub-packages. Moving functions to `ui/index`, `ui/tasks`, etc. would require exporting everything and introducing import cycles. Same-package file splitting is free.

## File Map

| File | ~Lines | Contents |
|------|--------|----------|
| `model.go` | 250 | Model struct, type/const declarations, message types, `New`, `NewSinglePath`, `Init`, state query helpers (`current`, `tabAvailable`, `defaultTab`, `artifactPath`, `contentHeight`, …) |
| `update.go` | 340 | `Update` — complete key and message switch |
| `viewport.go` | 130 | `loadViewport` — async glamour render dispatch |
| `index.go` | 380 | `handleTick`, `enterIndex`, `buildSpecOrder`, `buildIndexItems`, `refreshIndexViewport`, `renderIndexContent`, `renderActiveItem`, `renderArchivedItem`, `taskCounts`, `viewIndex`, `sameNames`, `sameStrings`, `extractRequirement`, `specSuffix` |
| `tasks.go` | 220 | `loadTaskItems`, `firstTaskIdx`, `moveCursorDown`, `moveCursorUp`, `doToggle`, `refreshTasksViewport`, `renderTasksContent`, `sectionProgress`, `progressBar`, `extractOpeningEscape`, `inlineMarkdown` |
| `view.go` | 200 | `View`, `renderHeader`, `renderTabBar`, `renderSpecSubnav`, `hasSpecSubnav`, `boxTop`, `boxBottom`, `boxInnerSep`, `addBorderSides`, `renderHelpBar`, `emptyView` |

## Risks / Trade-offs

- **Risk**: A function moved to the wrong file creates a misleading map. → Mitigation: follow the grouping in this doc exactly; verify with `go build` after each move.
- **Risk**: Merge conflicts if other branches modify `model.go` concurrently. → Mitigation: this branch should be short-lived; land it before branching for other UI work.
- **Trade-off**: Six files instead of one means more files to open. Accepted — each file is now small enough to read in full.
