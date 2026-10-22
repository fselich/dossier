# Requirements — drafted: artifact rows show stable identity and honest status

## Acceptance Criteria

- Each status-source artifact row renders its schema artifact ID as its primary label; when exactly one output file is present, its filename stem renders as dimmed secondary text and never replaces the ID. *(retires the Question: what is a row's stable identity?)*
- A row's primary label is identical before and after its output file is authored or removed.
- Discovery never marks an artifact authored: presence-only artifacts render through the stray or discovered badges, and no `[authored]` badge appears without status-source backing. *(retires the Question: does file presence ever imply authored?)*
- When a change's status call fails, its metadata shows `files N · status unavailable` — never a `planning` fraction — its rows badge `[discovered]`, and the diagnostic renders without losing the failing command text. *(retires the Questions: honest degraded rendering, and what may be truncated?)*
- Unclaimed discovered artifacts (strays) badge `[stray]`, sort after the Spec's step rows, are excluded from the `planning` count, and surface as `· N strays` when present. *(retires the Question: where do unclaimed files belong?)*
- Index rows, artifact tabs, and the index filter derive the artifact name from one shared labeling rule. *(retires the Question: which surfaces display the name, and must they agree?)*

## Scope

- `internal/ui/index.go` — label pair helper, `renderArtifactItem`, `renderArtifactStatusBadge`, `activePlanningText`, `activeDiagnosticText`.
- `internal/ui/view.go` — artifact tab labels.
- `internal/openspec/loader.go` — `loadDiscoveredArtifacts` status assignment.
- `internal/openspec/cli.go` — `EnrichChange` stray and degraded semantics.
- Tests: hierarchy, index-refresh, loader, and CLI fixtures.

## Out of Scope

- Renaming schemas or templates in user projects — vintage drift there is the user's to cure.
- Migrating, validating, or deleting stray artifacts in target projects.
- Changing the openspec CLI or its status JSON contract.
- Task-checkbox progress (`tasks x/y`).

## Guardrails

- Rendering only: no change to status-source statuses or the gate's grant semantics.
- Classic schemas whose output stem equals the artifact ID (e.g. `proposal`, `tasks`) gain no redundant secondary text.
- Index column alignment (label, dependency, badge) stays stable across rows and changes.
- Behavior changes land with updated tests in the same change; full suite, vet, build, and lint stay green.

## Human Validation Process

Run `dossier` against a project whose schema numbers its templates but not its artifact IDs (e.g. `~/galactica`) and one whose schema is not installed (e.g. `~/projects/kernel-rust`), open the Active Work panel, and observe: rows carry schema IDs as primary labels with filename stems dimmed after them; the schema-less change shows `files N · status unavailable` with `[discovered]` badges and a fully readable diagnostic; unclaimed old-vintage files render `[stray]` after the Spec's steps with `· N strays` in the metadata and do not inflate the `planning` fraction. That observed screen is VALID.
