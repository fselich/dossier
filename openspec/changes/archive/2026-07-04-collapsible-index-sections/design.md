## Context

The index view uses a flat `[]indexItem` list where each item is a concrete artifact (change, spec, requirement, archived change). Sections are rendered as plain styled headers in `renderIndexContent()` but are not themselves items — the cursor skips over them. Currently only specs have a toggle (expand/collapse requirements via `Space`). There is no `indexItemKind` for sections.

## Goals / Non-Goals

**Goals:**
- Make the three index sections (Active Changes, Specifications, Archived Changes) navigable by the cursor
- `Space` on a section header toggles collapse/expand; collapsed sections hide their child items
- `Space` on a spec still toggles requirement expansion (unchanged)
- Visual indicator on section headers shows state: `▼` when expanded, `▶` when collapsed
- Collapse state survives index rebuilds (tick polling, filter changes)
- Filtering respects collapsed sections (filtered items inside a collapsed section are still hidden)
- Minimum lines changed; reuse existing patterns (`ExpandedSpecs` → `CollapsedSections` map)

**Non-Goals:**
- Not changing the rendering style or layout of section headers beyond the collapse indicator
- Not changing cursor behavior beyond enabling section navigation
- Not adding animation or progressive disclosure effects

## Decisions

1. **New `indexKindSection` kind** — A fourth `indexItemKind` value (`indexKindSection`) lets sections participate in the flat item list alongside active changes, specs, and archived changes. This is the simplest change that makes sections navigable.

2. **`CollapsedSections map[indexKindSection]bool`** — A new field in `indexState`, akin to `ExpandedSpecs`, keyed by section identity (we can use `int` constants 0, 1, 2 for the three sections). Default is `false` (all sections expanded).

3. **`buildIndexItems` inserts section items** — The first item added is the "Active Changes" section, followed by its children (if not collapsed), then the "Specifications" section with its children, then "Archived Changes" with its children. This mirrors the current rendering order but now sections are real items.

4. **Cursor on sections** — `visibleItemIdx` and `visibleItemCount` already work on the flat list; no changes needed. `renderIndexContent` renders section items with a `▶`/`▼` prefix and the section name. The cursor `▶` marker aligns with the section header line.

5. **Space context-sensitivity in `updateIndex`** — The current `Space` handler checks `item.kind == indexKindSpec`. We add `item.kind == indexKindSection` before it. `Space` on a section toggles `CollapsedSections[idx]` and rebuilds. `Space` on a spec still toggles `ExpandedSpecs[idx]` as before. The order of checks doesn't matter since one item can't be both kinds.

6. **Filter + collapse interaction** — `applyFilter` already filters after `buildIndexItems`. If a section is collapsed, its children aren't in the item list at all, so filtering is a no-op for hidden items. The cursor resets naturally.

7. **Loading archive changes** — Currently `enterIndex` and `pollIndexMode` load `ArchiveChanges` separately. These don't change — section collapse is purely a presentation concern in the model.

8. **Help bar** — Updated from `Space: expand` to `Space: toggle section  Space: expand spec` (or similar concise phrasing).

## Risks / Trade-offs

1. **[Cursor position on rebuild]** When a section is toggled, the item list changes length and existing cursor positions may be invalid. → Mitigation: Clamp cursor to `max(0, visibleItemCount()-1)` after rebuild, same as existing pattern for filter/expand.

2. **[Confusion between Space for sections vs specs]** Having both actions on `Space` could be surprising. → Mitigation: This is context-sensitive and natural — the action is always "toggle the thing under the cursor". The help bar will clarify.

3. **[Section items change total visible count]** Specs and requirements inside a collapsed section become unreachable. → Mitigation: This is the intended behavior; the user deliberately collapsed the section.

4. **[Mouse click on section headers]** Mouse click handling (`indexItemAtContentLine`) needs to account for section items now taking a content line. → Mitigation: Verify `indexItemAtContentLine` handles `indexKindSection` — clicking a section header should do nothing special (or toggle it, but that adds complexity; initially: no-op).
