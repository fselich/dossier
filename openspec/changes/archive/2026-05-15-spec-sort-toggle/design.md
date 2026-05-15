## Context

The index view (`ModeIndex`) loads project specs into `m.projectSpecs []ProjectSpec`, sorted alphabetically by `LoadProjectSpecs()`. The 500 ms tick compares `m.projectSpecs` names against `diskSpecs` from `ListSpecNames()` (also alphabetical) to detect filesystem changes.

If `m.projectSpecs` were re-sorted in place, the tick comparison would always see a mismatch and trigger a reload every 500 ms, destroying any toggled sort order.

## Goals / Non-Goals

**Goals:**
- Toggle sort order of the Specifications section between full-name and suffix ordering
- Preserve cursor position across toggles
- Keep the 500 ms live-reload detection intact

**Non-Goals:**
- Persisting sort preference
- Sorting active or archived changes
- More than two sort modes

## Decisions

### Permutation slice, not in-place sort

`m.projectSpecs` stays in alphabetical order at all times. A new field `specOrder []int` holds a permutation of indices into `m.projectSpecs`. Default: `[0, 1, 2, …]`. Suffix sort: same indices, re-sorted by `specSuffix(m.projectSpecs[i].Name)`.

`buildIndexItems()` and `renderIndexContent()` iterate `m.specOrder` instead of `m.projectSpecs` directly. Because `indexItem.idx` still refers to the original slot in `m.projectSpecs`, all downstream logic (expand/collapse, requirement navigation, spec viewer, tick reload) is unaffected.

Alternative considered: re-sort `m.projectSpecs` in place. Rejected because it breaks the tick comparison.

### buildSpecOrder called inside buildIndexItems

`buildSpecOrder()` is called at the top of `buildIndexItems()` rather than at every call site. This keeps `m.specOrder` always in sync with `m.projectSpecs` and `m.specSortBySuffix` without requiring callers to remember the right sequence.

### Suffix extraction

```
specSuffix("archive-viewer")       → "viewer"
specSuffix("index-specs-section")  → "section"
specSuffix("path-arg")             → "arg"
specSuffix("tui-viewer")           → "viewer"  (ties broken by full name, stable sort)
```

`strings.LastIndex(name, "-")` — O(n), trivial.

### Help bar as sort indicator

The help bar shows the action that `s` will perform next:
- When sorted by name:   `… s: sort by suffix …`
- When sorted by suffix: `… s: sort by name …`

This gives the user both the key binding and the current state without adding noise to the section header.

## Risks / Trade-offs

- **Ties in suffix sort**: multiple specs share the same suffix (e.g., `archive-viewer`, `spec-detail-viewer`, `tui-viewer` all end in `viewer`). `sort.SliceStable` preserves their relative alphabetical order within ties — predictable, no special handling needed.
- **Cursor restoration**: after toggling, the cursor is restored by searching `m.indexItems` for an item matching the saved `kind`/`idx`/`reqIdx`. This is the same pattern already used by Space (expand/collapse). If the cursor was on an active or archived item (which don't reorder), it stays in place naturally.
- **expandedSpecs map**: keyed by original `m.projectSpecs` index — unaffected by the permutation.
