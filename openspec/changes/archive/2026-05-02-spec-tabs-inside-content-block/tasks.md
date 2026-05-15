## 1. Move the subnav to the content block

- [x] 1.1 In `View()` (`internal/ui/model.go`), remove the conditional block that adds `renderSpecSubnav` before `boxInnerSep`
- [x] 1.2 In `View()`, add a conditional block that adds `renderSpecSubnav` immediately after `boxInnerSep` and before `vp.View()`

## 2. Verification

- [x] 2.1 Build and open a change with multiple specs; confirm that the chips appear inside the content block, below the separator
- [x] 2.2 Verify that on scroll the subnav remains visible (it is a static row, not part of the viewport)
- [x] 2.3 Verify that on tabs other than `specs` the chip row is not shown
- [x] 2.4 Verify that the viewport does not overflow (height remains correct)
