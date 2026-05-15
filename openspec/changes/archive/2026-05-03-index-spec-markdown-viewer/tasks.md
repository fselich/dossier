## 1. Loader

- [x] 1.1 Add `Content string` field to `openspec.ProjectSpec` in `internal/openspec/loader.go`
- [x] 1.2 Read the content of `spec.md` in `LoadProjectSpecs()` and assign it to `Content`

## 2. TUI Model — types and state

- [x] 2.1 Add `indexKindSpec` to `indexItemKind` in `internal/ui/model.go`
- [x] 2.2 Add `ModeViewingSpec` to `Mode` in `internal/ui/model.go`
- [x] 2.3 Add `specViewerCursor int` field to `Model`

## 3. Index construction

- [x] 3.1 Include specs in `buildIndexItems()` with kind `indexKindSpec`

## 4. Navigation and Enter action in the index

- [x] 4.1 In the `enter` handler for `ModeIndex`: detect `indexKindSpec` and enter `ModeViewingSpec` (assign `specViewerCursor`, adjust `vp.Height`, call `loadViewport`)
- [x] 4.2 In `ModeViewingSpec`, handle `Esc` to return to `ModeIndex` (call `enterIndex` with cursor restored to the viewed spec)

## 5. Viewport and render

- [x] 5.1 In `loadViewport()`: add a branch for `ModeViewingSpec` that takes `m.projectSpecs[m.specViewerCursor].Content` and launches the async glamour render
- [x] 5.2 In `contentHeight()`: add a case for `ModeViewingSpec` with the same calculation as `ModeIndex`

## 6. Views — header and helpbar

- [x] 6.1 In `renderHeader()`: add a case for `ModeViewingSpec` that shows `<project>  ·  <spec-name>  [spec]`
- [x] 6.2 In `renderHelpBar()`: add a case for `ModeViewingSpec` with `j/k: scroll  Esc: index  q: quit`

## 7. Index view — cursor rendering in specs

- [x] 7.1 In `renderIndexContent()`: show the cursor (`▶`) on the active spec within the "Specs" section
