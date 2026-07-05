## Context

The git changes tab (`TabGit`) currently shows a selectable list of changed files with status codes. The user can navigate with `j/k` and open files with `Enter`/`e`, but there's no way to preview changes without leaving the TUI. This adds a diff toggle within the same tab.

Chroma (syntax highlighter) is already a transitive dependency via glamour. It supports `lexers.Match(filename)` for automatic language detection from file extension or filename pattern.

## Goals / Non-Goals

**Goals:**
- Press `d` on a file in the changes tab to see its diff.
- For tracked files (M, A, D, R, C): `git diff HEAD --color=always -- <file>` — standard git diff with ANSI colors.
- For untracked files (??): read file content, detect language via `lexers.Match(filename)`, render with chroma syntax highlighting, show with a `(new file: path)` header.
- Press `d` or `Esc` to return to the file list.
- Diff content is cached per file in `gitState`; cache cleared on git status change.
- `j/k` scroll the diff viewport normally.

**Non-Goals:**
- Syntax highlighting of code within diffs for tracked files. The `git diff --color=always` output already colors `+`/`-`/`@@` markers; full code syntax highlighting inside diffs requires combining lexers and is left for a future iteration.
- Side-by-side diff or word-level diff.
- Interactive staging/unstaging from the diff view.
- External diff tools (delta, difftastic) — no system dependencies.

## Decisions

### Dual rendering strategy: git diff vs chroma

| File type | Display method | Why |
|-----------|---------------|-----|
| Tracked (X != '?') | `git diff --color=always HEAD -- <file>` | Already colored, no parsing needed. Shows all changes (staged + unstaged) vs HEAD. |
| Untracked (??) | `os.ReadFile` + chroma via `lexers.Match(filename)` | No git history to diff against. File content with full syntax highlighting is more useful than an empty diff. |

### Diff command choice

`git diff HEAD -- <file>` shows both staged and unstaged changes compared to HEAD, which handles ` M` (unstaged), `M ` (staged), and `MM` (both) correctly. For deleted files (` D` / `D `), the file is gone from the worktree, so `git diff HEAD -- <file>` works because git can diff the HEAD version against nothing. For renamed files, use the new path from `FileStatus.Path`.

For untracked files, `git diff HEAD -- <file>` returns an error (path not in HEAD). We fall back to reading the file directly and highlighting with chroma.

### Toggle within git tab, not a new mode

The diff view is a sub-state of the git tab (`gitState.ShowingDiff bool`, `gitState.DiffContent string`, `gitState.DiffFile string`) rather than a new `Mode`. This avoids:
- Adding mode-switching boilerplate in `update.go`/`dispatchKey`
- Interacting with `prevMode` restore logic
- Duplicating the chrome/layout logic

When `ShowingDiff` is true:
- `renderGitContent()` returns the diff instead of the file list
- The help bar shows diff-specific shortcuts
- Key dispatch handles `d`/`Esc` to toggle back

### Cache invalidation

Diff content is cached in `gitState.DiffContent` for the current file. The cache is cleared when:
- Git status changes (`pollGitStatus` detects new files or status changes)
- User returns to list view (presses `d`/`Esc` from diff view)
- Cursor moves to a different file (diff is file-specific)

### Chroma integration for untracked files

For untracked files, we read the file contents and pass them to chroma:

```go
lexer := lexers.Match(filename)
if lexer == nil {
    lexer = lexers.Fallback
}
style := styles.Get("monokai")
// ... tokenise + format to ANSI
```

The output is ANSI-colored text that the viewport renders directly (same as glamour output). The style "monokai" is chosen for its dark background compatibility (same as the existing glamour "dark" style).

## Architecture

```
internal/ui/gitdiff.go (NEW):
  - computeDiff(m *Model) string        — runs git diff or reads+highlights file
  - highlightUntrackedFile(path string) string — chroma highlighting
  - toggleGitDiff(m *Model)             — toggles ShowingDiff, populates DiffContent

internal/ui/git.go (MODIFIED):
  + gitState.ShowingDiff bool
  + gitState.DiffContent string
  + gitState.DiffFile string           — tracks which file the diff is for
  + renderGitContent() → branch when ShowingDiff
  + pollGitStatus() → clear DiffContent on change
  + refreshGitViewport() → show diff when ShowingDiff

internal/ui/view.go (MODIFIED):
  + renderHelpBar() → diff mode text: "d/Esc: back  j/k: scroll  q: quit"

internal/ui/viewer.go (MODIFIED):
  + "d" key in TabGit → toggleGitDiff
  + Esc in TabGit when ShowingDiff → back to list
```

## Risks / Trade-offs

- **`git diff --color=always` ANSI codes**: The viewport handles ANSI codes (it's how glamour renders markdown), so this should work. However, ANSI width calculation may be slightly off if git uses different escape sequences than glamour/chroma. `lipgloss.Width()` handles standard ANSI sequences correctly.
- **Large diffs**: For large files with many changes, the diff output could be thousands of lines. The viewport handles this (it's scrollable). No performance concern — `git diff` is O(changes) and typically fast.
- **Chroma theme mismatch**: The "monokai" theme for chroma may look different from the "dark" theme used by glamour. This is acceptable — the diff view is a different context from markdown rendering.
- **Binary files**: `git diff --color=always` for binary files outputs "Binary files differ". This is handled gracefully by the viewport.
- **No diff for untracked empty files**: An empty file would show nothing useful. Edge case, negligible.

## Migration Plan

1. Add `gitState.ShowingDiff`, `DiffContent`, `DiffFile` fields.
2. Create `internal/ui/gitdiff.go` with `computeDiff`, `highlightUntrackedFile`, `toggleGitDiff`.
3. Modify `renderGitContent` to branch on `ShowingDiff`.
4. Add `d`/`Esc` key dispatch in `viewer.go`.
5. Update help bar in `view.go`.
6. Clear cache in `pollGitStatus`.
7. `make build && make lint && make test`.
