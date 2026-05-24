# openspec-root-path Specification

## Purpose
All public loader functions in `internal/openspec` accept an explicit root path parameter, decoupling them from `os.Getwd()` and enabling testability without global state manipulation.

## ADDED Requirements

### Requirement: Load project from explicit root path
The loader SHALL expose a `LoadFrom(root string) (*Project, error)` function that loads the openspec project structure from the given root directory path. It SHALL behave identically to `Load()` but use `root` instead of `os.Getwd()`.

#### Scenario: Valid root path with openspec directory
- **WHEN** `LoadFrom("/path/to/project")` is called and `/path/to/project/openspec/` exists
- **THEN** the function returns a `*Project` with changes loaded from `openspec/changes/` and no error

#### Scenario: Root path without openspec directory
- **WHEN** `LoadFrom("/path/to/empty")` is called and the directory has no `openspec/` subdirectory
- **THEN** the function returns an error containing "no openspec/ directory found"

### Requirement: Load project config from explicit root path
The loader SHALL expose a `LoadConfigFrom(root string) (ProjectConfig, error)` function that reads `openspec/config.yaml` from the given root path.

#### Scenario: Config file exists at root path
- **WHEN** `LoadConfigFrom("/path/to/project")` is called and `openspec/config.yaml` exists with valid YAML
- **THEN** the function returns the parsed `ProjectConfig` and nil error

#### Scenario: Config file missing at root path
- **WHEN** `LoadConfigFrom("/path/to/project")` is called and `openspec/config.yaml` does not exist
- **THEN** the function returns an empty `ProjectConfig` and nil error (missing config is not an error)

### Requirement: Load project specs from explicit root path
The loader SHALL expose a `LoadProjectSpecsFrom(root string) ([]ProjectSpec, error)` function that reads all spec subdirectories under `openspec/specs/` from the given root path, sorted alphabetically by name.

#### Scenario: Specs directory with multiple subdirectories
- **WHEN** `LoadProjectSpecsFrom("/path/to/project")` is called and `openspec/specs/` contains `auth/` and `api/` subdirectories each with a `spec.md`
- **THEN** the function returns two `ProjectSpec` entries sorted alphabetically, and nil error

#### Scenario: Specs directory does not exist
- **WHEN** `LoadProjectSpecsFrom("/path/to/project")` is called and `openspec/specs/` does not exist
- **THEN** the function returns nil slice and nil error (missing directory is not an error)

### Requirement: List functions accept explicit root path
The loader SHALL expose `*From(root string)` variants for all list functions:
- `ListChangeNamesFrom(root string) ([]string, error)`
- `ListArchiveChangesFrom(root string) ([]Change, error)`
- `ListArchiveNamesFrom(root string) ([]string, error)`
- `ListSpecNamesFrom(root string) ([]string, error)`

Each SHALL operate relative to the given `root` path instead of `os.Getwd()`.

#### Scenario: ListChangeNamesFrom with active changes
- **WHEN** `ListChangeNamesFrom("/path/to/project")` is called and `openspec/changes/` contains `feat-a/` and `feat-b/`
- **THEN** the function returns `["feat-a", "feat-b"]` and nil error

#### Scenario: ListArchiveChangesFrom with no archived changes
- **WHEN** `ListArchiveChangesFrom("/path/to/project")` is called and `openspec/changes/archive/` is empty or does not exist
- **THEN** the function returns nil slice and nil error
