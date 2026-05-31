## Why

`loadViewport()` constructs a new `glamour.NewTermRenderer` on every render call (4 sites in `viewport.go`), each with identical options. This adds unnecessary allocation overhead on every tab switch and resize. Creating the renderer once and reusing it eliminates the repeated construction.

## What Changes

- Add a `glamourRenderer *glamour.TermRenderer` field to `Model`.
- Initialize the renderer once in `New()`.
- Replace all 4 `glamour.NewTermRenderer(...)` calls in `loadViewport()` with `m.glamourRenderer`.
- Extract the `minWidth=80` constant and width clamping logic into a shared helper.

## Capabilities

<!-- No spec changes — implementation detail only. -->

## Impact

- `internal/ui/model.go`: New `glamourRenderer` field and initialization in `New()`.
- `internal/ui/viewport.go`: 4 `NewTermRenderer` calls replaced, min-width/clamp extracted.
- No behavioral changes; rendering output is identical.
