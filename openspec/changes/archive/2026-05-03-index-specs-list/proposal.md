## Why

The current index only shows active and archived changes. Project specs (`openspec/specs/`) are not visible anywhere in the TUI, forcing the user to leave the application to check which specs exist and what their requirements are.

## What Changes

- Add a new "Specs" section at the end of the `ModeIndex` view, after the "Archived" section.
- Each spec is shown with its name and the list of requirements it contains.
- Specs are not selectable or navigable in this phase (read-only/display only).
- The spec loader (`openspec/specs/`) is integrated into startup and polling to keep the list up to date.

## Capabilities

### New Capabilities

- `index-specs-section`: Section in the index view (`ModeIndex`) that shows project specs with their names and requirements, in read-only mode.

### Modified Capabilities

- `change-index`: The new "Specs" section is added to the existing index. No behavior changes in the already-specified "Active" and "Archived" sections.

## Impact

- `internal/openspec/loader.go`: Needs to load specs from `openspec/specs/` (the project specs directory, distinct from `openspec/changes/`).
- `internal/ui/model.go`: Model state extended with the spec list; `renderIndexContent` and `buildIndexItems` updated to include the new section.
- `internal/ui/styles.go`: Possibly new styles for the spec name and requirements in the index.
- No changes to the public API or the format of existing artifacts.
