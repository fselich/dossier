## 1. Theme

- [x] 1.1 Add `BaseText` field to `ThemeStyles` struct in `internal/ui/themes.go`
- [x] 1.2 Add `BaseText` initialization in `BuildStyles`: `lipgloss.NewStyle().Foreground(c.PrimaryFg)`
- [x] 1.3 Update `themes_test.go`: verify `BaseText.GetForeground() == PrimaryFg` for dark and light

## 2. Index rendering

- [x] 2.1 Apply `m.theme.Styles.BaseText.Render(paddedName)` in `renderActiveItem` for non-selected items
- [x] 2.2 Apply `m.theme.Styles.BaseText.Render(ch.Name) + pad` in `renderArchivedItem` for non-selected items
- [x] 2.3 Apply `m.theme.Styles.BaseText.Render(ps.Name)` in spec rendering inside `renderIndexContent` for non-selected items

## 3. Validation

- [x] 3.1 `make test` passes
- [x] 3.2 `make lint` passes
- [x] 3.3 Manual check: `--theme light` index items are readable when not selected
