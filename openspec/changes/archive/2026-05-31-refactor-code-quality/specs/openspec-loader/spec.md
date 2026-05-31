# openspec-loader Delta Specification

## Purpose
Loader functions are converted to methods on a `*Loader` struct with injected filesystem dependency. Backward-compatible wrapper functions preserve the existing public API.

## MODIFIED Requirements

### Requirement: Descubrir openspec desde CWD
The loader SHALL look for the `openspec/` directory relative to the current working directory on startup. The `Load()` function SHALL delegate to `LoadFrom(os.Getwd())` internally. If the directory does not exist or if `os.Getwd()` fails, the loader SHALL return an error. The implementation SHALL be a method on `*Loader` with a package-level wrapper delegating to a default `Loader` using the real filesystem.

#### Scenario: openspec presente
- **WHEN** `dossier` is run in a directory that contains `openspec/`
- **THEN** the loader loads the structure without error

#### Scenario: openspec ausente
- **WHEN** `dossier` is run in a directory without `openspec/`
- **THEN** the program terminates with an error message indicating the openspec directory was not found

### Requirement: Releer artifacts de un change en disco
The loader SHALL expose a method on `*Loader` that, given an already-loaded `Change`, rereads from disk the content of its artifacts (`proposal.md`, `design.md`, `tasks.md`, `specs/*/spec.md`) and returns a new `Change` with the updated content. If a file does not exist or cannot be read, the corresponding artifact SHALL be marked as absent without returning an error. A package-level wrapper function `ReloadChange(ch Change) Change` SHALL delegate to `defaultLoader.ReloadChange(ch)`.

#### Scenario: Contenido de tasks.md actualizado en disco
- **WHEN** `tasks.md` has been externally modified since the last load
- **THEN** the function returns a `Change` with `Tasks.Content` equal to the new file content

#### Scenario: Archivo eliminado entre recargas
- **WHEN** `design.md` existed in the previous load but has been deleted
- **THEN** the function returns a `Change` with `Design.Present = false` and `Design.Content = ""`

#### Scenario: Sin cambios en disco
- **WHEN** no file in the change has changed since the last load
- **THEN** the function returns a `Change` with the same content as the original

### Requirement: Error logging for non-critical loader calls
Callers of loader functions SHALL surface errors through the UI's error display mechanism (`m.errMsg`) rather than silent `log.Printf`. Loader functions themselves SHALL continue to return errors normally.

#### Scenario: Archive load fails
- **WHEN** `ListArchiveChangesFrom()` returns a non-nil error
- **THEN** the caller sets `m.errMsg` with the error message, displayed to the user and auto-cleared after 3 seconds

#### Scenario: Specs load fails
- **WHEN** `LoadProjectSpecsFrom()` returns a non-nil error
- **THEN** the caller sets `m.errMsg` with the error message, displayed to the user and auto-cleared after 3 seconds
