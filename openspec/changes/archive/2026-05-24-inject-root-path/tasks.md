## 1. Add `*From` variants to loader.go

- [x] 1.1 Implement `LoadFrom(root string) (*Project, error)` — loads project structure relative to `root`, refactor `Load()` to delegate to `LoadFrom(os.Getwd())`
- [x] 1.2 Implement `LoadConfigFrom(root string) (ProjectConfig, error)` — reads `openspec/config.yaml` from `root`, returns error on YAML parse failures instead of silent zero-value
- [x] 1.3 Implement `LoadProjectSpecsFrom(root string) ([]ProjectSpec, error)` — reads specs from `root/openspec/specs/`
- [x] 1.4 Implement `ListChangeNamesFrom(root string) ([]string, error)` — lists active change names from `root/openspec/changes/`
- [x] 1.5 Implement `ListArchiveChangesFrom(root string) ([]Change, error)` — loads archived changes from `root/openspec/changes/archive/`
- [x] 1.6 Implement `ListArchiveNamesFrom(root string) ([]string, error)` — lists archive dir names from `root`
- [x] 1.7 Implement `ListSpecNamesFrom(root string) ([]string, error)` — lists spec dir names from `root`

## 2. Update zero-argument wrapper signatures

- [x] 2.1 Change `LoadConfig()` signature to `(ProjectConfig, error)` — wrapper delegates to `LoadConfigFrom` with `os.Getwd()`
- [x] 2.2 Change `LoadProjectSpecs()` signature to `([]ProjectSpec, error)` — wrapper delegates to `LoadProjectSpecsFrom`
- [x] 2.3 Change `ListArchiveChanges()` signature to `([]Change, error)` — wrapper delegates to `ListArchiveChangesFrom`
- [x] 2.4 Change `ListArchiveNames()` signature to `([]string, error)` — wrapper delegates to `ListArchiveNamesFrom`
- [x] 2.5 Change `ListSpecNames()` signature to `([]string, error)` — wrapper delegates to `ListSpecNamesFrom`
- [x] 2.6 Change `ListChangeNames()` signature to `([]string, error)` — wrapper delegates to `ListChangeNamesFrom`

## 3. Update callers

- [x] 3.1 Update `cmd/dossier/main.go`: handle new error returns from `LoadConfig()`, store root path, pass to model constructor
- [x] 3.2 Add `root string` field to `ui.Model` struct
- [x] 3.3 Update `ui.New()` and `ui.NewSinglePath()` to accept and store `root`
- [x] 3.4 Update `internal/ui/index.go` polling code: switch `ListChangeNames()`, `ListSpecNames()`, `ListArchiveNames()` calls to `*From(m.root)` variants with error handling

## 4. Add tests for loader.go

- [x] 4.1 Create `internal/openspec/loader_test.go` with a `setupTestDir(t)` helper that creates a temp directory with the openspec structure
- [x] 4.2 Test `LoadFrom`: valid project with changes, missing `openspec/`, empty changes dir
- [x] 4.3 Test `LoadConfigFrom`: valid config, missing file (returns empty + nil), malformed YAML (returns error)
- [x] 4.4 Test `LoadProjectSpecsFrom`: specs dir with subdirectories, missing `specs/` dir, specs without `spec.md`
- [x] 4.5 Test `ListChangeNamesFrom` / `ListArchiveChangesFrom` / `ListArchiveNamesFrom` / `ListSpecNamesFrom`: with entries, empty, missing dirs
- [x] 4.6 Verify `go test -race -count=1 ./internal/openspec/` passes with no race conditions
