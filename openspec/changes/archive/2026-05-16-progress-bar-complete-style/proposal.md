## Why

Progress bars currently use the same green color whether 10% or 100% complete, giving no visual signal that a change is fully done. A distinct color when complete helps users distinguish at a glance between in-progress and fully completed work.

## What Changes

- Add a new `progressCompleteStyle` (cyan, terminal color `"14"`) to the style registry.
- All three progress bar render sites switch to `progressCompleteStyle` when `done == total`.
- No change to behavior when the bar is partially filled.

## Non-goals

- Changing the progress bar characters (`█`, `░`, `─`).
- Animating or otherwise changing the bar on completion.
- Differentiating partially-filled bars by percentage (e.g., warning at <30%).

## Capabilities

### New Capabilities

- `progress-bar-complete-style`: Visual differentiation of progress bars that have reached 100% completion, using cyan color instead of green.

### Modified Capabilities

<!-- none -->

## Impact

- `internal/ui/styles.go` — new style variable added.
- `internal/ui/view.go:127` — general progress bar in tab bar.
- `internal/ui/index.go:343` — per-change progress bar in the index.
- `internal/ui/tasks.go:188` — per-section progress bar in the tasks view.
