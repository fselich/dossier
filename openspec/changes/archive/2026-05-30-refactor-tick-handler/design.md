## Context

`handleTick()` in `internal/ui/index.go` runs every 500ms. The current logic has three distinct phases:

1. **Index mode** (lines 18-74): Check if active changes, archived changes, or specs changed on disk; if so, reload and rebuild the index. Otherwise, reload task counts for progress bars.
2. **Normal mode change detection** (lines 76-104): When not in single-path mode, detect if change list changed on disk; if so, reload project, preserve current change if it still exists.
3. **Normal mode content reload** (lines 106-154): Reload individual artifact content (tasks, proposal, design, specs) for the current change and update viewport if dirty.

These are three independent concerns sharing a single method body. Extracting them is a pure mechanical refactor.

## Goals / Non-Goals

**Goals:**
- Extract the three phases into named methods: `pollIndexMode()`, `pollNormalModeChanges()`, `pollNormalModeContent()`
- `handleTick()` retains only the guard (skip for archive/spec mode) and dispatches to extracted methods
- Exact same control flow, early returns, and nil checks preserved
- All existing tests continue to pass without modification

**Non-Goals:**
- Changing any logic, timing, or behavior
- Extracting the content reload sub-blocks (tasks vs proposal vs design vs specs) further
- Adding interfaces or making methods public (package-private is sufficient)

## Decisions

### Decision: Method boundaries follow existing early-return structure

The three phases are already separated by the early returns and the `if m.mode == ModeIndex` block. The extraction mirrors this:

```
pollIndexMode() → early return nil   (was line 74 return nil)
pollNormalModeChanges() → early return cmd (was line 101-102 / 103 return nil)
pollNormalModeContent() → return cmd  (was line 154 return nil)
```

`handleTick()` becomes:
```go
func (m *Model) handleTick() tea.Cmd {
    if m.mode == ModeViewingArchive || m.mode == ModeViewingSpec {
        return nil
    }
    if cmd := m.pollIndexMode(); cmd != nil || m.mode == ModeIndex {
        return cmd
    }
    if cmd := m.pollNormalModeChanges(); cmd != nil {
        return cmd
    }
    return m.pollNormalModeContent()
}
```

**Alternatives considered:**
- Function-level extraction (passing Model as arg) — rejected; methods on `*Model` are idiomatic and avoid copying the full struct.
- Single method for all normal mode polling — rejected; the change-list detection and content reload are different concerns with different early-return points.

### Decision: No new files

All extracted methods remain in `index.go` since they belong to the tick/timer concern that is already in this file. This avoids unnecessary file fragmentation.

## Risks / Trade-offs

- **Method ordering**: The three extracted methods reference `m.project`, `m.current()`, etc. No risk since they remain methods on `*Model`.
- **Diff size**: The extraction changes indentation across ~130 lines, making the diff larger than the conceptual change. Mitigation: review via side-by-side diff focusing on the dispatcher logic.
