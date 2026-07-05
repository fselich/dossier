## Why

The git changes tab has two UX issues:

1. **`Enter` and `e` open files in the editor**, but the user wants to see the diff first (as `d` does). Opening in the editor has the same effect as selecting the file from the file system — it's redundant with the index mode's file browsing.

2. **The changes tab is always selectable** even when the working tree is clean (showing "working tree clean"), unlike other tabs that are disabled when no artifact exists.

## What Changes

- `Enter` and `e` in the git changes tab now call `toggleGitDiff` (show diff) instead of `openGitFile` (open in editor).
- `tabAvailable(TabGit)` returns `false` when `len(m.gitState.Files) == 0`, disabling the tab (grayed out) for a clean working tree.
- Remove `openGitFile` function from `git.go` (dead code).

## Capabilities

### Modified Capabilities

- `git-status-tab`: `Enter`/`e` show diff view instead of opening in editor. Tab is disabled when working tree is clean.

## Impact

- `internal/ui/viewer.go`: change `"enter"` and `"e"` handlers for TabGit from `openGitFile` to `toggleGitDiff`.
- `internal/ui/model.go`: `tabAvailable(TabGit)` adds `&& len(m.gitState.Files) > 0`.
- `internal/ui/git.go`: remove `openGitFile` function.
- No other files.
