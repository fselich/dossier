## 1. Reduce polling interval

- [x] 1.1 Change `2*time.Second` to `500*time.Millisecond` in the `tea.Tick` call in `model.go`

## 2. Artifact presence detection

- [x] 2.1 In the tick handler in `model.go`, compare the `Present` flag of each artifact between the previous and new state; if any artifact transitions from absent to present, update `m.project` and recalculate the enabled tabs
- [x] 2.2 Verify that the enabled-tabs calculation (`tabEnabled`) is invoked when updating `m.project` in the tick, not only during initialisation

## 3. Change list reload from empty state

- [x] 3.1 At the start of the tick handler, if `len(m.project.Changes) == 0`, call `loader.LoadProject` and adopt the result if it returns at least one change

## 4. Immediate toggle update

- [x] 4.1 In `handleTaskToggle` in `model.go`, after successfully writing `tasks.md`, update `m.taskItems[cursor].Done` in memory immediately (without waiting for the next tick)
