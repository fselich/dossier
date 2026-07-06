## Context

`pollGitStatus()` in `internal/ui/git.go` is called every 500ms tick. When it detects any change in the git status output, it unconditionally clears the diff view:

```go
m.gitState.Files = files
m.gitState.ShowingDiff = false   // always cleared
m.gitState.DiffLines = nil
m.gitState.DiffFile = ""
m.gitState.ScrollX = 0
```

This means editing any file while viewing a diff for another file kicks the user back to the file list.

## Goals / Non-Goals

**Goals:**
- Preserve the diff view when the viewed file (`DiffFile`) has not changed on disk
- Only clear the diff when the viewed file itself is modified, removed, or has a status change

**Non-Goals:**
- No changes to the diff rendering, cursor navigation, or file list behavior
- No changes to the polling interval or mechanism

## Decisions

- **Approach**: After updating `m.gitState.Files`, check if we're showing a diff. If so, search for `DiffFile` in the new file list and compare `X`/`Y` status. If the file still exists with the same status, keep the diff. Otherwise, clear it.
- **Why not "always preserve"**: If the viewed file itself was edited, the diff is stale — we should clear it so the user can re-enter to see the current state.
- **Why not compare file content (hash)**: Overkill. `git status --porcelain` already tells us if something changed. If the status bytes are the same, the diff remains valid for the user's purposes.

## Risks / Trade-offs

- [Edge case] If the viewed file's content changes but its git status stays the same (e.g., unstaged modification → still unstaged modification), the diff won't refresh. Mitigation: the user can press `Esc` and re-enter the file to see the latest diff.
