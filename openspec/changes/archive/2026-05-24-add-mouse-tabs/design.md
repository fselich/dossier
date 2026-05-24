## Context

Dossier currently handles only keyboard input (`tea.KeyMsg`). Mouse events are not captured because `tea.WithMouseCellMotion()` is not enabled in `main.go`. Terminal wheel scrolling works today only because many terminals translate wheel events into `up`/`down` arrow key presses when mouse mode is off — a side effect, not an application feature.

Enabling mouse capture will let us add click targets (starting with the tab bar) but requires handling wheel events explicitly since the terminal will stop translating wheel → arrow keys.

## Goals / Non-Goals

**Goals:**
- Enable `tea.WithMouseCellMotion()` to capture mouse events
- Handle wheel events so scroll continues working (replace the terminal translation)
- Handle left-click on the tab bar row to switch tabs, mirroring the `1`/`2`/`3`/`4` keyboard shortcuts
- Disabled tabs (absent artifacts) are not clickable

**Non-Goals:**
- Clicking index items, archived items, or spec items
- Clicking spec subnav labels (when viewing specs tab)
- Right-click back navigation
- Drag-to-scroll, hover effects, or motion tracking
- `tea.WithMouseAllMotion()` (only needed for hover effects)

## Decisions

### 1. New file: `internal/ui/mouse.go`

Mouse handling gets its own file consistent with the existing pattern: `index.go`, `tasks.go`, `viewport.go` each handle a distinct concern. `mouse.go` will contain a single method:

```go
func (m Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd)
```

Called from the `tea.MouseMsg` case in `update.go`.

**Alternatives considered:**
- Inline in update.go: would make an already-large switch even larger.
- Mix into model.go: model.go should stay focused on data structures and initialization.

### 2. Wheel → viewport forwarding

```
MouseButtonWheelUp   → m.vp.LineUp(3)
MouseButtonWheelDown → m.vp.LineDown(3)
```

3 lines per wheel tick is the standard TUI convention (matches most terminal emulators' default). Applies to all modes since the viewport is always present.

In task mode: wheel scrolls the viewport content, not the task cursor. Keyboard (`j`/`k`) stays as the task cursor navigation mechanism. Wheel is a scroll device; `j`/`k` is a navigation device. This separation is intuitive.

In index mode: same pattern — wheel scrolls, `j`/`k` moves the index cursor.

### 3. Tab click coordinate math

The tab bar is always at Y=2 (only in `ModeNormal` and `ModeViewingArchive`, where the layout includes tabs). Each tab label has `Padding(0, 1)` → visual width = `len(label) + 2`. Labels are joined with one space between them. The `│` border occupies X=0.

```
X: 1         11 12      19 20 21     27 28 29    35
   │ proposal │ │ design │ │ specs │ │ tasks │ ...
   └─ 10 cols ┘ └─ 8 cols ┘ └─ 7 cols ┘ └─ 7 cols ┘
```

Algorithm:

```
x := 1  // first tab starts at X=1 (past the │ border)
for each tab T in [Proposal, Design, Specs, Tasks]:
    w := len(tabLabels[T]) + 2
    if clickX in [x, x+w-1] and m.tabAvailable(T):
        switch to tab T (same logic as pressing key "1"/"2"/"3"/"4")
        return
    x += w + 1   // label width + space separator
// click fell on a space, progress bar area, or disabled tab: no-op
```

All tab styles (`active`, `inactive`, `disabled`) use the same `Padding(0, 1)`. The `Bold(true)` on the active style does not change terminal cell width. So the coordinate math is deterministic regardless of which tab is currently active.

**Edge cases:**
- Click on disabled tab → no-op
- Click on currently active tab → reloads viewport (same as pressing the number key when already on that tab)
- Click in the progress bar area (X > 35) → no-op

### 4. Only left-clicks on press

Filter to `msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft` before doing any coordinate mapping. Release events and other buttons are ignored. This prevents double-firing actions on press+release and keeps the door open for right-click (future).

### 5. Mode gating

Tab click coordinate check runs only when the current mode has a tab bar:

```
if m.mode != ModeNormal && m.mode != ModeViewingArchive:
    skip tab click handling
```

In other modes (Index, ViewingSpec, ViewingConfig), the tab bar row does not exist in the layout and clicks at Y=2 land on other elements (inner separator or viewport).

## Risks / Trade-offs

- **Scroll regression on terminals without mouse support**: If a terminal doesn't support mouse escape sequences, enabling `WithMouseCellMotion` has no effect — wheel continues being translated to arrow keys by the terminal. Low risk, most modern terminals support mouse.

- **Text selection requires Shift**: When mouse capture is active, clicking and dragging without holding Shift is captured by the app instead of selecting text. This is standard TUI behavior and unavoidable with mouse-enabled TUIs.

- **Coordinate fragility**: The tab bar Y=2 and X ranges are hard-coded based on the current layout. If the layout changes (e.g., a new row is added above tabs), the Y offset must be updated. Mitigation: the layout is stable and documented in `tui-viewer` spec; any layout change would be a separate change with its own spec update.

- **Tab label width depends on padding**: If tab styles ever change their padding, the width math breaks. Mitigation: all three tab styles explicitly set `Padding(0, 1)` in `styles.go`; this is visible in a single file.
