## Context

The TUI stores the change artifacts in a `Project` structure loaded once at startup. There is no subsequent update mechanism. The main scenario motivating this change is an external agent (`opsx:apply`) that modifies `tasks.md` while the user has the TUI open.

Bubble Tea manages the event loop via messages. Periodic operations are implemented with `tea.Tick`, which emits a message every N seconds and re-registers itself from the handler, creating a continuous cycle.

## Goals / Non-Goals

**Goals:**
- Detect changes in the artifact files of the currently visible change
- Update the tasks view immediately and preserve the cursor
- Invalidate the render cache for proposal/design/specs when changes are detected
- No new dependencies

**Non-Goals:**
- Detection of new changes or archived changes
- Bidirectional synchronization (the user's toggle already writes to disk directly)
- Reload of the changes list

## Decisions

### 1. Polling with tea.Tick instead of fsnotify
`tea.Tick(2s)` fits naturally into the Bubble Tea event loop without additional goroutines or explicit channels. `fsnotify` would require a goroutine→channel→tea.Cmd bridge and adds an external dependency. At 2 seconds the latency is imperceptible for the main use case (agent marking tasks).

### 2. Comparison by content, not by mtime
Comparing the read content with the in-memory content detects real changes regardless of whether the mtime was touched (e.g. tools that rewrite without changing the content). The cost is reading the full file every 2s — for markdown files of a few KB this is negligible.

### 3. Cursor restoration by text, not by index
When `tasks.md` changes, indices can shift if tasks are added or removed. Saving the `Text` of the task under the cursor before the reload and searching for it in the new list keeps the position semantically correct. If the text disappears (task deleted), the cursor falls to the first available item.

### 4. Differentiated granularity by artifact type
- `tasks.md`: full reload (re-parse) because the interactive state depends on it
- `proposal/design/specs`: only invalidate `renderCache[tab]` — the glamour re-render occurs on demand when returning to that tab, without blocking the tick

## Risks / Trade-offs

- **Simultaneous toggle** → If the user presses Space at the exact moment a tick with external changes arrives, the toggle writes first and the reload overwrites with the state on disk (which includes the toggle). There is no data loss because the toggle already persisted before the reload.
- **File deleted between ticks** → If `tasks.md` disappears (change archived while the TUI is open), the reload simply will not find changes and the in-memory state remains stable. No crash.
