## Why

When opening a spec from the index (or switching changes/tabs/modes), the viewport briefly shows raw unstyled markdown before glamour finishes rendering. The flash is imperceptible with small documents but very noticeable with specs that have many requirements, where glamour takes ~100-200ms. This visual defect breaks the expectation that content is always shown with the active theme applied.

## What Changes

- Remove the `m.vp.SetContent(raw)` calls in `loadViewportForSpec`, `loadViewportForConfig`, and `loadViewportForArtifact` that push unstyled content into the viewport before the async render completes
- Clear the viewport (`m.vp.SetContent("")`) in those three methods while waiting for glamour's result
- The `m.loading` flag is already set correctly; final rendering already happens in the `specRenderedMsg`, `renderedMsg`, and `renderedConfigMsg` handlers

## Capabilities

### New Capabilities

(None — this is a bug fix with no new feature)

### Modified Capabilities

(None — specified behavior does not change; existing specs already require content to be shown "rendered as markdown")

## Impact

- `internal/ui/viewport.go`: remove 3 lines of `SetContent(raw)` and replace with `SetContent("")`
- No impact on APIs, dependencies, or observable behavior (only eliminates the visual flash)
