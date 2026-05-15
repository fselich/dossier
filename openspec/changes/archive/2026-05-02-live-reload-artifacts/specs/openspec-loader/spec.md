## ADDED Requirements

### Requirement: Re-read a change's artifacts from disk
The loader SHALL expose a function that, given an already-loaded `Change`, re-reads the content of its artifacts (`proposal.md`, `design.md`, `tasks.md`, `specs/*/spec.md`) from disk and returns a new `Change` with the updated content. If a file does not exist or cannot be read, the corresponding artifact SHALL be marked as absent without returning an error.

#### Scenario: tasks.md content updated on disk
- **WHEN** `tasks.md` has been modified externally since the last load
- **THEN** the function returns a `Change` with `Tasks.Content` equal to the new content of the file

#### Scenario: File deleted between reloads
- **WHEN** `design.md` existed in the previous load but has been deleted
- **THEN** the function returns a `Change` with `Design.Present = false` and `Design.Content = ""`

#### Scenario: No changes on disk
- **WHEN** no file of the change has changed since the last load
- **THEN** the function returns a `Change` with the same content as the original
