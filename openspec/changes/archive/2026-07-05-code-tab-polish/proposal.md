## Why

Two issues with the git tab:

1. The tab appears in archive mode (`ModeViewingArchive`) as disabled (grayed out), adding visual noise. It's irrelevant for archived changes since they have no working tree to monitor.

2. The label `changes` is ambiguous — it could refer to OpenSpec changes (the core concept of the app) or git file changes. `code` is clearer and shorter.

## What Changes

- Rename tab label from `changes` to `code`.
- Hide the git tab entirely in `ModeViewingArchive` mode (not just disabled).

## Capabilities

### Modified Capabilities

- `git-status-tab`: Tab label changed from `changes` to `code`. Tab is hidden in archive mode.

## Impact

- `internal/ui/model.go`: change `tabLabels[TabGit]` from `"changes"` to `"code"`.
- `internal/ui/view.go`: skip rendering TabGit in `renderTabBar` when mode is not ModeNormal.
