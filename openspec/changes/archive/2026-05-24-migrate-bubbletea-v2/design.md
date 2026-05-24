## Context

Bubble Tea v2 is a major API redesign. The core change is that `View()` returns a `tea.View` struct instead of `string`, making terminal features declarative rather than imperative. This change cascades into every file that touches the view layer.

The project currently has ~1300 lines of UI code across 9 files in `internal/ui/`. All of them will need changes.

## Goals / Non-Goals

**Goals:**
- Upgrade bubbletea, bubbles, lipgloss, and glamour to their v2 versions
- Adapt all UI code to the v2 API without changing behavior or adding features
- Fix the mouse-after-editor bug (implicit, via declarative mouse mode)
- Keep the existing layout, navigation, and interaction model identical

**Non-Goals:**
- Add new features (search, themes, etc.)
- Refactor Model into sub-models (item 2.3)
- Change the visual appearance
- Optimize performance beyond what v2 provides by default

## Decisions

### Decision 1: View() builds string first, wraps in tea.View

Instead of refactoring every rendering function to return `tea.View`, we keep the string-building approach and wrap at the top level:

```go
func (m Model) View() tea.View {
    content := m.buildLayoutString() // existing logic unchanged
    v := tea.NewView(content)
    v.AltScreen = true
    v.MouseMode = tea.MouseModeCellMotion
    v.BackgroundColor = m.theme.ViewBg
    return v
}
```

**Rationale**: Minimizes diff. The rendering helpers (`renderHeader`, `renderTabBar`, etc.) don't need to change. Only the top-level `View()` and view-modal functions (`viewConfig`, `viewIndex`) change.

**Alternative considered**: Refactor every rendering function to return `tea.View`. Rejected because it would touch 15+ functions with minimal benefit.

### Decision 2: Key matching via msg.String() (unchanged pattern)

v2 introduces `tea.KeyPressMsg` with `msg.Code` and `msg.Text`, but `msg.String()` still returns the key name. We keep the existing pattern:

```go
case tea.KeyPressMsg:
    switch msg.String() {
    case "q", "ctrl+c": ...
    case "enter": ...
    }
```

**Rationale**: Zero change to the key dispatch logic. The only thing that changes is `tea.KeyMsg` → `tea.KeyPressMsg`.

### Decision 3: Mouse handling — separate type switches

v2 splits mouse messages. The current `handleMouse(tea.MouseMsg)` takes a single type and dispatches on `msg.Type` and `msg.Button`. We replace it with a type switch in `Update()`:

```go
case tea.MouseClickMsg:  m.handleMouseClick(msg)
case tea.MouseWheelMsg:  m.handleMouseWheel(msg)
```

The existing `handleMouse()` is refactored into two methods. Behavior is preserved.

### Decision 4: tea.Exec for editor launch

`tea.ExecProcess` is removed in v2. `tea.Exec` is the direct replacement:

```go
// v1:
return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
    return editorReturnMsg{}
})

// v2:
return m, tea.Exec(cmd, func(err error) tea.Msg {
    return editorReturnMsg{}
})
```

### Decision 5: viewport getter/setter migration

bubbles v2 changes `viewport.Model` fields to methods. All direct field accesses become method calls:

```go
// v1:
m.vp.Width = m.width - 2
m.vp.SetContent(content)
m.vp.YOffset

// v2:
m.vp.SetWidth(m.width - 2)
m.vp.SetContent(content)
m.vp.YOffset()
```

`viewport.New(width, height)` → `viewport.New(viewport.WithWidth(w), viewport.WithHeight(h))`

### Decision 6: lipgloss v2 — NoColor and Color API check

lipgloss v2 removes `AdaptiveColor` (not used by dossier) but keeps `lipgloss.Color("12")` and `lipgloss.NoColor{}`. The `Theme` struct in `styles.go` uses these — should be compatible without changes. Verify at upgrade time.

### Decision 7: glamour v2 — transparent upgrade

glamour v2 has the same API as v1 (only import path changes). `glamour.NewTermRenderer` with options should work unchanged.

## Risks / Trade-offs

- **tea.View.BackgroundColor might not accept lipgloss.Color**: The field is `color.Color` from `image/color`. Needs conversion from lipgloss terminal color. Risk: Medium. Mitigation: Verify during upgrade; lipgloss colors may implement `color.Color` already.
- **Cursed Renderer behaves differently**: v2 has a completely new renderer. Risk: Low (rendering differences should be improvements). Mitigation: Visual testing after upgrade.
- **Glamour may not be compatible with v2 styling**: Risk: Low. Glamour produces styled strings — it doesn't interact with tea.View directly. Mitigation: Test markdown rendering after upgrade.
- **Missing dependencies in go.sum**: After changing import paths, `go mod tidy` may pull different dependency trees. Risk: Low. Mitigation: Run `go mod tidy` and `go build` as first step.

## Open Questions

- Does `tea.View.BackgroundColor` accept `lipgloss.NoColor{}` or `lipgloss.Color` directly, or is a conversion needed?
- Does the Cursed Renderer handle our box-drawing characters (`┌─┐`, `├─┤`, etc.) correctly?
