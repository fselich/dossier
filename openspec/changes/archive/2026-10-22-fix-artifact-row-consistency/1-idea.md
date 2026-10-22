# Idea: artifact rows render stable identity and honest status

## Pain

The Active Work panel reads as random. A row renames itself the moment its file is authored (`plan` becomes `5-plan`), so numbered and bare labels mix inside one change's step list. Changes whose schema is not installed render every file `[authored]` with `planning 4/4` — indistinguishable from validated progress — behind a diagnostic truncated mid-command. Stale old-vintage files wear the same `[authored]` badge as real step rows and inflate the planning fraction (`planning 6/9` where the true Spec stands at 3/6).

## Questions

- What is a row's stable identity — the schema artifact ID, the output filename stem, or a composition of both?
- Does a file's presence on disk ever imply its artifact is authored?
- How must a change render when its status cannot be loaded — what makes degraded state visibly honest?
- Where do unclaimed output files belong in the step ordering and the planning count?
- What may be truncated in diagnostics, and what must always remain readable?
- Which surfaces display the artifact's name (index rows, artifact tabs, filter) and must they agree?
