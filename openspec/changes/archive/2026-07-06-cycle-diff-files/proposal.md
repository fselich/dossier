## Why

When reviewing diffs in the "code" tab, switching between files requires three steps per file: exit diff, navigate to the next file, re-enter diff. For a review session with many changed files this is tedious. Keyboard shortcuts to cycle directly between diffs would enable fast, fluid review without leaving the diff view.

## What Changes

- Add keyboard shortcuts in diff view to jump to the previous/next file's diff
- When cycling, reload the diff for the target file and update the cursor in the underlying file list
- Update the help bar to show the new shortcuts
- Wrap around at the ends of the file list

## Capabilities

### New Capabilities

- (none)

### Modified Capabilities

- `git-status-tab`: Requirement "Diff view toggle within git changes tab" — add keyboard shortcuts for cycling between file diffs without returning to the file list

## Impact

- Changes in `internal/ui/viewer.go` (new key handlers in `updateViewer`)
- Changes in `internal/ui/gitdiff.go` (extract and reuse diff-loading logic)
- Changes in `internal/ui/view.go` (help bar)
