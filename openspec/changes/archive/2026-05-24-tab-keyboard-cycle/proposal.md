## Why

The user currently switches between artifact tabs with direct number keys (`1`-`4`) or mouse clicks. There is no way to cycle sequentially forward through tabs without knowing the number. A `Tab` key binding provides a complementary, relative navigation model that is natural for keyboard-driven explorers reviewing all artifacts of a change in order.

## What Changes

- `Tab` key cycles forward to the next available tab (skipping disabled tabs, wrapping from last to first)
- `Shift+Tab` cycles backward to the previous available tab
- Help bar updated to document the new shortcut

## Capabilities

### New Capabilities

None. This is a modification of existing tab navigation behavior.

### Modified Capabilities

- `tui-viewer`: The "Tabs de artifact" requirement gains `Tab`/`Shift+Tab` as additional navigation keys alongside `1`-`4`

## Impact

- `internal/ui/update.go`: Two new key cases (`"tab"` / `"shift+tab"`) calling helper to cycle tabs
- `internal/ui/model.go`: Tab cycle helper method
- `internal/ui/view.go`: Help bar string updated to include `Tab` shortcut

### Non-goals

- Focus-based UI model with tabindex
- Cycling through the specs subnav
- Including `ModeViewingConfig` (textarea) — Tab must still insert indentation there
