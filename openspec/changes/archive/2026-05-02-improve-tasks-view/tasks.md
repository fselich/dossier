## 1. Word wrap on all items

- [x] 1.1 Apply `.Width(m.width - 2)` to `taskPendingStyle` and `taskDoneStyle` in `renderTasksContent()` (currently only `taskCursorStyle` has Width)
- [x] 1.2 Verify that items without cursor also wrap in narrow terminals

## 2. Inline markdown renderer

- [x] 2.1 Define `codeStyle` and `boldStyle` in `internal/ui/styles.go`
- [x] 2.2 Implement function `inlineMarkdown(s string) string` in `internal/ui/model.go` with regex for `` `code` `` and `**bold**`
- [x] 2.3 Apply `inlineMarkdown()` to each item's text in `renderTasksContent()` before the lipgloss render

## 3. Global progress bar

- [x] 3.1 Calculate total and completed counts for all `KindTask` items of the change before the items loop in `renderTasksContent()`
- [x] 3.2 Render the global progress bar as the first content line using `progressBar()` with a larger width (10 blocks)
- [x] 3.3 Omit the bar if the change has no task-type items
