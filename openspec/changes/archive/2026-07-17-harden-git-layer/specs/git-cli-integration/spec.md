## ADDED Requirements

### Requirement: Robust porcelain status parsing

The git layer SHALL obtain working-tree status via `git status --porcelain=v1 -z -u` and parse NUL-separated entries. Parsing SHALL be byte-exact for any filename git can produce, including names containing spaces, newlines, quotes, and non-ASCII characters. For renamed or copied entries (`X` or `Y` is `R` or `C`), the entry token SHALL provide the new path and the following NUL token SHALL be consumed as the old path. Files under `openspec/` SHALL continue to be excluded.

#### Scenario: Filename containing a newline parses as one entry
- **WHEN** the working tree contains a modified file whose name contains a newline character
- **THEN** `Status` returns exactly one entry for that file with the full, unmodified path

#### Scenario: Rename produces old and new paths
- **WHEN** a file is renamed and staged (`R ` status)
- **THEN** the returned entry has `Path` set to the new name and `OldPath` set to the old name

#### Scenario: Standard statuses parse unchanged
- **WHEN** the working tree contains modified, added, untracked, and deleted files
- **THEN** `Status` returns entries with the same `X`, `Y`, `Path`, and `IsDeleted` values as the previous line-based parser

#### Scenario: Clean tree returns no entries
- **WHEN** the working tree is clean
- **THEN** `Status` returns an empty result and no error

### Requirement: WorkTreeRoot reports failure explicitly

`WorkTreeRoot` SHALL return `(string, error)`. When `git rev-parse --show-toplevel` fails, it SHALL return a non-nil error and SHALL NOT silently return the input path. The UI caller SHALL fall back to the openspec root when an error is returned.

#### Scenario: Inside a worktree
- **WHEN** `WorkTreeRoot` is called with a directory inside a git worktree
- **THEN** it returns the absolute worktree root and a nil error

#### Scenario: Outside a worktree
- **WHEN** `WorkTreeRoot` is called with a directory that is not inside a git worktree
- **THEN** it returns a non-nil error

### Requirement: Bounded execution time for git subprocesses

Every git subprocess invocation (`IsInsideWorkTree`, `WorkTreeRoot`, `Status`, and diff retrieval) SHALL be executed with a timeout of 2 seconds. When the timeout elapses, the invocation SHALL fail with an error and the process SHALL be terminated; the UI SHALL keep its last known state.

#### Scenario: Hung git does not freeze the UI
- **WHEN** a git invocation blocks for longer than the timeout
- **THEN** the call returns an error within approximately 2 seconds and the TUI remains responsive, keeping the previously displayed status

### Requirement: Test coverage for the git layer

The `internal/git` package SHALL have automated tests exercising real git repositories created in temporary directories, covering: modified, added, untracked, renamed, copied, and deleted files; the `openspec/` filter; unusual filenames; and worktree detection for repo and non-repo directories. Tests SHALL skip gracefully when the `git` binary is unavailable.

#### Scenario: Status table tests pass against a real repo
- **WHEN** the test suite runs on a machine with git installed
- **THEN** `Status`, `IsInsideWorkTree`, and `WorkTreeRoot` are verified against repositories built in `t.TempDir()`

#### Scenario: Git binary unavailable
- **WHEN** the test suite runs on a machine without git on PATH
- **THEN** git-layer tests are skipped, not failed
