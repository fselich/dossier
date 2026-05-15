## Context

OpenSpec organises work in `openspec/changes/<name>/` (active) and `openspec/changes/archive/` (completed). Each change has up to four artifacts: `proposal.md`, `design.md`, `tasks.md` and `specs/<capability>/spec.md`. Today there is no dedicated tool to browse and operate on them; the usual workflow is to open the files in a text editor.

The TUI runs globally from any project directory and is not part of the OpenSpec repository.

## Goals / Non-Goals

**Goals:**
- Discover and display active changes from `./openspec/` in CWD
- Render proposal, design and specs with readable formatting (markdown → ANSI)
- Expose the tasks from `tasks.md` as navigable and interactive items
- Write the `[ ]` ↔ `[x]` toggle directly to `tasks.md` without opening an editor

**Non-Goals:**
- Free-text editing (proposal, design, specs) — V2
- Archive navigation — V2
- View of canonical specs (`openspec/specs/`) — V2
- Real-time synchronisation of external file changes

## Decisions

### 1. Bubble Tea as TUI framework
Bubble Tea (Elm Architecture) is the de facto standard in the Go ecosystem for TUIs. Considered alternative: tview (immediate mode). Bubble Tea is chosen for the composability of its models and the charmbracelet ecosystem (lipgloss, glamour, bubbles) that covers all needs without ad-hoc code.

### 2. Glamour for markdown rendering
Glamour converts markdown to ANSI with word-wrap and theme support. Alternative: render raw markdown. Glamour is chosen because the readability of the proposal/design is one of the main values of the TUI.

### 3. Task toggle without a full markdown parser
Checkbox lines have a stable format (`- [ ] text` / `- [x] text`). The file is parsed line by line, a list of items with their line number is built, and the toggle only modifies that exact line in the file. A full markdown parser is not used because it is not needed and would add fragility.

### 4. Navigation: one main pane with artifact tabs
With the typical 1-3 active changes, a persistent side panel is not needed. The header shows `project · change [N/M]` and navigation between changes is done with `h`/`l`. The artifact tabs (`1`-`4`) occupy a fixed strip and the rest of the screen is content. Alternative: three-column layout. Discarded as unnecessarily complex given the typical volume of active changes.

### 5. Reading `config.yaml` from the project
The project name is inferred from the directory name (not from `config.yaml`), since the `context` field of config is optional and not always present.

## Risks / Trade-offs

- **Glamour + wide terminal** → On very wide terminals the render may produce long lines without breaks. Mitigation: fix glamour width to the real viewport width.
- **Concurrent writes to tasks.md** → If the user edits the file externally while the TUI is open, the in-memory state becomes out of sync. Mitigation: reload the file on focus return (V2); in V1 document the limitation.
- **Absent artifacts** → A change may not have all artifacts (e.g. no design.md). The TUI must disable the corresponding tabs and show a clear message.

## Open Questions

- Should the binary name be `spec-viewer` or something shorter like `sv`? Pending decision before publishing.
