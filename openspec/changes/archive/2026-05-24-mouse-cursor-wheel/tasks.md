## 1. Modify wheel handling in mouse.go

- [x] 1.1 Replace `m.vp.LineUp(3)` / `m.vp.LineDown(3)` in `handleMouse()` with mode-aware dispatch: index mode moves index cursor, tasks tab moves task cursor, all else scrolls viewport

## 2. Verification

- [x] 2.1 Manual test: wheel up/down moves index cursor in `ModeIndex`, viewport auto-follows
- [x] 2.2 Manual test: wheel up/down moves task cursor in `TabTasks` / `ModeNormal`, viewport auto-follows
- [x] 2.3 Manual test: wheel scrolls viewport 3 lines in all non-cursor views (proposal, design, specs, config, archive)
- [x] 2.4 Manual test: wheel at first index item with up does nothing (no crash)
- [x] 2.5 Manual test: wheel at last index item with down does nothing (no crash)
