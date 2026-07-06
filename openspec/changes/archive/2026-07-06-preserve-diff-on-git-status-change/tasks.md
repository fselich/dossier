## 1. Core Implementation

- [x] 1.1 In `pollGitStatus()`, after updating the file list, skip clearing the diff if `ShowDiff` is true and the viewed file (`DiffFile`) still exists in the new list with the same git status (`X`, `Y`)

## 2. Verification

- [x] 2.1 Run `make test` to confirm all tests pass
- [x] 2.2 Run `make lint` to confirm no linting issues
