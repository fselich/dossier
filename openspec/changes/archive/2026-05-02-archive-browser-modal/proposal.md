## Why

There is no way to consult already-archived changes from the TUI — once archived they are invisible. Being able to review the history (proposal, design, specs, tasks) without leaving the viewer adds valuable context when working on new related changes.

## What Changes

- New key `a` that opens a modal listing the archived changes.
- The modal shows the clean change name (without the date prefix) and the date in compact format as visual metadata.
- Selecting an archived change from the modal opens a read-only mode where its artifacts can be navigated with the same keys as active changes (1-4, j/k).
- Keys `e` and `Space` are disabled in archive mode (read-only).
- `Esc` from the archive viewer returns to the modal; `Esc` from the modal returns to normal state. `a` from the archive viewer is equivalent to `Esc` (returns to the modal).
- The header shows `[archivo]` instead of `[N/M]` when viewing an archived change.

## Capabilities

### New Capabilities

- `archive-picker`: Modal for selecting archived changes with keyboard navigation.
- `archive-viewer`: Read-only display mode for an archived change, integrated into the existing viewer.

### Modified Capabilities

_(none — active changes do not change behaviour)_

## Impact

- `internal/ui/model.go`: new `mode` field (or similar), key logic for `a` and `Esc` by mode, modal rendering and helpbar adaptation.
- `internal/openspec/loader.go`: new function to list and load changes from `openspec/changes/archive/`.
- No changes to the openspec CLI or to on-disk artifacts.
