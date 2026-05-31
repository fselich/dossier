## Why

`LoadConfigFrom` returns `nil` on _all_ errors (line 126), masking non-IsNotExist failures like permission errors. Additionally, `ListArchiveChangesFrom` and `LoadProjectSpecsFrom` errors are silently discarded with `_` in `index.go` and `model.go`, making failures invisible to the user.

## What Changes

- Fix `LoadConfigFrom` to return the error on non-IsNotExist failures instead of `nil`.
- Replace `_` error discarding with a logged/warned error in `index.go` and `model.go` for `ListArchiveChangesFrom` and `LoadProjectSpecsFrom` calls.

## Capabilities

### Modified Capabilities

- `openspec-loader`: Fix error propagation — `LoadConfigFrom` returns the actual error on read failures (not nil), and callers in the UI no longer discard errors silently.

## Impact

- `internal/openspec/loader.go`: `LoadConfigFrom` line 126 changes from `return ProjectConfig{}, nil` to `return ProjectConfig{}, err`.
- `internal/ui/index.go`: Lines 66–67 and 159–161 replace `_` with proper error handling.
- `internal/ui/model.go`: Lines 130–131 replace `_` with proper error handling.
- Existing callers that check `err != nil` will now correctly receive non-nil errors for permission/IO failures.
