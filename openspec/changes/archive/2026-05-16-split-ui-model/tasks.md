## 1. Preparation

- [x] 1.1 Verify the project builds cleanly: `go build ./...`
- [x] 1.2 Run existing tests as a baseline: `go test ./...`

## 2. Create new files (empty, with package declaration)

- [x] 2.1 Create `internal/ui/update.go` with `package ui`
- [x] 2.2 Create `internal/ui/viewport.go` with `package ui`
- [x] 2.3 Create `internal/ui/index.go` with `package ui`
- [x] 2.4 Create `internal/ui/tasks.go` with `package ui`
- [x] 2.5 Create `internal/ui/view.go` with `package ui`

## 3. Move functions — one file at a time

- [x] 3.1 Move `Update` to `update.go`; verify build
- [x] 3.2 Move `loadViewport` to `viewport.go`; verify build
- [x] 3.3 Move index functions to `index.go`: `handleTick`, `viewIndex`, `enterIndex`, `specSuffix`, `buildSpecOrder`, `buildIndexItems`, `refreshIndexViewport`, `renderIndexContent`, `renderActiveItem`, `renderArchivedItem`, `taskCounts`, `sameNames`, `sameStrings`, `extractRequirement`; verify build
- [x] 3.4 Move task functions to `tasks.go`: `loadTaskItems`, `firstTaskIdx`, `moveCursorDown`, `moveCursorUp`, `doToggle`, `refreshTasksViewport`, `renderTasksContent`, `sectionProgress`, `progressBar`, `extractOpeningEscape`, `inlineMarkdown`, and the `var` block for the regex/escape vars; verify build
- [x] 3.5 Move rendering functions to `view.go`: `View`, `renderHeader`, `renderTabBar`, `renderSpecSubnav`, `hasSpecSubnav`, `boxTop`, `boxBottom`, `boxInnerSep`, `addBorderSides`, `renderHelpBar`, `emptyView`; verify build

## 4. Clean up model.go

- [x] 4.1 Confirm `model.go` now contains only: type/const declarations, message types, `Model` struct, `New`, `NewSinglePath`, `Init`, and state query helpers (`current`, `firstAvailableTab`, `tabAvailable`, `defaultTab`, `artifactPath`, `currentArchive`, `contentHeight`)
- [x] 4.2 Remove any leftover blank lines or orphaned imports from `model.go`

## 5. Verify

- [x] 5.1 `go build ./...` passes with no errors
- [x] 5.2 `go test ./...` passes (same result as baseline)
- [x] 5.3 `go vet ./...` passes
- [x] 5.4 Manually run the TUI and exercise each mode (index, normal, archive, spec viewer) to confirm no regressions
