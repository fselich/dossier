## 1. Upgrade Go modules

- [x] 1.1 Change `go.mod` imports: `bubbletea` → `charm.land/bubbletea/v2`, `bubbles` → `charm.land/bubbles/v2`, `lipgloss` → `charm.land/lipgloss/v2`, `glamour` → `charm.land/glamour/v2`
- [x] 1.2 Run `go mod tidy` to resolve v2 dependency trees
- [x] 1.3 Update all `.go` file import statements to use v2 import paths

## 2. Adapt cmd/dossier/main.go

- [x] 2.1 Remove `tea.WithAltScreen()` and `tea.WithMouseCellMotion()` from `tea.NewProgram`
- [x] 2.2 Update `Init()` usage if needed (v2 `Init` signature unchanged)

## 3. Adapt Model and View signature

- [x] 3.1 Change `View() string` → `View() tea.View` on `Model`
- [x] 3.2 Wrap existing layout string in `tea.NewView(content)` and set `AltScreen = true`, `MouseMode = tea.MouseModeCellMotion`, `BackgroundColor` from theme
- [x] 3.3 Update `viewConfig()` and `viewIndex()` renamed to `viewConfigContent()`/`viewIndexContent()` returning string, wrapped in `View()`
- [x] 3.4 Remove `renderWithBackground()` — background handled by `tea.View.BackgroundColor`

## 4. Adapt key handling (update.go)

- [x] 4.1 Replace `tea.KeyMsg` with `tea.KeyPressMsg` in the type switch
- [x] 4.2 Verify all `msg.String()` key names still match (space returns `"space"` instead of `" "` in v2)
- [x] 4.3 Replace `tea.ExecProcess` with `tea.Exec` for editor launch (kept `tea.ExecProcess` — it exists in v2)

## 5. Adapt mouse handling (mouse.go)

- [x] 5.1 Replace `tea.MouseMsg` type switch with separate cases for `tea.MouseClickMsg`, `tea.MouseWheelMsg`
- [x] 5.2 Refactor `handleMouse(tea.MouseMsg)` into `handleMouseClick(tea.MouseClickMsg)` and `handleMouseWheel(tea.MouseWheelMsg)`
- [x] 5.3 Mouse button constants (`tea.MouseLeft`, `tea.MouseWheelUp`, `tea.MouseWheelDown`) — compatible with v2

## 6. Adapt viewport usage (bubbles v2)

- [x] 6.1 Change `viewport.New(width, height)` → `viewport.New(viewport.WithWidth(w), viewport.WithHeight(h))`
- [x] 6.2 Replace all `vp.Width` / `vp.Height` field access with `vp.SetWidth()` / `vp.Width()` / `vp.SetHeight()` / `vp.Height()` / `vp.SetYOffset()` / `vp.YOffset()`
- [x] 6.3 `LineUp/LineDown` → `ScrollUp/ScrollDown`. Updated all files: `model.go`, `viewport.go`, `tasks.go`, `index.go`, `update.go`, `mouse.go`

## 7. Adapt lipgloss v2

- [x] 7.1 `lipgloss.Color("12")` returns `color.Color`; `lipgloss.NoColor{}` implements `color.Color` — compatible with v2
- [x] 7.2 `tea.View.BackgroundColor` field type is `color.Color`; `Theme.ViewBg` changed from `lipgloss.TerminalColor` to `color.Color`

## 8. Adapt glamour v2

- [x] 8.1 Update glamour import in `viewport.go` to `charm.land/glamour/v2`
- [x] 8.2 `glamour.NewTermRenderer` API unchanged — compatible

## 9. Validation

- [x] 9.1 Run `go build ./...` — compiles without errors
- [x] 9.2 Run `go vet ./...` — passes
- [x] 9.3 Run `go test -race ./...` — all tests pass
- [x] 9.4 Manual smoke test: launch TUI, navigate tabs, scroll, use mouse, open editor (`e`), verify mouse works after editor
