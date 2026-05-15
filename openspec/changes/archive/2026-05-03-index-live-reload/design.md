## Context

`handleTick()` fires every 500 ms. Its first line is an early return for `ModeViewingArchive`, `ModeIndex`, and `ModeViewingSpec` — none of these modes do any polling. For `ModeIndex` this is the bug: the three lists that populate the index (active changes, archived changes, project specs) are loaded once in `enterIndex()` and never re-read.

The existing pattern for ModeNormal already does cheap detection: `ListChangeNames()` returns only directory names (no file reads) and compares them against the in-memory list. The same pattern can be applied to archives and specs.

## Goals / Non-Goals

**Goals:**
- Index refreshes automatically within one 500 ms tick when an active change, archived change, or spec appears or disappears on disk
- Detection is cheap: directory listing only, no file reads on every tick
- Cursor is preserved when the item it points to still exists after reload; reset to 0 only on out-of-bounds

**Non-Goals:**
- Refreshing archive or spec _content_ on tick (only list membership is tracked)
- Polling while in `ModeViewingArchive` or `ModeViewingSpec`

## Decisions

### 1. Two new loader functions: `ListArchiveNames()` and `ListSpecNames()`

Both return `[]string` of directory names — identical pattern to `ListChangeNames()`. Kept in `loader.go` alongside the existing function.

**Alternative**: inline the `os.ReadDir` calls in `handleTick`. Rejected — loader owns disk access; model should not reach into the filesystem directly.

### 2. Detection by name slice comparison, not count

Comparing full name slices catches renames and out-of-order additions, not just count changes. Uses the existing `sameNames` helper (already defined in `model.go`).

For archives and specs a parallel `sameStrings` helper is added (archive names are plain strings, not `openspec.Change`).

**Alternative**: compare lengths only. Rejected — would miss simultaneous add+remove (same count, different content).

### 3. On change detected: full reload + rebuild

When any of the three lists differs, all three are reloaded together (same as `enterIndex()` minus the mode/cursor reset), `buildIndexItems()` is called, cursor is clamped to `len(indexItems)-1`, and `refreshIndexViewport()` is called.

**Alternative**: partial reload of only the changed list. Rejected — `buildIndexItems()` interleaves all three lists; a partial reload would require tracking offsets. Not worth the complexity.

### 4. Cursor preservation

After rebuild, if `m.indexCursor >= len(m.indexItems)`, clamp to `max(0, len-1)`. No attempt to re-find the previously selected item by identity — the index lists are positional, not identity-based, and the cursor landing on a nearby item is acceptable UX.

## Risks / Trade-offs

- [Three `ReadDir` calls per tick in ModeIndex] → Each is a single syscall on a small directory. At 500 ms cadence this is negligible. Polling stops the moment the user leaves `ModeIndex`.
- [Cursor jump on reload] → If an item is removed above the cursor position, the cursor shifts. Acceptable; mirrors behaviour in ModeNormal when a change disappears.
