## Why

`internal/ui/model.go` has grown to 1527 lines — 75% of the entire codebase — mixing eight distinct responsibilities in a single file. Navigation and cognitive load suffer: finding any function requires scrolling or symbol search through an undifferentiated wall of code.

## What Changes

- Split `internal/ui/model.go` into six focused files within the same `ui` package
- No API changes, no new interfaces, no behavior changes — purely file reorganization

### New files

- `model.go` (~250 lines): Model struct, type declarations, message types, `New`, `NewSinglePath`, `Init`, core state query helpers
- `update.go` (~340 lines): `Update` function — the full Elm-style switch, kept together
- `viewport.go` (~130 lines): `loadViewport` — async render dispatch
- `index.go` (~380 lines): all `ModeIndex` logic — `handleTick`, `buildIndexItems`, `renderIndexContent`, `renderActiveItem`, `renderArchivedItem`
- `tasks.go` (~220 lines): all task tab logic — `moveCursor`, `doToggle`, `renderTasksContent`, `progressBar`
- `view.go` (~200 lines): shared rendering — `View`, `renderHeader`, `renderTabBar`, box helpers, `renderHelpBar`

`styles.go` is untouched.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

None. This change touches only code organization, not observable behavior or requirements.

## Impact

- `internal/ui/` package only
- No callers change (`cmd/dossier/main.go` imports the package, not specific files)
- All existing tests remain valid without modification

## Non-goals

- Changing behavior, key bindings, or rendering output
- Introducing sub-packages or new abstractions
- Refactoring the `Update` function internals
