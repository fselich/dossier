## 1. Tooling

- [x] 1.1 Add `.golangci.yml` with standard linters
- [x] 1.2 Complete `Makefile` with `test`, `lint`, and `fmt` targets
- [x] 1.3 Run `make fmt` and `make lint` to verify clean baseline

## 2. Filesystem interface

- [x] 2.1 Define `fileSystem` interface in `internal/ui/fs.go` with ReadFile, WriteFile, ReadDir, Stat
- [x] 2.2 Define `osFS` adapter in `internal/openspec/` satisfying the `fileSystem` interface
- [x] 2.3 Create `Loader` struct in `internal/openspec/` with injected `fileSystem` field
- [x] 2.4 Convert `LoadFrom`, `LoadConfigFrom`, `LoadProjectSpecsFrom`, `ListChangeNamesFrom`, `ListArchiveChangesFrom`, `ListArchiveNamesFrom`, `ListSpecNamesFrom`, `LoadFromPath`, `ReloadChange`, `ToggleTask` to methods on `*Loader`
- [x] 2.5 Add backward-compatible package-level wrapper functions delegating to `defaultLoader`
- [x] 2.6 Wire `Loader` in `cmd/dossier/main.go` and `internal/ui/model.go`
- [x] 2.7 Run `go test -race ./...` to verify no regressions

## 3. Code style fixes

- [x] 3.1 Eliminate naked return in `taskCounts` (`internal/ui/index.go:491`), use explicit `return done, total`
- [x] 3.2 Preallocate `specs` slice in `loader.go:132`: `make([]ProjectSpec, 0, len(entries))`
- [x] 3.3 Preallocate `items` slice in `tasks.go:31`: `make([]TaskItem, 0, len(lines))`
- [x] 3.4 Preallocate `parts` slice in `view.go:91`: `make([]string, 0, 4)`
- [x] 3.5 Add comment documenting intentionally ignored error at `loader.go:326` (`_ = yaml.Unmarshal(...)`)
- [x] 3.6 Extract named layout constants in `model.go:310-327` replacing magic numbers for `contentHeight()`
- [x] 3.7 Run `go test -race ./...` to verify no regressions

## 4. DRY: merge reload logic

- [x] 4.1 Extract `mergeReloadedChange(fresh Change)` method in `internal/ui/`, used by both `editorReturnMsg` handler and `pollNormalModeContent`
- [x] 4.2 Replace duplicated code in `update.go:65-93` and `index.go:132-176` with calls to `mergeReloadedChange`
- [x] 4.3 Run `go test -race ./...` to verify no regressions

## 5. Per-mode handler files

- [x] 5.1 Extract `updateViewer()` from `handleKeyPress` for ModeNormal + ModeViewingArchive into new `internal/ui/viewer.go`
- [x] 5.2 Extract `updateIndex()` from `handleKeyPress` into `internal/ui/index.go`
- [x] 5.3 Extract `updateSpec()` from `handleKeyPress` into new `internal/ui/spec.go`
- [x] 5.4 Extract `updateConfig()` from `handleKeyPress` into new `internal/ui/config.go`
- [x] 5.5 Refactor `update.go` `Update()` into thin dispatcher: global keys + `switch m.mode` delegating to per-mode methods
- [x] 5.6 Remove monolithic `handleKeyPress()` from `update.go`
- [x] 5.7 Move `handleMouseWheel` and `handleMouseClick` to `mouse.go` (already there)

## 6. Keymap + help bar integration (skipped — low priority, current help bar works)

- [ ] ~~6.1 Add `help` bubble to `Model` struct dependencies~~ (skipped)
- [ ] ~~6.2 Replace manual `renderHelpBar()` with `m.help.ShortHelpView()`~~ (skipped)
- [ ] ~~6.3 Update per-mode update methods to enable/disable `key.Binding` values~~ (skipped)
- [ ] ~~6.4 Run `go test -race ./...` and `make lint` to verify~~ (skipped)

## 7. Error surfacing

- [x] 7.1 Replace `log.Printf` for archive/spec load errors in `index.go` with `m.errMsg` assignments
- [x] 7.2 Replace `log.Printf` for archive/spec load errors in `model.go` (`New` constructor) with `m.errMsg` assignments
- [x] 7.3 Verify errors display temporarily in help bar and auto-clear after 3 seconds

## 8. Final verification

- [x] 8.1 Run full test suite: `make test`
- [x] 8.2 Run full lint: `make lint`
- [x] 8.3 Manual smoke test
: run `./dossier` in a project with openspec, verify all keybindings, mode transitions, and error display
