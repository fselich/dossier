# openspec-loader Specification

## Purpose
Loads and exposes the openspec project structure from the working directory: active and archived changes, their artifacts (`proposal`, `design`, `tasks`, `specs`) and the project specs, with support for on-demand reloading.

## Requirements


### Requirement: Descubrir openspec desde CWD
The loader SHALL look for the `openspec/` directory relative to the current working directory on startup. The `Load()` function SHALL delegate to `LoadFrom(os.Getwd())` internally. If the directory does not exist or if `os.Getwd()` fails, the loader SHALL return an error.

#### Scenario: openspec presente
- **WHEN** `dossier` is run in a directory that contains `openspec/`
- **THEN** the loader loads the structure without error

#### Scenario: openspec ausente
- **WHEN** `dossier` is run in a directory without `openspec/`
- **THEN** the program terminates with an error message indicating the openspec directory was not found

### Requirement: Listar changes activos
The loader SHALL consider as active changes all direct subdirectories of `openspec/changes/` except `archive/`. It SHALL ignore any loose files at that level.

#### Scenario: Múltiples changes activos
- **WHEN** `openspec/changes/feat-a/` and `openspec/changes/feat-b/` exist
- **THEN** the loader returns both as active changes

#### Scenario: Solo archive
- **WHEN** the only subdirectory is `openspec/changes/archive/`
- **THEN** the loader returns an empty list of active changes

#### Scenario: Sin directorio changes
- **WHEN** `openspec/changes/` does not exist
- **THEN** the loader returns an empty list of active changes without error

### Requirement: Leer metadatos del change
The loader SHALL read `openspec/changes/<name>/.openspec.yaml` to obtain the creation date (`created`). If the file does not exist, it SHALL use an empty string for the date.

#### Scenario: Metadatos presentes
- **WHEN** `.openspec.yaml` contains `created: 2026-05-01`
- **THEN** the change exposes `Created = "2026-05-01"`

### Requirement: Cargar artifacts disponibles
The loader SHALL attempt to read each of the four artifacts (`proposal.md`, `design.md`, `tasks.md`, `specs/<cap>/spec.md`) for each change. It SHALL mark as absent any artifact whose file does not exist, without returning an error.

#### Scenario: Artifact ausente
- **WHEN** `design.md` does not exist in the change
- **THEN** the `design` artifact is marked as absent and produces no error

#### Scenario: Todos los artifacts presentes
- **WHEN** `proposal.md`, `design.md`, `tasks.md`, and at least one `specs/*/spec.md` exist
- **THEN** all artifacts are loaded with their content

### Requirement: Inferir nombre del proyecto
The loader SHALL use the name of the current working directory as the project name, without reading `openspec/config.yaml`.

#### Scenario: Directorio con nombre simple
- **WHEN** the CWD is `/home/user/my-project`
- **THEN** the project name is `"my-project"`

### Requirement: Releer artifacts de un change en disco
The loader SHALL expose a function that, given an already-loaded `Change`, rereads from disk the content of its artifacts (`proposal.md`, `design.md`, `tasks.md`, `specs/*/spec.md`) and returns a new `Change` with the updated content. If a file does not exist or cannot be read, the corresponding artifact SHALL be marked as absent without returning an error.

#### Scenario: Contenido de tasks.md actualizado en disco
- **WHEN** `tasks.md` has been externally modified since the last load
- **THEN** the function returns a `Change` with `Tasks.Content` equal to the new file content

#### Scenario: Archivo eliminado entre recargas
- **WHEN** `design.md` existed in the previous load but has been deleted
- **THEN** the function returns a `Change` with `Design.Present = false` and `Design.Content = ""`

#### Scenario: Sin cambios en disco
- **WHEN** no file in the change has changed since the last load
- **THEN** the function returns a `Change` with the same content as the original

### Requirement: Listar nombres de changes archivados
The loader SHALL expose a `ListArchiveNames()` function that returns the names and an error. It SHALL return only the names of the subdirectories of `openspec/changes/archive/`, sorted from most recent to oldest, without reading any files inside them. If the directory does not exist, it SHALL return an empty list and nil error. If a read error occurs, the error SHALL be propagated.

#### Scenario: Orden descendente garantizado
- **WHEN** `openspec/changes/archive/` contains directories with different date prefixes
- **THEN** `ListArchiveNames()` returns the names in descending order and nil error

#### Scenario: Directorio archive ausente
- **WHEN** `openspec/changes/archive/` does not exist
- **THEN** `ListArchiveNames()` returns an empty list and nil error

### Requirement: Listar nombres de specs del proyecto
The loader SHALL expose a `ListSpecNames()` function that returns the names and an error. It SHALL return only the names of the subdirectories of `openspec/specs/`, sorted alphabetically, without reading any files inside them. If the directory does not exist, it SHALL return nil and nil error. If a read error occurs for another reason, the error SHALL be propagated.

#### Scenario: Orden alfabético garantizado
- **WHEN** `openspec/specs/` contains three subdirectories in non-alphabetical creation order
- **THEN** `ListSpecNames()` returns the names in ascending alphabetical order and nil error

#### Scenario: Directorio specs ausente
- **WHEN** `openspec/specs/` does not exist
- **THEN** `ListSpecNames()` returns nil and nil error

### Requirement: LoadConfig returns error on failure
`LoadConfig()` and `LoadConfigFrom()` SHALL return `(ProjectConfig, error)`. If the config file does not exist, they SHALL return an empty config and nil error. If the file exists but cannot be read for another reason (e.g., permission denied, corrupt filesystem), they SHALL return the error. If the YAML is malformed, they SHALL return an error.

#### Scenario: Config file with valid YAML
- **WHEN** `LoadConfig()` is called and `openspec/config.yaml` contains valid YAML
- **THEN** the function returns the parsed `ProjectConfig` and nil error

#### Scenario: Config file with invalid YAML
- **WHEN** `LoadConfig()` is called and `openspec/config.yaml` contains malformed YAML
- **THEN** the function returns an empty `ProjectConfig` and a non-nil error

#### Scenario: Config file read failure (not IsNotExist)
- **WHEN** `LoadConfigFrom()` is called and `openspec/config.yaml` exists but `os.ReadFile` fails with a non-IsNotExist error (e.g., permission denied)
- **THEN** the function returns an empty `ProjectConfig` and the error from `os.ReadFile`

### Requirement: Error logging for non-critical loader calls
Callers of `ListArchiveChangesFrom` and `LoadProjectSpecsFrom` SHALL not discard errors silently. When these functions return a non-nil error, the caller SHALL log it so the failure is visible.

#### Scenario: Archive load fails
- **WHEN** `ListArchiveChangesFrom()` returns a non-nil error
- **THEN** the caller logs the error to stderr and continues with an empty archive list

### Requirement: LoadProjectSpecs returns error on failure
`LoadProjectSpecs()` SHALL return `([]ProjectSpec, error)`. If the `openspec/specs/` directory does not exist, it SHALL return nil and nil error. If a read error occurs for another reason, the error SHALL be propagated.

#### Scenario: Specs directory readable
- **WHEN** `LoadProjectSpecs()` is called and `openspec/specs/` exists with subdirectories
- **THEN** the function returns the list of specs and nil error

#### Scenario: Specs directory missing
- **WHEN** `LoadProjectSpecs()` is called and `openspec/specs/` does not exist
- **THEN** the function returns nil and nil error

### Requirement: ListArchiveChanges returns error on failure
`ListArchiveChanges()` SHALL return `([]Change, error)`. If the archive directory does not exist, it SHALL return nil and nil error. If a read error occurs, the error SHALL be propagated.

#### Scenario: Archive directory exists
- **WHEN** `ListArchiveChanges()` is called and `openspec/changes/archive/` exists
- **THEN** the function returns the list and nil error

### Requirement: ListChangeNames returns error on failure
`ListChangeNames()` SHALL return `([]string, error)`. If the changes directory does not exist, it SHALL return nil and nil error. If a read error occurs for another reason, the error SHALL be propagated.

#### Scenario: Changes directory exists with entries
- **WHEN** `ListChangeNames()` is called and `openspec/changes/` contains active change subdirectories
- **THEN** the function returns the sorted list and nil error

#### Scenario: Changes directory missing
- **WHEN** `ListChangeNames()` is called and `openspec/changes/` does not exist
- **THEN** the function returns nil and nil error

### Requirement: Ordenar cambios activos por fecha de creación
`LoadFrom()` SHALL sort active changes by their `created` date in descending order (newest first). Changes without a `created` date SHALL be placed after all dated changes and sorted alphabetically by name (stable sort as tiebreaker). Changes with equal `created` dates SHALL keep their relative input order (stable sort).

#### Scenario: Cambios con fechas distintas
- **WHEN** two active changes have `created` dates `2026-05-01` and `2026-05-10`
- **THEN** the change from `2026-05-10` appears before the change from `2026-05-01` in the returned list

#### Scenario: Cambios sin fecha van al final
- **WHEN** one active change has no `created` date and another has `2026-05-01`
- **THEN** the dated change appears first and the undated change appears after, sorted alphabetically

#### Scenario: Cambios con la misma fecha
- **WHEN** two changes have the same `created` date
- **THEN** their relative order is stable (preserved from directory listing)
