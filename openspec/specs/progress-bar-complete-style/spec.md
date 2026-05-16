## Purpose

Defines the visual styling of progress bars in the TUI, distinguishing between in-progress and fully complete bars using distinct colors.

## Requirements

### Requirement: Progress bar complete style
The TUI SHALL render all progress bars using two distinct colors: cyan (`"6"`) for in-progress bars and green (`"2"`) for bars that have reached 100% completion (`done == total`).

#### Scenario: Bar is fully complete
- **WHEN** a progress bar's `done` count equals its `total` count (and `total > 0`)
- **THEN** the filled portion SHALL be rendered using `progressCompleteStyle` (green, color `"2"`)

#### Scenario: Bar is partially complete
- **WHEN** a progress bar's `done` count is less than its `total` count
- **THEN** the filled portion SHALL be rendered using `progressDoneStyle` (cyan, color `"6"`)

#### Scenario: Bar applies to all render sites
- **WHEN** any of the three progress bar render sites (global tab bar, change index, per-section tasks view) reaches 100%
- **THEN** each SHALL independently apply `progressCompleteStyle` to its filled blocks
