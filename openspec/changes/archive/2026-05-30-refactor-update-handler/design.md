## Context

`Update()` in `internal/ui/update.go` handles `tea.KeyPressMsg` with a 283-line `switch msg.String()` (lines 106-388). The switch interleaves mode guards with action blocks for multiple modes. For example, `j`/`k` handle `ModeIndex` and the default case, `h`/`l` handle `ModeViewingSpec` and `ModeNormal`, and `q`/`esc`/`enter` each have their own mode checks. This makes it hard to tell at a glance which keys are active in a given mode.

Additionally, many cases share the same 2-line boilerplate:

```go
m.vp.SetHeight(m.contentHeight())
return m, m.loadViewport()
```

This appears ~15 times and obscures the actual state change logic.

## Goals / Non-Goals

**Goals:**
- Extract four mode-specific key handler methods: `handleIndexModeKeys()`, `handleNormalModeKeys()`, `handleArchiveModeKeys()`, `handleSpecModeKeys()`
- Extract `commitStateChange()` helper for the `vp.SetHeight` + `loadViewport` boilerplate
- `Update()` delegates to the appropriate handler based on `m.mode` after the shared keys (`q`, `i`, `esc` which affect multiple modes)
- Zero behavior change: exact same key bindings, exact same return values

**Non-Goals:**
- Refactoring non-keypress message handlers (WindowSizeMsg, renderedMsg, tickMsg, etc.)
- Changing any key binding
- Splitting mode handlers into sub-methods

## Decisions

### Decision: Top-level mode dispatch in Update()

The Update() keypress switch becomes:

```go
case tea.KeyPressMsg:
    switch msg.String() {
    case "q", "ctrl+c":
        // shared handler (quits or exits config)
    case "i":
        // shared handler (opens config)
    case "a":
        // shared handler (enters index)
    case "esc":
        // shared handler (multi-mode back navigation)
    case "enter":
        // shared handler (index-only)
    default:
        switch m.mode {
        case ModeIndex:
            return m.handleIndexModeKeys(msg)
        case ModeNormal:
            return m.handleNormalModeKeys(msg)
        case ModeViewingArchive:
            return m.handleArchiveModeKeys(msg)
        case ModeViewingSpec:
            return m.handleSpecModeKeys(msg)
        }
    }
```

Keys that apply to multiple modes (`q`, `ctrl+c`, `i`, `a`, `esc`, `enter`) stay in Update() as shared concerns. Mode-specific keys move to their respective handlers.

### Decision: commitStateChange() helper

```go
func (m *Model) commitStateChange() (Model, tea.Cmd) {
    m.vp.SetHeight(m.contentHeight())
    return *m, m.loadViewport()
}
```

Returns `Model` (value) not `tea.Model` to avoid interface boxing at each call site; the caller can append `, tea.Cmd` if needed.

**Alternatives considered:**
- Returning `(tea.Model, tea.Cmd)` — rejected; forces `return m` (not `return &m`) which is inconsistent with method receivers.
- Not extracting the helper — rejected; ~15 call sites benefit from a single point of change.

### Decision: No new files

All extracted methods remain in `update.go`. They are pure refactoring of the existing `Update()` method.

## Risks / Trade-offs

- **Mode-specific methods returning Model/tea.Cmd**: The existing code returns `(tea.Model, tea.Cmd)` from `Update()`. The extracted handlers follow the same signature. No risk.
- **Diff size**: Will be large because the switch cases move. Mitigation: review via side-by-side diff or commit the extraction in logical chunks.
