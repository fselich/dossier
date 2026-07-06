## Why

`git status --porcelain` by default collapses untracked files in new directories into a single `?? dir/` entry, hiding individual filenames. This makes the "code" tab show directory placeholders instead of actual files, forcing the user to leave the TUI to discover what was created.

## What Changes

- Pass `-u` (or `--untracked-files=all`) to `git status --porcelain` in `internal/git/git.go:34`
- The "code" tab file list will show each untracked file individually (e.g., `?? src/Domain/Cache/CacheInterface.php`) instead of collapsing the directory (`?? src/Domain/Cache/`)

## Capabilities

### New Capabilities

- (none)

### Modified Capabilities

- `git-status-tab`: Requirement "List of changed files with status indicators" — change the git status command to use individual untracked file listing (`-u` flag) instead of the default collapsed directory mode

## Impact

- Single line change in `internal/git/git.go`
- No API changes, no new dependencies
- The diff viewer for untracked files (`untrackedFileDiffLines`) already handles individual file paths correctly — no changes needed there
