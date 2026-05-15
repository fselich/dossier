## Why

The TUI can only be opened in "all active changes" mode. There is no way to inspect the artifacts of an archived change without unarchiving it first or opening the files manually. Passing a path as an argument allows using the viewer as a lookup tool for any change, active or archived.

## What Changes

- `main.go` reads `os.Args[1]` if provided and interprets it as the path to a change directory
- When a path is passed, the loader uses it directly instead of scanning `openspec/changes/`
- The UI does not change: the model receives a `Project` with a single change, navigates between tabs as usual
- If the path does not exist or does not contain `.openspec.yaml`, the binary prints an error and exits
- Without argument: current behavior unchanged

## Capabilities

### New Capabilities

- `path-arg`: Invocation mode with explicit path to a change directory

### Modified Capabilities

(none — the UI requirements do not change)

## Impact

- Only affects `cmd/spec-viewer/main.go` and `internal/openspec/loader.go`
- No new dependencies
- Compatible with the current behavior (without argument = no changes)
