## 1. Tests for internal/openspec (loader functions)

- [x] 1.1 Create `internal/openspec/loader_test.go` with `setupProjectDir(t, changes)` helper that scaffolds a tempdir with `openspec/changes/<name>/.openspec.yaml` structure
- [x] 1.2 Test `LoadFrom`: project with multiple active changes, empty changes dir, missing `openspec/` dir (returns error)
- [x] 1.3 Test `LoadFromPath`: valid change path, nonexistent path (returns error), path without `.openspec.yaml` (returns error)
- [x] 1.4 Test `LoadProjectSpecs`: specs dir with subdirectories containing `spec.md`, empty specs dir, specs dir missing (returns nil + error)
- [x] 1.5 Test `ListChangeNames`: with active changes, empty directory, directory missing
- [x] 1.6 Test `ListArchiveChanges`: with dated archive directories, empty archive, missing archive dir
- [x] 1.7 Test `ReloadChange`: file modified on disk produces updated content, file deleted produces absent artifact
- [x] 1.8 Test `LoadConfig`: valid YAML, missing file (empty + nil error), malformed YAML (returns error)

## 2. Tests for internal/openspec (tasks functions)

- [x] 2.1 Test `ParseTasks`: content with sections + pending + done tasks, content with only sections (no tasks), empty content, content with mixed checkbox formats
- [x] 2.2 Test `ToggleTask`: toggle pending → done writes to disk, toggle done → pending writes to disk, idx out of range returns nil, read-only file returns error
- [x] 2.3 Test `FindCursorByText`: text found in items, text not found (returns first task), only sections (returns 0)

## 3. Tests for internal/ui (rendering helpers)

- [x] 3.1 Create `internal/ui/view_test.go`
- [x] 3.2 Test `extractRequirement`: requirement name found in content (returns block), name not found (returns empty), last requirement in document, requirement with no following header
- [x] 3.3 Test `buildIndexItems`: with active changes + specs + archived, empty index, sort order verification
- [x] 3.4 Test `firstAvailableTab`: change with all tabs, change with only proposal and tasks, change with no artifacts
- [x] 3.5 Test `renderTasksContent`: with task cursor at valid position (renders cursor indicator), with all tasks done (shows 100%), empty task list

## 4. Validation

- [x] 4.1 Run `go test -race -count=1 ./internal/...` and ensure all tests pass with no race conditions
- [x] 4.2 Run `go test -coverprofile=coverage.out ./internal/... && go tool cover -func=coverage.out` and verify openspec > 60%, ui > 40%
