## 1. Message type and path helper

- [x] 1.1 Define `editorReturnMsg` in `model.go` (empty type, used to signal the editor's return)
- [x] 1.2 Add helper `(m *Model) artifactPath() string` that returns the path of the active artifact file based on the current tab (`proposal.md`, `design.md`, `tasks.md`, first `spec.md` from `specs/`)

## 2. Launch the editor

- [x] 2.1 Add case `"e"` in the `Update` key handler: if `m.tabAvailable(m.tab)`, build `exec.Command(editor, path)` where `editor` is `$EDITOR` or `"vi"` as fallback, and return `tea.ExecProcess(cmd, func(err error) tea.Msg { return editorReturnMsg{} })`

## 3. Reload on return from editor

- [x] 3.1 Add case `editorReturnMsg` in `Update`: call `ReloadChange` on the current change, apply the changes to the model (same as the tick handler — tasks re-parse with cursor restore, markdown cache invalidation), and call `m.loadViewport()` to refresh the view immediately

## 4. Help bar

- [x] 4.1 Update `renderHelpBar()` in `model.go` to include `e: edit` in both variants of the help text (tasks tab and markdown tab)
