## 1. Data model

- [x] 1.1 Add `type NamedSpec struct { Name, Content string }` to `internal/openspec/loader.go`
- [x] 1.2 Add field `SpecFiles []NamedSpec` to the `Change` struct in `loader.go`
- [x] 1.3 Update `loadSpecs` so that in addition to building the concatenated `Artifact`, it also populates and returns `[]NamedSpec` — update `Load` and `ReloadChange` to assign `ch.SpecFiles`

## 2. UI model state

- [x] 2.1 Add field `specIdx int` to `Model` in `model.go`
- [x] 2.2 Reset `specIdx = 0` in the change-switch handlers (`h`/`l`) and when entering TabSpecs from another tab

## 3. Key 3 with dual behaviour

- [x] 3.1 Update the `"3"` case in `Update`: if `m.tab != TabSpecs`, switch to TabSpecs (current behaviour); if `m.tab == TabSpecs`, do `m.specIdx = (m.specIdx + 1) % len(specFiles)` and invalidate the TabSpecs cache before calling `loadViewport()`

## 4. Rendering the selected spec

- [x] 4.1 Update `loadViewport` for the `TabSpecs` case: use `ch.SpecFiles[m.specIdx].Content` instead of `ch.Specs.Content`; guard if `len(ch.SpecFiles) == 0`
- [x] 4.2 Update the tick handler: when `fresh.Specs` differs from `ch.Specs`, also update `m.project.Changes[m.changeIdx].SpecFiles` with the new `NamedSpec` values

## 5. Navigation sub-bar

- [x] 5.1 Add method `renderSpecSubnav() string` in `model.go` that generates the chip row using `tabActiveStyle` for the active spec and `tabInactiveStyle` for the others
- [x] 5.2 Insert the sub-bar in `View()` between the tab bar and the lower tab bar separator, only when `m.tab == TabSpecs && len(specFiles) > 0`

## 6. Content height

- [x] 6.1 Update `contentHeight()` to return `m.height - 8` when `m.tab == TabSpecs && len(specFiles) > 0`, and `m.height - 7` in all other cases
- [x] 6.2 Ensure that `m.vp.Height` is recalculated with `m.contentHeight()` when switching tabs (already happens in the numeric key handlers — verify that the `3` handler also does it)
