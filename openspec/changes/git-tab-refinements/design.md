## Context

The git changes tab currently opens files in the editor on `Enter`/`e` via `openGitFile()`, while `d` shows the diff view via `toggleGitDiff()`. This is inconsistent — the index mode's file list uses `Enter` to navigate into a change, not open files. The diff view is the more useful in-TUI action.

Additionally, the tab is always visible when inside a git repo, even when `len(m.gitState.Files) == 0`. Other tabs (proposal, design, specs, tasks) are disabled when their artifact is absent. The git tab should follow the same convention.

## Goals / Non-Goals

**Goals:**
- `Enter` and `e` in git tab call `toggleGitDiff()` instead of `openGitFile()`.
- Tab is disabled (grayed out, not selectable) when working tree is clean.
- Remove `openGitFile` function to avoid dead code.
- Tab becomes enabled automatically when `pollGitStatus` detects new files.

**Non-Goals:**
- Adding new file-opening workflows from other parts of the UI.
- Changing index mode behavior.
- Changing the `d` key behavior (it already calls `toggleGitDiff`).

## Decisions

### `Enter` and `e` call `toggleGitDiff`

Simple change in `viewer.go`: both keys now delegate to the same function.

```go
case "enter":
    if m.tab == TabGit {
        m.toggleGitDiff()
    }
case "e":
    if m.tab == TabGit {
        m.toggleGitDiff()
    }
```

### `tabAvailable` checks file count

```go
case TabGit:
    return m.isGitRepo && m.mode == ModeNormal && len(m.gitState.Files) > 0
```

When `pollGitStatus` runs and detects files, `m.gitState.Files` becomes non-empty, `tabAvailable(TabGit)` returns `true`, and `renderTabBar` shows the tab as active. The transition happens on the next tick (max 500ms).

### Remove `openGitFile`

The function in `git.go` used `os.Getenv`, `exec.Command`, `tea.ExecProcess`, and `filepath.Join`. With no callers, the function and its unused imports are removed. The import cleanup is handled by `make fmt` (`goimports -w .`).

## Risks / Trade-offs

- **User can no longer open files from the git tab**: This is intentional. The user can still open files from the index mode or from the file system. The git tab is for browsing changes and viewing diffs.
- **Tab flickers on transition**: When the working tree switches from clean to dirty, the tab transitions from disabled to enabled. This is consistent with how other tabs behave when their artifacts appear/disappear during polling.
