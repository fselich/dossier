## Why

Expanding a spec in the index view collapses automatically after ~500ms. The live-reload tick compares in-memory spec and archive name lists against fresh `ReadDir` results, but the in-memory lists are sorted (`LoadProjectSpecs` sorts alphabetically, `ListArchiveChanges` sorts descending) while the `List*Names()` polling functions return names in raw filesystem order. The comparison always fails — triggering a full reload that resets `expandedSpecs` every tick.

## What Changes

- `ListSpecNames()` returns names sorted alphabetically, matching `LoadProjectSpecs()` sort order
- `ListArchiveNames()` returns names sorted descending, matching `ListArchiveChanges()` sort order

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

_(none — internal polling correctness fix, no user-visible spec behavior changes)_

## Impact

- `internal/openspec/loader.go`: two one-line sort additions in `ListSpecNames()` and `ListArchiveNames()`
