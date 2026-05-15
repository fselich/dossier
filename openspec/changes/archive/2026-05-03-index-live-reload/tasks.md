## 1. Loader — cheap detection functions

- [x] 1.1 Add `ListArchiveNames() []string` in `internal/openspec/loader.go` — lists only subdirectory names from `openspec/changes/archive/`
- [x] 1.2 Add `ListSpecNames() []string` in `internal/openspec/loader.go` — lists only subdirectory names from `openspec/specs/`

## 2. Model — comparison helper

- [x] 2.1 Add function `sameStrings(a, b []string) bool` in `internal/ui/model.go` to compare slices of archive and spec names

## 3. handleTick — polling in ModeIndex

- [x] 3.1 Remove the early-return for `ModeIndex` from the initial guard in `handleTick()`, leaving only `ModeViewingArchive` and `ModeViewingSpec`
- [x] 3.2 Add `if m.mode == ModeIndex` branch in `handleTick()` that compares `ListChangeNames()` against the current names in `m.project.Changes`, `ListArchiveNames()` against the names in `m.archiveChanges`, and `ListSpecNames()` against the names in `m.projectSpecs`
- [x] 3.3 If any change is detected: reload `m.project` with `openspec.Load()`, `m.archiveChanges` with `openspec.ListArchiveChanges()`, and `m.projectSpecs` with `openspec.LoadProjectSpecs()`, call `buildIndexItems()`, clamp the cursor, and call `refreshIndexViewport()`
