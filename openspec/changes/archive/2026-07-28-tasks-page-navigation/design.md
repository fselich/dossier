## Context

The Tasks tab renders a flat `[]openspec.TaskItem` (mixing `KindSection` and `KindTask`) as styled text in `renderTasksContent` (`internal/ui/tasks.go`). Item heights are not uniform: a `KindSection` consumes two extra blank/spacing lines around its header, and a `KindTask` line can wrap to multiple rows once styled with `Width(contentWidth)`. The cursor (`m.tasks.Cursor`) is an *item index*, not a line offset; `moveCursorDown`/`moveCursorUp` just increment/decrement it, and `refreshTasksViewport` re-derives the cursor's pixel line by re-rendering and recording `cursorLine` when the loop reaches `i == m.tasks.Cursor`.

Every other tab (Proposal, Design, Specs, Git diff) treats `pgdown`/`pgup` as `viewport.PageDown()`/`PageUp()` — jump by one screen's worth of lines. The Tasks tab currently ignores both keys entirely, in both `ModeNormal` (live, editable cursor) and `ModeViewingArchive` (read-only, plain-scroll) — see `viewer.go:197-236`, which no-ops the `TabTasks` case unconditionally for both keys.

## Goals / Non-Goals

**Goals:**
- `pgdown`/`pgup` in `ModeViewingArchive` + `TabTasks` scroll the raw viewport by a page (parity with `j`/`k`'s existing `vp.ScrollDown/Up(1)` branch there).
- `pgdown`/`pgup` in `ModeNormal` + `TabTasks` move the cursor by a screen's worth of *rendered lines*, not a fixed item count — so the jump feels the same size regardless of how many sections or wrapped lines are on screen.
- No duplication of the line-height-counting logic that already exists in `renderTasksContent`.

**Non-Goals:**
- Git tab's non-diff file list keeps its current `pgdown`/`pgup` no-op.
- No change to single-step `j`/`k`/`up`/`down` cursor behavior, section-header reachability, or toggle semantics.
- No change to the definition of a "page" (still `m.vp.Height()` lines, same as other tabs).

## Decisions

**Decision: extend `renderTasksContent`'s existing line-counting loop to report every item's start line, instead of adding a parallel "line lookup" function.**

`renderTasksContent(m *Model) (string, int)` currently tracks a local `line` counter and only captures it into `cursorLine` when `i == m.tasks.Cursor`. Change the signature to `(string, int, []int)`, appending `itemLines[i] = line` unconditionally inside the existing loop (before each item's height is added). `cursorLine` becomes `itemLines[m.tasks.Cursor]` — a straight simplification, not new logic.

Item heights (and thus every entry in `itemLines`) do not depend on which item is currently the cursor — the cursor only swaps `▶ ` for `  ` in the rendered prefix, and both are the same rendered width — so `itemLines` is valid for computing a page-jump target regardless of where the cursor currently sits.

Alternative considered: a fixed item-count jump (`Cursor = min(Cursor + vp.Height(), len-1)`), reusing `moveCursorDown` N times. Rejected — it does not correspond to a real page: dense screens (many sections, one extra line each) would overshoot by more than a screen, sparse screens (few wrapped long lines) would undershoot. Extending the render loop instead makes the jump length match `viewport.PageDown()` exactly, and costs one extra `[]int` allocation on a keypress-driven, non-hot path.

**Decision: persist `itemLines` on `taskState` (`ItemLines []int`), refreshed every `refreshTasksViewport()` call, rather than recomputing it inline inside the `pgdown`/`pgup` key handler.**

`refreshTasksViewport` already re-renders on every task mutation, tab switch, resize, and tick-triggered reload (via `loadViewport` → `loadViewportForTasks`), so `m.tasks.ItemLines` is guaranteed fresh by the time the next keypress is dispatched. The `pgdown`/`pgup` handlers in `viewer.go` become pure lookups against already-computed state, keeping `updateViewer`'s key-dispatch switch (a value-receiver method) free of rendering side effects.

**Decision: page-jump target selection — nearest item at or past `cursorLine ± vp.Height()`, clamped to `[0, len(Items)-1]`.**

`moveCursorPageDown(pageHeight int)`: scan `itemLines` for the last index `i` such that `itemLines[i] <= itemLines[Cursor] + pageHeight`; if none advance past the current cursor (e.g., last screen is shorter than a page), clamp to `len(Items)-1`. `moveCursorPageUp` is the mirror: last index `i` such that `itemLines[i] <= itemLines[Cursor] - pageHeight`, clamped to `0`. No special-casing for landing on a `KindSection` item — normal `up`/`down` navigation already allows resting on section headers (see `tasks-cursor-navigation` spec), and toggle (`space`) already rejects non-`KindTask` items, so page-jump reuses that existing invariant rather than adding a `lastTaskIdx()` skip.

## Risks / Trade-offs

- [Risk] `itemLines` becomes stale if a future code path mutates `m.tasks.Items` or resizes the viewport without calling `refreshTasksViewport()` first, causing a page-jump to land on the wrong item for one keypress. → Mitigation: `refreshTasksViewport` is already the sole place that renders task content; no new call site is introduced that skips it. Existing gotcha #8 in AGENTS.md flags `renderCache` invalidation the same way — same discipline applies here.
- [Risk] Extremely small viewports (`vp.Height()` near 0, e.g. terminal resized very short) make page-jump degenerate to near single-step. → Mitigation: acceptable; matches how `viewport.PageDown()` already degrades for other tabs at the same terminal size.
</content>
