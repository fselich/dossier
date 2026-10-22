# Requirements — refined: research-driven steering of the drafted sheet

## Acceptance Criteria

- Every artifact row from the status source renders its schema artifact ID as its primary label; when exactly one output file is present, its filename stem renders as dimmed secondary text after the ID, and is suppressed when the stem equals the ID.
- The primary label is byte-identical across all rendering surfaces — index rows, artifact tabs, and index filter matches — before and after its output file is authored or removed.
- No `[authored]` badge renders without status-source backing; discovery-derived presence renders only through the `[stray]` or `[discovered]` badges.
- A change whose status call failed renders exactly `files N · status unavailable` in its metadata — never a `planning` fraction — with rows badged `[discovered]` and a diagnostic whose failing-command text remains readable end to end (elided from the middle, never the command).
- Unclaimed discovered artifacts render `[stray]` after the Spec's steps in alpha order, excluded from the `planning` count, with ` · N strays` appended to the metadata when present.
- **Final claim (Tasks judges this):** running `dossier` against a numbered-template schema project and a schema-less project shows stable ID-primary rows, an honest degraded header, and uncounted strays — the observed screen is VALID.

## Deltas from Drafted

- **Sharpened** — "one shared labeling rule" becomes byte-identical primary labels across rows, tabs, and filter (sharpens the identity and surfaces Questions).
- **Sharpened** — degraded header pinned to the exact string `files N · status unavailable`; diagnostic elision pinned to middle-only, command text never lost (sharpens the degraded-rendering and truncation Questions).
- **Sharpened** — stray suffix pinned to ` · N strays`; stray ordering pinned to alpha after the Spec's steps (Research confirmed `EnrichChange` already appends unclaimed artifacts in that order — no new ordering logic).
- **Added** — the label-stability criterion gains its measurable form: primary label equals the schema artifact ID in every render, on every surface.
- **Removed** — none. **Relaxed** — none.

## Human Validation Process

unchanged from drafted
