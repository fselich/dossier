## Context

The collapsible index sections feature was recently implemented. Two issues surfaced in review:

1. **Enter bug**: The Enter handler in `updateIndex` has no explicit case for `indexKindSection`. It falls through to the archive navigation block which uses `item.idx` (a section constant 0/1/2) as an `ArchiveCursor`, causing index-out-of-bounds or navigation to the wrong archive.

2. **Visual indicator noise**: Section headers show `▼` (expanded) or `▶` (collapsed) as a prefix before the title. The expanded state should be visually clean — the foldable nature is already discoverable via cursor navigation and Space key. The collapsed state needs a subtle cue at the end, not a prominent prefix.

## Goals / Non-Goals

**Goals:**
- Enter on any section header is a silent no-op (no crash, no navigation)
- Expanded section headers show no indicator — just name and count
- Collapsed section headers show `…` (unicode ellipsis) in `helpStyle` after the count, separated by a space
- Cursor `▶` on the left is unchanged — it marks position, not collapse state

**Non-Goals:**
- Not changing cursor behavior, Space toggle, or any other key
- Not changing the click handler (clicking a section already just moves the cursor, which is correct)
- Not modifying styles beyond the indicator

## Decisions

1. **Enter guard** — Add `if item.kind == indexKindSection { return m, nil }` before the archive fallthrough in `updateIndex` (around line 756). Simplest fix, no side effects. Also add same guard in `clickIndexItem` for consistency (it's already a no-op via fallthrough, but explicit is better).

2. **Indicator logic in `renderIndexContent`** — Replace the current ternary:
   ```
   indicator := "▼" / "▶"
   header := fmt.Sprintf("%s %s (%d)", indicator, sectionName, totalCount)
   ```
   With:
   ```
   header := fmt.Sprintf("%s (%d)", sectionName, totalCount)
   if isCollapsed {
       header += " " + helpStyle.Render("…")
   }
   ```
   When the cursor is on the section, the cursor mark `▶` still appears on the left as before — that's the cursor indicator, not the collapse indicator.

3. **Unicode ellipsis** — Single character `…` (U+2026, HORIZONTAL ELLIPSIS). Rendered in `helpStyle` (gray, color 8) — same as the archive date and other secondary text.

## Risks / Trade-offs

1. **[Narrow terminals]** The ellipsis adds 2 cells (space + character). Negligible risk since section headers are short.
2. **[Discoverability]** Without a visible indicator on expanded sections, users might not know sections are foldable. → Mitigation: The help bar shows `Space: toggle`, and the cursor can land on section headers — these cues remain. The ellipsis on collapsed sections also reveals the mechanic.
