## Context

The TUI is implemented as a Bubble Tea model with a single `Model` struct. Navigation state is implicit: there is no `mode` field; logic is derived from `m.tab` and flags like `m.singlePath`. To add a modal and a read-only mode we need to make that state explicit.

Active changes are loaded in `loader.go:Load()`, which explicitly skips the `archive/` directory. Archived changes have the `YYYY-MM-DD-` prefix in their directory name.

## Goals / Non-Goals

**Goals:**
- Add a `mode` field to the Model that governs key behaviour and rendering.
- Implement the archived-change selection modal as an additional rendering layer.
- Load archived changes on demand (only when the picker is opened).
- Show the clean name + date in the modal, and `[archivo]` in the viewer header.
- Disable `e` and `Space` in `ViewingArchive` mode.

**Non-Goals:**
- Persisting which archived change was being viewed between sessions.
- Unarchiving (returning a change to active) from the TUI.
- Live-reload of archived changes (they are immutable).
- Pagination of the archived list (scroll with j/k is sufficient).

## Decisions

**Add an explicit `Mode` type to the Model**

Alternative: continue deriving state from boolean flags. Discarded because the picker+viewer combination would require 2–3 additional flags with fragile interactions. A `Mode` enum (`ModeNormal`, `ModeArchivePicker`, `ModeViewingArchive`) makes transitions explicit and key logic clean.

**Load archived changes on demand, not at startup**

Alternative: load all archived changes at startup alongside active changes. Discarded because the archive directory can grow indefinitely and would slow startup. They are loaded only when the user presses `a`, and cached in `m.archiveChanges` to avoid re-reading disk on each modal open.

**Render the modal as an overlay in `View()`**

The modal is rendered on top of the normal content using `lipgloss.Place` to centre it. The normal content continues to render beneath — visually the modal floats above it. Alternative: replace the entire view with the modal. Discarded because losing the visual context of the active change worsens user orientation.

**Clear `renderCache` when entering and leaving `ModeViewingArchive`**

The cache uses `Tab` as its key. An archived change's content and an active change's content would collide on the same keys. Clearing on context switch (active ↔ archived) avoids showing content from the previous context.

**Extract the clean archive name with `name[11:]`**

The `YYYY-MM-DD-` prefix always occupies 11 characters. The display date is parsed from the first 10 characters and formatted as `DD Mon`.

## Risks / Trade-offs

- [If the directory name format changes] The slice `name[11:]` would produce incorrect names → Mitigation: a helper function `parseArchiveName(dir string) (name, date string)` that validates the format before slicing; if it does not match, it returns the full name.
- [Very long list of archived changes] The modal has no internal scroll of its own — it uses the `j/k` cursor with a scrollable viewport. With 50+ archived changes this may be uncomfortable → acceptable for now, can be improved later.
