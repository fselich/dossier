## Context

The TUI uses a `tea.Tick(2*time.Second)` message to poll disk state. On each tick, `loader.LoadProject` re-reads artifact files and compares content hashes to detect changes. Tab availability is computed once at startup from the initial load result and never updated; if an artifact is absent at launch and appears later, the corresponding tab stays disabled until restart. The task toggle writes to disk but does not update the in-memory Done flag directly, so the progress counter in the tab bar lags by up to one tick.

Relevant code lives in `internal/ui/model.go` (tick handler, tab rendering, task toggle) and `internal/openspec/loader.go` (project loading).

## Goals / Non-Goals

**Goals:**
- Detect new artifacts appearing on disk in under 500 ms
- Enable previously-disabled tabs when their artifact becomes present
- Reload the change list when the TUI starts with zero active changes and a change is created externally
- Reflect task toggle state in the progress counter immediately (no tick lag)

**Non-Goals:**
- Detecting artifact deletions (present → absent) — tabs remain enabled once unlocked
- Watching subdirectory tree changes with inotify/fsnotify — polling is sufficient given the 500 ms target
- Animated or transition effects on tab unlock

## Decisions

**D1: 500 ms poll interval instead of 2 s**

The proposal targets "less than half a second" reactivity. 500 ms is the simplest value that satisfies that. Halving to 250 ms gives diminishing returns and doubles syscall pressure; 1 s is noticeably laggy for interactive use.

Alternatives considered:
- `fsnotify` file-system events: zero-lag but adds a dependency, cross-platform complexity, and requires fd management. Overkill for a local CLI tool.
- 250 ms: Marginal improvement over 500 ms, doubles poll rate for no user-visible benefit.

**D2: Presence detection by comparing `artifact.Present` flags, not content hashes**

The existing tick handler compares `artifact.Hash` values to decide whether to re-render. Adding a separate pass that checks `oldArtifact.Present != newArtifact.Present` covers the absent→present transition without touching the hash logic. The two checks remain orthogonal.

**D3: Reload change list when `len(m.project.Changes) == 0`**

When the TUI starts with no active changes (e.g., before the first `openspec new change` call), `m.project` holds an empty slice. On the next tick after a change directory is created, `LoadProject` is called again, and if it now returns at least one change the model adopts it. This avoids a special "watching for new change" state machine; the existing poll loop handles it.

**D4: Immediate Done-flag update on task toggle**

Currently `handleTaskToggle` writes the updated file and waits for the next tick to re-read it. Instead, update `m.taskItems[cursor].Done` in memory right after the write succeeds. This makes the progress counter in the tab bar update on the same frame as the keystroke. The next tick will re-read and confirm (or correct) the disk state, which is a safe eventual-consistency guarantee.

## Risks / Trade-offs

[Risk: CPU usage increase from 2 s → 500 ms polling] → Mitigation: The poll body is a directory read + stat calls on a handful of files. At 500 ms the overhead is negligible on any modern laptop. Acceptable trade-off.

[Risk: Race between TUI toggle write and tick re-read] → Mitigation: D4 ensures the UI reflects the toggle immediately. If a concurrent external edit overwrites the file between toggle and next tick, the tick corrects the in-memory state. No data loss; at worst a brief visual flicker.

[Risk: Tab unlock is one-way] → Mitigation: This is intentional (Non-Goal). A tab that becomes available stays available for the session; re-disabling would be surprising UX. If an artifact is deleted mid-session the tab shows stale content until restart.

## Migration Plan

No migration needed. The change is limited to two files; the binary is rebuilt and replaces the existing one. No config changes, no data migrations.

Rollback: revert the two files, rebuild.

## Open Questions

None. The scope is narrow and all decisions are resolved.
