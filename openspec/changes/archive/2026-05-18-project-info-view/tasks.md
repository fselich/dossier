## 1. Config loading

- [x] 1.1 Add `ProjectConfig` struct to `internal/openspec/loader.go` with `Context string` and `Rules map[string][]string` fields
- [x] 1.2 Implement `LoadConfig() ProjectConfig` in `internal/openspec/loader.go` that reads and parses `openspec/config.yaml`
- [x] 1.3 Add `projectConfig openspec.ProjectConfig` field to `Model` in `internal/ui/model.go`
- [x] 1.4 Pass loaded config into `Model` from `cmd/dossier/main.go`

## 2. New mode and state

- [x] 2.1 Add `ModeViewingConfig Mode` constant in `internal/ui/model.go`
- [x] 2.2 Add `prevMode Mode` field to `Model` to restore caller's mode on exit
- [x] 2.3 Add `renderedConfigMsg` type (carries Glamour output) alongside existing `renderedMsg` / `specRenderedMsg`

## 3. Rendering

- [x] 3.1 Implement `configToMarkdown(cfg openspec.ProjectConfig) string` helper in `internal/ui/view.go` that builds the `## Context` + `### <key>` rules markdown
- [x] 3.2 Implement `viewConfig()` in `internal/ui/view.go` reusing the `viewIndex` box layout
- [x] 3.3 Extend `renderHeader()` to return `<project-name>  ·  project config` when `m.mode == ModeViewingConfig`
- [x] 3.4 Extend `renderHelpBar()` to return `j/k: scroll  i/Esc: back  q: quit` when `m.mode == ModeViewingConfig`
- [x] 3.5 Wire `viewConfig()` into the `View()` dispatch at the top of `view.go`

## 4. Key handling

- [x] 4.1 In `update.go`, handle `i` in `ModeIndex`: set `prevMode`, enter `ModeViewingConfig`, trigger async Glamour render
- [x] 4.2 In `update.go`, handle `i` in `ModeNormal`: same as above
- [x] 4.3 In `update.go`, handle `Esc` and `q` in `ModeViewingConfig`: restore `m.mode = m.prevMode`, reset viewport
- [x] 4.4 Handle `renderedConfigMsg` in `Update()` to set viewport content

## 5. Tests

- [x] 5.1 Add unit test for `LoadConfig()` covering: valid file, missing file, empty context, missing rules
- [x] 5.2 Add unit test for `configToMarkdown()` covering: context output, rules grouping, empty config
