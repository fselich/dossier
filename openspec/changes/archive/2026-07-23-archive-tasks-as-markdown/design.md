## Context

The tasks tab (`TabTasks`) currently has two rendering paths:

```
loadViewport() dispatch (viewport.go:23)
├── ModeNormal + TabTasks     → loadViewportForTasks()  → interactive (cursor, toggle, progress bars)
├── ModeViewingArchive + TabTasks → loadViewportForTasks()  → interactive (same path!)  ← BUG SOURCE
└── any other tab             → loadViewportForArtifact() → glamour markdown rendering
```

The interactive path calls `renderTasksContent()`, which builds a custom view with:
- Section headers with progress bars (`─── 5/8`)
- Task items with `▶` cursor marks and `[x]`/`[ ]` checkboxes
- Inline markdown styling for code and bold

This is appropriate for `ModeNormal` where toggling is the primary interaction. In `ModeViewingArchive`, it:
1. Introduces a visual inconsistency (tasks look different from proposal/design/specs)
2. Requires special guards in `viewer.go` to block `space` (toggle) and `e` (editor)
3. Requires `loadTaskItems()` to be called from multiple entry points (`index.go`, `mouse.go`)
4. Is the root cause of a regression where cursor movement blanks the viewport

The previous fix (9debbb4) patched symptoms rather than addressing the structural mismatch.

## Goals / Non-Goals

**Goals:**
- Tasks tab in archive mode renders via glamour (same as proposal/design/specs tabs)
- j/k scroll the viewport instead of moving a cursor
- Remove archive-specific guards and special cases
- Eliminate the regression at its root cause

**Non-Goals:**
- Do not change `ModeNormal` tasks behavior
- Do not remove `renderTasksContent()` — it remains used for active changes
- Do not change how archived changes appear in the index

## Decisions

### Decision 1: Change `loadViewport()` dispatch, not add another mode check

**Chosen:** Remove `ModeViewingArchive` from the tasks case in `loadViewport()`, letting it fall through to `loadViewportForArtifact()`.

```go
// Before (viewport.go:23):
case m.tab == TabTasks && (m.mode == ModeNormal || m.mode == ModeViewingArchive):
    return m.loadViewportForTasks()

// After:
case m.tab == TabTasks && m.mode == ModeNormal:
    return m.loadViewportForTasks()
```

**Why:** `loadViewportForArtifact()` already handles `TabTasks` (line 167-168: `raw = ch.Tasks.Content`). It renders via glamour with caching. This is a one-line change that fixes the root cause.

**Alternative considered:** Add a `m.mode == ModeViewingArchive` guard inside `loadViewportForTasks()` to call `loadViewportForArtifact()`. Rejected because it adds complexity to a function that shouldn't need archive awareness.

### Decision 2: Guard j/k TabTasks handlers with mode check instead of blocking at dispatch

**Chosen:** Add `m.mode == ModeNormal` guard to the j/down and k/up TabTasks cases in `updateViewer()`.

```go
// Before:
case "j", "down":
    case TabTasks:
        m.moveCursorDown()
        m.refreshTasksViewport()

// After:
case "j", "down":
    case TabTasks:
        if m.mode == ModeNormal {
            m.moveCursorDown()
            m.refreshTasksViewport()
        }
```

**Why:** With Decision 1, archive mode tasks use the glamour path. If j/k still triggered `moveCursorDown()` + `refreshTasksViewport()`, it would overwrite the glamour viewport with interactive content — a visual glitch. The guard makes archive mode fall through to the `default` case (`m.vp.ScrollDown(1)`).

**Alternative considered:** Remove TabTasks from the j/k switch and handle it separately only for ModeNormal. Rejected as more invasive.

### Decision 3: Remove dead archive guards, not keep them

**Chosen:** Remove the `if m.mode == ModeViewingArchive { return m, nil }` guards for `space` and `e` in `updateViewer()`.

**Why:** With Decision 2, when in archive mode + TabTasks, the j/k handlers fall through to the default scroll case. The `space` and `e` handlers are already after the j/k handlers, but since archive+tasks now uses the default path (not the TabTasks path), `space` with TabTasks in archive mode would enter ... let me trace: actually `space` handler checks `m.mode == ModeViewingArchive` FIRST (line 232), then checks `m.tab == TabTasks`. Wait, line 232's check is mode-level, not tab-level. In archive mode, space always returns nil regardless of tab. The e handler at line 250 also has `m.mode == ModeViewingArchive` check (tab-agnostic). These guards are still needed because e and space handlers check mode first, not tab. Keep these.

**Correction:** The archive guards for `space` and `e` are tab-agnostic (they return nil for ALL tabs in archive mode), so they remain correct. We only remove the unnecessary complexity from the tasks-specific path.

### Decision 4: Unify help bar for all archive tabs

**Chosen:** Remove the `if m.tab == TabTasks` special case in the archive help bar.

```go
// Before:
if m.mode == ModeViewingArchive {
    if m.tab == TabTasks {
        return "...j/k: navigate..."
    }
    return "...j/k: scroll..."
}

// After:
if m.mode == ModeViewingArchive {
    return "...j/k: scroll..."
}
```

**Why:** All tabs now scroll in archive mode. No need for a special case.

## Risks / Trade-offs

- **[Low] Loss of progress bar visualization** — Archived tasks rendered as markdown won't show progress bars. This is acceptable because archived data is frozen and read-only; progress information is stale.
- **[Low] Loss of cursor-marker visuals** — No `▶` cursor in archived tasks. Acceptable since there's no action to take on any item.
- **[None] Behavioral change** — j/k now scrolls instead of navigating in archived tasks. This matches every other archive tab and is more intuitive.
