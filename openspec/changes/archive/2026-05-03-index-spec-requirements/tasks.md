## 1. Loader — add RequirementNames to ProjectSpec

- [x] 1.1 Add `RequirementNames []string` to the `ProjectSpec` struct in `internal/openspec/loader.go`
- [x] 1.2 In `LoadProjectSpecs()`, extract the name of each requirement (substring after `### Requirement: `, trimmed) in the same loop that counts `RequirementCount`, and store it in `RequirementNames`

## 2. Model — new kind and expand state

- [x] 2.1 Add `indexKindRequirement` to the `indexItemKind` enum in `internal/ui/model.go`
- [x] 2.2 Add field `reqIdx int` to the `indexItem` struct (used for `indexKindRequirement`)
- [x] 2.3 Add field `expandedSpecs map[int]bool` to the `Model` struct
- [x] 2.4 Initialize `expandedSpecs` to an empty map in `enterIndex()`

## 3. buildIndexItems — flatten requirements when spec is expanded

- [x] 3.1 In `buildIndexItems()`, after inserting the `indexKindSpec` item, if `m.expandedSpecs[i]` is true, insert one `indexKindRequirement` item for each name in `m.projectSpecs[i].RequirementNames`

## 4. Space — toggle expand/collapse

- [x] 4.1 In the key handler for `ModeIndex`, add a case `" "` (Space): if the item under the cursor is `indexKindSpec`, toggle `expandedSpecs[item.idx]`; if the cursor would be out of bounds after rebuild (due to collapsing with cursor inside), move it to the spec item; call `buildIndexItems()` and `refreshIndexViewport()`
- [x] 4.2 If the cursor is on an `indexKindRequirement` and the user presses Space, ignore (do nothing)

## 5. renderIndexContent — draw indented requirements

- [x] 5.1 In `renderIndexContent()`, after drawing the spec item, if `m.expandedSpecs[i]` is true, iterate `m.projectSpecs[i].RequirementNames` and draw each name with indentation (4 spaces), cursor marker if applicable, and without the `N requirements` column

## 6. Enter on requirement — open viewer with scroll target

- [x] 6.1 Add field `jumpToRequirement string` to `specRenderedMsg` (empty = no jump)
- [x] 6.2 In the `Enter` handler in `ModeIndex`, add a case for `indexKindRequirement`: save the requirement name, enter `ModeViewingSpec`, launch the render goroutine passing the requirement name as target
- [x] 6.3 In the render goroutine, after rendering with glamour, if `jumpToRequirement` is not empty, scan the rendered output line by line (stripping ANSI with regex `\x1b\[[0-9;]*m`) looking for a line that contains the name; store the found line index in `jumpToLine int` of the message
- [x] 6.4 In the `specRenderedMsg` handler, if `jumpLine > 0`, call `m.viewport.SetYOffset(jumpLine)` after setting the viewport content

## 7. Esc from viewer opened via requirement — return to the requirement item

- [x] 7.1 Ensure that when entering `ModeViewingSpec` from an `indexKindRequirement`, `m.specViewerCursor` points to the correct index in `indexItems` (the requirement item), so that the existing `Esc` restores the cursor to that item
