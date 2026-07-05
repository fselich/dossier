## Context

The git tab is part of `ModeNormal` — it shows working-tree changes. In `ModeViewingArchive`, the tab bar still includes a disabled (grayed out) git tab because `renderTabBar` iterates over all tabs and shows disabled ones. Other disabled tabs in archive mode (e.g., a missing proposal) make sense because they reference the archive change's artifacts. The git tab doesn't — it has no relation to archived changes.

Additionally, the label `changes` is easily confused with the OpenSpec concept of "changes" (a named modification directory). `code` is unambiguous.

## Goals / Non-Goals

**Goals:**
- Git tab label becomes `code` in the tab bar and in `tabLabels`.
- Git tab is completely hidden in `ModeViewingArchive` (not rendered at all), not just disabled.

## Decisions

### Rename in tabLabels array

One-line change in `model.go`:

```go
var tabLabels = [tabCount]string{"proposal", "design", "specs", "tasks", "code"}
```

This affects: the tab bar text, keyboard shortcut help (shows `5` to access it), and any label-based code. No other changes needed.

### Hide in archive mode

In `renderTabBar`, the loop that builds tab parts skips TabGit entirely when the mode is not `ModeNormal`:

```go
for t := Tab(0); t < tabCount; t++ {
    if t == TabGit && m.mode != ModeNormal {
        continue
    }
    // ... existing rendering
}
```

This is cleaner than adding mode checks to `tabAvailable`, which would affect other logic. The help bar already handles archive mode with its own block and doesn't mention the git tab.

When the user switches from normal mode to archive mode, `renderTabBar` is called with the new mode and the tab disappears. When switching back, it reappears.

## Risks / Trade-offs

- **Tab labels change**: Any external tool or user script that reads the tab bar text "changes" would need updating. This is a cosmetic change with no API impact.
- **Hidden tabs in archive mode**: The tab count shown in help bar is static (`"1-4"` or `"1-5"`). In archive mode, it says `"1-4"` (since `m.isGitRepo` is still true, but `tabCount` is always 5). Actually, the help bar for archive mode is hardcoded as `"1-4/Tab: artifact"` and doesn't depend on `isGitRepo`. So there's no inconsistency.
