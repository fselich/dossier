## Context

The viewport rendering pipeline currently sets raw markdown content immediately via `m.vp.SetContent(raw)` before dispatching the async glamour render. This causes one frame of unstyled text to be visible before the themed output arrives via `specRenderedMsg`/`renderedMsg`/`renderedConfigMsg`.

The `m.loading` flag already exists and is set to `true` during this window, but neither the viewport nor the view checks it.

Three functions share this pattern:

- `loadViewportForSpec()` — spec viewer (most visible: large specs)
- `loadViewportForConfig()` — config viewer
- `loadViewportForArtifact()` — normal artifact tabs (proposal, design, specs)

All three follow the same sequence: `m.loading = true` → `m.vp.SetContent(raw)` → `return async glamour render`.

## Goals / Non-Goals

**Goals:**
- Eliminate the flash of unstyled markdown in all three rendering paths
- Preserve existing behavior for all message handlers, glamour caching, loading flag, editor return flow

**Non-Goals:**
- No loading indicator or spinner (would be its own flash; viewport stays blank for < 1 frame)
- No change to glamour rendering behavior, render cache, or error handling
- No change to the viewport scroll position logic in message handlers

## Decisions

### Replace `SetContent(raw)` with `SetContent("")`

The fix is a two-line diff per function:

```
- m.loading = true
- m.vp.SetContent(raw)
+ m.loading = true
+ m.vp.SetContent("")
```

Rationale:
- The message handlers (`update.go:32,40,52`) already call `m.vp.SetContent(msg.content)` when the async render completes — they are the sole correct place to set content.
- Clearing to `""` is invisible: the viewport shows nothing for the gap between the old and new content. In practice glamour completes within the same tick for small content, and within ~100ms for large specs — not long enough to warrant a "Loading..." placeholder.
- The `m.loading = true` flag is preserved as-is; it's already set and used by the message handlers to know not to re-trigger a render.

**Alternatives considered:**
- **Show "Loading..."**: rejected — for small content the flash of "Loading..." would be more visually jarring than a blank viewport. For large content, the blank viewport still looks better.
- **Keep old content until render completes**: rejected — switching from change A to change B would show A's content briefly, which is misleading.
- **Make glamour synchronous**: rejected — would block the TUI event loop during rendering.

## Risks / Trade-offs

- [Risk] Viewport is blank for ~100ms on large specs — acceptable trade-off vs showing raw markdown
- [Risk] If glamour fails to render (renderer is nil), the content gets set to raw md by the fallback goroutine — this still works because the raw content is delivered via the message handler, just without the flash
