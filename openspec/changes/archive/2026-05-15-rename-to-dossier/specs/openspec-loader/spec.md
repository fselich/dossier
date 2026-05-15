## MODIFIED Requirements

### Requirement: Descubrir openspec desde CWD
The loader SHALL look for the `openspec/` directory relative to the current working directory on startup. If it does not exist, it SHALL terminate with a clear error message.

#### Scenario: openspec presente
- **WHEN** `dossier` is run in a directory that contains `openspec/`
- **THEN** the loader loads the structure without error

#### Scenario: openspec ausente
- **WHEN** `dossier` is run in a directory without `openspec/`
- **THEN** the program terminates with the message `"No openspec/ directory found in current directory"`
