## 1. Git diff computation

- [x] 1.1 Create `internal/ui/gitdiff.go` with `highlightUntrackedFile(path string) string` using chroma's `lexers.Match`, "monokai" style, terminal16m formatter
- [x] 1.2 Add `computeDiff` that for tracked files runs `git diff --color=always` with appropriate flags and for untracked calls `highlightUntrackedFile`
- [x] 1.3 Handle error cases: file read failure, git diff failure (show error or fallback message)

## 2. Git state changes

- [x] 2.1 Add `ShowingDiff bool`, `DiffContent string`, `DiffFile string` to `gitState` in `model.go`
- [x] 2.2 Add `toggleGitDiff` in `gitdiff.go` that toggles `ShowingDiff`, populates `DiffContent` and `DiffFile`, and refreshes viewport
- [x] 2.3 Clear `DiffContent` and set `ShowingDiff = false` in `pollGitStatus` when git status changes in `git.go`

## 3. Rendering diff content

- [x] 3.1 Modify `renderGitContent` in `git.go` to branch: if `ShowingDiff`, render `DiffContent` with a file header
- [x] 3.2 Add `renderDiffContent` in `git.go` showing diff header + content
- [x] 3.3 Update `renderHelpBar` in `view.go` for diff mode: "d/Esc: back  j/k: scroll"

## 4. Key dispatch

- [x] 4.1 Add `case "d"` in `updateViewer` for TabGit: toggle diff view
- [x] 4.2 Add Esc handling in TabGit when `ShowingDiff` → back to list
- [x] 4.3 `h`/`l` still work from diff view; `j`/`k` scroll viewport in diff view, navigate list otherwise

## 5. Verify

- [x] 5.1 `make build` succeeds
- [x] 5.5 `make lint` passes
- [x] 5.6 `make test` passes
- [x] 5.2 Press `d` on modified tracked file shows colored diff
- [x] 5.3 Press `d` on untracked file shows syntax-highlighted content
- [x] 5.4 Press `d` or `Esc` returns to file list
