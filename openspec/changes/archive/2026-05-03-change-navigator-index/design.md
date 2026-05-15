## Context

The TUI currently manages archived selection via a modal overlay (`ModeArchivePicker`) that overlaps the main content visually. This causes two problems: the modal appears visually disconnected from the TUI chrome (its own borders, centered positioning with a dark background), and it conceptually separates active changes from archived ones, forcing the user into two distinct navigation flows.

The relevant code lives in `internal/ui/model.go`. The current modes are `ModeNormal`, `ModeArchivePicker`, and `ModeViewingArchive`. The modal is rendered in `renderArchivePicker()` using `lipgloss.Place` with an overridden background.

## Goals / Non-Goals

**Goals:**
- Replace `ModeArchivePicker` with `ModeIndex`: a full-screen view that lists active changes (with progress) and archived changes (with date)
- Simplify the navigation model: `a`/`Esc` from a change view → index; `Esc` from index → quit; `q` from anywhere → quit
- Keep `ModeViewingArchive` with no changes to its internal behaviour

**Non-Goals:**
- Changing the behaviour of an archived change view (same tabs, same interaction)
- Adding search or filtering to the index
- Showing the index automatically on startup when there are no active changes

## Decisions

### D1: ModeIndex is full-screen, not a modal

**Decision**: `ModeIndex` occupies the full screen with the same TUI chrome (borders, header, helpbar), instead of an overlay.

**Alternatives considered**: Keep the overlay but redesign it visually. Discarded because the underlying problem is that the modal interrupts the visual flow — a redesign does not solve that.

**Rationale**: Full-screen integrates the index as just another "page" of the TUI, consistent with the existing chrome.

---

### D2: Flat list with visual section separation

**Decision**: The cursor navigates through a flat list (active + archived in order). A section separator line visually divides the two groups but is not a selectable item.

**Alternatives considered**: Cursor separated by section (Tab to jump between sections). Discarded for unnecessary complexity.

**Rationale**: Simple `j`/`k` works well. The user knows active changes are at the top and archived at the bottom.

---

### D3: Esc from ModeNormal goes to the index (new behaviour)

**Decision**: `Esc` in `ModeNormal` now opens `ModeIndex`. Previously it was a no-op.

**Rationale**: Symmetry with the inverse flow (`Esc` from the index quits). The user can enter and leave the index with `Esc` regardless of where they are.

---

### D4: Index cursor does not persist between openings

**Decision**: Each time `ModeIndex` is opened, the cursor is positioned on the first active item (or the first archived item if there are no active ones).

**Rationale**: Simplicity. The index is a navigator, not a persistent view.

## Risks / Trade-offs

- [Esc in ModeNormal is new] → Users accustomed to the old behaviour may be surprised. Mitigation: the helpbar in ModeNormal will include `Esc: index`.
- [ModeArchivePicker removed] → Dead code to clean up (`renderArchivePicker`, `modal*` styles). No functional risk.
- [Loading archived changes when opening the index] → With 10 current archived changes this is imperceptible. With a large N it could be noticeable. Acceptable for now.
