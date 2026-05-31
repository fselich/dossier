## ADDED Requirements

### Requirement: ExtractRequirement extracts a single requirement block
The loader SHALL expose `ExtractRequirement(content string, name string) string` which extracts the content of a single requirement block from a spec's full markdown content. The function SHALL return the text from `### Requirement: <name>` to the next `### Requirement:` (or end of content). If the requirement name is not found, it SHALL return an empty string.

#### Scenario: Requirement found in content
- **WHEN** `ExtractRequirement` is called with valid content containing the named requirement
- **THEN** the function returns the requirement block text excluding the header line

#### Scenario: Requirement not found
- **WHEN** `ExtractRequirement` is called with a name not present in the content
- **THEN** the function returns an empty string

### Requirement: ConfigToMarkdown formats project config
The loader SHALL expose `ConfigToMarkdown(cfg ProjectConfig) string` which formats a `ProjectConfig` struct into a human-readable markdown string for display.

#### Scenario: Config with all fields populated
- **WHEN** `ConfigToMarkdown` is called with a fully populated `ProjectConfig`
- **THEN** the function returns a markdown-formatted string representing the config
