## MODIFIED Requirements

### Requirement: Cargar artifacts disponibles
The loader SHALL use an internal `loadChangeFromPath(dirPath, entryName string, isArchived bool) (Change, error)` helper to build a `Change` struct from a directory path and entry name. This helper SHALL be called by `LoadFrom`, `LoadFromPath`, and `ListArchiveChangesFrom` to avoid duplicate construction logic. The helper SHALL read `.openspec.yaml` for the creation date, scan artifact subdirectories, and return the populated `Change`. If `.openspec.yaml` does not exist, it SHALL use an empty string for the date. Each of the four artifacts (`proposal.md`, `design.md`, `tasks.md`, `specs/<cap>/spec.md`) SHALL be marked as absent if the file does not exist, without returning an error.

#### Scenario: Artifact ausente
- **WHEN** `design.md` does not exist in the change
- **THEN** the `design` artifact is marked as absent and produces no error

#### Scenario: Todos los artifacts presentes
- **WHEN** `proposal.md`, `design.md`, `tasks.md`, and at least one `specs/*/spec.md` exist
- **THEN** all artifacts are loaded with their content

#### Scenario: Helper called from all loading functions
- **WHEN** `LoadFrom`, `LoadFromPath`, or `ListArchiveChangesFrom` builds a Change
- **THEN** each SHALL delegate to `loadChangeFromPath` for the common construction logic
