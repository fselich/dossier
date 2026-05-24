# openspec-loader Delta Specification

## MODIFIED Requirements

### Requirement: Descubrir openspec desde CWD
The loader SHALL look for the `openspec/` directory relative to the current working directory on startup. If it does not exist or if `os.Getwd()` fails, the loader SHALL return an error. The `Load()` function SHALL delegate to `LoadFrom(os.Getwd())` internally.

#### Scenario: openspec presente
- **WHEN** `dossier` is run in a directory that contains `openspec/`
- **THEN** the loader loads the structure without error

#### Scenario: openspec ausente
- **WHEN** `dossier` is run in a directory without `openspec/`
- **THEN** the program terminates with the message `"No openspec/ directory found in current directory"`

### Requirement: LoadConfig returns error on failure
`LoadConfig()` SHALL return `(ProjectConfig, error)`. If `os.Getwd()` fails, the error SHALL be propagated. If the config file does not exist, it SHALL return an empty config and nil error. If the YAML is malformed, it SHALL return an error.

#### Scenario: Config file with valid YAML
- **WHEN** `LoadConfig()` is called and `openspec/config.yaml` contains valid YAML
- **THEN** the function returns the parsed `ProjectConfig` and nil error

#### Scenario: Config file with invalid YAML
- **WHEN** `LoadConfig()` is called and `openspec/config.yaml` contains malformed YAML
- **THEN** the function returns an empty `ProjectConfig` and a non-nil error

### Requirement: LoadProjectSpecs returns error on failure
`LoadProjectSpecs()` SHALL return `([]ProjectSpec, error)`. If `os.Getwd()` fails or the `openspec/specs/` directory cannot be read, the error SHALL be propagated.

#### Scenario: Specs directory readable
- **WHEN** `LoadProjectSpecs()` is called and `openspec/specs/` exists with subdirectories
- **THEN** the function returns the list of specs and nil error

#### Scenario: os.Getwd fails
- **WHEN** `LoadProjectSpecs()` is called but `os.Getwd()` returns an error
- **THEN** the function returns nil and the error from `os.Getwd()`

### Requirement: ListArchiveChanges returns error on failure
`ListArchiveChanges()` SHALL return `([]Change, error)`. Read errors on the archive directory SHALL be propagated.

#### Scenario: Archive directory exists
- **WHEN** `ListArchiveChanges()` is called and `openspec/changes/archive/` exists
- **THEN** the function returns the list and nil error

#### Scenario: Archive directory does not exist
- **WHEN** `ListArchiveChanges()` is called and `openspec/changes/archive/` does not exist
- **THEN** the function returns nil slice and nil error (missing archive is not an error)

### Requirement: ListArchiveNames returns error on failure
`ListArchiveNames()` SHALL return `([]string, error)`. Read errors on the archive directory SHALL be propagated.

#### Scenario: Archive directory exists with entries
- **WHEN** `ListArchiveNames()` is called and `openspec/changes/archive/` contains subdirectories
- **THEN** the function returns the sorted list and nil error

### Requirement: ListSpecNames returns error on failure
`ListSpecNames()` SHALL return `([]string, error)`. Read errors on the specs directory SHALL be propagated.

#### Scenario: Specs directory exists with entries
- **WHEN** `ListSpecNames()` is called and `openspec/specs/` contains subdirectories
- **THEN** the function returns the sorted list and nil error

### Requirement: ListChangeNames returns error on failure
`ListChangeNames()` SHALL return `([]string, error)`. Read errors on the changes directory SHALL be propagated.

#### Scenario: Changes directory exists with entries
- **WHEN** `ListChangeNames()` is called and `openspec/changes/` contains active change subdirectories
- **THEN** the function returns the sorted list and nil error
