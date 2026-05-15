## 1. Separators

- [x] 1.1 Add `separatorStyle` in `internal/ui/styles.go` (foreground color "0", no bold)
- [x] 1.2 Add method `renderSeparator() string` in `model.go` that returns `strings.Repeat("─", m.width)` with `separatorStyle`
- [x] 1.3 Insert the two separators in `View()`: one between tabBar and viewport, another between viewport (or globalProgressBar) and helpBar
- [x] 1.4 Update `contentHeight()` to subtract 2 extra lines (one per separator)
