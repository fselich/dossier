## 1. Rename entry point directory

- [x] 1.1 Rename `cmd/spec-viewer/` to `cmd/specview/` using `git mv`

## 2. Clean up project root

- [x] 2.1 Delete untracked binary `main` from the project root
- [x] 2.2 Delete untracked binary `sv` from the project root

## 3. Add Makefile

- [x] 3.1 Create `Makefile` at the project root with `build`, `install`, and `clean` targets

## 4. Update .gitignore

- [x] 4.1 Add `specview`, `main`, and `sv` to `.gitignore` to prevent committing compiled binaries

## 5. Verify

- [x] 5.1 Run `make build` and confirm `./specview` is produced
- [x] 5.2 Run `make install` and confirm `specview` is available from PATH
- [x] 5.3 Run `make clean` and confirm `./specview` is removed
