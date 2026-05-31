## Context

Three progress bar implementations exist across the codebase with identical logic: fill N `█` blocks with active color, pad with `·` for remaining width, and switch color to green when done==total. Two `viewContent` functions in view.go are also duplicates rendering a single component's content.

## Goals / Non-Goals

**Goals:**
- One parameterized `renderProgressBar(done, total, width int) string` function
- All call sites use it directly (no wrappers)
- One `viewContent` function replacing the two identical ones

**Non-Goals:**
- Changing bar characters, colors, or visual appearance
- Extracting anything beyond the listed functions

## Decisions

**Decision: Put `renderProgressBar` in `tasks.go`**
- `tasks.go` already owns the `progressBar` function — the most natural home for the unified version
- Other files already import from the `ui` package and can call it directly as a package-level function

**Decision: Drop wrappers — call sites call `renderProgressBar` directly**
- `renderTabBar` becomes just `renderProgressBar(globalDone, globalTotal, width)` inline
- `renderActiveItem` becomes `renderProgressBar(done, total, width)` inline
- The old `progressBar` in tasks.go calls `renderProgressBar` (it has extra logic for styling the right portion)

**Decision: Merge `viewIndexContent`/`viewConfigContent` into `viewContent`**
- Both functions are identical: they render a title header + content body with a separator
- The merged function takes title, body, and optional footer parameters

## Risks / Trade-offs

- [Risk: unintentional visual change] → Mitigation: compare before/after renders; the logic is a pure move with no behavioral changes
