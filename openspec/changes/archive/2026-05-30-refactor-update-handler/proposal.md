## Why

`Update()` in `update.go` contains a 283-line `tea.KeyPressMsg` switch that handles keypresses for all modes in a single monolithic block. This makes it difficult to navigate and understand which keys are active in each mode. Extracting mode-specific handlers and a common boilerplate helper clarifies the structure without any behavior change.

## What Changes

- Extract `handleIndexModeKeys()`: handles keypresses when `m.mode == ModeIndex`
- Extract `handleNormalModeKeys()`: handles keypresses when `m.mode == ModeNormal`
- Extract `handleArchiveModeKeys()`: handles keypresses when `m.mode == ModeViewingArchive`
- Extract `handleSpecModeKeys()`: handles keypresses when `m.mode == ModeViewingSpec`
- Extract `commitStateChange()`: helper wrapping `m.vp.SetHeight()` + `return m, m.loadViewport()` boilerplate
- `Update()` delegates key handling to the appropriate mode method

## Capabilities

### New Capabilities

*(none — pure internal refactor)*

### Modified Capabilities

*(none — no spec-level behavior changes)*

## Non-goals

- Changing any key binding or behavior
- Refactoring non-keypress message handlers (WindowSizeMsg, tickMsg, etc.)
- Extracting sub-methods beyond the four mode handlers and one helper

## Impact

- `internal/ui/update.go`: keypress switch replaced with delegation to 4 mode-specific methods and 1 helper
- No API changes, no new dependencies
