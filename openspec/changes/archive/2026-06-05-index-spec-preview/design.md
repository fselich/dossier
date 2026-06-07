## Context

The index view (`ModeIndex`) currently shows three sections (active changes, specifications, archived changes) inside a scrollable viewport. Users browsing specs see only the spec name and requirement count — no purpose text. Checking what a spec is about requires pressing `Enter` to open the full spec viewer, then `Esc` to return to the index. This friction discourages exploration.

The existing chrome layout places a helpbar at the bottom, separated from the viewport by an inner separator:

```
header
───
viewport (scrollable)
───
helpbar
```

## Goals / Non-Goals

**Goals:**
- Add a fixed 1-line preview bar between the viewport and helpbar in `ModeIndex`
- Show the selected spec's name and purpose text (or nothing if cursor is on a non-spec item)
- Text truncation with `…` when exceeding available width

**Non-Goals:**
- Adding the preview bar to other modes (normal, archive, spec viewer)
- Interactive elements in the preview bar (display only)
- Markdown rendering in the preview bar (plain text only)

## Decisions

### Fixed chrome row vs inside viewport
- **Chosen**: Fixed chrome row between viewport and helpbar
- **Why**: Always visible regardless of scroll position; doesn't compete with scrollable content; follows the existing pattern of the helpbar
- **Alternative considered**: Inside viewport as a footer — would scroll away, defeating the purpose

### Always-reserved vs conditional height
- **Chosen**: Always reserve 1 line, even when no spec is selected
- **Why**: Avoids re-calculating viewport height on every cursor movement; simpler implementation; no jarring layout shifts
- **Trade-off**: Costs 1 terminal row when viewing non-spec items

### Text extraction from markdown
- **Chosen**: Parse `## Purpose` heading, extract plain text up to next `##` heading, strip markdown syntax
- **Why**: Self-contained in the `Content` string; no new parsing library needed; consistent with existing `ExtractRequirement()` approach
- **Alternative considered**: Adding a `Purpose` field to `ProjectSpec` — over-engineering for this use case

### Truncation behavior
- **Chosen**: Word-wrap at available width minus `len(name) + 3` (for ` ┊ ` separator), truncate to 1 line with `…`
- **Why**: Keeps the bar to exactly 1 line; `…` is the universal truncation signal

## Risks / Trade-offs

- **[Risk]** Purpose text may be empty or missing `## Purpose` heading → bar shows only the spec name
- **[Risk]** Long spec names + long purpose text = little room for the purpose → mitigated by truncation
- **[Risk]** ContentHeight must be reduced by 1 in `ModeIndex` → `WindowSizeMsg` handler already recalculates, so terminal resize works naturally
