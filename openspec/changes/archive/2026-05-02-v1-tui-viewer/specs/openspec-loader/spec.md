## ADDED Requirements

### Requirement: Discover openspec from CWD
The loader SHALL search for the `openspec/` directory relative to the current working directory on startup. If it does not exist, SHALL exit with a clear error message.

#### Scenario: openspec present
- **WHEN** `spec-viewer` is run in a directory that contains `openspec/`
- **THEN** the loader loads the structure without error

#### Scenario: openspec absent
- **WHEN** `spec-viewer` is run in a directory without `openspec/`
- **THEN** the program exits with the message `"No openspec/ directory found in current directory"`

### Requirement: List active changes
The loader SHALL consider as active changes all direct subdirectories of `openspec/changes/` except `archive/`. SHALL ignore any loose files at that level.

#### Scenario: Multiple active changes
- **WHEN** `openspec/changes/feat-a/` and `openspec/changes/feat-b/` exist
- **THEN** the loader returns both as active changes

#### Scenario: Archive only
- **WHEN** the only subdirectory is `openspec/changes/archive/`
- **THEN** the loader returns an empty list of active changes

#### Scenario: No changes directory
- **WHEN** `openspec/changes/` does not exist
- **THEN** the loader returns an empty list of active changes without error

### Requirement: Read change metadata
The loader SHALL read `openspec/changes/<name>/.openspec.yaml` to obtain the creation date (`created`). If the file does not exist, SHALL use an empty string for the date.

#### Scenario: Metadata present
- **WHEN** `.openspec.yaml` contains `created: 2026-05-01`
- **THEN** the change exposes `Created = "2026-05-01"`

### Requirement: Load available artifacts
The loader SHALL attempt to read each of the four artifacts (`proposal.md`, `design.md`, `tasks.md`, `specs/<cap>/spec.md`) for each change. SHALL mark as absent any artifact whose file does not exist, without returning an error.

#### Scenario: Absent artifact
- **WHEN** `design.md` does not exist in the change
- **THEN** the `design` artifact is marked as absent and produces no error

#### Scenario: All artifacts present
- **WHEN** `proposal.md`, `design.md`, `tasks.md` and at least one `specs/*/spec.md` exist
- **THEN** all artifacts are loaded with their content

### Requirement: Infer project name
The loader SHALL use the name of the current working directory as the project name, without reading `openspec/config.yaml`.

#### Scenario: Directory with simple name
- **WHEN** the CWD is `/home/user/my-project`
- **THEN** the project name is `"my-project"`
