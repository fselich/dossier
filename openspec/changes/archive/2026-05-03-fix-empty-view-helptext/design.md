## Context

`emptyView()` in `internal/ui/model.go` is rendered when `ModeNormal` has no active changes. It currently shows a hardcoded `helpStyle.Render("\n  q: salir")`. All other help strings in the app were translated in the `ui-english-strings` change; this one was missed.

The `a` key already works in `ModeNormal` to enter `ModeIndex` (same as the handler in `Update`), so adding `a: index` to the helptext is accurate without any code change beyond the string.

## Goals / Non-Goals

**Goals:**
- Help text in `emptyView()` is in English and consistent with the rest of the UI
- User knows they can open the index from the empty state

**Non-Goals:**
- Any structural change to `emptyView()` layout
