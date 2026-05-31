## 1. Extract Helper

- [x] 1.1 Extract `loadChangeFromPath(dirPath, entryName string, isArchived bool) (Change, error)` in `loader.go` from the duplicated construction logic
- [x] 1.2 Refactor `LoadFrom` to call `loadChangeFromPath` (replace lines 86–97)
- [x] 1.3 Refactor `LoadFromPath` to call `loadChangeFromPath` (replace lines 324–335)
- [x] 1.4 Refactor `ListArchiveChangesFrom` to call `loadChangeFromPath` (replace lines 206–214)

## 2. Verification

- [x] 2.1 Verify tests pass — run `go test ./internal/openspec/...`
