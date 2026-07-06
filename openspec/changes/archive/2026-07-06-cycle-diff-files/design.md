## Context

The diff view is entered from the file list by pressing `Enter`, `d`, or `e` on a file. Currently there is no way to switch to a different file's diff without first returning to the file list (`Esc`/`d`), navigating (`j`/`k`), and re-entering (`Enter`/`d`).

State in `gitState`:
- `Files []git.FileStatus` — file list
- `Cursor int` — current position in file list (synced with which file's diff is showing)
- `ShowingDiff bool`
- `DiffLines []DiffLine`
- `DiffFile string`
- `ScrollX int`

## Goals / Non-Goals

**Goals:**
- Keyboard shortcuts to move to the previous/next file's diff while staying in diff view
- Wrap around at list boundaries
- Update help bar to show new shortcuts

**Non-Goals:**
- No changes to how diffs are loaded or rendered
- No changes to the file list navigation (`j`/`k` in list mode)

## Decisions

- **Keys**: Use `[` (previous file) and `]` (next file) when in diff view
  - Chosen over `ctrl+pgup`/`ctrl+pgdown` because they're one-key chords, terminal-portable, and follow the convention of `[`/`]` for prev/next navigation (common in pagers like `less`)
  - `ctrl+pgup`/`ctrl+pgdown` are not universally recognized across terminal emulators
- **Mechanism**: Extract the "load diff for file at cursor" logic from `toggleGitDiff()` into a reusable method `loadDiffForFile(cursor int)`. The cycling handler moves `m.gitState.Cursor` by ±1 (skipping deleted files, wrapping), then calls this method.
- **Cursor sync**: The underlying file cursor moves with cycling, so exiting the diff returns to the correct file in the list.

## Risks / Trade-offs

- [Conflict] `[` and `]` are currently unassigned in diff view and in normal mode. No conflict risk.
- [Diff freshness] When cycling to a new file, the diff is loaded fresh. When cycling back to a previous file, it's reloaded from disk — no stale diff risk.
