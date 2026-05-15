## Context

BubbleTea manages the terminal in raw mode throughout the entire TUI lifetime. To hand control over to an external editor, the BubbleTea program must be suspended in a coordinated way: restore the terminal to normal mode, launch the editor, wait for it to finish, and retake control. BubbleTea has native support for this via `tea.ExecProcess`, which encapsulates exactly that cycle.

The file to edit depends on the active tab: `proposal.md`, `design.md`, `tasks.md`, or the `specs/` directory (for specs, the first `spec.md` found is opened, or the directory, depending on what the editor supports).

## Goals / Non-Goals

**Goals:**
- Pressing `e` on any tab with an available artifact opens `$EDITOR` with the corresponding file
- When the editor is closed, the TUI resumes and forces an immediate reload of the edited artifact
- Fallback to `vi` if `$EDITOR` is not defined
- The help bar shows `e: edit` when an artifact is available in the active tab

**Non-Goals:**
- Support for non-terminal editors (GUI editors); if the user has one, the behavior is undefined but does not crash
- Editing multiple files at once
- Detection of whether the file was actually modified (always reloads)

## Decisions

**D1: `tea.ExecProcess` instead of direct `os/exec`**

`tea.ExecProcess(cmd, callback)` is the official BubbleTea mechanism for handing the terminal over to a subprocess. It handles the suspension/resumption of raw mode, the redraw after returning, and delivers a `tea.Msg` with the process exit error. Using `os/exec` directly without this wrapper would corrupt the terminal state.

Discarded alternative: `tea.Suspend`/`tea.Resume` — these are for system signals (SIGTSTP), not editor subprocesses.

**D2: Resolve the file path by active tab**

Direct mapping: `TabProposal` → `proposal.md`, `TabDesign` → `design.md`, `TabTasks` → `tasks.md`, `TabSpecs` → the first `spec.md` inside `specs/`. If there is none (specs empty but `Present` true), `specs/` is opened as a directory — editors like vim handle this as netrw.

**D3: Forced immediate reload after returning from the editor**

The `tea.ExecProcess` callback returns a `tea.Msg`. An `editorReturnMsg` type is defined that, when processed in `Update`, calls `ReloadChange` on the current change and updates the state (equivalent to the tick handler but executed at the exact moment of return). This guarantees that the updated content is visible immediately, without waiting for the next 500 ms tick.

**D4: Fallback `vi`**

`os.Getenv("EDITOR")` returns `""` if not defined. In that case `"vi"` is used since it is the lowest common denominator available on any Unix system. No additional configuration is offered: users who want another editor must define `$EDITOR`.

## Risks / Trade-offs

[Risk: GUI editor (VS Code, etc.) in `$EDITOR`] → The TUI yields the terminal but the editor opens in a separate window. On return, the TUI resumes correctly. The user might not notice that the editor opened in another window. No active mitigation; it is the user's responsibility to have `$EDITOR` pointing to a terminal editor.

[Risk: Editor exits with error (e.g., file locked)] → The callback receives the error and can display it in `m.errMsg`. The TUI continues functioning.

[Risk: Specs with multiple subdirectories] → If there are several subdirectories in `specs/`, only the first is opened. This is a rare case in the current flow (one capability per change) and simplifies the implementation.

## Migration Plan

No migration. The change adds a new keybinding (`e`) that does not collide with any existing one. No data or file format changes.

## Open Questions

None.
