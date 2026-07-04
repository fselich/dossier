## Why

The collapsible index sections have two rough edges: pressing Enter on a section header silently navigates to a wrong view (or crashes), and the visual indicator is noisy — showing `▼`/`▶` before every section title when the expanded state should just be clean. These small issues make the feature feel unpolished.

## What Changes

- **Enter on section header** → no-op (currently falls through to archive navigation, which is a bug)
- **Collapse indicator** replaced: expanded sections show no indicator; collapsed sections show a trailing unicode ellipsis `…` in muted color (`helpStyle`)
- **Expanded sections** become visually clean — just the section name and count

## Capabilities

### New Capabilities

- (none)

### Modified Capabilities

- `collapsible-sections`: Modified "Collapse state is visually indicated" requirement (change indicator format); added requirement that Enter on section header does nothing

## Impact

- `internal/ui/index.go`: renderIndexContent indicator logic (lines 347-352), Enter handler guard (line ~756)
- No changes to model, styles, or tests
