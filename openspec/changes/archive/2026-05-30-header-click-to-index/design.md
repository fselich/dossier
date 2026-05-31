## Context

The header row (screen Y=1) displays `project · change-name [N/M]` in `ModeNormal` and `ModeViewingArchive`. It's rendered by `renderHeader()` in `view.go` and currently has no mouse click handler.

## Goals / Non-Goals

**Goals:**
- Left-click on header in `ModeNormal` → enter index
- Left-click on header in `ModeViewingArchive` → enter index
- Ignore header click in `ModeIndex`, `ModeViewingSpec`, `ModeViewingConfig`

**Non-Goals:**
- Right-click or middle-click on header
- Visual indicator that header is clickable

## Decisions

### Decision: One extra check in handleMouseClick

In `handleMouseClick`, before the tab bar check (Y==2), add:

```
if msg.Y == 1 && (m.mode == ModeNormal || m.mode == ModeViewingArchive) {
    m.enterIndex()
    return m, nil
}
```

`enterIndex()` already handles all initialization (loads specs, archives, builds items, sets cursor). No new logic needed.

## Risks / Trade-offs

- The header always looks the same regardless of clickability. Users discover this by trying. Could add a subtle hint later.
