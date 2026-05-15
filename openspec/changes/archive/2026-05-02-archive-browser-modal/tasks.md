## 1. Loader — loading archived changes

- [x] 1.1 Add function `parseArchiveName(dir string) (name, date string)` in `loader.go` that extracts the clean name and date (`DD Mon`) from the `YYYY-MM-DD-` prefix; if the format does not match, returns the full name and an empty date
- [x] 1.2 Add function `ListArchiveChanges() []Change` in `loader.go` that reads `openspec/changes/archive/`, loads each subdirectory as a `Change` using `parseArchiveName` for the `Name` field, and returns the list sorted by directory name (most recent first)

## 2. Model — explicit state

- [x] 2.1 Define type `Mode int` with constants `ModeNormal`, `ModeArchivePicker`, `ModeViewingArchive` in `model.go`
- [x] 2.2 Add fields to the `Model` struct: `mode Mode`, `archiveChanges []openspec.Change`, `archiveCursor int`
- [x] 2.3 Add helper function `(m *Model) enterArchivePicker()` that loads `archiveChanges` if empty, resets `archiveCursor` to 0 and sets `mode = ModeArchivePicker`

## 3. Keyboard — mode transitions

- [x] 3.1 In the `tea.KeyMsg` handler, add case `"a"`: if `ModeNormal` → `enterArchivePicker()`; if `ModeViewingArchive` → `ModeArchivePicker`
- [x] 3.2 Add case `"esc"`: if `ModeArchivePicker` → `ModeNormal`; if `ModeViewingArchive` → `ModeArchivePicker`
- [x] 3.3 Add cases `"j"` / `"k"` in `ModeArchivePicker`: move `archiveCursor` within the bounds of `archiveChanges`
- [x] 3.4 Add case `"enter"` in `ModeArchivePicker`: clear `renderCache`, load the selected archived change as the display context and set `mode = ModeViewingArchive`
- [x] 3.5 Disable `"e"` and `" "` (space) when `mode == ModeViewingArchive`
- [x] 3.6 Disable `"h"` / `"l"` (change switching) when `mode == ModeViewingArchive`

## 4. Rendering — viewer in archive mode

- [x] 4.1 Adapt `renderHeader()` to use the clean archived change name and show `[archivo]` instead of `[N/M]` when `mode == ModeViewingArchive`
- [x] 4.2 Adapt `renderHelpBar()` to show `1-4: artifact   j/k: scroll   a/Esc: volver` when `mode == ModeViewingArchive`
- [x] 4.3 Adapt `loadViewport()` / `current()` to read from the selected archived change when `mode == ModeViewingArchive` instead of `project.Changes[changeIdx]`

## 5. Rendering — picker modal

- [x] 5.1 Implement `renderArchivePicker() string` that builds the modal content: title, list of items with clean name + date, `>` cursor, and internal helpbar
- [x] 5.2 In `View()`, when `mode == ModeArchivePicker`, overlay the modal on top of the normal content using `lipgloss.Place`

## 6. Verification

- [x] 6.1 Compile without errors
- [x] 6.2 Press `a` from an active change: confirm the modal appears with archived changes, clean names and dates
- [x] 6.3 Navigate with `j/k` in the modal and select an archived change with `Enter`: confirm the viewer shows its artifacts with `[archivo]` in the header
- [x] 6.4 Confirm that `e` and `Space` do nothing in archive mode
- [x] 6.5 Confirm that `Esc` from the viewer returns to the modal, and `Esc` from the modal returns to normal
- [x] 6.6 Confirm that `a` from the archive viewer returns to the modal
