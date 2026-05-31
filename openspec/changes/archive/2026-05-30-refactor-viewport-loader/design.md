## Context

`loadViewport()` is a 155-line method in `viewport.go` that handles 5 distinct modes. Each mode has its own content source, rendering strategy (sync vs. async glamour), and edge cases. The function structure is:

```
if ModeIndex → sync refresh
if ModeViewingConfig → config markdown + glamour
if ModeViewingSpec → extract requirement + glamour (focus vs full)
if TabTasks + ModeNormal → sync task refresh
else → cache check + glamour (proposal/design/specs)
```

All callers (`Update`, `handleTick`, `editorReturnMsg`, mouse handlers) call `loadViewport()` uniformly.

## Goals / Non-Goals

**Goals:**
- Extract each mode branch into its own named method on `*Model`
- Reduce `loadViewport()` to a simple delegation switch
- Preserve exact behavior — no logic changes, no ordering changes

**Non-Goals:**
- Change the rendering logic for any mode
- Add new test coverage (separate change: `add-core-tests`)
- Introduce interfaces or abstractions over modes
- Change the glamour configuration

## Decisions

**Decision 1: Extract methods returning `tea.Cmd`**

Each extracted method returns `tea.Cmd` (same as `loadViewport()`), keeping the caller interface unchanged. The delegation is:

```go
func (m *Model) loadViewport() tea.Cmd {
    if !m.vpReady { return nil }
    switch m.mode {
    case ModeIndex:
        return m.loadViewportForIndex()
    case ModeViewingConfig:
        return m.loadViewportForConfig()
    case ModeViewingSpec:
        return m.loadViewportForSpec()
    default:
        if m.tab == TabTasks {
            return m.loadViewportForTasks()
        }
        return m.loadViewportForArtifact()
    }
}
```

**Decision 2: Method ordering by use-frequency**

Methods ordered as: Index → Config → Spec → Tasks → Artifact. This matches the original if-else ordering in `loadViewport()`.

**Decision 3: No new file**

All extracted methods stay in `viewport.go` to keep related viewport-loading logic together. No new package or file.

**Alternative considered**: Use a strategy map `map[Mode]func() tea.Cmd`. Rejected — the `TabTasks` check depends on `m.tab`, not just `m.mode`, so a simple mode map is insufficient.

## Risks / Trade-offs

- **[Risk]** Extraction could accidentally change call order or early-return semantics → **Mitigation**: copy-paste exact blocks, verify tests pass after each extraction
- **[Risk]** Functions become slightly longer than the original `loadViewport()` if measuring file-level → **Mitigation**: total file lines increase modestly (method signatures), but each function is single-responsibility and ~20-40 lines
