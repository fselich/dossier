## 1. Reorder indexItems construction

- [x] 1.1 In `internal/ui/model.go` (~line 593), move the block that adds `indexKindArchived` after the block that adds `indexKindSpec` and `indexKindRequirement`, so that the order is: active → spec/requirement → archived.

## 2. Reorder index rendering

- [x] 2.1 In `internal/ui/model.go` (~line 629), move the "Archived Changes" rendering block (~line 652) so that it appears after the "Specifications" rendering block (~line 680).

## 3. Update tests

- [x] 3.1 Find tests that verify the order of sections or items in `ModeIndex` and update the expected indexes/order to reflect `Active → Specifications → Archived`.
