## 1. Line number tracking in diff parser

- [x] 1.1 Add `OldNum int` and `NewNum int` fields to `DiffLine` struct
- [x] 1.2 Update `parseDiff` to track `oldNum`/`newNum` from hunk headers (`@@ -old,count +new,count @@`)
- [x] 1.3 Update `parseDiffLine` to increment numbers per line type (context: both++, added: newNum++, removed: oldNum++)

## 2. State refactor

- [x] 2.1 Replace `DiffContent string` with `DiffLines []DiffLine` and `ScrollX int` in `gitState` (model.go)
- [x] 2.2 Update `toggleGitDiff` to store `DiffLines` instead of `DiffContent`
- [x] 2.3 Update `pollGitStatus` in git.go to clear `DiffLines` and `ScrollX`

## 3. Scroll-aware rendering

- [x] 3.1 Update `renderDiff` to accept `width int` and `scrollX int`, truncate raw content before highlighting
- [x] 3.2 Update `renderDiffContent` in git.go to pass `m.width` and `m.gitState.ScrollX` to `renderDiff`
- [x] 3.3 Display line numbers (4-char format) in diff view

## 4. Scroll keys

- [x] 4.1 Add `→` and `←` key handlers in diff view: adjust `ScrollX` by ±10, clamp to >=0
- [x] 4.2 `h`/`l` scroll horizontally in diff view, navigate changes otherwise
- [x] 4.3 `d`/`Esc` reset `ScrollX = 0` when returning to file list
- [x] 4.4 Regular `j`/`k` vertical scroll still works in diff view

## 5. Help bar

- [x] 5.1 Update help bar for diff mode to show: "d/Esc: back  j/k: vertical  h/l: ←→ horizontal"

## 6. Verify

- [x] 6.1 `make build` succeeds
- [x] 6.6 `make lint` passes
- [x] 6.7 `make test` passes
- [x] 6.2 Long lines scroll horizontally with `h`/`l`
- [x] 6.3 Line numbers display correctly in diff view
- [x] 6.4 Scroll resets when returning to file list
- [x] 6.5 Regular vertical scroll still works
