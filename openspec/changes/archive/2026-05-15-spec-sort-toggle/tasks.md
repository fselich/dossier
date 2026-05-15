## 1. Model changes

- [x] 1.1 Add `specSortBySuffix bool` field to `Model` in `internal/ui/model.go`
- [x] 1.2 Add `specOrder []int` field to `Model` in `internal/ui/model.go`

## 2. Sort logic

- [x] 2.1 Add `specSuffix(name string) string` helper (returns substring after last `-`, or full name if no `-`)
- [x] 2.2 Add `(m *Model) buildSpecOrder()` method: builds `m.specOrder` as identity permutation, then `sort.SliceStable` by `specSuffix` when `m.specSortBySuffix` is true

## 3. Wire buildSpecOrder into item building

- [x] 3.1 Call `m.buildSpecOrder()` at the top of `buildIndexItems()` so `m.specOrder` is always in sync
- [x] 3.2 Change the spec iteration loop in `buildIndexItems()` from `for i, ps := range m.projectSpecs` to `for _, i := range m.specOrder` (emit `indexItem{kind: indexKindSpec, idx: i}`)
- [x] 3.3 Change the spec iteration loop in `renderIndexContent()` from `for i, ps := range m.projectSpecs` to `for _, i := range m.specOrder` with `ps := m.projectSpecs[i]`

## 4. Key handler

- [x] 4.1 Add `case "s":` in the `tea.KeyMsg` switch, guarded by `m.mode == ModeIndex`
- [x] 4.2 Inside the `s` handler: save current item (`kind`, `idx`, `reqIdx`), toggle `m.specSortBySuffix`, call `m.buildIndexItems()`, restore cursor by searching `m.indexItems`, call `m.refreshIndexViewport()`

## 5. Help bar

- [x] 5.1 Update `renderHelpBar()` for `ModeIndex`: when `m.specSortBySuffix` is false show `s: sort by suffix`; when true show `s: sort by name`
