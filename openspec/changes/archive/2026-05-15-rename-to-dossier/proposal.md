## Why

The project name `spec-viewer` / binary `specview` was a working title that described the tool mechanically. The new name `dossier` better captures the concept: a structured collection of documents about a unit of work, which maps directly to what an OpenSpec change is.

## What Changes

- Binary renamed from `specview` to `dossier`
- Entry point directory renamed from `cmd/specview/` to `cmd/dossier/`
- Go module path updated from `github.com/fselich/dossier` to `github.com/fselich/dossier`
- All internal import paths updated accordingly
- `Makefile` updated to build and install `dossier`
- `.gitignore` updated to ignore the `dossier` binary
- `README.md` and `README.es.md` updated throughout
- Existing OpenSpec specs updated to remove all references to `spec-viewer` and `specview`

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `build-tooling`: Binary name and entry point directory change from `specview`/`cmd/specview/` to `dossier`/`cmd/dossier/`
- `path-arg`: Binary name in all scenarios changes from `spec-viewer` to `dossier`
- `openspec-loader`: Binary name in all scenarios changes from `spec-viewer` to `dossier`

## Impact

- `go.mod`: module declaration
- `cmd/dossier/main.go`: import paths (previously `cmd/specview/main.go`)
- `internal/ui/model.go`: import path
- `Makefile`: BIN and CMD variables
- `.gitignore`: ignored binary name
- `README.md`, `README.es.md`: all name references
- `openspec/specs/build-tooling/spec.md`: binary name and paths in requirements
- `openspec/specs/path-arg/spec.md`: binary name in all scenarios
- `openspec/specs/openspec-loader/spec.md`: binary name in scenarios
