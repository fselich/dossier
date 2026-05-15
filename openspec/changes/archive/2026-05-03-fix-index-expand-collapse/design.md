## Context

`handleTick()` in ModeIndex calls three `List*Names()` functions and compares their output against the in-memory lists to detect changes. The comparison uses `sameStrings()`, which is order-sensitive. The in-memory lists come from `LoadProjectSpecs()` (alphabetical sort) and `ListArchiveChanges()` (reverse sort), but `ListSpecNames()` and `ListArchiveNames()` return raw `os.ReadDir` order. On Linux, `ReadDir` returns filesystem order (inode/hash order for ext4), which does not match alphabetical or reverse-alphabetical. The mismatch causes `sameStrings` to return false on every tick, triggering a full reload that resets `expandedSpecs`.

## Goals / Non-Goals

**Goals:**
- `ListSpecNames()` and `ListArchiveNames()` return names in the same order as the in-memory lists they are compared against

**Non-Goals:**
- Changing the comparison from order-sensitive to set-based (unnecessary once sort orders match)

## Decisions

### Sort `ListSpecNames()` alphabetically and `ListArchiveNames()` descending

Add `sort.Strings(names)` to `ListSpecNames()` before returning, and `sort.Sort(sort.Reverse(sort.StringSlice(names)))` to `ListArchiveNames()` before returning. These mirror the existing sort calls in `LoadProjectSpecs()` and `ListArchiveChanges()` respectively.

`ListChangeNames()` is not affected: `Load()` iterates `ReadDir(changesDir)` without sorting, and `ListChangeNames()` also uses `ReadDir` without sorting, so both are consistently in filesystem order.
