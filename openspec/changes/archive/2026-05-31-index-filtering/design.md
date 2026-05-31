## Context

The index view (`ModeIndex`) renders a flat list of `[]indexItem` grouped into three sections: Active Changes, Specifications, and Archived Changes. Navigation is linear (j/k). Users who know what they're looking for must scroll through all items.

The filtering feature adds a vim-style "/" prompt at the bottom of the screen. As the user types, items are filtered in real-time via case-insensitive substring matching.

## Goals / Non-Goals

**Goals:**
- Press "/" in `ModeIndex` to open a filter prompt (replaces help bar while typing)
- Real-time case-insensitive substring filtering of active changes, specs, requirements, and archived changes
- `Enter` confirms the filter (stays applied), `Esc` clears it or cancels typing
- When no items match a section, show "No items match '<query>'" in that section
- Esc is contextual: filter active → clear filter; no filter → quit
- Cursor is preserved when filter doesn't remove the current item

**Non-Goals:**
- Fuzzy matching (simple substring only)
- Filtering by status (done/total counts), date, or section
- Persisting filter after leaving `ModeIndex`
- Regex or prefix matching modes
- Search history

## Decisions

### Decision 1: FilterIndices indirection layer over Items

Filtered items are tracked via `FilterIndices []int` — a slice of indices into the full `Items` slice. When filter is active, the cursor indexes into `FilterIndices` rather than `Items` directly.

```
Items:       [A, B, C, D, E]        (always the full list)
Filter:      "foo"
FilterIndices: [0, 2, 4]            (A, C, E match)
Cursor:      1                       → Items[FilterIndices[1]] = Items[2] = C
```

**Why**: Keeps `buildIndexItems()` unchanged (always builds the full list). Filter just adds a view layer. On rebuild (tick/Space), filter is re-applied automatically.

**Alternative rejected**: Rebuilding Items as filtered — would lose the full list, making filter clearing and cursor preservation harder.

### Decision 2: FilterText + FilterActive as two-axis state

| FilterText | FilterActive | Meaning |
|---|---|---|
| `""` | false | No filter (default) |
| `"foo"` | false | Filter applied, browsing |
| `"foo"` | true | Actively editing query |

Separating "is there a filter" from "is the user typing" avoids conflating the two. During typing, keypresses go to filter editing. After Enter, the filter stays active but keys return to navigation.

### Decision 3: Help bar doubles as filter prompt

While `FilterActive` is true, `renderHelpBar()` shows `/` + query + cursor. After confirming, the help bar shows normal bindings with a `[/query]` indicator.

**Why**: No extra chrome or layout changes. The help bar is already the bottom line inside the box frame.

### Decision 4: Contextual Esc priority

```
Esc → FilterActive? → Yes: cancel typing (go to Filtered or Browsing)
     → FilterText != ""? → Yes: clear filter (go to Browsing)
     → Quit
```

**Why**: Single key does the "least destructive" thing at each level. Users who just pressed "/" accidentally can Esc out. Users who want to remove the filter press Esc once (from Browsing). Users who want to quit press Esc twice, or once if no filter.

### Decision 5: No-match sections show inline message

Each section (Active, Specs, Archives) independently checks if any filtered items exist in it. If a section has zero matching items, it shows "No items match '<query>'" in place of that section's items. Sections with no items at all (unfiltered) keep their existing "No active changes" / "No specifications available" / "No archived changes" messages.

### Decision 6: Mouse click maps through FilterIndices

`indexItemAtContentLine` already walks items in render order. It needs to skip non-matching items (same as `renderIndexContent`). The returned index is the raw `Items` index, which `handleMouseClick` then maps back through `FilterIndices` if filtering is active.

## Risks / Trade-offs

- **Line counting coupling**: Both `renderIndexContent` and `indexItemAtContentLine` must apply the same filter logic. If one diverges, mouse clicks will map to wrong items. Mitigation: extract a shared `isItemVisible(idx int) bool` helper.

- **Cursor jump on filter change**: When the user types, items may appear/disappear, causing the cursor to jump. Mitigation: clamp cursor to `len(FilterIndices)-1` after each filter change. If the current item still matches, keep its position.

- **Spec expansion + filter**: Expanding a spec adds requirement items. `buildIndexItems` is called, then filter is re-applied. The cursor re-targets the same spec item in the new filtered list. Requirements matching the filter will also appear.
