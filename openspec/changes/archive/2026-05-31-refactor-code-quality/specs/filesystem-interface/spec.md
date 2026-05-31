# filesystem-interface Specification

## Purpose

A `fileSystem` interface in `internal/ui/` decouples filesystem access for testability. The `openspec` loader converts to a struct with an injected filesystem dependency, while preserving backward-compatible package-level functions.

## ADDED Requirements

### Requirement: FileSystem interface defined in consumer package
The UI package SHALL define a `fileSystem` interface with methods `ReadFile`, `WriteFile`, `ReadDir`, and `Stat`, matching the signatures of `os.ReadFile`, `os.WriteFile`, `os.ReadDir`, and `os.Stat` respectively. The interface SHALL NOT be exported.

#### Scenario: Interface compiles and is satisfied by os operations
- **WHEN** a struct embeds `*os.File` or wraps `os` functions
- **THEN** the struct satisfies the `fileSystem` interface

#### Scenario: In-memory filesystem satisfies interface
- **WHEN** a test provides a fake filesystem with in-memory data
- **THEN** it satisfies the `fileSystem` interface and can be injected

### Requirement: Loader struct with injected filesystem
The `openspec` package SHALL expose a `Loader` struct with an unexported `fileSystem` field. All existing loader functions (`LoadFrom`, `LoadConfigFrom`, `LoadProjectSpecsFrom`, `ListChangeNamesFrom`, `ListArchiveChangesFrom`, `ListArchiveNamesFrom`, `ListSpecNamesFrom`, `ReloadChange`, `ToggleTask`, `LoadFromPath`) SHALL become methods on `*Loader`.

#### Scenario: Loader created with real filesystem
- **WHEN** `NewLoader(osFS{})` is called with a real OS filesystem adapter
- **THEN** the loader behaves identically to the previous package-level functions

#### Scenario: Loader created with fake filesystem
- **WHEN** `NewLoader(fakeFS{})` is called with an in-memory filesystem
- **THEN** the loader reads and writes through the fake without touching the real disk

### Requirement: Backward-compatible wrapper functions
The `openspec` package SHALL provide package-level wrapper functions with the same signatures as before (`LoadFrom`, `LoadConfigFrom`, etc.) that delegate to a default `Loader` using the real OS filesystem.

#### Scenario: Existing callers compile unchanged
- **WHEN** existing code calls `openspec.LoadFrom(root)` without a Loader
- **THEN** it compiles and behaves identically to before

#### Scenario: Wrapper delegates to default Loader
- **WHEN** `openspec.LoadFrom("/some/path")` is called
- **THEN** the call delegates to `defaultLoader.LoadFrom("/some/path")` using the real filesystem

### Requirement: OS filesystem adapter
The `openspec` package SHALL define an unexported `osFS` struct that wraps `os.ReadFile`, `os.WriteFile`, `os.ReadDir`, and `os.Stat` to satisfy the `fileSystem` interface. This is the default injected dependency.

#### Scenario: osFS.ReadFile delegates to os.ReadFile
- **WHEN** `osFS{}.ReadFile("existing.txt")` is called
- **THEN** it returns the same bytes and error as `os.ReadFile("existing.txt")`
