## Why

The `internal/git` package is read-only today but fragile: porcelain parsing breaks on filenames with newlines, `WorkTreeRoot` silently returns a wrong path on error, git subprocesses have no timeout (a hung git freezes the UI polling every 500ms), and the package has zero test coverage. Upcoming git mutations (stage/unstage) will act on parsed paths, so a mis-parsed path would mutate the wrong file — the layer must be hardened first.

## What Changes

- Switch `git.Status` to `git status --porcelain=v1 -z -u` and parse NUL-separated entries (renames/copies arrive as two consecutive NUL entries instead of ` -> `).
- Change `WorkTreeRoot` to return `(string, error)` instead of silently falling back to the input path; caller (`internal/ui/model.go`) decides how to degrade. **BREAKING** (internal API only).
- Add a timeout (~2s) to all git subprocess invocations (`IsInsideWorkTree`, `WorkTreeRoot`, `Status`, and the diff invocations in `internal/ui/gitdiff.go`) so a hung git cannot freeze the UI.
- Add a test suite for `internal/git` using real git repos in `t.TempDir()`: modified/added/untracked/renamed/copied/deleted files, `openspec/` filtering, filenames with spaces and newlines, short lines, non-repo directories.

## Capabilities

### New Capabilities

- `git-cli-integration`: behavior of the git subprocess layer — robust porcelain parsing (NUL-separated, rename pairs, unusual filenames), explicit error propagation from `WorkTreeRoot`, and bounded execution time for all git invocations.

### Modified Capabilities

<!-- none: git-status-tab requirements (what the tab shows) are unchanged; this change hardens how the data is obtained -->

## Non-goals

- No git mutations (stage/unstage/commit) — that is a follow-up change (`git-stage-unstage`).
- No changes to the git tab UI, key bindings, or rendering.
- No resolution of the `code` vs `changes` tab label inconsistency (PROPUESTAS.md 1.8).
- No async diff computation (PROPUESTAS.md 2.4); timeouts only bound worst-case blocking.

## Impact

- `internal/git/git.go`: parsing rewrite, signature change, timeouts.
- `internal/git/git_test.go`: new file.
- `internal/ui/model.go`: adapt to `WorkTreeRoot` returning an error.
- `internal/ui/gitdiff.go`: timeout on `git diff` invocations.
- AGENTS.md gotcha #10 (porcelain parsing notes) needs updating after the `-z` switch.
