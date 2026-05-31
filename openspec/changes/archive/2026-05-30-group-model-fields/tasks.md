## 1. Define IndexState struct

- [x] 1.1 Define `indexState` type in `model.go` with fields: `Items []indexItem`, `Cursor int`, `ExpandedSpecs map[int]bool`, `SortBySuffix bool`, `Order []int`, `ArchiveChanges []openspec.Change`, `ArchiveCursor int`
- [x] 1.2 Add `index indexState` field to `Model`, remove the flat fields (`indexItems`, `indexCursor`, `expandedSpecs`, `specSortBySuffix`, `specOrder`, `archiveChanges`, `archiveCursor`)

## 2. Define SpecViewerState struct

- [x] 2.1 Define `specViewerState` type in `model.go` with fields: `Cursor int`, `JumpTarget string`, `FocusMode bool`, `ReqCursor int`
- [x] 2.2 Add `specViewer specViewerState` field to `Model`, remove the flat fields (`specViewerCursor`, `specJumpTarget`, `specFocusMode`, `specReqCursor`)

## 3. Define TaskState struct

- [x] 3.1 Define `taskState` type in `model.go` with fields: `Items []openspec.TaskItem`, `Cursor int`
- [x] 3.2 Add `tasks taskState` field to `Model`, remove the flat fields (`taskItems`, `taskCursor`)

## 4. Update Model struct

- [x] 4.1 Verify `Model` compiles with the three new embedded fields and all old flat fields removed
- [x] 4.2 Update `New()` and `NewSinglePath()` to reference new field paths

## 5. Update all field references

- [x] 5.1 Update `index.go`: all `m.indexItems` → `m.index.Items`, `m.indexCursor` → `m.index.Cursor`, etc.
- [x] 5.2 Update `update.go`: all affected field references
- [x] 5.3 Update `mouse.go`: all affected field references
- [x] 5.4 Update any other files in `internal/ui/` that reference the flat fields (view.go, etc.)

## 6. Verify tests pass

- [x] 6.1 Run `go build ./...` to confirm no compilation errors
- [x] 6.2 Run `go test ./internal/ui/...` and confirm all tests pass
