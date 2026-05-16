## Context

The TUI renders progress bars in three locations: the global tab bar (`view.go`), the change index (`index.go`), and the per-section task view (`tasks.go`). All three share two style variables from `styles.go`: `progressDoneStyle` (green) for filled blocks and `progressEmptyStyle` (gray) for empty blocks. There is no visual distinction between a bar at 80% and one at 100%.

## Goals / Non-Goals

**Goals:**
- Introduce a `progressCompleteStyle` (cyan) applied when `done == total` across all three render sites.
- Keep the change minimal and self-contained inside the `ui` package.

**Non-Goals:**
- Changing progress bar characters or layout.
- Any animation or transition effect on completion.
- Theming or user-configurable colors.

## Decisions

**Single new style variable vs. modifying existing style**

Adding `progressCompleteStyle` as a separate variable in `styles.go` is preferred over making `progressDoneStyle` dynamic. Keeping styles as package-level `var` declarations is the established pattern in this file, and a separate variable makes the intent explicit at the call site without adding any indirection.

**Color assignment: cyan (`"6"`) for in-progress, green (`"2"`) for complete**

Cyan (`"6"`) is used for the in-progress state and green (`"2"`) for the complete state. Bright variants (`"14"`, `"10"`) were considered but the standard variants provide sufficient contrast without being visually dominant alongside other UI elements.

**Conditional at render time, not at model level**

The `done == total` check already exists at all three render sites (to clamp `filled` to `barSpace`). Reusing that same condition to select the style adds zero new state and no changes outside the `ui` package.

## Risks / Trade-offs

- **Terminal color support**: Colors `"6"` and `"2"` are standard ANSI colors available in virtually all terminals. This is the same constraint that all existing styles have — no regression.
- **Future theme support**: Hardcoded color values are consistent with the rest of the codebase. If theming is introduced later, all styles will need migration equally.

## Migration Plan

No migration needed. Pure additive UI change with no data model or file format impact. Rollback is reverting the three render-site edits and removing the style variable.
