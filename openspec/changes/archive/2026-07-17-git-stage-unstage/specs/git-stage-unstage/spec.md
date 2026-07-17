## ADDED Requirements

### Requirement: Stage/unstage toggle with `s`

The TUI SHALL toggle the staged state of the file under the cursor when the user presses `s` in the git tab file list. If the file has unstaged changes (worktree status `Y != ' '`, including untracked `??`), the file SHALL be staged via `git add -- <path>`. Otherwise (fully staged), the file SHALL be unstaged. The `s` key SHALL do nothing inside the diff view and when the working tree is clean.

#### Scenario: Stage a modified file
- **WHEN** the cursor is on a file with status ` M` and the user presses `s`
- **THEN** the file is staged and its status becomes `M·` in the list

#### Scenario: Stage an untracked file
- **WHEN** the cursor is on a file with status `??` and the user presses `s`
- **THEN** the file is staged and appears as added (`A·`)

#### Scenario: Unstage a fully staged file
- **WHEN** the cursor is on a file with status `M·` (no worktree changes) and the user presses `s`
- **THEN** the file is unstaged and its status becomes ` M`

#### Scenario: Mixed state stages the worktree side
- **WHEN** the cursor is on a file with status `MM` and the user presses `s`
- **THEN** the worktree changes are staged and the status becomes `M·`

#### Scenario: Stage a deleted file
- **WHEN** the cursor is on a file with status ` D` and the user presses `s`
- **THEN** the deletion is staged and the status becomes `D·`

#### Scenario: Unstage a staged rename
- **WHEN** the cursor is on a renamed entry with status `R·` and the user presses `s`
- **THEN** both the old and new paths are unstaged and the list reflects the resulting statuses

#### Scenario: s inactive in diff view
- **WHEN** the diff view is showing and the user presses `s`
- **THEN** nothing happens

#### Scenario: s inactive on clean working tree
- **WHEN** the working tree is clean and the user presses `s`
- **THEN** nothing happens

### Requirement: Unstage works in a repository without commits

Unstaging SHALL work in a repository with no commits (unresolvable `HEAD`). The TUI SHALL first attempt `git reset -q HEAD -- <paths>` and, on failure, fall back to `git rm --cached -q -- <paths>`.

#### Scenario: Unstage in a fresh repository
- **WHEN** the repository has no commits, a file is staged (`A·`), and the user presses `s`
- **THEN** the file is unstaged and appears as untracked (`??`)

### Requirement: Immediate refresh with cursor preserved by path

After a stage or unstage action, the TUI SHALL re-run `git status` and refresh the file list immediately, without waiting for the next 500ms poll tick. The cursor SHALL be restored to the entry whose path matches the file that was acted on; if that path no longer exists in the list, the cursor SHALL be clamped to the nearest valid index.

#### Scenario: Status updates without tick delay
- **WHEN** the user presses `s` on a file
- **THEN** the list shows the new XY status immediately

#### Scenario: Cursor stays on the same file
- **WHEN** the user stages a file in the middle of the list
- **THEN** after the refresh the cursor remains on that file's path

#### Scenario: Cursor falls back when the path disappears
- **WHEN** unstaging a rename splits the entry into two different paths
- **THEN** the cursor is clamped to a valid index and no crash occurs

### Requirement: Mutation errors shown in help bar

When a stage or unstage command fails, the TUI SHALL display a one-line error message in the help bar. The message SHALL be cleared on the next successful stage/unstage action, on any subsequent key press in the git tab, or when polling observes a status change. Failures SHALL NOT crash the TUI or corrupt the file list.

#### Scenario: git add failure surfaces a message
- **WHEN** `git add` fails (e.g. another process holds `index.lock`)
- **THEN** the help bar shows an error message and the file list keeps its previous state

#### Scenario: Message clears on next key press
- **WHEN** an error message is visible and the user presses any key in the git tab
- **THEN** the message is removed and the normal key hints return

### Requirement: Stage/unstage key hint in help bar

The help bar in the git tab file list SHALL include a hint for the `s` key (e.g. `s stage/unstage`).

#### Scenario: Hint visible in file list
- **WHEN** the git tab shows the file list with at least one file
- **THEN** the help bar includes the `s` hint
