## 1. Extract mode-specific load methods

- [x] 1.1 Extract `loadViewportForIndex()` from `loadViewport()` (ModeIndex branch, sync refresh)
- [x] 1.2 Extract `loadViewportForConfig()` from `loadViewport()` (ModeViewingConfig branch, configToMarkdown + glamour)
- [x] 1.3 Extract `loadViewportForSpec()` from `loadViewport()` (ModeViewingSpec branch, requirement extraction + glamour)
- [x] 1.4 Extract `loadViewportForTasks()` from `loadViewport()` (TabTasks + ModeNormal branch, sync task rendering)
- [x] 1.5 Extract `loadViewportForArtifact()` from `loadViewport()` (cache check + glamour for proposal/design/specs)

## 2. Wire up delegation

- [x] 1.6 Update `loadViewport()` to delegate to the appropriate mode-specific method

## 3. Verify

- [x] 1.7 Run `go test ./internal/...` and verify all existing tests pass
