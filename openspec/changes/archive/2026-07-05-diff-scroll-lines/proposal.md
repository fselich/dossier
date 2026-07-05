## Why

The diff view shows code with syntax highlighting, but long lines overflow the viewport edge. There's no way to see the rest of a wide line. Adding horizontal scroll solves this cleanly — truncate raw content at scroll offset, then apply chroma highlighting. This avoids the complexity of wrapping ANSI-colored text.

Additionally, line numbers provide orientation within the diff, showing which lines in the file are being changed.

## What Changes

- Parse diff lines with `OldNum`/`NewNum` line numbers (tracked from hunk headers, like differ does).
- Store parsed `[]DiffLine` in `gitState` instead of just a content string; regenerate `DiffContent` on render.
- Add `ScrollX int` to `gitState` for horizontal scroll offset.
- Render truncates raw content at `ScrollX` rune offset *before* chroma highlighting, so ANSI codes are never split.
- Keys `→`/`l` and `←`/`h` scroll horizontally by 10 runes. `d`/`Esc` reset scroll to 0 and return to file list.
- Display line numbers (4-char wide) in diff view.

## Capabilities

### Modified Capabilities

- `git-status-tab`: Diff view now includes line numbers and horizontal scroll. Parsed diff lines are stored as structured data, not a pre-rendered string.

## Impact

- Refactor `internal/ui/gitdiff.go`: store parsed diff in `gitState.DiffLines`; render on-the-fly with scroll offset.
- Modify `internal/ui/model.go`: `gitState.DiffLines []DiffLine`, `gitState.ScrollX int` replace `DiffContent string`.
- Modify `internal/ui/git.go`: `renderDiffContent` uses scroll-aware rendering. `pollGitStatus` clears `DiffLines` and `ScrollX`.
- Modify `internal/ui/viewer.go`: `→`/`←`/`h`/`l` scroll in diff view. `d`/`Esc` reset scroll.
- Modify `internal/ui/view.go`: update help bar for scroll keys.
