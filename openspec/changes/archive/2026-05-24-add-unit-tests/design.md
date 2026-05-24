## Context

Current test coverage is nearly zero (2 test files, ~174 lines). The `inject-root-path` change (P0) changes function signatures across the `openspec` package — doing this without tests risks silent breakage.

The `openspec` package is purely functional (file I/O + parsing), making it straightforward to test with `t.TempDir()`. The `ui` package mixes rendering (string output) with Bubble Tea state — rendering helpers can be tested as pure functions, but full `Update`/`View` cycle tests are out of scope for now.

## Goals / Non-Goals

**Goals:**
- Test all public functions in `internal/openspec/loader.go` using `t.TempDir()` with filesystem scaffolding
- Test rendering helpers in `internal/ui/` that are pure functions (no Bubble Tea runtime)
- Achieve >60% coverage in `openspec`, >40% in `ui`
- Use table-driven tests (idiomatic Go pattern)
- Use `t.Parallel()` where possible to keep test suite fast

**Non-Goals:**
- Bubble Tea integration tests (require `tea.NewProgram`, out of scope for unit tests)
- Viewport rendering tests (depend on terminal dimensions and glamour, not unit-testable)
- Key dispatch tests (require full Update cycle infrastructure)
- 100% coverage (diminishing returns; focus on high-risk functions from mejoras.md 1.1)

## Decisions

### Decision 1: `t.TempDir()` scaffolding pattern

Each test creates a temporary directory with the exact openspec structure needed, then calls the function under test. This avoids the race conditions of `os.Chdir` in parallel tests.

Helper pattern:
```go
func setupProjectDir(t *testing.T, changes []string) string {
    t.Helper()
    root := t.TempDir()
    // Create openspec/ directory structure
    os.MkdirAll(filepath.Join(root, "openspec", "changes"), 0755)
    for _, name := range changes {
        dir := filepath.Join(root, "openspec", "changes", name)
        os.MkdirAll(dir, 0755)
        os.WriteFile(filepath.Join(dir, ".openspec.yaml"), []byte("created: 2026-05-24"), 0644)
    }
    return root
}
```

### Decision 2: Test rendering helpers as pure functions, not Bubble Tea components

Functions like `extractRequirement`, `renderTasksContent`, `buildIndexItems` take inputs and return strings or data structures — no Bubble Tea runtime needed. This keeps tests fast and deterministic.

Full `Update()` method tests (key presses, mode transitions, viewport rendering) will be added later when the model is broken into sub-models (item 2.3).

### Decision 3: Table-driven tests

All tests follow Go's idiomatic table-driven pattern:
```go
tests := []struct {
    name    string
    input   ...
    want    ...
    wantErr bool
}{...}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // test body
    })
}
```

## Risks / Trade-offs

- **Tests may break on function signature changes**: The `inject-root-path` change modifies signatures. Risk: tests written against current signatures need updating. Mitigation: implement `inject-root-path` first, then write tests against new signatures. OR write tests now and update them during the refactor. Decision deferred to implementation order.
- **UI tests fragile**: Rendering output depends on Lipgloss color codes and exact spacing. Risk: cosmetic changes break tests. Mitigation: test structure (presence of sections, requirement names) not exact ANSI sequences. Use `strings.Contains` over exact match.
