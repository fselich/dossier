## 1. Theme definition

- [x] 1.1 Create `internal/ui/themes.go` with `Theme` struct (`Name`, `GlamourStyle`, `ChromaStyle`, `ViewBg`) and `ThemeColors` struct (placeholder for phase 2)
- [x] 1.2 Define three built-in themes (`dark`, `light`, `dracula`) as a `map[string]Theme`
- [x] 1.3 Add `DefaultTheme()` function returning the `dark` theme

## 2. CLI flag parsing

- [x] 2.1 Replace manual `os.Args` parsing in `cmd/dossier/main.go` with Go `flag` package (`--theme`, `--version`, `--help`, positional `[path]`)
- [x] 2.2 Look up theme by `--theme` flag value (case-insensitive), exit with error listing valid themes if not found
- [x] 2.3 Pass resolved `Theme` to `ui.New()` and `ui.NewSinglePath()`

## 3. Model integration

- [x] 3.1 Update `Model.theme` to use the `Theme` type from `themes.go` (remove or repurpose existing `Theme struct` in `model.go`)
- [x] 3.2 Update `New()` and `NewSinglePath()` signatures to accept `Theme` parameter
- [x] 3.3 Update test helper `newTestModel()` in `view_test.go` to pass default theme

## 4. Glamour threading

- [x] 4.1 Update `ensureRenderer()` to use `m.theme.GlamourStyle` instead of hardcoded `"dark"`
- [x] 4.2 Handle glamour renderer creation failure gracefully (return early, let callers render raw markdown as fallback)

## 5. Chroma threading

- [x] 5.1 Replace `chromaOnce`/`chromaStyle` globals with `chromaStyleCache map[string]*chroma.Style` and `getChromaStyle(name string)` function
- [x] 5.2 Add `chromaStyleName string` parameter to `highlightLine()`, `renderDiff()` (package functions)
- [x] 5.3 Thread `m.theme.ChromaStyle` through `refreshGitViewport()` → `renderDiffContent()` → `renderDiff()` → `highlightLine()`

## 6. View background

- [x] 6.1 Confirm `tea.View.BackgroundColor` is correctly set from `m.theme.ViewBg` in `View()` (existing code in `model.go:237`)
- [x] 6.2 Ensure nil `ViewBg` (e.g., from default-initialized Theme{}) does not break rendering — terminal default background used as fallback

## 7. Tests

- [x] 7.1 Test that built-in theme map contains expected themes with correct glamour/chroma style names
- [x] 7.2 Test `getChromaStyle()` returns correct style for valid name and falls back for invalid name
- [x] 7.3 Test `ensureRenderer()` uses theme's glamour style (verify renderer is created with correct style)
- [x] 7.4 Test `--theme` flag parsing: valid name, invalid name, default, case-insensitive
- [x] 7.5 Run `make test` and `make lint` to verify no regressions
