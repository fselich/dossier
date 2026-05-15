## Context

The index (`ModeIndex`) lists active, archived, and project specs. Archived changes can be opened with Enter, which takes you to `ModeViewingArchive` where their content is displayed in the viewport. Specs are currently purely decorative: `indexItems` only contains active and archived items, and the cursor never reaches the specs.

The codebase already has the full "open index item in read mode" pattern implemented for archived items: a new `Mode` enum (`ModeViewingArchive`), a separate cursor, and `loadViewport` that loads the correct content based on the mode. This change replicates that pattern for specs.

## Goals / Non-Goals

**Goals:**
- Specs navigable in the index (cursor with `j`/`k`, same as active and archived items)
- `Enter` on a spec opens `ModeViewingSpec` with the content of `spec.md` rendered as markdown via glamour
- `j`/`k` scrolls the viewport; `Esc` returns to the index
- `spec.md` is read once when entering `ModeViewingSpec` (no hot-reload)

**Non-Goals:**
- Editing the spec
- Navigation between multiple files of the spec (tab subnav)
- Hot-reload of content while viewing
- View of specs from a change (that is the existing TabSpecs)

## Decisions

### 1. New `indexItemKind`: `indexKindSpec`

`buildIndexItems` adds one `indexItem{kind: indexKindSpec, idx: i}` per spec in `m.projectSpecs`. This integrates specs into the existing navigation flow without special cursor code.

**Discarded alternative**: keep specs outside of `indexItems` and manage the cursor manually. More code, more error-prone.

### 2. New `Mode`: `ModeViewingSpec`

Like `ModeViewingArchive`, `ModeViewingSpec` uses the existing viewport and glamour to render markdown. A field `specViewerCursor int` is added that points to the currently viewed spec in `m.projectSpecs`.

**Discarded alternative**: reuse `ModeViewingArchive` with a flag. Semantically incorrect and breaks `renderHeader` and `renderHelpBar`.

### 3. Content loading in `LoadProjectSpecs`

`openspec.ProjectSpec` adds a `Content string` field with the raw content of `spec.md`. It is read in `LoadProjectSpecs` alongside the requirement count.

**Discarded alternative**: read the file in `model.go` when entering `ModeViewingSpec`. Mixes responsibilities; the loader already opens the file.

### 4. Async render same as the rest of markdown

`loadViewport` detects `m.mode == ModeViewingSpec`, takes `m.projectSpecs[m.specViewerCursor].Content` and launches the glamour render in the background with a `renderedMsg`. No per-tab cache is needed (specs don't have tabs); the same `renderCache` is used with a dummy key or directly `vp.SetContent` when receiving `renderedMsg`.

## Risks / Trade-offs

- [Large specs] → The glamour render is asynchronous, there is no blocking. The "Loading..." placeholder already exists.
- [Synchrony of `projectSpecs`] → If the user edits a spec while viewing it, it will not update (explicit non-goal). Acceptable for v1.
- [Index with many specs] → The cursor and the index viewport scroll already handle this correctly once `indexKindSpec` is in `indexItems`.
