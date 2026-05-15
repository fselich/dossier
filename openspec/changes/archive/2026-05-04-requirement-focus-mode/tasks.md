## 1. Model state

- [x] 1.1 Add fields `specFocusMode bool` and `specReqCursor int` to `Model` in `internal/ui/model.go`

## 2. Extraction helper

- [x] 2.1 Implement `extractRequirement(raw, name string) string` that returns the markdown block of the specified requirement (from `### Requirement: <name>` to the next `### Requirement:` or EOF)

## 3. Entry from the index

- [x] 3.1 In the `indexKindRequirement` → `Enter` handler, set `specFocusMode = true` and `specReqCursor = item.reqIdx` in addition to the already existing `specJumpTarget`
- [x] 3.2 In the `indexKindSpec` → `Enter` handler, ensure `specFocusMode = false` and `specReqCursor = 0` are cleared

## 4. Rendering in focus mode

- [x] 4.1 In the `loadViewport()` branch for `ModeViewingSpec`, if `specFocusMode == true`, call `extractRequirement` with the current `specJumpTarget` and pass the extracted block to glamour instead of the full `Content`
- [x] 4.2 Remove the `jumpLine` search when `specFocusMode == true` (no scroll needed)

## 5. h/l navigation in focus mode

- [x] 5.1 In the `case "h"` handler, add a branch for `m.mode == ModeViewingSpec && m.specFocusMode`: decrement `specReqCursor` with wrap, update `specJumpTarget`, and call `loadViewport()`
- [x] 5.2 In the `case "l"` handler, add an equivalent branch to increment `specReqCursor` with wrap

## 6. Header and HelpBar

- [x] 6.1 In `renderHeader()`, if `m.mode == ModeViewingSpec && m.specFocusMode`, show `<project>  ·  <spec-name>  ·  Req N/M`
- [x] 6.2 In the helpbar, if `m.mode == ModeViewingSpec && m.specFocusMode`, show `h/l: prev/next req  Esc: index  q: quit`

## 7. Verification

- [x] 7.1 Navigate from the index to a requirement and confirm that only that requirement is shown
- [x] 7.2 Press `l` and `h` and confirm that the viewport changes to the next/previous requirement
- [x] 7.3 Confirm that the header shows `Req N/M` correctly
- [x] 7.4 Confirm that `Esc` returns to the index with the cursor on the correct requirement
- [x] 7.5 Confirm that opening a spec from its index item (not a requirement) still shows the full spec
