## MODIFIED Requirements

### Requirement: Invocación con ruta explícita
The binary SHALL accept an optional positional argument: the path to a change directory. If provided, it SHALL load only that change and display it in the TUI without scanning `openspec/changes/`. If not provided, the behavior SHALL be identical to the current behavior.

#### Scenario: Invocación sin argumentos
- **WHEN** the user runs `./dossier` with no arguments
- **THEN** the TUI starts in normal mode loading all active changes from `openspec/changes/`

#### Scenario: Invocación con ruta a change activo
- **WHEN** the user runs `./dossier ./openspec/changes/mi-change`
- **THEN** the TUI shows only the artifacts of `mi-change`

#### Scenario: Invocación con ruta a change archivado
- **WHEN** the user runs `./dossier ./openspec/changes/archive/2026-05-02-mi-change`
- **THEN** the TUI shows only the artifacts of the archived change, with tab navigation the same as in normal mode

### Requirement: Validación de ruta
If a path is provided and does not correspond to a valid change, the binary SHALL print a descriptive error message to stderr and exit with exit code 1, without opening the TUI.

#### Scenario: Ruta inexistente
- **WHEN** the user runs `./dossier ./ruta/que/no/existe`
- **THEN** the binary prints `"error: path not found: ./ruta/que/no/existe"` and exits with code 1

#### Scenario: Ruta sin .openspec.yaml
- **WHEN** the user runs `./dossier ./algún/directorio` and that directory does not contain `.openspec.yaml`
- **THEN** the binary prints `"error: not a valid change directory (missing .openspec.yaml)"` and exits with code 1
