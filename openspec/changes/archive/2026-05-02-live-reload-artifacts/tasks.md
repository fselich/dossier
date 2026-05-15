## 1. openspec-loader: reload function

- [x] 1.1 Add function `ReloadChange(ch Change) Change` in `loader.go` that re-reads proposal, design, tasks and specs from disk
- [x] 1.2 Verify that absent artifacts are left with `Present=false` without error

## 2. tui-viewer: polling cycle

- [x] 2.1 Add `tickMsg` and `fileChangedMsg` as message types in `model.go`
- [x] 2.2 Emit `tea.Tick(2s)` from `Init()` to start the cycle
- [x] 2.3 In the `tickMsg` handler, call `ReloadChange` for the active change
- [x] 2.4 If `tasks.md` changed: re-parse items, restore cursor by text, refresh viewport if active tab is tasks
- [x] 2.5 If proposal/design/specs changed: update `ch.Content` in memory and invalidate the corresponding `renderCache[tab]`
- [x] 2.6 Re-emit `tea.Tick(2s)` at the end of the handler to maintain the cycle

## 3. tasks-toggle: cursor restoration by text

- [x] 3.1 Add function `findCursorByText(items []TaskItem, text string) int` that returns the index of the first task item with that text, or the first available item if not found
- [x] 3.2 Use `findCursorByText` in the tasks reload to restore the cursor position

## 4. Manual verification

- [x] 4.1 Open the TUI in spec-viewer, edit `tasks.md` externally and verify that the view updates in ~2s
- [x] 4.2 Verify that the cursor stays on the same task after the reload
- [x] 4.3 Open the proposal tab, edit `proposal.md` externally, return to the tab and verify that it re-renders
