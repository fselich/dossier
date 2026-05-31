## Why

The codebase lacks testability (no interfaces abstracting the filesystem), has a monolithic 270-line key handler that mixes concerns across five UI modes, and is missing linter configuration and standard Makefile targets. A structured code review identified 12 actionable improvements spanning architecture, style, tooling, and patterns. Addressing them now ensures maintainability as the project grows.

## What Changes

- **Split `handleKeyPress` into per-mode `update*()` functions**, each in its own file (`viewer.go`, `index.go`, `spec.go`, `config.go`), with `update.go` becoming a thin dispatcher
- **Introduce a `fileSystem` interface** in the UI package for filesystem operations, decoupling `openspec` loader functions from direct `os` calls
- **Add `.golangci.yml`** with errcheck, staticcheck, goimports, misspell, unconvert, unparam
- **Complete the `Makefile`** with `test`, `lint`, and `fmt` targets
- **Eliminate naked return** in `taskCounts` (`internal/ui/index.go:491`)
- **Preallocate slices** where size is known at allocation time (`loader.go:132`, `tasks.go:31`, `view.go:91`)
- **Document intentionally ignored error** in `loader.go:326` (`_ = yaml.Unmarshal(...)`)
- **Use named layout constants** instead of magic numbers (`m.height - 6`, `m.height - 7`)
- **DRY the reload-merge logic** duplicated between `editorReturnMsg` handler and `pollNormalModeContent`
- **Surface poll errors** via `m.errMsg` instead of silent `log.Printf`

## Capabilities

### New Capabilities
- `per-mode-handlers`: Each UI mode has its own update handler in a dedicated file, with `update.go` acting as a thin dispatcher. Keys are defined via `key.Binding` structs (keymap pattern) instead of raw strings.
- `filesystem-interface`: A `fileSystem` interface in `internal/ui/` decouples filesystem access for testability. The `openspec` loader converts to a struct with an injected filesystem dependency.
- `lint-and-makefile`: Project tooling includes `.golangci.yml` with standard linters and a complete `Makefile` with `test`, `lint`, and `fmt` targets.

### Modified Capabilities
- `tui-viewer`: Key handling is restructured from a monolithic switch to per-mode delegation. No user-facing behavior changes.
- `openspec-loader`: Functions become methods on `*Loader` struct. Public API names preserved via wrapper functions. No behavior changes.

## Impact

- Affected code: `internal/ui/` (all files, new per-mode files), `internal/openspec/` (loader.go), `cmd/dossier/` (main.go wiring), `Makefile`, new `.golangci.yml`
- No external API changes — `main.go` wiring is internal-only
- No dependency changes — pure refactor within existing dependencies
- Tests: existing tests continue to pass; new interface enables future in-memory testing without `t.TempDir()`
