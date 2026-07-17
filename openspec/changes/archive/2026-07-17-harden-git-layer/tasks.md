## 1. Subprocess helper and timeouts

- [x] 1.1 Add a `runGit(dir string, args ...string) ([]byte, error)` helper in `internal/git` using `exec.CommandContext` with a 2s `context.WithTimeout` (created and cancelled locally)
- [x] 1.2 Migrate `IsInsideWorkTree` and `Status` to the helper
- [x] 1.3 Apply the same timeout policy to the `git diff` invocations in `internal/ui/gitdiff.go`

## 2. WorkTreeRoot error propagation

- [x] 2.1 Change `WorkTreeRoot` to return `(string, error)` and remove the silent fallback to the input path
- [x] 2.2 Update the call site in `internal/ui/model.go` to fall back to the openspec root explicitly when an error is returned

## 3. Porcelain -z parsing

- [x] 3.1 Switch `Status` to `git status --porcelain=v1 -z -u` and parse NUL-separated entries (3-byte `XY ` prefix + path)
- [x] 3.2 Handle rename/copy entries: consume the following NUL token as `OldPath`
- [x] 3.3 Preserve the `openspec/` filter and `IsDeleted` derivation; remove the obsolete ` -> ` splitting and `TrimRight` line handling

## 4. Tests for internal/git

- [x] 4.1 Add test helpers: init a real repo in `t.TempDir()` (with user config), skip when git is not on PATH
- [x] 4.2 Table-driven `Status` tests: modified, added, untracked, renamed, copied, deleted, `openspec/` filter, clean tree
- [x] 4.3 `Status` tests for unusual filenames: spaces, non-ASCII, newline (best-effort create + skip if the filesystem rejects it)
- [x] 4.4 Tests for `IsInsideWorkTree` and `WorkTreeRoot` (repo root, subdirectory, non-repo → error)

## 5. Verification and docs

- [x] 5.1 Run `make test`, `make lint`, `go vet ./...`
- [x] 5.2 Update AGENTS.md gotcha #10 (porcelain parsing notes) to reflect `-z` parsing and the timeout helper
