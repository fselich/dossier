## 1. Change Enter/e to show diff

- [x] 1.1 In `viewer.go`, change `case "enter":` for TabGit from `openGitFile` to `toggleGitDiff`
- [x] 1.2 In `viewer.go`, change `case "e":` for TabGit from `openGitFile` to `toggleGitDiff`

## 2. Disable tab when no files

- [x] 2.1 In `model.go`, update `tabAvailable(TabGit)` to add `&& len(m.gitState.Files) > 0`

## 3. Remove dead code

- [x] 3.1 Remove `openGitFile` function from `git.go`
- [x] 3.2 Run `make fmt` to clean up unused imports

## 4. Verify

- [x] 4.5 `make build` succeeds
- [x] 4.6 `make lint` passes
- [x] 4.7 `make test` passes
- [x] 4.1 `Enter` on a file shows diff view
- [x] 4.2 `e` on a file shows diff view
- [x] 4.3 Tab is disabled when working tree is clean
- [x] 4.4 Tab becomes enabled when files appear
