# lint-and-makefile Specification

## Purpose

Project tooling includes `.golangci.yml` with standard linters and a complete `Makefile` with `test`, `lint`, and `fmt` targets.

## ADDED Requirements

### Requirement: .golangci.yml with standard linters
The project root SHALL contain a `.golangci.yml` file configuring at minimum: `errcheck`, `gosimple`, `govet`, `ineffassign`, `staticcheck`, `unused`, `gofmt`, `goimports`, `misspell`, `unconvert`, and `unparam`. It SHALL enable `check-type-assertions` for `errcheck` and `check-shadowing` for `govet`.

#### Scenario: Linter config file exists
- **WHEN** the project root is inspected
- **THEN** `.golangci.yml` is present with the specified linters enabled

#### Scenario: golangci-lint runs without errors on clean code
- **WHEN** `golangci-lint run ./...` is executed on the refactored codebase
- **THEN** all enabled linters pass without warnings

### Requirement: Makefile with test, lint, and fmt targets
The `Makefile` SHALL include `.PHONY` targets: `build` (existing), `install` (existing), `test` (running `go test -race -cover ./...`), `lint` (running `golangci-lint run ./...`), `fmt` (running `goimports -w .`), and `clean` (existing).

#### Scenario: make test runs all tests with race detection
- **WHEN** `make test` is executed
- **THEN** all packages are tested with `-race` and `-cover` flags

#### Scenario: make lint checks all packages
- **WHEN** `make lint` is executed
- **THEN** `golangci-lint run ./...` is invoked

#### Scenario: make fmt formats all Go files
- **WHEN** `make fmt` is executed
- **THEN** all Go source files are formatted with `goimports -w .`
