## 1. Extract handleIndexModeKeys

- [x] 1.1 Extract `handleIndexModeKeys(msg tea.KeyPressMsg) (Model, tea.Cmd)` from `Update()`: move `j`/`k`/`space`/`s`/`enter` cases that apply when `m.mode == ModeIndex`
- [x] 1.2 Use `commitStateChange()` where applicable for viewport height + reload boilerplate

## 2. Extract handleNormalModeKeys

- [x] 2.1 Extract `handleNormalModeKeys(msg tea.KeyPressMsg) (Model, tea.Cmd)` from `Update()`: move `h`/`l`/`j`/`k`/`space`/`e`/`1`-`4`/`tab`/`shift+tab` cases that apply when `m.mode == ModeNormal`
- [x] 2.2 Use `commitStateChange()` where applicable for viewport height + reload boilerplate

## 3. Extract handleArchiveModeKeys

- [x] 3.1 Extract `handleArchiveModeKeys(msg tea.KeyPressMsg) (Model, tea.Cmd)` from `Update()`: move key handlers that apply when `m.mode == ModeViewingArchive` (tab switching keys)
- [x] 3.2 Use `commitStateChange()` where applicable

## 4. Extract handleSpecModeKeys

- [x] 4.1 Extract `handleSpecModeKeys(msg tea.KeyPressMsg) (Model, tea.Cmd)` from `Update()`: move `h`/`l` cases that apply when `m.mode == ModeViewingSpec`
- [x] 4.2 Use `commitStateChange()` where applicable

## 5. Extract commitStateChange helper

- [x] 5.1 Add `commitStateChange() (Model, tea.Cmd)` method that wraps `m.vp.SetHeight(m.contentHeight())` + `return *m, m.loadViewport()`
- [x] 5.2 Replace all ~15 instances of the boilerplate across mode handlers with calls to `commitStateChange()`

## 6. Update Update() to delegate

- [x] 6.1 Replace the mode-specific key cases in `Update()` with a mode dispatch switch: `ModeIndex → handleIndexModeKeys`, `ModeNormal → handleNormalModeKeys`, `ModeViewingArchive → handleArchiveModeKeys`, `ModeViewingSpec → handleSpecModeKeys`
- [x] 6.2 Keep shared keys (`q`, `ctrl+c`, `i`, `a`, `esc`, `enter`) at the top level in `Update()`

## 7. Verify tests pass

- [x] 7.1 Run `go build ./...` to confirm no compilation errors
- [x] 7.2 Run `go test ./internal/ui/...` and confirm all tests pass
