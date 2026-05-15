## 1. Fix the return order of ListSpecNames and ListArchiveNames

- [x] 1.1 In `ListSpecNames()` in `internal/openspec/loader.go`, add `sort.Strings(names)` before the final `return names`, so the order matches that of `LoadProjectSpecs()` (ascending alphabetical)
- [x] 1.2 In `ListArchiveNames()` in `internal/openspec/loader.go`, add `sort.Sort(sort.Reverse(sort.StringSlice(names)))` before the final `return names`, so the order matches that of `ListArchiveChanges()` (descending)
