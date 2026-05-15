## 1. Go module and source

- [x] 1.1 Rename `go.mod` module path to `github.com/fselich/dossier`
- [x] 1.2 Rename `cmd/specview/` directory to `cmd/dossier/`
- [x] 1.3 Update import paths in `cmd/dossier/main.go`
- [x] 1.4 Update import path in `internal/ui/model.go`
- [x] 1.5 Verify `go build ./...` succeeds

## 2. Build tooling

- [x] 2.1 Update `Makefile` BIN variable from `specview` to `dossier`
- [x] 2.2 Update `Makefile` CMD variable from `./cmd/specview` to `./cmd/dossier`
- [x] 2.3 Update `.gitignore` entry from `specview` to `dossier`

## 3. Documentation

- [x] 3.1 Update all references in `README.md`
- [x] 3.2 Update all references in `README.es.md`

## 4. Project-level specs

- [x] 4.1 Update `openspec/specs/build-tooling/spec.md`: rename requirement "Binary named specview" to "Binary named dossier" and replace all `specview`/`cmd/specview/` references
- [x] 4.2 Update `openspec/specs/path-arg/spec.md`: replace `./spec-viewer` with `./dossier` in all scenarios
- [x] 4.3 Update `openspec/specs/openspec-loader/spec.md`: replace `spec-viewer` with `dossier` in scenarios
