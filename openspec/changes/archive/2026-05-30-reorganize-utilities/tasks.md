## 1. Move Functions to openspec Package

- [x] 1.1 Move `extractRequirement` from `ui/viewport.go` to `openspec/loader.go` as exported `ExtractRequirement`
- [x] 1.2 Move `configToMarkdown` from `ui/view.go` to `openspec/loader.go` as exported `ConfigToMarkdown`
- [x] 2.1 Update `ui/viewport.go` to import and call `openspec.ExtractRequirement`, remove local function
- [x] 2.2 Update `ui/view.go` to import and call `openspec.ConfigToMarkdown`, remove local function
- [x] 3.1 Verify tests pass — run `go test ./...`
