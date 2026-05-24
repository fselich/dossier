## Why

All dependencies in `go.mod` are marked `// indirect`, including those imported directly by the code (`bubbletea`, `glamour`, `lipgloss`, `bubbles`, `yaml.v3`). This is incorrect — `go mod tidy` should mark directly imported modules without the `// indirect` comment. The fix is trivial but important for tooling correctness and dependency auditing.

## What Changes

- Run `go mod tidy` to correctly classify direct vs indirect dependencies
- `go.mod` will show `bubbletea`, `glamour`, `lipgloss`, `bubbles`, and `yaml.v3` as direct dependencies (without `// indirect`)
- No code changes, no breaking changes

## Capabilities

_No functional changes. This is tooling hygiene only._

## Impact

- `go.mod`: dependency classification corrected, no version changes expected
- `go.sum`: potentially updated with correct checksums
