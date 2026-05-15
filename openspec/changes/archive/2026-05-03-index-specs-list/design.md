## Context

The TUI has a `ModeIndex` view that currently shows two sections: "Active" and "Archived". Project specs live in `openspec/specs/<name>/spec.md` and contain requirements formatted as `### Requirement: <name>`. These have no representation in the TUI.

The loader (`internal/openspec/loader.go`) already knows how to load specs from individual changes (`loadSpecs`) but has no function that reads `openspec/specs/` at the project level.

The UI model (`internal/ui/model.go`) builds `indexItems` only from active and archived items; the `renderIndexContent` render iterates that list.

## Goals / Non-Goals

**Goals:**
- Show a "Specs" section in `ModeIndex`, below "Archived", with each project spec and its number of requirements.
- Specs are loaded once when entering `ModeIndex` (same as archived items).
- No navigation or selection required.
- Both sections (Archived and Specs) use a two-column aligned layout: name on the left, secondary data on the right.

**Non-Goals:**
- Navigation, selection, or opening of specs from the index.
- Editing specs from the TUI.
- Polling for changes in `openspec/specs/` during the session.
- Showing the full content of requirements (only their names).

## Decisions

### 1. New `ProjectSpec` type in the loader

A `ProjectSpec` type is added with fields `Name string` and `RequirementCount int`, and a `LoadProjectSpecs() []ProjectSpec` function that reads `openspec/specs/`, counts `### Requirement:` lines in each `spec.md`, and returns the list sorted by name.

**Discarded alternative**: store the full list of requirement names. The display only needs the count; storing the strings would carry data that no one consumes.

### 2. `projectSpecs []openspec.ProjectSpec` field in `Model`

Added to the model. Loaded in `enterIndex()`, same as `archiveChanges`, to avoid paying the disk cost on each render.

**Discarded alternative**: load in `Init()`. Would be earlier than necessary and complicates startup without active changes.

### 3. Static section in `renderIndexContent`

The "Specs" section is added at the end of `renderIndexContent` as plain text, outside of `indexItems`. No new `indexKind` is created because specs are not navigable in this phase.

**Discarded alternative**: add `indexKindSpec` but mark them as non-selectable. Adds unnecessary state complexity when the behavior is purely for display.

### 4. Two-column layout in render

In both the "Archived" and "Specs" sections, the render computes the maximum width of names in the list before iterating, and applies padding to align the right column. No fixed width is used to avoid wasting space with short lists.

### 5. Requirement count in the loader

`strings.HasPrefix` line by line over `spec.md` to detect `### Requirement: `. No full Markdown parsing; sufficient given the canonical format of the specs.

## Risks / Trade-offs

- [Specs with non-standard format] → The line parser is tolerant: lines that do not begin with `### Requirement:` are ignored. Worst case: an empty spec in the list.
- [Load in `enterIndex` without polling] → If someone creates a new spec during the session, it will not appear until the user leaves and re-enters the index. Acceptable given the explicit Non-Goal.
- [No dedicated scroll for the Specs section] → The index viewport is already scrollable, so long sections are accessible with `j/k`.
