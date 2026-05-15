## Context

All UI-visible strings in `internal/ui/model.go` are in Spanish. The change is a pure text replacement with no logic changes.

## Goals / Non-Goals

**Goals:**
- Replace every Spanish string literal in the UI with its English equivalent
- Update the five specs that explicitly specify string values in their requirements

**Non-Goals:**
- Internationalisation infrastructure (no i18n framework, no string extraction)
- Translating spec document prose (only requirement-normative string values change)

## Decisions

### No string constants file

All strings stay inline in `model.go` as literals. Extracting them to a constants file would be over-engineering for a tool with a single locale.

## Risks / Trade-offs

- [Spec drift] Any future spec that quotes a helpbar string verbatim must use the new English value. Mitigated by updating all five affected specs as part of this change.
