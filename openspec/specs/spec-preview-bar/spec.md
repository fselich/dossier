# spec-preview-bar Specification

## Purpose
Provides a persistent preview bar in the index view that displays the purpose text of the currently selected specification, allowing users to browse specs without leaving the index.

## Requirements

### Requirement: Purpose text extraction
The system SHALL provide a function `ExtractPurpose(content string) string` that extracts the plain text from between the `## Purpose` heading and the next `##` heading (or end of file) in a spec's raw markdown content. Markdown formatting (bold, italic, links) SHALL be stripped from the extracted text. If no `## Purpose` heading is found, the function SHALL return an empty string.

#### Scenario: Purpose section present
- **WHEN** the content contains `## Purpose\nDefines the layout...\n\n## Requirements`
- **THEN** `ExtractPurpose` returns `"Defines the layout..."`

#### Scenario: No Purpose heading
- **WHEN** the content has no `## Purpose` heading
- **THEN** `ExtractPurpose` returns `""`

#### Scenario: Purpose at end of file
- **WHEN** `## Purpose` is the last heading before EOF
- **THEN** `ExtractPurpose` returns all text from `## Purpose` to EOF (excluding the heading)

### Requirement: Text truncation with ellipsis
The preview bar text SHALL be truncated with `…` if the combined string (spec name + ` ┊ ` + purpose) exceeds the available content width (`m.width - 2`). The truncation SHALL preserve the spec name and the ` ┊ ` separator, truncating only the purpose portion.

#### Scenario: Short purpose fits
- **WHEN** the purpose text is short enough to fit in the available width
- **THEN** the full text is shown without truncation

#### Scenario: Long purpose truncated
- **WHEN** the combined text exceeds the available width
- **THEN** the purpose portion is shortened and `…` is appended at the end

### Requirement: Empty bar when no spec selected
When no spec is selected (cursor on active/archived item), the preview bar SHALL be hidden entirely (no gap in the chrome).

#### Scenario: Empty bar hidden
- **WHEN** the cursor is on an active change
- **THEN** the spec preview bar is not rendered, leaving no gap between viewport and helpbar
