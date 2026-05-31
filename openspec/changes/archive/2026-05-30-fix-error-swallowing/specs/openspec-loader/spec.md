## MODIFIED Requirements

### Requirement: LoadConfig returns error on failure
`LoadConfig()` and `LoadConfigFrom()` SHALL return `(ProjectConfig, error)`. If the config file does not exist, they SHALL return an empty config and nil error. If the file exists but cannot be read for another reason (e.g., permission denied, corrupt filesystem), they SHALL return the error. If the YAML is malformed, they SHALL return an error.

#### Scenario: Config file read failure (not IsNotExist)
- **WHEN** `LoadConfigFrom()` is called and `openspec/config.yaml` exists but `os.ReadFile` fails with a non-IsNotExist error (e.g., permission denied)
- **THEN** the function returns an empty `ProjectConfig` and the error from `os.ReadFile`

## ADDED Requirements

### Requirement: Error logging for non-critical loader calls
Callers of `ListArchiveChangesFrom` and `LoadProjectSpecsFrom` SHALL not discard errors silently. When these functions return a non-nil error, the caller SHALL log it so the failure is visible.

#### Scenario: Archive load fails
- **WHEN** `ListArchiveChangesFrom()` returns a non-nil error
- **THEN** the caller logs the error to stderr and continues with an empty archive list
