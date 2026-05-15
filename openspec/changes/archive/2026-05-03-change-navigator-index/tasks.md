## 1. New ModeIndex mode

- [x] 1.1 Add `ModeIndex` to the `Mode` constant in `model.go`
- [x] 1.2 Add fields `indexCursor int` and `indexItems []indexItem` to the `Model` struct (where `indexItem` references an active or archived change)
- [x] 1.3 Implement `enterIndex()`: load archived changes if not already loaded, build the flat item list (active + archived), position cursor at 0

## 2. Index rendering

- [x] 2.1 Implement `renderIndex()`: full screen with TUI chrome (borders, header, helpbar)
- [x] 2.2 Render "Active" section with progress bar `[█░] N/M` per item (reuse `progressBar` logic)
- [x] 2.3 Render "Archived" section with `DD/MM/YYYY  name` format
- [x] 2.4 Visually highlight the item under the cursor with `▶` and active style
- [x] 2.5 Implement index-specific helpbar: `j/k: navigate  Enter: open  Esc: quit`

## 3. Navigation in the index

- [x] 3.1 Handle `j`/`k` in `ModeIndex`: move cursor through the flat list (skipping section separators)
- [x] 3.2 Handle `Enter` in `ModeIndex`: if active → `ModeNormal` with correct `changeIdx`; if archived → `ModeViewingArchive`
- [x] 3.3 Handle `Esc` in `ModeIndex`: `tea.Quit`
- [x] 3.4 Handle `q` in `ModeIndex`: `tea.Quit` (already covered by the global `q` case)

## 4. Transitions from change views

- [x] 4.1 In `ModeNormal`: `a` → `enterIndex()` (replaces the call to `enterArchivePicker()`)
- [x] 4.2 In `ModeNormal`: `Esc` → `enterIndex()` (previously a no-op)
- [x] 4.3 In `ModeViewingArchive`: `Esc` → `enterIndex()` (previously went to `ModeArchivePicker`)
- [x] 4.4 In `ModeViewingArchive`: `a` → `enterIndex()` (previously went to `ModeArchivePicker`)

## 5. Remove ModeArchivePicker

- [x] 5.1 Remove `ModeArchivePicker` from the `Mode` constant
- [x] 5.2 Remove `renderArchivePicker()` and its call in `View()`
- [x] 5.3 Remove `modal*` styles from `styles.go` that are no longer used
- [x] 5.4 Remove struct fields used only by `ModeArchivePicker` (`archiveCursor`)

## 6. Update helpbars

- [x] 6.1 Update `renderHelpBar()` in `ModeNormal` (tasks): include `Esc: index`
- [x] 6.2 Update `renderHelpBar()` in `ModeNormal` (other tabs): include `Esc: index`
- [x] 6.3 Update `renderHelpBar()` in `ModeViewingArchive`: change `a/Esc: back` to `a/Esc: index`

## 7. Verification

- [x] 7.1 Compile without errors (`go build ./...`)
- [x] 7.2 Verify full flow: Normal → index → select archived → ViewingArchive → index
- [x] 7.3 Verify flow: Normal → index → select active → Normal with correct change
- [x] 7.4 Verify that `Esc` from the index quits the application
- [x] 7.5 Verify that `q` works from all modes
