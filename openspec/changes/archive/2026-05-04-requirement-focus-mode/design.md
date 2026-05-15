## Context

`ModeViewingSpec` today has one rendering path: glamour renders the full spec markdown, then `loadViewport()` searches the rendered ANSI output for the requirement name and calls `SetYOffset()` to scroll there. Two fields drive this: `specViewerCursor int` (which spec) and `specJumpTarget string` (which requirement to scroll to).

When entering from an `indexKindRequirement` item, the target requirement is buried in the full spec. There is no visual distinction between the target requirement and its neighbours.

`ProjectSpec` already carries `RequirementNames []string` and `Content string` (raw markdown), so all the data needed for focused rendering is already loaded.

## Goals / Non-Goals

**Goals:**
- When entering `ModeViewingSpec` from a requirement item, render only that requirement's block.
- `h` / `l` in focus mode navigate to the previous / next requirement within the same spec.
- Header shows `Req N/M` counter in focus mode.
- HelpBar reflects focus-mode controls.
- Entering from a spec item (not a requirement item) continues to render the full spec unchanged.

**Non-Goals:**
- Toggle between focus mode and full spec while viewing (out of scope).
- Dimming or greying out non-target requirements in the full spec.
- Any changes to the data loader or on-disk format.

## Decisions

**Two new fields on `Model`: `specFocusMode bool` and `specReqCursor int`**

`specFocusMode` gates the focused rendering path in `loadViewport()`. `specReqCursor` is the index into `projectSpecs[specViewerCursor].RequirementNames` and drives `h/l` navigation and the `N/M` counter. Both are set when entering from an `indexKindRequirement` item and cleared when entering from a spec item.

_Alternative_: repurpose `specJumpTarget` — skip adding new fields and derive focus mode from `specJumpTarget != ""`. Rejected: `specJumpTarget` is also set for the scroll-to-line path and conflating the two makes the logic harder to reason about. Explicit fields are clearer.

**Extract requirement block from raw markdown before glamour**

A helper `extractRequirement(raw, name string) string` scans `Content` line by line:
1. Finds the line `### Requirement: <name>`.
2. Collects lines until the next `### Requirement:` prefix or EOF.
3. Returns the block as-is to be rendered by glamour.

This is pure string processing on the raw markdown — no ANSI manipulation, no post-processing of glamour output. The extracted block renders correctly because it is valid markdown on its own.

_Alternative_: post-process glamour's ANSI output to dim lines outside the target. Rejected: glamour output contains complex, interleaved ANSI sequences; injecting dim codes reliably is brittle and would break with style changes.

**`h` / `l` in `ModeViewingSpec` with `specFocusMode == true` navigate requirements**

`h` decrements `specReqCursor` (wraps at 0), `l` increments (wraps at len-1), updates `specJumpTarget` to the new name, and calls `loadViewport()`. The existing `h/l` handler for `ModeNormal` (change navigation) is guarded by `m.mode == ModeNormal`, so there is no conflict.

**`Esc` from focus mode returns to index** — existing behaviour, no change needed. The index cursor is already preserved.

## Risks / Trade-offs

- [Requirement name not found in content] → `extractRequirement` returns an empty string → `loadViewport()` falls back to `"(spec not available)"`. Acceptable edge case; requirement names come from the same `Content` string so mismatch is unlikely.
- [Specs with a single requirement] → `h/l` wraps to the same requirement. Harmless; counter shows `1/1`.
- [`specFocusMode` persisting across sessions] — not applicable; state is in-memory only.

## Migration Plan

Pure additive in-memory change. No data migration. No backwards-compatibility concerns.
