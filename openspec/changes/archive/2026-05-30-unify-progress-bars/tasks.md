## 1. Extract Unified Progress Bar

- [x] 1.1 Extract `renderProgressBar(done, total, width int, filledChar, emptyChar string) string` in `tasks.go`
- [x] 1.2 Update `renderTabBar()` in `view.go` to call `renderProgressBar`
- [x] 1.3 Update `renderActiveItem()` in `index.go` to call `renderProgressBar`
- [x] 1.4 Update `progressBar()` in `tasks.go` to call `renderProgressBar`

## 2. Merge Duplicate View Functions

- [x] 2.1 Merge `viewIndexContent` and `viewConfigContent` into `viewContentWithChrome()`

## 3. Verification

- [x] 3.1 `go build` and `go test` pass — all rendering unchanged
