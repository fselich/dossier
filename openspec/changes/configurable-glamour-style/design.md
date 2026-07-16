## Context

Dossier renders change artifacts in a terminal UI where Markdown prose comes from Glamour and the surrounding chrome comes from Lip Gloss styles. Light terminal palettes exposed two separate contrast problems:

- Glamour used the hardcoded `dark` style for Markdown artifacts.
- Dossier chrome used foreground colors that can become white or near-white on light backgrounds.

This change keeps absent-artifact behavior unchanged: if an artifact file is missing, its tab remains disabled and unselectable. This `design.md` file is present so the design tab is available for this OpenSpec change.

## Decisions

- Read `DOSSIER_GLAMOUR_STYLE` when constructing the cached Glamour renderer.
- Default to Glamour's `dark` standard style when the environment variable is unset or empty.
- Keep the renderer cache scoped to the active process and render width.
- Use light-palette-safe foregrounds for available inactive tabs, task cursor markers, and pending task rows.
- Keep the task cursor marker outside section and task row styles so its reset sequence cannot cancel the row's intended style.
- Keep checked and unchecked task checkboxes inside their task row style.
- Cap the tab-bar progress indicator at a compact width instead of consuming all remaining columns.

## Tradeoffs

- Invalid `DOSSIER_GLAMOUR_STYLE` values are passed to Glamour; renderer creation already handles errors by falling back to raw content.
- No in-app theme picker is added.
- Terminal palettes can still vary, but the selected ANSI colors avoid the specific white-on-white cases seen on light palettes.

## Validation

- OpenSpec strict validation.
- Focused unit coverage for Glamour style selection, tab progress width, and task focus style boundaries.
- Full Go race test suite.
- Manual pty checks with `DOSSIER_GLAMOUR_STYLE=light` and `NO_COLOR` unset.
