# Tasks

Step states: `[ ]` Pending, `[-]` Running/interrupted/failed, `[x]` Complete.

## Task 1 — Stable artifact identity across surfaces

> Spec: fix-artifact-row-consistency

- [x] Step 1 — Failing label tests: ID primary, stem dimmed secondary, suppressed when stem equals ID; tabs and filter agree
      start: 2026-10-22T14:15Z · elapsed: 3m · check: `go test ./internal/ui/ -run 'TestRenderArtifactRows'` — red for the right reasons (stem-as-label at column 4)
      paths: internal/ui/label_consistency_test.go
      commit: 75a7086 — fix-artifact-row-consistency: Task 1 Step 1 — label tests [x]
- [x] Step 2 — Refactor to `artifactLabels(primary, secondary)`; update rows, width, filter, tabs; retire stem-pinned test
      start: 2026-10-22T14:25Z · elapsed: 8m · check: `go test ./internal/ui/ ./internal/openspec/` — green (only gate-wrapper git test fails, exec 126, environmental)
      paths: internal/ui/index.go, internal/ui/index_refresh_test.go
      commit: 76e8a52 — fix-artifact-row-consistency: Task 1 Step 2 — artifactLabels refactor [x]

## Task 2 — Strays distinct and uncounted

> Spec: fix-artifact-row-consistency

- [x] Step 1 — Failing stray tests: `[stray]` badge, alpha after the Spec's steps, excluded count, ` · N strays` suffix
      start: 2026-10-22T14:40Z · elapsed: 4m · check: `go test ./internal/ui/ -run 'TestRenderStrayArtifactsDistinctAndUncounted'` — red for the right reasons (strays badge [authored] today)
      paths: internal/ui/label_consistency_test.go, internal/openspec/loader.go
      commit: 5f5f81d — fix-artifact-row-consistency: Task 2 Step 1 — stray tests [x]
- [x] Step 2 — Source-only planning count and stray rendering in `activePlanningText` / badge
      start: 2026-10-22T14:50Z · elapsed: 5m · check: `go test ./internal/ui/ ./internal/openspec/` — green (only environmental git-wrapper test fails, exec 126)
      paths: internal/ui/index.go, internal/openspec/cli.go
      commit: f296d65 — fix-artifact-row-consistency: Task 2 Step 2 — stray rendering [x]

## Task 3 — Honest degraded mode

> Spec: fix-artifact-row-consistency

- [x] Step 1 — Failing degraded tests: `[discovered]` badges, `files N · status unavailable`, middle-elided diagnostic
      start: 2026-10-22T15:05Z · elapsed: 4m · check: `go test ./internal/ui/ -run 'TestRenderDegraded|TestDiscoveryNever'` — red for the right reasons (planning 0/4 + tail-truncated diagnostic + discovery marking done)
      paths: internal/ui/label_consistency_test.go
      commit: 765cd43 — fix-artifact-row-consistency: Task 3 Step 1 — degraded tests [x]
- [x] Step 2 — Discovery stops hardcoding done; degraded header and diagnostic rendering
      start: 2026-10-22T15:15Z · elapsed: 8m · check: `go test ./internal/ui/ ./internal/openspec/` — green (only environmental git-wrapper test fails, exec 126); archived changes stamp done so History keeps authored rows
      paths: internal/openspec/loader.go, internal/ui/index.go, internal/ui/label_consistency_test.go
      commit: aa535e9 — fix-artifact-row-consistency: Task 3 Step 2 — honest degraded rendering [x]

## Task 4 — Verification

> Spec: fix-artifact-row-consistency

- [x] Step 1 — Full gate: gofmt, vet, build, lint, full suite with race
      start: 2026-10-22T15:30Z · elapsed: 10m · check: gofmt/vet/build/lint clean; `go test ./...` green except 7 git-spawning tests (internal/git ×6, TestGitSUnstageStaged) intercepted by the gate's git wrapper (exec 126, untouched paths) — receipts deferred to the human's ungated terminal; race clean of DATA RACE
      paths: internal/ui/hierarchy_test.go, internal/ui/label_consistency_test.go, internal/openspec/loader.go, internal/ui/index.go
      commit: 683cd32 — fix-artifact-row-consistency: Task 4 Step 1 — full verification [x]
- [x] Human validation — run dossier against `~/galactica` and `~/projects/kernel-rust`, observe the refined Requirements' final claim
      response: "VALID" — literal reply in session · outcome: VALID · at: 2026-10-22T14:40Z
