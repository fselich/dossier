## Why

The git tab is read-only: to stage or unstage a file the user must leave the TUI. Adding a stage/unstage toggle on the file under the cursor keeps the review-then-stage loop inside dossier. This builds on `harden-git-layer`, which makes parsed paths safe to use as mutation arguments.

## What Changes

- New key `s` in the git tab file list: toggles staged state of the file under the cursor (`git add` / unstage).
- Toggle semantics by status: if the file has unstaged changes (`Y != ' '` or untracked) → stage; otherwise → unstage. Unstaging a staged rename resets both old and new paths.
- Cursor no longer skips deleted files: the cursor can land on ` D` entries so they can be staged; `Enter`/`e`/`d` remain no-ops on deleted files. **BREAKING** for the `git-status-tab` cursor requirement (fixes PROPUESTAS.md 1.7 as a side effect).
- Immediate status refresh after a mutation (no 500ms tick wait), preserving the cursor by file path.
- Ephemeral error message in the help bar when a git mutation fails (e.g. `index.lock` held).
- `s` is inactive inside the diff view (list only).

## Capabilities

### New Capabilities

- `git-stage-unstage`: staging and unstaging files from the git tab, including toggle semantics, immediate refresh, cursor preservation, and error feedback.

### Modified Capabilities

- `git-status-tab`: the "Cursor navigation skips deleted files" requirement changes — the cursor SHALL be able to land on deleted files (open/diff actions remain no-ops on them).

## Non-goals

- No commit flow (message popup, `git commit`) — possible future change.
- No partial/hunk staging (`git add -p`).
- No staging of files under `openspec/` (they stay filtered out of the list).
- No `s` action inside the diff view; the existing behavior that a changed XY closes the diff view stays as is.
- No general error/notification system beyond the git tab help-bar message.

## Impact

- `internal/git/git.go`: new mutation functions (stage/unstage) using the hardened `runGit` helper; tests in `internal/git/git_test.go`.
- `internal/ui/viewer.go`: `s` key handling in `TabGit`.
- `internal/ui/git.go`: cursor logic (`moveGitCursor*`, `clampGitCursor`) stops skipping deleted files; cursor preservation by path; immediate refresh path.
- `internal/ui/view.go`: help bar shows `s stage/unstage` hint and the ephemeral error message.
- AGENTS.md gotcha #11 (cursor skips deleted files) needs updating.
- Depends on `harden-git-layer` being implemented first.
