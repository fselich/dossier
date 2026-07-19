## 1. Viewport rendering fix

- [x] 1.1 In `loadViewportForSpec()`, replace `m.vp.SetContent(raw)` with `m.vp.SetContent("")` at `internal/ui/viewport.go:92`
- [x] 1.2 In `loadViewportForConfig()`, replace `m.vp.SetContent(raw)` with `m.vp.SetContent("")` at `internal/ui/viewport.go:64`
- [x] 1.3 In `loadViewportForArtifact()`, replace `m.vp.SetContent(raw)` with `m.vp.SetContent("")` at `internal/ui/viewport.go:176`

## 2. Verification

- [x] 2.1 Run `make test` and confirm no regressions
- [x] 2.2 Run `make lint` and confirm clean
- [x] 2.3 Manual smoke test: open a spec with many requirements from the index, confirm no raw markdown flash
- [x] 2.4 Manual smoke test: open config view (`i` from index), confirm no raw markdown flash
- [x] 2.5 Manual smoke test: switch between change tabs, confirm no raw markdown flash
