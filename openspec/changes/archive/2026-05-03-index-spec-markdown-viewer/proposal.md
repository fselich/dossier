## Why

The "Specs" section of the index shows project specs but does not allow interacting with them — they are purely decorative. Adding selection and markdown viewing allows consulting the content of a spec directly from the index, without leaving the tool.

## What Changes

- Specs in the "Specs" section of the index become navigable with `j`/`k`
- Pressing `Enter` on a spec opens a read-only viewing mode with the content of `spec.md` rendered as markdown
- The view uses the existing viewport with `j`/`k` scroll; `Esc` returns to the index
- No editing, no tabs, no subnav — read-only and scroll only

## Capabilities

### New Capabilities

- `spec-detail-viewer`: Read-only viewing mode for the content of a project spec (`openspec/specs/<name>/spec.md`) rendered as markdown, accessible from the index.

### Modified Capabilities

- `index-specs-section`: The "Specs not selectable in the index" requirement is reversed: specs become navigable items in `indexItems`, with cursor, selection, and `Enter` action.

## Impact

- `internal/ui/model.go`: new `indexItemKind` (`indexKindSpec`), new `Mode` field (`ModeViewingSpec`), updates to `buildIndexItems`, `renderIndexContent`, `Update`, and `renderHelpBar`
- `internal/openspec/loader.go`: reading the content of `spec.md` in `LoadProjectSpecs` (the function already exists, just add the `Content` field)
- `openspec/specs/index-specs-section/spec.md`: delta that revokes the non-selectability requirement
