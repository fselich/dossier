## ADDED Requirements

### Requirement: List archived change names
The loader SHALL expose a function `ListArchiveNames()` that returns only the names of subdirectories in `openspec/changes/archive/`, without reading any files inside them. If the directory does not exist or is empty, it SHALL return an empty list without error.

#### Scenario: Archives present
- **WHEN** `openspec/changes/archive/` contains two subdirectories
- **THEN** `ListArchiveNames()` returns the two names without reading any files inside them

#### Scenario: Absent archive directory
- **WHEN** `openspec/changes/archive/` does not exist
- **THEN** `ListArchiveNames()` returns an empty list without error

### Requirement: List project spec names
The loader SHALL expose a function `ListSpecNames()` that returns only the names of subdirectories in `openspec/specs/`, without reading any files inside them. If the directory does not exist or is empty, it SHALL return an empty list without error.

#### Scenario: Specs present
- **WHEN** `openspec/specs/` contains three subdirectories
- **THEN** `ListSpecNames()` returns the three names without reading any files inside them

#### Scenario: Absent specs directory
- **WHEN** `openspec/specs/` does not exist
- **THEN** `ListSpecNames()` returns an empty list without error
