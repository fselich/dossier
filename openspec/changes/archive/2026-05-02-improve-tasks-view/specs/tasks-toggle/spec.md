## ADDED Requirements

### Requirement: Inline markdown rendering in task items
The TUI SHALL convert inline markdown marks present in each task's text to ANSI styles before rendering the item with lipgloss. The supported patterns are `` `code` `` (backtick) and `**bold**` (double asterisk).

#### Scenario: Task with a code snippet
- **WHEN** a task item's text contains `` `func main()` ``
- **THEN** the snippet is displayed with the visual code style (distinct background or colour) in the TUI

#### Scenario: Task with bold text
- **WHEN** a task item's text contains `**importante**`
- **THEN** the word is displayed in bold in the TUI

#### Scenario: Multiple snippets in the same task
- **WHEN** an item's text contains several `` `code` `` or `**bold**` fragments separated from each other
- **THEN** each fragment is rendered with its corresponding style independently

#### Scenario: Task without inline markdown
- **WHEN** an item's text contains no backticks or double asterisks
- **THEN** the text is displayed unchanged, with no visual artefacts
