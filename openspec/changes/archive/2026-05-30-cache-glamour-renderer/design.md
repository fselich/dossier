## Context

`loadViewport()` constructs a `glamour.TermRenderer` 4 times with identical options (`glamour.WithStandardStyle("dark")`, `glamour.WithWordWrap(width)`). The width is clamped to a minimum of 80 in each of those 4 sites. Creating the renderer involves parsing the style definition — cheap but unnecessary to repeat.

## Goals / Non-Goals

**Goals:**
- Initialize one `glamour.TermRenderer` per process and reuse it.
- Extract width clamping (`max(width, 80)`) into a helper.
- Keep rendering output identical.

**Non-Goals:**
- Making the renderer style or theme configurable.
- Adding renderer lifecycle (close/recreate).

## Decisions

**Add `glamourRenderer *glamour.TermRenderer` to `Model` and initialize in `New()`**
The renderer is constructed with fixed options at startup. At call sites, `glamour.NewTermRenderer(...)` is replaced with `m.glamourRenderer`. Word wrap width is dynamic per call, so the renderer must support runtime width changes. Glamour's `TermRenderer` is stateless for rendering and accepts width via options; since width is dynamic, we pass it via `WithWordWrap` at each `.Render()` call — wait, Glamour's API bakes width into the renderer at construction. This means we may need to recreate when width changes, or accept a fixed width.

**Actually: Glamour width is fixed at construction** — the `WithWordWrap` option sets the renderer's internal width. Recreating on every resize negates the benefit. Instead, set the renderer's width to the Viewport's width via a setter or accept the current width at render time. Check Glamour API: `TermRenderer` has `Render(input string) (string, error)` — no width parameter. Width is constructor-only.

**Revised approach**: Cache the `glamourRenderer` but recreate it when `width` changes. Track `lastRendererWidth` and only call `NewTermRenderer` when it differs. This preserves the caching benefit across tab switches (most common case) while still supporting resize.
_Alternative_: always recreate below a threshold — no, the "recreate on width change" approach is strictly better.

**Extract `minWidth=80` constant and `clampWidth(w int) int` helper**
The pattern `width := m.width - 2; if width < 20 { width = 80 }` appears 4 times. Extract to a helper:
```go
const minTermWidth = 80

func clampWidth(w int) int {
    if w < minTermWidth {
        return minTermWidth
    }
    return w
}
```
Note: the existing code clamps at 20 (not 80) — `clampWidth` uses 80 as the floor.

**Wait, re-reading**: the existing code uses `if width < 20 { width = 80 }` — this sets width to 80 when it's below 20. We'll clean this up to simply use `max(m.width-2, minWidth)`.

## Risks / Trade-offs

- **Stale renderer width after resize** → Mitigation: track `lastRendererWidth` and recreate when width changes.
- **Glamour renderer not goroutine-safe** → Mitigation: all render calls are dispatched via Bubble Tea's Cmd system (single goroutine), so no concurrent access.
