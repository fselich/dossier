## Why

The repository has a release workflow (triggered on `v*` tags) but no CI workflow for pushes to `main` or pull requests. There is no automated verification that the code compiles, tests pass, or passes vet checks. All three other P0 items (`inject-root-path`, `add-unit-tests`, `fix-go-mod`) will benefit from automated CI validation.

## What Changes

- Add `.github/workflows/ci.yml` that runs on `push` to `main` and `pull_request` to `main`
- Runs `go test -race -coverprofile=coverage.out ./...` on every push/PR
- Runs `go vet ./...` as a separate step
- Reports coverage summary via `go tool cover -func=coverage.out`
- Uses Go 1.25, matching the `go.mod`

## Capabilities

_No functional changes. This is CI infrastructure only._

## Impact

- New file: `.github/workflows/ci.yml`
- No code changes
- GitHub Actions minutes consumed per run (~30s)
