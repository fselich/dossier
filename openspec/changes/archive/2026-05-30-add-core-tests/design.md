## Context

The project has 3 test files (`view_test.go`, `tasks_test.go`, `loader_test.go`) totaling ~700 lines. The `view_test.go` covers `extractRequirement`, `IndexItemAtContentLine`, `FirstAvailableTab`, `BuildIndexItems`, `ClickInIndex`, `HeaderClick`, and `RenderTasksContent`. No tests exist for `Update`, `doToggle`, `loadViewport`, `handleTick`, or `renderTabBar`.

## Goals / Non-Goals

**Goals:**
- Add smoke tests for `Update()` covering all keybindings in their primary modes
- Add `doToggle()` tests verifying disk writes and error cases
- Add `loadViewport()` tests verifying correct mode dispatch
- Add `handleTick()` tests verifying reload behavior
- Add `renderTabBar()` tests verifying output for all tab states

**Non-Goals:**
- 100% branch coverage or exhaustive edge-case testing
- Mocking the filesystem in `doToggle()` — use a real temp directory
- Testing glamour rendering output (it's visual, tested manually)
- Testing mouse wheel/motion handlers (covered by existing mouse tests)

## Decisions

**Decision 1: Use real temp directories for disk-backed tests**

`doToggle()` writes to disk via `openspec.ToggleTask()`. Tests create a temp directory with a `tasks.md` fixture, build a `Model` pointing at it, call `doToggle()`, then verify the file was modified. No mock filesystem needed.

**Decision 2: Smoke tests, not exhaustive**

For `Update()` keypress tests, verify the mode/tab/cursor changes expected for each keybinding in its primary mode. Don't test every edge case — these are regression guards, not unit-level correctness proofs.

**Decision 3: Tests live in `view_test.go`**

Keep all UI tests in `internal/ui/view_test.go` to match the existing convention. No new test files.

**Decision 4: `loadViewport()` tests verify mode dispatch**

Build a `Model` in each mode, call `loadViewport()`, verify the returned `tea.Cmd` is nil or non-nil as expected, and that viewport content was set correctly.

**Decision 5: `handleTick()` tests verify poll behavior**

Build a `Model` on a temp directory, write a change to disk after initialization, call `handleTick()`, verify the model's `project.Changes` was updated.

## Risks / Trade-offs

- **[Risk]** Temp directory tests may be slow → **Mitigation**: tests only create tiny fixtures (1-2 line markdown files), so overhead is minimal
- **[Risk]** `doToggle()` tests may leave temporary files → **Mitigation**: use `t.TempDir()` which `go test` automatically cleans up
