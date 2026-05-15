# openspec-loader Specification

## Purpose
Loads and exposes the openspec project structure from the working directory: active and archived changes, their artifacts (`proposal`, `design`, `tasks`, `specs`) and the project specs, with support for on-demand reloading.

## Requirements


### Requirement: Descubrir openspec desde CWD
The loader SHALL look for the `openspec/` directory relative to the current working directory on startup. If it does not exist, it SHALL terminate with a clear error message.

#### Scenario: openspec presente
- **WHEN** `dossier` is run in a directory that contains `openspec/`
- **THEN** the loader loads the structure without error

#### Scenario: openspec ausente
- **WHEN** `dossier` is run in a directory without `openspec/`
- **THEN** the program terminates with the message `"No openspec/ directory found in current directory"`

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
The loader SHALL expose a `ListArchiveNames()` function that returns only the names of the subdirectories of `openspec/changes/archive/`, sorted from most recent to oldest (descending alphabetical order by directory name), without reading any files inside them. If the directory does not exist or is empty, it SHALL return an empty list without error.

#### Scenario: Orden descendente garantizado
- **WHEN** `openspec/changes/archive/` contains directories with different date prefixes
- **THEN** `ListArchiveNames()` returns the names in descending order, matching the order of `ListArchiveChanges()`

#### Scenario: Directorio archive ausente
- **WHEN** `openspec/changes/archive/` does not exist
- **THEN** `ListArchiveNames()` returns an empty list without error

### Requirement: Listar nombres de specs del proyecto
The loader SHALL expose a `ListSpecNames()` function that returns only the names of the subdirectories of `openspec/specs/`, sorted alphabetically, without reading any files inside them. If the directory does not exist or is empty, it SHALL return an empty list without error.

#### Scenario: Orden alfabético garantizado
- **WHEN** `openspec/specs/` contains three subdirectories in non-alphabetical creation order
- **THEN** `ListSpecNames()` returns the names in ascending alphabetical order

#### Scenario: Directorio specs ausente
- **WHEN** `openspec/specs/` does not exist
- **THEN** `ListSpecNames()` returns an empty list without error
