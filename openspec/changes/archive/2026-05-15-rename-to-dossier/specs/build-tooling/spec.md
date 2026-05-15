## MODIFIED Requirements

### Requirement: Binary named dossier
The project SHALL produce a binary named `dossier`. The entry point directory SHALL be `cmd/dossier/` so that `go install` names the binary by convention.

#### Scenario: go install produces correct binary name
- **WHEN** the developer runs `go install ./cmd/dossier/`
- **THEN** a binary named `dossier` is placed in `$GOPATH/bin`

#### Scenario: Binary is executable from PATH
- **WHEN** `$GOPATH/bin` is in `$PATH` and `make install` has been run
- **THEN** `dossier` is available as a command from any directory

### Requirement: Makefile build target
The project SHALL provide a `make build` target that compiles the application into a local binary named `dossier` in the project root.

#### Scenario: Local build
- **WHEN** the developer runs `make build`
- **THEN** a `dossier` binary is created in the project root

### Requirement: Makefile install target
The project SHALL provide a `make install` target that compiles and installs the binary to `$GOPATH/bin`.

#### Scenario: System install
- **WHEN** the developer runs `make install`
- **THEN** `go install ./cmd/dossier/` is executed and `dossier` is available system-wide

### Requirement: Makefile clean target
The project SHALL provide a `make clean` target that removes compiled binaries from the project root.

#### Scenario: Cleanup removes local binary
- **WHEN** the developer runs `make clean`
- **THEN** the `dossier` binary in the project root is deleted (if present)

### Requirement: No stale binaries in repository root
The project root SHALL NOT contain committed or untracked compiled binaries. A `.gitignore` entry SHALL prevent accidental commits of compiled output.

#### Scenario: Compiled binaries are ignored by git
- **WHEN** the developer builds the project
- **THEN** `git status` does not show `dossier`, `main`, or `sv` as untracked files
