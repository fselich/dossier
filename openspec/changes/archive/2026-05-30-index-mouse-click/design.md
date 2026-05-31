## Context

The index view (`ModeIndex`) currently supports mouse wheel scrolling (moves cursor) and keyboard navigation (j/k, Enter, Space), but left-clicking on an item does nothing. The tab bar already supports click-to-switch via `handleMouseClick`. This design extends that same function to handle `ModeIndex`.

## Goals / Non-Goals

**Goals:**
- Left-click on an index item moves the cursor to that item
- Left-click on the already-selected item performs its primary action (Enter → open change/spec, Space → toggle spec expansion)
- Clicks outside the viewport content area or on section headers / blank lines do nothing
- Works correctly with viewport scroll offset (YOffset)

**Non-Goals:**
- Double-click support (single-click semantics only)
- Right-click or middle-click handling
- Click on section headers or blank lines
- Touch / gesture support

## Decisions

### Decision: Screen Y → content line → index item via line-counting

The viewport occupies screen rows 3 to `3 + vp.Height() - 1`. The viewport may be scrolled (`YOffset > 0`).

```
contentLine = msg.Y - 3 + m.vp.YOffset()
```

A new method `indexItemAtContentLine(line int) (idx int, found bool)` iterates `m.indexItems` while mirroring the line-counting logic from `renderIndexContent()`. This avoids storing redundant position data and keeps the single source of truth in the rendering code.

**Alternatives considered:**
- Pre-compute line positions in `indexItem` struct — rejected because it couples data to layout; needs recomputation on every change.
- Screen-relative index (row 3 = first item) — fails when viewport is scrolled.

### Decision: Click-to-select-then-click-to-open

- Click on uncursor item → move cursor there, no action
- Click on cursor item → perform action (Enter for active/archived/requirement, Space for spec toggle)

This matches standard TUI mouse conventions and feels safe (accidental clicks just move the cursor, only deliberate clicks on the already-selected item trigger actions).

## Risks / Trade-offs

- **Viewport content area calculation**: Depends on the hardcoded row offset (3 = boxTop + header + innerSep). If the index layout changes, this offset must be updated. Mitigation: add a named constant or helper function.
- **Section header line counting**: The `indexItemAtContentLine` function must stay in sync with `renderIndexContent()`. Mitigation: place them adjacent in the source file and add a comment referencing the coupling.
- **Spec expanded state**: When a spec is expanded, requirement items appear. The line count changes. `indexItemAtContentLine` reads `m.expandedSpecs` at call time, so the mapping is always current.
