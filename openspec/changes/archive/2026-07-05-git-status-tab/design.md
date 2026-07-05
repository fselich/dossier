## Context

Dossier's UI has four tabs in `ModeNormal`: proposal, design, specs, tasks. Each tab either renders Markdown via glamour (proposal, design, specs) or a custom selectable list (tasks). The git tab is a new selectable-list tab — it doesn't render Markdown, it renders a cursor-driven file list.

The `Tab` enum is `iota`-based with `tabCount` as the sentinel. Adding a fifth tab means bumping `tabCount` from 4 to 5 and adding `TabGit = 4`.

## Goals / Non-Goals

**Goals:**
- Detect git repo at startup; store `isGitRepo bool` in `Model`.
- Tab `changes` (key `5`) only available when `isGitRepo`.
- Parse `git status --porcelain` output into `[]FileStatus`.
- Render as a cursor-driven list with status codes colored by type.
- `j/k` navigate; cursor skips deleted files (` D`, `D `).
- `Enter` or `e` opens selected file in `$EDITOR` with absolute path.
- Renamed files (`R `) show `old → new`; open the new path.
- Poll status on the existing 500ms tick (only when TabGit is active).
- Show file count in tab bar: `changes (N)`.

**Non-Goals:**
- Show the tab when there are zero OpenSpec changes.
- Staged vs unstaged split view (we show the raw `XY` code).
- `git diff`, `git log`, `git blame`, or any other git subcommand.
- Interactive staging/unstaging of files.
- Custom git root detection beyond the project root.

## Decisions

### `git status --porcelain` over go-git

The porcelain format is stable, well-documented, and requires zero dependencies. `go-git` adds ~20 transitive packages for a single API call we can implement in ~30 lines of parsing code. Since git must be installed for the repo to exist, the binary dependency is guaranteed.

### Status codes: raw `XY` pair

`git status --porcelain` outputs `XY path` where X is the index status and Y is the worktree status. We display the raw two-character code (e.g. ` M`, `??`, `MM`) rather than friendly labels. This is compact, unambiguous, and familiar to git users.

### Cursor skips deleted files

Deleted files (` D` / `D `) have no file to open. The cursor navigates past them; pressing Enter on a deleted line is a no-op. They are rendered dimmed to indicate non-interactivity.

### Renamed files: open the new path

`git status --porcelain` outputs renames as `R  old → new`. We parse the arrow, show both paths, and use the new path for the editor command.

### Filter out OpenSpec directory

Files under `openspec/` are already tracked and navigable through the existing tabs (proposal, design, specs, tasks). Showing them in the git tab would be redundant. `Status()` filters out any path with prefix `openspec/`.

### Polling on existing tick

`git status --porcelain` on a typical repo takes <15ms. The 500ms tick is already polling for OpenSpec changes; adding git status to the same tick loop adds negligible overhead. We only run `git status` when `m.isGitRepo && m.mode == ModeNormal` (regardless of current tab, so switching tabs shows fresh data immediately).

### Editor: absolute path

`exec.Command(editor, filepath.Join(m.root, file.Path))` — always absolute, works regardless of the user's cwd.

## Architecture

```
internal/git/git.go:
  - IsInsideWorkTree(root string) bool
  - FileStatus{XY, Path string, IsDeleted bool}
  - Status(root string) ([]FileStatus, error)

internal/ui/model.go:
  + TabGit = 4, tabCount = 5
  + tabLabels[tabCount] = {..., "changes"}
  + isGitRepo bool
  + gitState struct { Files []git.FileStatus; Cursor int }
  + tabAvailable(): case TabGit → m.isGitRepo

internal/ui/view.go:
  + renderTabBar(): include TabGit if isGitRepo
  + renderHelpBar(): when tab==TabGit → "j/k: navigate  Enter/e: edit  1-5: tabs"
  + contentHeight(): reduce by 1 when git tab needs subnav? No — it's a simple list, fits in viewport.

internal/ui/viewer.go:
  + key "5" → TabGit
  + j/k in TabGit → move cursor (skip deleted)
  + Enter/e in TabGit → open file in editor (same pattern as artifactPath editor)

internal/ui/viewport.go:
  + loadViewport(): case TabGit → loadGitViewport (synchronous, no glamour)

internal/ui/mouse.go:
  + handleMouseClick: Y=2 includes new tab width
  + handleMouseWheel: TabGit scrolls cursor

internal/ui/git.go (NEW):
  + refreshGitViewport(): render + set content + scroll-to-cursor
  + renderGitContent(): builds the selectable list string
  + pollGitStatus(): called from handleTick → runs git.Status if isGitRepo

internal/ui/styles.go:
  + gitStatusModified (yellow), gitStatusAdded (green), gitStatusDeleted (gray),
    gitStatusRenamed (cyan), gitStatusUntracked (green)
```

## Risks / Trade-offs

- **git binary not in PATH**: `exec.LookPath("git")` returns error → `isGitRepo` stays false, tab never appears. Graceful degradation.
- **Filernames with spaces/special chars**: `git status --porcelain` is space-delimited. Parsing must split on first space, not `strings.Fields`. Paths may contain spaces — the format always has exactly 2 chars + space prefix, so we can safely slice `[3:]`.
- **Large repos with thousands of changed files**: `git status --porcelain` scales linearly. For a monorepo with 10k files it could take 100-200ms. Acceptable for 500ms polling. If it becomes a problem, we'd add a `--untracked-files=normal` or debounce.
- **Renamed file parsing**: `R  old → new` — we split on ` -> ` (with spaces). If a filename literally contains ` -> `, parsing would be wrong. This is extremely unlikely and git itself uses this delimiter in porcelain output.

## Migration Plan

1. Create `internal/git/git.go` with `IsInsideWorkTree` and `Status`.
2. Add `TabGit` to enum, `isGitRepo` to `Model`, detect in `New()`.
3. Create `internal/ui/git.go` with rendering + viewport logic.
4. Wire up tab visibility (`tabAvailable`), key dispatch (`viewer.go`), mouse handling, help bar.
5. Add git polling to `handleTick`.
6. `make build && ./dossier` to verify in a git repo.
7. Test: clean repo (shows "working tree clean"), modified files, deleted files (cursor skips), renamed files.
