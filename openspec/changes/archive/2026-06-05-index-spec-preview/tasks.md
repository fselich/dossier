## 1. Helper function

- [x] 1.1 Add `ExtractPurpose(content string) string` to the `openspec` package that extracts plain text between `## Purpose` and the next `##` heading
- [x] 1.2 Add unit tests for `ExtractPurpose` covering: purpose present, purpose absent, purpose at EOF

## 2. Chrome adjustments

- [x] 2.1 Add `chromeSpecPreview = 2` constant in `model.go` and update `contentHeight()` to include it in `ModeIndex`
- [x] 2.2 Add `renderSpecPreview()` function in `view.go` that returns the formatted bar text (or empty string)
- [x] 2.3 Insert the preview bar between viewport and helpbar in `viewContentWithChrome()`
- [x] 2.4 Verify the bar renders correctly in the index layout

## 3. Cursor interaction

- [x] 3.1 In `renderSpecPreview()`, determine the currently selected spec by examining the cursor item kind; extract its purpose and format the bar
- [x] 3.2 Handle the three cursor states: spec item (show), requirement item (show parent spec), other (empty)
- [x] 3.3 Ensure truncation with `…` when the combined name + separator + purpose exceeds available width
