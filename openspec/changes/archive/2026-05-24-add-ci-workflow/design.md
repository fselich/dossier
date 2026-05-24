## Context

The release workflow (`.github/workflows/release.yml`) already uses `actions/setup-go@v5` with Go 1.25. The CI workflow should follow the same patterns for consistency.

Current CI gaps: no automated test runs, no lint/vet checks, no coverage tracking on PRs.

## Goals / Non-Goals

**Goals:**
- Run tests on every push to `main` and every PR
- Run `go vet` for static analysis
- Show coverage summary in CI logs

**Non-Goals:**
- Coverage gates (blocking PRs on coverage thresholds) — can be added later
- Linter integration (`.golangci.yml` from item 1.5 is a separate P1 change)
- Code coverage reporting to external services (Codecov, Coveralls)

## Decisions

**Decision: Single job, multiple steps**

A single `test` job with sequential steps keeps the workflow simple:
1. Checkout
2. Setup Go 1.25
3. `go vet ./...`
4. `go test -race -coverprofile=coverage.out ./...`
5. `go tool cover -func=coverage.out`

No need for a matrix build since it's a single-platform Go tool.

**Decision: `-race` flag enabled**

The race detector adds overhead but catches data races early. The project is small enough that the overhead is negligible.

## Risks / Trade-offs

- **Tests may fail initially**: The `inject-root-path` change may break existing test expectations. Mitigation: implement `inject-root-path` and `add-unit-tests` before enabling CI, or ensure CI passes on `main` before merging.
- **Race detector on CI**: GitHub Actions runners may be slower with `-race`. Mitigation: acceptable trade-off for early detection.
