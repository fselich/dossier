# view-background Specification

## Purpose
Defines how views render a configurable solid background color that fills the entire terminal viewport, covering all whitespace areas including gaps between styled segments, padding areas, and empty vertical space below the box frame.

## Requirements

### Requirement: Background color fills entire viewport
The system SHALL render a solid background color across the full terminal window (width × height) when a background color is configured for the active view. The background SHALL cover all areas visible to the user: box borders, content text, inter-element whitespace, and any empty area below the box frame. No area of the terminal-default background SHALL be visible within the rendered view.

#### Scenario: Solid fill with content shorter than viewport
- **WHEN** the configured background color is set to ANSI color 234 (dark grey) and the viewport content is shorter than the terminal height
- **THEN** the area below the box frame renders as full-width background-colored lines down to the bottom of the terminal

#### Scenario: No gaps between styled segments on same line
- **WHEN** a line contains multiple styled segments (e.g., a border character, then plain spaces, then a section header, then padding, then another border character) and a background color is configured
- **THEN** every character cell on that line renders with the configured background color, including the spaces between styled segments

#### Scenario: Explicit element background preserved
- **WHEN** an element has its own explicit background color (e.g., the cursor-highlighted index item with dark blue background) and a view background color is configured
- **THEN** the element retains its own background color, and areas after the element revert to the view background color

### Requirement: Terminal default fallback
When the view background color is nil or unset (either because the theme defines no background or no theme is loaded), the system SHALL render exactly as it does today: all areas not covered by styled segments display the terminal emulator's default background. No post-processing or wrapping SHALL occur.

#### Scenario: No background configured
- **WHEN** the view background color is nil or unset
- **THEN** the rendering pipeline passes through the view string unchanged, producing the same visual output as the current implementation

### Requirement: Background is configurable per view
The system SHALL support a per-view background color configuration via a Theme struct on the Model. The background color SHALL be a `color.Color` value. Setting it to `nil` SHALL disable the background (terminal default fallback). The background color SHALL be populated from the active built-in theme selected via `--theme` flag.

#### Scenario: Different views have different backgrounds
- **WHEN** a future theme configures the index view with one background color and the normal view with another
- **THEN** switching between modes displays the corresponding background color for each view

#### Scenario: Built-in theme provides background color
- **WHEN** the user selects a built-in theme (e.g., `--theme light`)
- **THEN** the view renders with the theme's defined background color

### Requirement: Rendering pipeline is reusable
The background fill pipeline SHALL be implemented as a single shared method on the Model, callable from any view-rendering function (`viewIndex`, `viewConfig`, `View`, etc.). Each view SHALL only need to call the pipeline with its rendered content string.

#### Scenario: Multiple views use the same pipeline
- **WHEN** `viewIndex()`, `viewConfig()`, and the main `View()` method all call the shared background render method
- **THEN** each view renders with the configured background color using the same logic, without code duplication
