## Context

`loadSpecs` concatenates all `spec.md` files from subdirectories of `specs/` into a single `Artifact.Content` separated by `---`. The TUI treats specs as a monolithic artifact: a single viewport, with no notion of how many specs there are or their individual names.

The change adds a navigation level inside the `specs` tab without touching the main tab structure or the polling mechanism.

## Goals / Non-Goals

**Goals:**
- Show a sub-bar with one chip per spec when the active tab is `specs`
- The `3` key cycles through specs if already on the `specs` tab; otherwise, switches to it
- The viewport shows one spec at a time, individually rendered by glamour
- `contentHeight` adjusts dynamically: −1 when the sub-nav is visible

**Non-Goals:**
- Changing the on-disk change detection mechanism (the concatenated `Artifact.Content` remains the polling hash)
- Support for nested specs (one subdirectory level, same as now)
- Persisting `specIdx` across change switches

## Decisions

**D1: `[]NamedSpec` in `Change` alongside the existing `Artifact`**

`SpecFiles []NamedSpec` is added to `Change`. `loadSpecs` populates it at the same time as the `Artifact`. The concatenated `Artifact.Content` remains intact so the polling tick handler can continue comparing hashes without changes. `NamedSpec` is a simple struct `{Name, Content string}`.

Discarded alternative: replacing `Artifact` with something richer — would break the tick handler and the polling model with no real benefit.

**D2: `specIdx int` in `Model`, resets to 0 when switching changes or tabs**

The index resets when the user switches changes (`h`/`l`) or when entering the specs tab from another tab. This avoids incoherent states (e.g., `specIdx = 2` for a change with only 1 spec).

**D3: `3` key with dual behaviour**

If `m.tab != TabSpecs`: switch to TabSpecs (current behaviour).
If `m.tab == TabSpecs`: `specIdx = (specIdx + 1) % len(specFiles)` and reload viewport.

If there is only one spec, the cycle does nothing visible (same spec). No special guard is needed.

**D4: Conditional `contentHeight`**

`contentHeight()` returns `m.height - 8` when `m.tab == TabSpecs && len(specFiles) > 0`, and `m.height - 7` in all other cases. The viewport is resized when switching tabs (already happens today in the `1`/`2`/`3`/`4` handlers).

**D5: Render cache by `(tab, specIdx)` — discarded, cache by tab is sufficient**

Glamour renders the selected spec when `specIdx` changes. The `TabSpecs` cache is invalidated when cycling (same as with content changes). The cache only stores the last rendered spec for this tab — sufficient for typical use (the user rarely goes back within specs).

## Risks / Trade-offs

[Risk: Cache miss when returning to a previous spec] → The user notices a brief load ("Loading…") when going back. Acceptable given the typical size of specs.

[Risk: `contentHeight` changes when entering/leaving specs] → The viewport already resizes on every tab change; this change is consistent with the existing pattern.

## Migration Plan

No migration. The change is additive — `SpecFiles` is populated just like `Artifact`, and `Artifact.Content` continues to exist. Older binaries simply ignore `SpecFiles`.

## Open Questions

None.
