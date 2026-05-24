# CI Infrastructure

## ADDED Requirements

### Requirement: Automated test execution on push and PR
The CI workflow SHALL run `go test -race -coverprofile=coverage.out ./...` on every push to `main` and every pull request targeting `main`.

#### Scenario: Push to main
- **WHEN** code is pushed to the `main` branch
- **THEN** the CI workflow triggers and runs all tests

#### Scenario: Pull request opened
- **WHEN** a pull request targeting `main` is opened or updated
- **THEN** the CI workflow triggers and runs all tests

### Requirement: Static analysis in CI
The CI workflow SHALL run `go vet ./...` before tests to catch compilation and static analysis issues.

#### Scenario: vet passes
- **WHEN** `go vet ./...` runs and finds no issues
- **THEN** the CI continues to the test step

### Requirement: Coverage summary in CI logs
The CI workflow SHALL generate and display a coverage summary using `go tool cover -func=coverage.out` after tests complete.

#### Scenario: Tests produce coverage data
- **WHEN** `go test -coverprofile=coverage.out` completes
- **THEN** `go tool cover -func=coverage.out` outputs per-function coverage percentages in the CI logs
