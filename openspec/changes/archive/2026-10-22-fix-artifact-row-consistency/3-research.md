# Research — fix-artifact-row-consistency

## Findings

- `artifactDisplayLabel` (`internal/ui/index.go:1007`) returns the single present output's DisplayName (the filename stem) when an artifact has exactly one present output, else the artifact ID — row identity therefore depends on disk state. (source: code)
- Discovery (`loadDiscoveredArtifacts`, `internal/openspec/loader.go`) hardcodes `Status: ArtifactStatusDone` for every discovered artifact; the `[discovered]` badge branch (`index.go:1057`) is unreachable for discovery-derived rows. (source: code)
- `activePlanningText` (`internal/ui/index.go:866`) computes done/total over all artifacts regardless of source — strays count. (source: code)
- `EnrichChange` (`internal/openspec/cli.go`) lays status-source artifacts first, then appends unclaimed discovery artifacts — strays already land after the Spec's steps in alpha order. (source: code)
- A failed status call appends the error to `Change.Diagnostic` and leaves discovery-only artifacts in place; the row renderer truncates the diagnostic (`! openspec status --cha…`). (source: `internal/ui/enrichment.go`, `index.go`)
- The same label helper feeds index rows (`index.go:1030`), artifact tabs (`internal/ui/view.go:103`), and the index filter (`index.go:250`). (source: code)
- `internal/ui/index_refresh_test.go:469` pins the current stem-label behavior (`artifactDisplayLabel(tasks) == "4-tasks"`). (source: tests)
- Observed on real projects (2026-10-22): `~/galactica` `fix-startup-race` rendered `planning 6/9` with three old-vintage strays (`2-proposition.md`, `2-requirements-a.md`, `3-proposal.md`); `ramona-voice-extension` showed `tasks [blocked]` beside `4-tasks [authored]`. `~/projects/kernel-rust` `honest-debug` mixed stem `1-idea` with bare IDs; three `ramona-plan` changes degraded to all-`[authored]` `planning 4/4` behind a truncated diagnostic — schema not installed user-level. (source: live sessions)

## Options

- Label pair from one helper: primary = artifact ID, secondary = single present output's stem (suppressed when stem equals ID) — shared by rows, tabs, and filter.
- Status split by source: discovery never yields done; strays (unclaimed after enrichment) and degraded rows (status call failed) distinguished at the change level.
- Count split: `planning` fraction over status-source artifacts only, `· N strays` suffix when present, degraded header `files N · status unavailable`.
- Diagnostic rendering: elide from the middle, never the command text.

## Evidence & prices — as of 2026-10-22

- Focused scope: four source files (`index.go`, `view.go`, `loader.go`, `cli.go`) plus four test files; estimate 3–5 hours.
- Cost of status quo: any project mixing schema vintages shows inflated progress and identity-flipping rows (two real projects today).
- Width price of secondary text: ~8–14 columns in the label column; elided before the primary label in narrow panes.

## Alternatives

- ID-primary with dimmed stem secondary
  · Keeps schema identity canonical while preserving the file mapping on screen.
  · Adds a second label column; classic schemas must suppress it to avoid noise.
  · Pick when rows must stay stable across authoring and schemas number their templates.

- Stem-primary always (generalize today's rule)
  · No renderer change beyond making the existing rule unconditional.
  · Row identity still flips on file create/delete; unnumbered IDs vanish from screens.
  · Pick when screens must mirror disk filenames exactly and schemas never renumber.

- ID-only (drop stems from rows entirely)
  · Simplest rows; zero width cost; no identity flip.
  · Loses the file-name mapping users rely on when stems differ from IDs.
  · Pick when schemas guarantee stems match IDs or filenames are never user-relevant.
