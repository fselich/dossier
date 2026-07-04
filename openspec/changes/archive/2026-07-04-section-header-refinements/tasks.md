## 1. Enter guard on section headers

- [x] 1.1 Add `if item.kind == indexKindSection { return m, nil }` guard in `updateIndex` before the archive fallthrough (around line 756)
- [x] 1.2 Add same guard in `clickIndexItem` (mouse.go) for consistency (currently a no-op via switch fallthrough, make explicit)

## 2. Visual indicator refinements

- [x] 2.1 In `renderIndexContent`, change section header rendering: no indicator when expanded, `…` (unicode ellipsis) in `helpStyle` at end when collapsed
- [x] 2.2 Remove the `▼`/`▶` character prefix from the header string; collapse indicator is now a suffix, not a prefix

## 3. Tests

- [x] 3.1 Add test: Enter on section header returns without navigation or crash
- [x] 3.2 Update existing collapse visual tests to reflect new indicator format (ellipsis, no marker when expanded)
