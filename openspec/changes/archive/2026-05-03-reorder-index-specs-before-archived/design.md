## Context

The index (`ModeIndex`) builds a list of `indexItems` by concatenating three blocks in order: active, archived, and specs. The current order is `Active → Archived → Specifications`. The proposal swaps the position of the last two blocks so that specs appear between active changes and archived changes.

## Goals / Non-Goals

**Goals:**
- Change the build order of `indexItems` so that specs appear before archived changes.
- Update the `change-index` and `index-specs-section` specs to reflect the new order.

**Non-Goals:**
- Changes to navigation (j/k still works the same, only the physical order changes).
- Changes to the visual format of any section.
- Changes to the loading or live-reload logic.

## Decisions

**Decision: minimal change to concatenation order**

The construction of `indexItems` is a sequential concatenation. It is enough to move the specs block before the archived block. No new abstraction or additional refactoring is required.

## Risks / Trade-offs

- [Low risk] Existing tests that verify the order of items in the index will fail and need to be updated → the adjustment is trivial (change the expected index).
