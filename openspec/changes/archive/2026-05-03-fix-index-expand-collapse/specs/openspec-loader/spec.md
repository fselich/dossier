## MODIFIED Requirements

### Requirement: List project spec names
The loader SHALL expose a function `ListSpecNames()` that returns only the names of subdirectories in `openspec/specs/`, sorted alphabetically, without reading any files inside them. If the directory does not exist or is empty, it SHALL return an empty list without error.

#### Scenario: Alphabetical order guaranteed
- **WHEN** `openspec/specs/` contains three subdirectories in non-alphabetical creation order
- **THEN** `ListSpecNames()` returns the names in ascending alphabetical order

#### Scenario: Absent specs directory
- **WHEN** `openspec/specs/` does not exist
- **THEN** `ListSpecNames()` returns an empty list without error

### Requirement: List archived change names
The loader SHALL expose a function `ListArchiveNames()` that returns only the names of subdirectories in `openspec/changes/archive/`, sorted from most recent to oldest (descending alphabetical order by directory name), without reading any files inside them. If the directory does not exist or is empty, it SHALL return an empty list without error.

#### Scenario: Descending order guaranteed
- **WHEN** `openspec/changes/archive/` contains directories with different date prefixes
- **THEN** `ListArchiveNames()` returns the names in descending order, matching the order of `ListArchiveChanges()`

#### Scenario: Absent archive directory
- **WHEN** `openspec/changes/archive/` does not exist
- **THEN** `ListArchiveNames()` returns an empty list without error
