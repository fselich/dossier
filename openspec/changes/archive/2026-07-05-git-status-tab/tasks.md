## 1. Git package

- [x] 1.1 Create `internal/git/git.go` with `IsInsideWorkTree(root string) bool` and `WorkTreeRoot(root string) string`
- [x] 1.2 Add `FileStatus` struct (XY, Path, IsDeleted bool) and `Status(root string) ([]FileStatus, error)` parsing `git status --porcelain`
- [x] 1.3 Handle renamed files: parse `old -> new` arrow, populate Path with new path
- [x] 1.4 Handle deleted files: set `IsDeleted = true` for ` D` and `D ` statuses
- [x] 1.5 Filter out paths with `openspec/` prefix (already tracked by existing tabs)

## 2. Model changes

- [x] 2.1 Add `TabGit = 4` to Tab enum, bump `tabCount` to 5, add `"changes"` to `tabLabels`
- [x] 2.2 Add `isGitRepo bool`, `gitRoot string`, and `gitState` struct with `Files []git.FileStatus` and `Cursor int` to Model
- [x] 2.3 Detect git repo in `New()` via `git.IsInsideWorkTree(root)`, store in `isGitRepo`
- [x] 2.4 Update `tabAvailable()` with `case TabGit: return m.isGitRepo`
- [x] 2.5 Update `defaultTab()` to respect new tabCount (no change needed — it already iterates to `tabCount`)

## 3. Git tab rendering

- [x] 3.1 Create `internal/ui/git.go` with `renderGitContent()` returning content string + cursor line
- [x] 3.2 Style status codes: M=yellow, A/??=green, R=cyan, D=gray (dimmed)
- [x] 3.3 Renamed: show `old → new` format
- [x] 3.4 Clean repo: show `(working tree clean)`
- [x] 3.5 Add `refreshGitViewport()` mirroring `refreshTasksViewport` pattern
- [x] 3.6 Add status colors to `internal/ui/styles.go`

## 4. Tab bar and help bar

- [x] 4.1 Update `renderTabBar()` to include TabGit label with count when `isGitRepo`: `changes (N)`
- [x] 4.2 Update tab bar width calculation to account for 5th tab in mouse click handler
- [x] 4.3 Update `renderHelpBar()` for TabGit: show navigation + editor hints

## 5. Key dispatch

- [x] 5.1 Add `case "5"` in `updateViewer()` to switch to TabGit
- [x] 5.2 Handle `j/k` in TabGit: move cursor, skip deleted files
- [x] 5.3 Handle `Enter` and `e` in TabGit: open file in `$EDITOR` with absolute path
- [x] 5.4 Ensure `h/l` (change navigation) still works in TabGit

## 6. Viewport loading

- [x] 6.1 Add `case m.tab == TabGit` in `loadViewport()` dispatching to `loadGitViewport`

## 7. Mouse handling

- [x] 7.1 Update tab bar click Y=2 coordinate calc for 5 tabs in `handleMouseClick`
- [x] 7.2 Handle mouse wheel scrolling in TabGit (move cursor, like tasks)

## 8. Polling

- [x] 8.1 Add `pollGitStatus()` to internal/ui/git.go: runs `git.Status(m.root)` if `isGitRepo`
- [x] 8.2 Call `pollGitStatus()` from `handleTick()` when `m.mode == ModeNormal`
- [x] 8.3 On status change, refresh git viewport if TabGit is active

## 9. Verify

- [x] 9.1 `make build` succeeds
- [x] 9.5 `make lint` passes
- [x] 9.6 `make test` passes
- [x] 9.2 Navigate git tab: j/k skip deleted, Enter opens file in editor
- [x] 9.3 Tab bar shows `changes (N)` count
- [x] 9.4 Clean working tree shows `(working tree clean)`
