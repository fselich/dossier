## MODIFIED Requirements

### Requirement: Progress bar complete style
The TUI SHALL render all progress bars using a single unified function `renderProgressBar(done, total, width int) string`. The function SHALL use two distinct colors: cyan (`"6"`) for in-progress bars and green (`"2"`) for bars that have reached 100% completion (`done == total`). All render sites SHALL call this function directly.

#### Scenario: Bar is fully complete
- **WHEN** a progress bar's `done` count equals its `total` count (and `total > 0`)
- **THEN** the filled portion SHALL be rendered using `progressCompleteStyle` (green, color `"2"`)

#### Scenario: Bar is partially complete
- **WHEN** a progress bar's `done` count is less than its `total` count
- **THEN** the filled portion SHALL be rendered using `progressDoneStyle` (cyan, color `"6"`)

#### Scenario: Unified function used by all render sites
- **WHEN** any of the three progress bar render sites (global tab bar, change index, per-section tasks view) needs a progress bar
- **THEN** it SHALL call `renderProgressBar(done, total, width)` with the appropriate parameters
- **THEN** each SHALL independently apply `progressCompleteStyle` to its filled blocks when done equals total
