## Why

Dossier renders markdown with Glamour's hardcoded `dark` standard style. On light terminal themes, that can produce low-contrast text and make proposals, designs, and specs difficult to read.

## What Changes

- Read `DOSSIER_GLAMOUR_STYLE` when constructing the cached Glamour renderer.
- Keep the existing `dark` style as the default when the environment variable is unset.
- Add a focused unit test for default and override selection.

## Non-goals

- No UI setting or config file is added.
- No changes are made to Lipgloss chrome, task rendering, or git diff highlighting.
- No validation layer for custom Glamour style names is added.

## Capabilities

### Modified Capabilities

- `tui-viewer`: Markdown rendering style can be selected by environment variable.

## Impact

- **Affected code**: `internal/ui/viewport.go`, `internal/ui/viewport_test.go`
- **Dependencies**: none
- **User behavior**: existing behavior remains unchanged unless `DOSSIER_GLAMOUR_STYLE` is set before launching `dossier`
