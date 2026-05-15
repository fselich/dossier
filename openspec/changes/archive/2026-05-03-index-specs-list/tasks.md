## 1. Loader: project specs support

- [x] 1.1 Add the `ProjectSpec` type with fields `Name string` and `RequirementCount int` in `internal/openspec/loader.go`
- [x] 1.2 Implement `LoadProjectSpecs() []ProjectSpec` that reads `openspec/specs/`, counts requirements from each `spec.md`, and returns the list sorted alphabetically
- [x] 1.3 Verify that `LoadProjectSpecs()` returns an empty list (without error) if the directory does not exist or is empty

## 2. UI Model: specs state

- [x] 2.1 Add field `projectSpecs []openspec.ProjectSpec` to the `Model` struct in `internal/ui/model.go`
- [x] 2.2 Call `openspec.LoadProjectSpecs()` inside `enterIndex()` and assign it to `m.projectSpecs`, same as is done with `archiveChanges`

## 3. Render: Specs section in the index

- [x] 3.1 Add at the end of `renderIndexContent` the "Specs" section with header `sectionStyle.Render("Specs")`
- [x] 3.2 If `m.projectSpecs` is empty, show `helpStyle.Render("  No specs available")`
- [x] 3.3 For each spec, render a line with the name and the number of requirements in `helpStyle`
- [x] 3.4 Verify that specs are not included in `indexItems` and are not reachable by the cursor
