## Why

Seven functions in `internal/openspec` hardcode `os.Getwd()`, coupling domain logic to global process state. This prevents testing with arbitrary paths, forces test code to use `os.Chdir` (a race condition hazard in parallel tests), and blocks usage of the package from a library or API context.

## What Changes

- Add `From(root string)` variants for all seven functions that currently call `os.Getwd()` internally:
  - `LoadFrom(root string) (*Project, error)`
  - `LoadConfigFrom(root string) ProjectConfig`
  - `LoadProjectSpecsFrom(root string) []ProjectSpec`
  - `ListArchiveChangesFrom(root string) []Change`
  - `ListArchiveNamesFrom(root string) []string`
  - `ListSpecNamesFrom(root string) []string`
  - `ListChangeNamesFrom(root string) []string`
- Refactor existing zero-argument functions to call their `*From` counterparts with `os.Getwd()`, preserving the public API. **BREAKING**: return `error` from functions that currently swallow it (`LoadConfig`, `LoadProjectSpecs`, `ListArchiveChanges`, `ListArchiveNames`, `ListSpecNames`, `ListChangeNames`).
- Update call sites in `cmd/dossier/main.go` and `internal/ui/` to use the `*From` variants directly when the root path is known.

## Capabilities

### New Capabilities

- `openspec-root-path`: All openspec loader functions accept an explicit root path parameter instead of relying on `os.Getwd()`.

### Modified Capabilities

- `openspec-loader`: Changed function signatures — zero-argument functions now return errors instead of silently swallowing them.

## Impact

- `internal/openspec/loader.go`: All seven functions modified; six new `*From` variants added.
- `cmd/dossier/main.go`: Switch from `Load()`/`LoadConfig()` to direct `*From` calls with known `os.Getwd()` value.
- `internal/ui/index.go` (and any polling code): Switch list-refresh calls to `*From` variants.
- Tests benefit immediately: can test with `t.TempDir()` without `os.Chdir`.
