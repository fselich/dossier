## 1. Cache glamour renderer

- [x] 1.1 Add `glamourRenderer *glamour.TermRenderer` and `lastRendererWidth int` fields to `Model` struct in `internal/ui/model.go`
- [x] 1.2 Initialize `glamourRenderer` once in `New()` and extract `minWidth = 80` constant + `clampWidth()` helper in `internal/ui/viewport.go`
- [x] 1.3 Replace all 4 `glamour.NewTermRenderer(...)` calls in `loadViewport()` with cached `m.glamourRenderer`, recreating only when width changes
