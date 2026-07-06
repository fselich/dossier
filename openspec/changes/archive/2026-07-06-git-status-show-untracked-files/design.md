## Context

The "code" tab uses `git status --porcelain` to list changed files. By default, git collapses untracked files in brand-new directories into a single `?? dir/` entry. This means files like `src/Domain/Cache/CacheInterface.php` show as `?? src/Domain/Cache/`, hiding individual filenames and preventing the user from selecting or diffing specific files.

## Goals / Non-Goals

**Goals:**
- Expand collapsed untracked directories into individual file entries in the "code" tab file list

**Non-Goals:**
- No changes to the diff viewer, git polling, cursor navigation, or any other part of the UI
- No changes to the `openspec/` directory filtering

## Decisions

- **Approach**: Pass `-u` (short form of `--untracked-files=all`) to `git status --porcelain`
  - `--untracked-files=normal` (default): git shows untracked directories as `?? dir/` without listing contents
  - `--untracked-files=all`: git lists every untracked file individually as `?? dir/file.go`
  - Using `--untracked-files=normal` would not solve the problem
  - **Rejected alternatives**: Post-processing the collapsed directories with `os.ReadDir` to expand them — adds complexity, duplicates git's own functionality, and would need to handle `.gitignore` correctly. Letting git do the expansion is simpler and more correct.

## Risks / Trade-offs

- [Noise] Changes with many generated files in untracked directories could produce longer file lists. Mitigation: users are unlikely to generate hundreds of files before committing, and the file list already handles large lists via viewport scrolling.
- [No risk] The `-u` flag is a stable, widely supported git option (since git 1.7.0). No compatibility concerns.
