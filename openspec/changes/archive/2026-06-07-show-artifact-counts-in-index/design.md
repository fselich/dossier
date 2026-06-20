## Context

The index view renders three sections: active changes, specifications, archived changes. Currently, section titles are plain labels. The proposal requests adding item counts to each title.

## Goals / Non-Goals

**Goals:**
- Show section item count in each section title
- Update `renderIndexContent()` in `internal/ui/index.go`

**Non-Goals:**
- No changes to data loading, filtering, layout, cursor interaction, or other UI elements
- Counts always reflect total items (not filtered subset)

## Decisions

- **Count source**: Use `len(m.project.Changes)` (active), `len(m.projectSpecs)` (specifications), `len(m.index.ArchiveChanges)` (archived) — the same variables already used for the empty-state guards
- **Format**: `"Section Name (N)"` using `fmt.Sprintf` — consistent across all three sections
- **Zero state**: When a section has 0 items, the empty-state message ("No active changes", etc.) is shown instead, so no count ever reads "(0)" — this matches existing behavior

## Risks / Trade-offs

No significant risks. The change is purely cosmetic and localized to `renderIndexContent()`.
