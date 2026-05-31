## Context

`extractRequirement` (in `ui/viewport.go`) extracts a single requirement block from spec markdown. `configToMarkdown` (in `ui/view.go`) converts a `ProjectConfig` to a markdown-formatted string for display. Both operate on openspec content structures, not UI concerns — they belong in the openspec package.

## Goals / Non-Goals

**Goals:**
- Move `extractRequirement` to `openspec/loader.go` as an exported function
- Move `configToMarkdown` to `openspec/loader.go` as an exported function
- Update UI package callers to use `openspec.ExtractRequirement` and `openspec.ConfigToMarkdown`

**Non-Goals:**
- Changing function signatures or behavior
- Moving any other UI functions

## Decisions

**Decision: Export both functions (uppercase)**
- They're called from the `ui` package — must be exported to be accessible across packages
- `ExtractRequirement(content string, name string) string`
- `ConfigToMarkdown(cfg ProjectConfig) string`

**Decision: Put both in `loader.go`**
- `loader.go` already handles openspec content processing (loading, parsing, listing)
- Both functions deal with openspec data structures — natural fit
- Avoids creating a new file for just two functions

## Risks / Trade-offs

- [Risk: import cycle] → Mitigation: the `openspec` package does not import `ui` — moving functions into openspec only creates a dependency from ui→openspec, which already exists
