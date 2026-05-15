## 1. Loader

- [x] 1.1 Add `LoadFromPath(path string) (*Project, error)` in `loader.go`: verify that the directory exists and contains `.openspec.yaml`, build a `Change` with `Name = filepath.Base(path)`, `Path = path`, read its artifacts with the existing functions (`loadFile`, `loadSpecs`), and return a `Project` with `Name = filepath.Base(filepath.Dir(path))` and that single change

## 2. UI Model

- [x] 2.1 Add field `singlePath bool` to `Model` in `model.go`
- [x] 2.2 Add alternative constructor `NewSinglePath(project *Project) Model` that calls `New` and sets `singlePath = true`
- [x] 2.3 In `handleTick`, if `m.singlePath`, skip the `ListChangeNames` / `sameNames` block and proceed directly to `ReloadChange` for the current change

## 3. Main

- [x] 3.1 In `main.go`, if `len(os.Args) > 1`, call `openspec.LoadFromPath(os.Args[1])` instead of `openspec.Load()`, and create the model with `ui.NewSinglePath(project)`; if there is an error, print it and exit with code 1
