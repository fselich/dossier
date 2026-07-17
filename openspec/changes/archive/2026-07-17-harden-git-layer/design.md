## Context

`internal/git` shells out to git and parses `git status --porcelain -u` line-by-line. Three weaknesses:

- Filenames containing `\n` (legal in git) split into multiple bogus entries; ` -> ` rename splitting is heuristic.
- `WorkTreeRoot` (git.go:24) returns the *input* path when `git rev-parse` fails, so `internal/ui/model.go` may build absolute editor paths from a wrong root with no signal.
- All `exec.Command` calls run without a deadline; `Status` runs on every 500ms tick, so a hung git (slow NFS, credential prompt) freezes the UI permanently.

The package has no tests. A follow-up change (`git-stage-unstage`) will run `git add`/`git restore --staged` on parsed paths, so parsing correctness becomes safety-critical.

## Goals / Non-Goals

**Goals:**
- Parsing that is byte-exact for any filename git can produce.
- Errors surfaced to callers instead of silent wrong values.
- Every git subprocess bounded in time.
- Test coverage for the whole package against real git repos.

**Non-Goals:**
- Git mutations (next change).
- Async diff computation or Chroma performance (PROPUESTAS.md 2.4).
- Error *display* in the UI (help-bar messaging arrives with `git-stage-unstage`); here callers keep degrading gracefully, just explicitly.

## Decisions

### D1: `--porcelain=v1 -z` over v2

`git status --porcelain=v1 -z -u` keeps the existing two-byte `XY` model (minimal churn in `FileStatus` and the UI) while making parsing exact: entries are NUL-separated, no quoting, and for renames/copies the *new* path is in the entry and the *old* path follows as the next NUL token. Porcelain v2 gives more data we don't need and would require a full parser rewrite.

Parsing rules:
- Split output on `\0`, drop the trailing empty token.
- Each entry: `XY ` (3 bytes) + path.
- If `X` or `Y` is `R`/`C`, consume the following token as `OldPath`.
- Keep the `openspec/` prefix filter and the `IsDeleted` derivation unchanged.

### D2: `WorkTreeRoot` returns `(string, error)`

Caller in `model.go` runs once at startup. On error it falls back to the openspec root (current de-facto behavior) — but now the decision is visible at the call site and testable. Alternative (keep silent fallback) rejected: the follow-up change uses `gitRoot` as `-C` argument for mutations.

### D3: Timeout via `exec.CommandContext` + `context.WithTimeout(2s)`

AGENTS.md says "no context.Context", but that guideline targets plumbing contexts through APIs. Here the context is created and cancelled inside each function (`defer cancel()`); no signatures change. A 2s budget is generous for local repos and small enough that a hung poll recovers within a few ticks. Applied to: `IsInsideWorkTree`, `WorkTreeRoot`, `Status`, and the `git diff` calls in `internal/ui/gitdiff.go` (extract a small `runGit(dir, args...)` helper in `internal/git` so the timeout policy lives in one place).

### D4: Tests use real git repos in `t.TempDir()`

Matches the repo convention (no mock filesystem). Helper initializes a repo with `git init` + `git -c user.name=... -c user.email=... commit`. Tests skip with `t.Skip` if git is not on PATH. Table-driven cases for `Status`; separate tests for `IsInsideWorkTree` / `WorkTreeRoot` (repo, non-repo, subdirectory).

## Risks / Trade-offs

- [`-z` output has no quoting, but also no newline terminator on the last entry] → split on `\0` and ignore empty tokens; test with trailing/absent NUL.
- [2s timeout may abort legitimately slow git on cold NFS] → poll retries every 500ms; a timed-out `Status` returns an error which the poller already ignores, keeping the last known state.
- [Filenames with `\n` in tests may misbehave on some filesystems] → guard that specific case with a best-effort create + skip.
- [Signature change on `WorkTreeRoot` touches UI init] → single call site; compile error makes it impossible to miss.
