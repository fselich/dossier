## Why

When viewing a file diff in the "code" tab, any change to any other file on disk (detected by the 500ms polling) clears the diff and returns the user to the file list. This is disruptive when focusing on a specific diff, as unrelated edits elsewhere dismiss the view.

## What Changes

- In `pollGitStatus()`, when a diff is being shown (`ShowingDiff == true`), do NOT unconditionally clear it on every git status change
- Only clear the diff if the file being viewed (`DiffFile`) was itself modified, removed, or had its status changed
- If the viewed file is unchanged, keep the diff view stable and only update the underlying file list

## Capabilities

### New Capabilities

- (none)

### Modified Capabilities

- `git-status-tab`: Requirement "Diff view toggle within git changes tab" — the "Diff content invalidated on status change" scenario narrows to only clear when the viewed file itself changes, not on any unrelated file change

## Impact

- Single function change in `internal/ui/git.go` (`pollGitStatus()`)
- No API changes, no new dependencies
- The diff content is already a snapshot — keeping it stable when the viewed file is untouched is safe
