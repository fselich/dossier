## Why

Bubble Tea v2 is a major release with a declarative rendering model, improved keyboard/mouse handling, and a rewritten renderer. The v1 version has a bug where mouse tracking is lost after launching an external editor (`e` key), because `ReleaseTerminal`/`RestoreTerminal` don't save/restore mouse state. In v2, mouse mode is declared in `View()` and re-applied on every frame, fixing this by design.

Additionally, v2 unlocks capabilities that will be valuable for future features: native clipboard, shift+enter key combos, real cursor support in viewport, horizontal scrolling, synchronized updates (no flickering), and color downsampling.

## What Changes

- **BREAKING**: Upgrade `bubbletea` from `v1.3.10` to `v2.0.6`, changing import path to `charm.land/bubbletea/v2`
- **BREAKING**: Upgrade `bubbles` from `v1.0.0` to `v2.1.0`, changing import path to `charm.land/bubbles/v2`
- **BREAKING**: Upgrade `lipgloss` from `v1.1.1` to `v2.0.3`, changing import path to `charm.land/lipgloss/v2`
- **BREAKING**: Upgrade `glamour` from `v1.0.0` to `v2.0.0`, changing import path to `github.com/charmbracelet/glamour/v2`
- **BREAKING**: `View()` returns `tea.View` struct instead of `string`. Alt screen, mouse mode, and other terminal options are set declaratively via fields.
- **BREAKING**: `tea.KeyMsg` replaced by `tea.KeyPressMsg` and `tea.KeyReleaseMsg`. `msg.String()` still works for key matching.
- **BREAKING**: `tea.MouseMsg` split into `tea.MouseClickMsg`, `tea.MouseWheelMsg`, `tea.MouseMotionMsg`, and `tea.MouseReleaseMsg`.
- **BREAKING**: `tea.ExecProcess` removed; use `tea.Exec` (non-blocking with callback) or `tea.ExecProcess` (blocking via `tea.Sequence`).
- **BREAKING**: `viewport.New(width, height)` → `viewport.New(viewport.WithWidth(w), viewport.WithHeight(h))`. Width/Height become getter/setter methods.
- Remove `tea.WithAltScreen()` and `tea.WithMouseCellMotion()` from program initialization.
- Editor return restores mouse state correctly (declarative model).

## Capabilities

### Modified Capabilities

- `tui-viewer`: `View()` signature changes from `string` to `tea.View`. Layout and behavior remain unchanged. Alt screen and background color are set via `tea.View` fields instead of program options.
- `mouse-navigation`: Mouse message types split (`MouseMsg` → `MouseClickMsg`, `MouseWheelMsg`, `MouseMotionMsg`). Mouse mode is set declaratively via `v.MouseMode = tea.MouseModeCellMotion` instead of `tea.WithMouseCellMotion()`.
- `editor-launch`: Editor launch uses `tea.Exec` instead of `tea.ExecProcess`. On return, mouse tracking is preserved because it's declared in `View()` and re-applied every frame.

## Impact

- `go.mod`: 4 dependency upgrades with new import paths
- `cmd/dossier/main.go`: Program initialization, `tea.View` return
- `internal/ui/model.go`: `View()` signature, mouse mode declaration, `tea.KeyMsg` → `tea.KeyPressMsg`, `Init()` updated
- `internal/ui/update.go`: All key/mouse message type switches rewritten (~400 lines)
- `internal/ui/view.go`: `View()` returns `tea.View` struct, all rendering helpers updated
- `internal/ui/viewport.go`: `glamour` import, viewport getter/setter methods
- `internal/ui/mouse.go`: `tea.MouseMsg` → split message types
- `internal/ui/tasks.go`: viewport method calls
- `internal/ui/index.go`: viewport method calls
- `internal/ui/styles.go`: `lipgloss` v2 API changes (if any)
