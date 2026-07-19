## 1. Load task items on archive entry

- [x] 1.1 Call `loadTaskItems()` in `updateIndex` Enter handler for archived items (`index.go:762`), before `commitStateChange`
- [x] 1.2 Call `loadTaskItems()` in `clickIndexItem` for archived items (`mouse.go:152`), before `loadViewport`

## 2. Use task list viewport for archive tasks tab

- [x] 2.1 Extend `loadViewport` condition in `viewport.go:23` to include `ModeViewingArchive`: `case m.tab == TabTasks && (m.mode == ModeNormal || m.mode == ModeViewingArchive)`

## 3. Block Space and e in archive mode

- [x] 3.1 Add `ModeViewingArchive` guard in `updateViewer` for `Space` handler (`viewer.go:231-234`)
- [x] 3.2 Add `ModeViewingArchive` guard in `updateViewer` for `e` handler (`viewer.go:247-259`)

## 4. Make helpbar tab-aware in archive mode

- [x] 4.1 Update `renderHelpBar` in `view.go` to show `j/k: navigate` when on TabTasks in archive mode, `j/k: scroll` on other tabs

## 5. Tests

- [x] 5.1 Add test for task navigation in archived changes (verify tasks.Items loaded and j/k works)
- [x] 5.2 Add test for Space in archive mode (verify no toggle)
- [x] 5.3 Add test for e in archive mode (verify no editor)
