## Context

After `harden-git-layer`, `internal/git` parses status exactly (`-z`) and runs every subprocess with a 2s timeout via a shared `runGit` helper. The git tab (`internal/ui/git.go`, `viewer.go`) renders `gitState.Files` with an index-based cursor that skips `IsDeleted` entries, and polls status every 500ms. Decisions below were settled in exploration with the user.

## Goals / Non-Goals

**Goals:**
- One key (`s`) toggles staged state of the file under the cursor, list view only.
- Deleted files become selectable so their deletion can be staged.
- Mutations reflect in the list immediately, cursor stays on the same file (by path).
- Failures surface as an ephemeral help-bar message.

**Non-Goals:**
- Commit flow, hunk staging, staging `openspec/` files.
- Keeping the diff view open across an XY change (existing close-on-change behavior stays).
- A general notification system.

## Decisions

### D1: Toggle direction — worktree side wins

For file status `XY`: if `Y != ' '` (unstaged changes exist, including `??` untracked) → **stage** (`git add -- <path>`); else (`X != ' '`, fully staged) → **unstage**. This resolves the ambiguous `MM` case toward staging, which matches the common "review then stage everything" loop. `git add` on a ` D` entry stages the deletion — no special case needed.

### D2: Unstage command — `git reset -q HEAD --` with `git rm --cached` fallback

`git restore --staged` and `git reset HEAD` both fail in a repo with no commits (unresolvable HEAD). Simplest robust approach: run `git reset -q HEAD -- <paths>`; if it fails, fall back to `git rm --cached -q -- <paths>` (correct for the no-HEAD case, where every staged file is newly added). No upfront HEAD detection, one extra call only on the rare path. For staged renames (`OldPath != ""`), pass both old and new paths; the entry may split into ` D` + `??` afterwards, which the refresh handles naturally.

### D3: Cursor lands on deleted files

Remove the `IsDeleted` skip from `moveGitCursorDown/Up` and `clampGitCursor` (they become simple wrapping/clamping). `Enter`/`e`/`d` keep their existing no-op guard on deleted files; `[`/`]` in diff view keep skipping deleted entries (a diff view cannot show them). This also fixes the "all files deleted" dead-end (PROPUESTAS.md 1.7).

### D4: Synchronous mutation + immediate refresh, cursor preserved by path

`s` runs the mutation and re-runs `git.Status` in the update loop (each call bounded by the 2s timeout from `harden-git-layer`). Rationale: `git add` on a single path is fast; an async Cmd would add message plumbing for little gain. After refresh, the cursor is restored by looking up the previous `Path` in the new list (mirroring the `FindCursorByText` pattern for tasks); if the path is gone (e.g. rename split), fall back to clamping the old index.

### D5: Error feedback — ephemeral help-bar message

New `gitState.ErrMsg string` rendered in the help bar (replacing the key hints while present). Set on mutation failure with a trimmed one-line git error; cleared on the next successful `s`, on any git-tab key press, or by the next poll tick that observes a status change. No timers — reuse existing redraw points.

### D6: New API in `internal/git`

`Stage(root string, paths ...string) error` and `Unstage(root string, paths ...string) error`, built on `runGit`. UI decides direction (D1) from the `FileStatus` it already has; the git package stays policy-free.

## Risks / Trade-offs

- [Synchronous mutation can block up to 2s on a pathological repo] → bounded by timeout; acceptable for an explicit user action, unlike the background poll.
- [Toggling `MM` always stages; user may want to unstage] → documented behavior; unstage becomes available on the next press once `Y == ' '`.
- [Cursor on deleted file changes muscle memory around `Enter`/`d`] → no-op guards keep it safe; help bar shows `s` as the available action.
- [Rename unstage splits one entry into two] → cursor-by-path falls back to index clamp; list is consistent after refresh.
- [Help-bar message could be overwritten by unrelated redraw] → message lives in model state, not in a transient render; cleared only at defined points (D5).
