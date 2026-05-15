## Context

The project was developed under the working name `spec-viewer` with the binary `specview`. The Go module path is `github.com/fselich/dossier` and the entry point lives at `cmd/specview/`. The name appears in: the module declaration, two internal import paths, the Makefile, `.gitignore`, both README files, and three project-level spec files.

The code changes have already been applied. This design documents the approach and serves as the record for the spec updates.

## Goals / Non-Goals

**Goals:**
- All user-facing references to `specview` / `spec-viewer` replaced with `dossier`
- All internal Go module references consistent with the new module path
- Project-level specs accurate (no stale binary name in requirements or scenarios)
- Archived changes left untouched — they are historical records

**Non-Goals:**
- Renaming the repository directory on disk (can be done separately when pushing to GitHub)
- Updating archived change artifacts — archives are immutable by convention
- Changing any runtime behavior

## Decisions

**Rename `cmd/specview/` to `cmd/dossier/`** — Go's `go install` derives the binary name from the last path segment of the package. Renaming the directory is sufficient; no additional flags or configuration needed.

**Update module path in `go.mod`** — The module path must match the new repository URL. All internal imports use this path as a prefix, so a single find-and-replace across the codebase covers them all. `go.sum` does not need changes because no external dependencies changed.

**Update specs in place** — The three affected project specs (`build-tooling`, `path-arg`, `openspec-loader`) have requirements that reference the binary name. These are updated directly rather than creating delta specs, because the requirement intent is unchanged — only the name within the requirement text changes.

**Leave archived changes untouched** — Archives represent completed historical work. Updating them would corrupt the audit trail.

## Risks / Trade-offs

[Stale binary on PATH] → Any developer who installed `specview` via `go install` will need to run `go install github.com/fselich/dossier/cmd/dossier@latest` and remove the old binary manually. No automated migration.

[Module path mismatch until repo is renamed] → Until the GitHub repository is renamed to `dossier`, the module path `github.com/fselich/dossier` will not resolve via `go install`. Local builds with `make build` are unaffected.
