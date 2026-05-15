## Context

The index cursor model is a flat `[]indexItem` slice built by `buildIndexItems()`. Each item has a `kind` (active/archived/spec) and an `idx` pointing into the corresponding data slice. `j`/`k` move an integer cursor over this slice, and `Enter` dispatches on the item kind.

`ProjectSpec` currently has `Name`, `RequirementCount`, and `Content`. Requirement names are not stored — only the count is derived from scanning `spec.md` for `### Requirement: ` lines.

The spec viewer (`ModeViewingSpec`) receives no scroll hint on entry; it always starts at the top.

## Goals / Non-Goals

**Goals:**
- Pressing `Space` on a spec item toggles its requirements inline as navigable sub-items
- `j`/`k` flow naturally through requirements without special casing
- `Enter` on a requirement opens `ModeViewingSpec` scrolled to that requirement
- Expand state survives cursor movement but resets when leaving and re-entering `ModeIndex`

**Non-Goals:**
- Persisting expand state across sessions
- Expanding archived change items (only specs get requirements)
- Pixel-perfect scroll positioning (best-effort line match is acceptable)

## Decisions

### 1. Flatten requirements into `indexItems` as `indexKindRequirement`

`indexItem` gains a `reqIdx int` field. When a spec is expanded, `buildIndexItems()` inserts one `indexKindRequirement` item per requirement immediately after the spec item. Collapsing removes them by rebuilding the slice.

**Alternative**: Two-level cursor (spec cursor + sub-cursor). Rejected — it would require special-casing `j`/`k`, clamp logic, and viewport tracking. The flat model reuses all existing cursor infrastructure unchanged.

### 2. Expand state: `expandedSpecs map[int]bool` keyed by spec index

`Model` gains `expandedSpecs map[int]bool`. `Space` on an `indexKindSpec` item toggles `expandedSpecs[item.idx]`, then calls `buildIndexItems()` and `refreshIndexViewport()`. The cursor is clamped after rebuild to handle the case where collapsing removes the item currently under the cursor.

`enterIndex()` initialises `expandedSpecs` to an empty map, so expand state resets when re-entering the index.

**Alternative**: Keyed by spec name (string). Rejected — spec index is sufficient and avoids a string lookup on every rebuild.

### 3. `RequirementNames []string` added to `ProjectSpec`

`LoadProjectSpecs()` already iterates `spec.md` line by line counting `### Requirement: ` prefixes. Extracting the name (the substring after the prefix, trimmed) in the same loop adds zero extra I/O.

### 4. Scroll-to-requirement via line search in rendered output

`specRenderedMsg` gains an optional `jumpLine int` (0 = no jump). When opening `ModeViewingSpec` from a `indexKindRequirement` item, the render goroutine receives the requirement name and, after glamour renders, scans the rendered string line by line for a line that contains the requirement name (stripped of ANSI escape sequences). The first matching line index is returned as `jumpLine`. On receiving `specRenderedMsg`, if `jumpLine > 0`, `viewport.SetYOffset(jumpLine)` is called.

**Alternative**: Count lines in raw markdown as a proxy. Rejected — glamour adds blank lines around headings and may reflow text, making the raw offset unreliable. Searching the rendered output is precise.

**Alternative**: Render per-requirement sections separately. Rejected — breaks glamour's document-level styling (heading hierarchy, consistent spacing).

ANSI stripping uses a simple regex `\x1b\[[0-9;]*m` applied per line before the name search.

### 5. Cursor preservation on expand/collapse

After `buildIndexItems()`, the new cursor position is calculated by finding the item that was under the cursor before the rebuild. For `indexKindSpec`, the item is found by matching `idx`. For `indexKindRequirement`, the item is found by matching both `idx` and `reqIdx`. If not found (e.g., collapsed the expanded spec while cursor was on a requirement), cursor falls back to the spec item for that `idx`. If out of bounds, clamp to `max(0, len-1)`.

## Risks / Trade-offs

- [ANSI search in rendered output] → glamour's heading rendering may include decorations (e.g., bold, color) interleaved with the text. Stripping ANSI codes before substring search should handle this reliably for the `### Requirement: <name>` pattern.
- [Cursor position after collapse] → If the cursor is on a requirement when the user collapses the spec, it snaps to the spec item. This is expected and matches standard tree-collapse UX.
- [buildIndexItems() called on every Space] → This is a linear operation over a small slice. No performance concern.
