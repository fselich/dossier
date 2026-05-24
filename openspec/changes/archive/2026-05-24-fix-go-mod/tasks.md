## 1. Fix go.mod

- [x] 1.1 Run `go mod tidy` to reclassify direct vs indirect dependencies
- [x] 1.2 Verify direct imports (`bubbletea`, `glamour`, `lipgloss`, `bubbles`, `yaml.v3`) appear without `// indirect` in `go.mod`
- [x] 1.3 Run `go build ./...` to verify compilation succeeds
- [x] 1.4 Run `go vet ./...` to verify no issues
