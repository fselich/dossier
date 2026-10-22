# Plan

## Estimate

3–5 h · high confidence — small, well-located change; every call site already identified and evidence-gathered.

## Skills

Go · bubbletea/lipgloss TUI rendering · dossier's loader/CLI/index architecture (session-discovered call sites: `index.go` label/badge/fraction, `view.go` tabs, `loader.go` discovery, `cli.go` enrichment).

## Workflow

**ID-primary with dimmed stem secondary** — rows must stay stable across authoring and dossier renders real schemas that number templates but not IDs (both observed live). Stem-primary lost: row identity still depends on disk state — the core pain. ID-only lost: drops the file-name mapping users rely on whenever stems differ from IDs (the `spec` schema does, by design).

## Steps

1. Add failing label tests: rows keep the artifact ID primary with the stem dimmed secondary, suppressed when stem equals ID; tabs and filter agree (updates `index_refresh_test.go:469`).
2. Refactor `artifactDisplayLabel` → `artifactLabels(artifact) (primary, secondary)`; update `renderArtifactItem` (dimmed, padded secondary), `buildArtifactIDWidth`, filter (`index.go:250`), tabs (`view.go:103`).
3. Add failing stray tests: `[stray]` badge, alpha order after the Spec's steps, excluded from count, ` · N strays` suffix; implement via source-only `activePlanningText` count.
4. Add failing degraded tests: `[discovered]` badges, `files N · status unavailable` header, middle-elided diagnostic; implement change-level degraded state from the enrichment error.
5. Full verification: `gofmt`, `go vet ./...`, `go build ./...`, `golangci-lint run`, `go test ./...` (race on `internal/ui`).

## Verification

Deterministic: focused `go test ./internal/ui/ ./internal/openspec/` green after each step; full suite, vet, build, lint green at the end. Human validation (the one human-only act): run `dossier` against `~/galactica` (numbered templates) and `~/projects/kernel-rust` (schema absent), observe the Active Work panel per the refined Requirements' final claim, reply VALID.

## History

Prices as of 2026-10-22 (Research evidence: two live projects, seven code facts). Branch `fresh-openspec`; schema `spec` v2 (user-sourced).
