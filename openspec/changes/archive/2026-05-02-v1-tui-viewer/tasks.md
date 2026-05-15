## 1. Project setup

- [x] 1.1 Initialize Go module (`go mod init github.com/user/spec-viewer`)
- [x] 1.2 Add dependencies: `bubbletea`, `lipgloss`, `glamour`, `bubbles`, `yaml.v3`
- [x] 1.3 Create directory structure: `cmd/spec-viewer/`, `internal/openspec/`, `internal/ui/`

## 2. openspec-loader

- [x] 2.1 Implement detection of `./openspec/` in CWD and error if not found
- [x] 2.2 Implement listing of active changes (subdirectories of `changes/` except `archive/`)
- [x] 2.3 Implement reading of `.openspec.yaml` per change (`created` field)
- [x] 2.4 Implement artifact content loading per change (proposal, design, specs) with absent flag
- [x] 2.5 Implement project name inference from directory name

## 3. tasks-toggle

- [x] 3.1 Implement `tasks.md` line-by-line parser: task and section items with line number
- [x] 3.2 Implement j/k navigation between task items (skipping sections)
- [x] 3.3 Implement on-disk toggle: modify only the corresponding line in `tasks.md`
- [x] 3.4 Implement write error handling with temporary message in the TUI

## 4. tui-viewer — base structure

- [x] 4.1 Create main Bubble Tea model with state: current change, active tab, viewport
- [x] 4.2 Implement layout: header (1 line) + tab bar (1 line) + content + help bar (1 line)
- [x] 4.3 Implement h/l navigation between changes with wrap
- [x] 4.4 Implement tab selection with keys 1-4 (disabling tabs for absent artifacts)
- [x] 4.5 Implement welcome screen when there are no active changes

## 5. tui-viewer — content rendering

- [x] 5.1 Integrate glamour for markdown rendering in proposal, design and specs tabs
- [x] 5.2 Fix glamour width to the real viewport width
- [x] 5.3 Integrate `bubbles/viewport` for j/k scrolling in markdown tabs
- [x] 5.4 Implement tasks view with cursor (▶), sections, checkboxes and progress bar per section

## 6. tui-viewer — help bar and styles

- [x] 6.1 Implement contextual help bar (different for tasks tab vs others)
- [x] 6.2 Define lipgloss styles: header, active/inactive/disabled tab, task cursor, section, progress bar
- [x] 6.3 Implement clean exit with q and Ctrl+C

## 7. Integration and manual testing

- [x] 7.1 Test with the openspec of this same project (`spec-viewer`)
- [x] 7.2 Test with the openspec of `lapo` (multiple archived changes, canonical specs)
- [x] 7.3 Test the case of a change without all artifacts (disabled tab)
- [x] 7.4 Test task toggle and verify that the file is written correctly
- [x] 7.5 Test in a terminal with different widths (80, 120, 200 cols)
