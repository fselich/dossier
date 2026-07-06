## 1. Core Implementation

- [x] 1.1 Extract diff-loading logic from `toggleGitDiff()` into a `loadDiffForFile(cursor int)` method
- [x] 1.2 Add `[` and `]` key handlers in `updateViewer()` for diff-view cycling that move `m.gitState.Cursor` (skipping deleted, wrapping) and call `loadDiffForFile`
- [x] 1.3 Update help bar in `view.go` to show `[/]: prev/next file` in diff view

## 2. Verification

- [x] 2.1 Run `make test` to confirm all tests pass
- [x] 2.2 Run `make lint` to confirm no linting issues
