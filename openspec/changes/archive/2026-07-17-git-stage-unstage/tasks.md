## 1. Git mutation API (internal/git)

- [x] 1.1 Add `Stage(root string, paths ...string) error` using the `runGit` helper (`git add -- <paths>`)
- [x] 1.2 Add `Unstage(root string, paths ...string) error`: `git reset -q HEAD -- <paths>` with fallback to `git rm --cached -q -- <paths>` on failure
- [x] 1.3 Tests in `internal/git/git_test.go`: stage modified/untracked/deleted, unstage staged file, unstage in a repo without commits, unstage a staged rename (both paths)

## 2. Cursor lands on deleted files

- [x] 2.1 Remove the `IsDeleted` skip from `moveGitCursorDown`, `moveGitCursorUp`, and `clampGitCursor` (simple wrap/clamp); keep `[`/`]` diff cycling skipping deleted files
- [x] 2.2 Verify `Enter`/`e`/`d` no-op guards on deleted files still hold; add a guard for `d` if missing
- [x] 2.3 Update UI tests covering cursor movement over deleted entries

## 3. `s` key toggle in the git tab

- [x] 3.1 Handle `s` in `updateViewer` for `TabGit` (list view only, no-op in diff view and on clean tree)
- [x] 3.2 Implement toggle direction: `Y != ' '` or untracked → `Stage`; else → `Unstage` (pass both paths for renames)
- [x] 3.3 After mutation, re-run `git.Status`, update `gitState.Files`, and refresh the viewport immediately
- [x] 3.4 Preserve cursor by path after refresh; clamp to a valid index when the path disappears

## 4. Help bar feedback

- [x] 4.1 Add `gitState.ErrMsg`; set it with a trimmed one-line error when a mutation fails
- [x] 4.2 Render `ErrMsg` in the help bar while present; clear it on next successful action, any git-tab key press, or a polled status change
- [x] 4.3 Add `s stage/unstage` hint to the git tab help bar (file list view)

## 5. Tests, docs, verification

- [x] 5.1 UI tests for the `s` toggle: stage, unstage, mixed `MM`, deleted file, error path (repo with held `index.lock`)
- [x] 5.2 Update AGENTS.md gotcha #11 (cursor no longer skips deleted files) and mention the `s` key
- [x] 5.3 Run `make test`, `make lint`, `go vet ./...`
