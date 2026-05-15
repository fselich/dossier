## Why

The project lacks a standard build workflow: the binary has no consistent name, compiled artifacts accumulate in the repo root, and there is no single command to install the tool system-wide. This makes day-to-day development friction-heavy.

## What Changes

- Rename `cmd/spec-viewer/` to `cmd/specview/` so `go install` produces a binary named `specview` by convention.
- Delete loose compiled binaries (`main`, `sv`) from the project root.
- Add a `Makefile` with `build`, `install`, and `clean` targets as the canonical way to compile and install the app.

## Capabilities

### New Capabilities

- `build-tooling`: A Makefile providing `make build`, `make install`, and `make clean` targets for compiling and installing the `specview` binary.

### Modified Capabilities

<!-- None — no existing spec-level behavior changes. -->

## Impact

- `cmd/spec-viewer/` directory renamed to `cmd/specview/` (import path changes accordingly).
- Two untracked binaries (`main`, `sv`) removed from the project root.
- New `Makefile` added at the project root.
- No runtime behavior changes.
