## 1. Fix error propagation

- [x] 1.1 Fix `LoadConfigFrom` in `internal/openspec/loader.go:126` — change `return ProjectConfig{}, nil` to `return ProjectConfig{}, err` for non-IsNotExist read failures
- [x] 1.2 Fix error swallowing in `internal/ui/index.go` and `internal/ui/model.go` — replace `_` with `log.Printf` for `ListArchiveChangesFrom` and `LoadProjectSpecsFrom` calls (lines 66–67, 130–131, 159–161)
- [x] 1.3 Run `go test ./internal/openspec/` and `go test ./internal/ui/` to verify existing tests pass and new error paths are covered
