## Context

The entry point lives at `cmd/spec-viewer/main.go`. Go's toolchain names the installed binary after the containing directory, so `go install` currently produces `spec-viewer`, not `specview`. Two stale compiled binaries (`main`, `sv`) sit in the project root from earlier manual builds. There is no `Makefile`.

The user has `~/go/bin` in `$PATH`, so `go install` is the right installation mechanism.

## Goals / Non-Goals

**Goals:**
- Binary installed as `specview` via `go install ./cmd/specview/`
- `make build` → local `./specview` binary
- `make install` → `~/go/bin/specview`
- `make clean` → removes compiled artifacts from the repo root
- Remove `main` and `sv` from the repo

**Non-Goals:**
- Cross-platform packaging (`.deb`, `.rpm`, Homebrew formula)
- Version embedding in the binary
- CI/CD pipeline changes

## Decisions

**Rename `cmd/spec-viewer/` → `cmd/specview/`**
Go convention: the directory name under `cmd/` becomes the binary name. Renaming is the idiomatic fix rather than overriding with `-o` everywhere.
_Alternative_: keep the directory, always pass `-o specview` — adds repetition and surprises contributors.

**Makefile over a shell script**
`make` is universally available on Linux/macOS, has a clean target syntax, and is the de-facto standard for Go projects with simple build needs.
_Alternative_: `Taskfile` (go-task) — overkill for three targets; adds a dependency.

**`make install` calls `go install`, not `cp`**
`go install` handles caching, cross-compilation flags, and writes directly to `$GOPATH/bin`. Copying a locally-built binary would bypass that.

**Delete `main` and `sv` with `make clean`, not via git**
These are untracked files (not committed), so a simple `rm -f` in `clean` is sufficient. A `.gitignore` entry for compiled binaries would prevent recurrence.

## Risks / Trade-offs

- **Broken internal references to `cmd/spec-viewer`** → Mitigation: grep for any Makefile, CI, or doc references before renaming (unlikely given the project is early-stage).
- **`main` or `sv` are in use** → Mitigation: they are untracked binaries; removing them has no git impact and is trivially reversible by rebuilding.

## Migration Plan

1. `git mv cmd/spec-viewer cmd/specview`
2. `rm main sv` (untracked, safe to delete)
3. Add `Makefile`
4. Add compiled binary names to `.gitignore`
5. Verify: `make install && specview`
