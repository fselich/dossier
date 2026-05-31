## MODIFIED Requirements

### Requirement: Mouse event capture

*(No requirement change. Internal restructuring only — `loadViewport()` is split into mode-specific methods. Behavior remains identical.)*

### Requirement: Wheel scrolling

*(No requirement change. Wheel handlers call `loadViewport()` unchanged.)*

### Requirement: Tab selection via left-click

*(No requirement change. Click handlers call `loadViewport()` unchanged.)*

### Requirement: Index item selection via left-click

*(No requirement change. Click handlers call `loadViewport()` unchanged.)*

### Requirement: Header click navigates to index

*(No requirement change. Header click triggers `enterIndex()` which calls `refreshIndexViewport()`, not `loadViewport()`.)*
