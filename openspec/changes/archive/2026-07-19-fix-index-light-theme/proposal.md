## Why

In the index view, non-selected item names (active changes, archived changes, spec names) are rendered as plain unstyled text. They rely on the terminal's default foreground color, which may not contrast with the viewport background set by the active theme. In light theme (`ViewBg = #ffffff`), unstyled text can become invisible or illegible, leaving only selected items (which use `IndexActive` with a blue background) readable.

## What Changes

- Add a `BaseText` style to `ThemeStyles` that applies `PrimaryFg` (no bold, no background)
- Apply `BaseText` to non-selected item names in `renderActiveItem`, `renderArchivedItem`, and spec name rendering in `renderIndexContent`
- Non-selected item names now use an explicit foreground color from the active theme

## Capabilities

### Modified Capabilities
- `theme-system`: add `BaseText` style to `ThemeStyles`; use it in index view for non-selected item names

## Impact

- `internal/ui/themes.go`: one new field in `ThemeStyles`, one line in `BuildStyles`
- `internal/ui/index.go`: 3 lines changed (`renderActiveItem`, `renderArchivedItem`, spec rendering)
- `internal/ui/themes_test.go`: verify `BaseText` foreground equals `PrimaryFg`
- `internal/ui/view_test.go`: may need adjusted assertions if tests check index output text
