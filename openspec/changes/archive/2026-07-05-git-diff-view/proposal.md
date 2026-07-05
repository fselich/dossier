## Why

The git changes tab shows which files are modified, but doesn't show what changed inside them. To see the actual diff the user must open the file in the editor or switch to a terminal. Pressing `d` on a file should show its diff inline, keeping the workflow inside the TUI.

## What Changes

- Pressing `d` on a file in the git changes tab toggles a diff view within the same tab.
- For tracked files: `git diff --color=always` produces a colored line-level diff with ANSI codes.
- For untracked files: file contents rendered with chroma syntax highlighting via `lexers.Match(filename)`.
- Pressing `d` or `Esc` in diff view returns to the file list.
- Diff content is cached and invalidated when git status changes.

## Capabilities

### Modified Capabilities

- `git-status-tab`: The changes tab now supports a diff view triggered by the `d` key, showing `git diff --color=always` for tracked files or syntax-highlighted file contents for untracked files.

## Impact

- New file `internal/ui/gitdiff.go` with diff computation, chroma highlighting for untracked files, and viewport refresh.
- Modified files: `internal/ui/git.go` (add `diffState` to gitState, toggle logic in rendering/viewport), `internal/ui/view.go` (help bar for diff mode), `internal/ui/viewer.go` (key dispatch for `d`/`Esc` in git tab).
- No new dependencies — chroma is already a transitive dependency via glamour.
- No changes to OpenSpec domain logic.
