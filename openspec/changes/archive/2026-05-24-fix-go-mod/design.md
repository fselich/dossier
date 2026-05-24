## Context

`go.mod` was likely generated incorrectly (possibly edited by hand or a Go version upgrade that didn't run `go mod tidy`). Go 1.25 requires explicitly declaring direct dependencies without `// indirect`.

## Goals / Non-Goals

**Goals:**
- Correctly classify direct imports as direct dependencies in `go.mod`

**Non-Goals:**
- Update dependency versions (out of scope)
- Add or remove dependencies

## Decisions

**Decision:** Run `go mod tidy` to let Go automatically detect direct imports and fix the classification.

This is the standard, zero-risk approach. No manual edits to `go.mod`.

## Risks / Trade-offs

- **Minimal risk of version changes**: `go mod tidy` may upgrade indirect dependencies if there are version conflicts. Mitigation: run `go build ./...` after to verify compilation.
