## Context

The TUI currently switches tabs with direct number keys (`1`-`4`) or mouse clicks. Each key maps to a specific tab; there is no concept of "next" or "previous" tab. This is efficient for known layouts but doesn't support a review workflow where the user wants to scan through all artifacts of a change sequentially.

Bubble Tea with `tea.WithAltScreen()` (already used) can receive `KeyTab` and `KeyShiftTab` events as `"tab"` and `"shift+tab"` respectively.

## Goals / Non-Goals

**Goals:**
- `Tab` advances to the next available tab, skipping disabled ones, wrapping around
- `Shift+Tab` moves to the previous available tab
- Works in `ModeNormal` and `ModeViewingArchive` only
- Help bar reflects the new shortcut

**Non-Goals:**
- No change to `ModeViewingConfig` (Tab must remain for indentation)
- No cycling through specs subnav items
- No focus system or tabindex model
- `1`-`4` and mouse clicks remain unchanged

## Decisions

### Cycle algorithm

A method `nextAvailableTab(Tab, delta int) Tab` on the Model that:

1. Starts from the current tab + delta (clamped to 0..tabCount-1 with modulo)
2. Iterates through all tab positions to find the first available one
3. Returns the found tab (or the original if none are available — edge case guard)

```
        ┌──────────┐
Start ◄─│ current  │──► delta = +1 (Tab)
        │ + delta  │──► delta = -1 (Shift+Tab)
        └────┬─────┘
             │
             ▼
    ┌────────────────┐    yes   ┌─────────────┐
    │ tabAvailable() │─────────▶│ return tab  │
    │                │          └─────────────┘
    └───────┬────────┘
            │ no
            ▼
    ┌────────────────┐
    │ advance by     │
    │ delta (wrap)   │──────► loop (max tabCount iterations)
    └────────────────┘
```

Infinite loop guard: stop after `tabCount` iterations; return unchanged.

### Key matching approach

Add two cases to the existing `switch msg.String()` in `update.go`. Follow the existing pattern (string comparison, no `key.Binding` types used elsewhere). Bubble Tea normalizes these:

- Tab → `"tab"`
- Shift+Tab → `"shift+tab"` (terminal-dependent; silently ignored if unsupported)

### Mode guard

Only activate in `ModeNormal` and `ModeViewingArchive`. The current key dispatch flows through a `switch m.mode` block (line ~96 in update.go) so placement is natural — add to those two mode cases.

### Help bar update

The relevant help strings in `view.go` are:

- Non-tasks tabs: `"h/l: change  1-4: artifact  j/k: scroll  e: edit  Esc: index  q: quit"`
- Tasks tab: `"h/l: change  1-4: artifact  j/k: navigate  Space: toggle  e: edit  Esc: index  q: quit"`

Change to:
- `"h/l: change  1-4/Tab: artifact  j/k: scroll  e: edit  Esc: index  q: quit"`
- `"h/l: change  1-4/Tab: artifact  j/k: navigate  Space: toggle  e: edit  Esc: index  q: quit"`

## Risks / Trade-offs

| Risk | Mitigation |
|------|-----------|
| `Shift+Tab` not supported in some terminals | Only Tab is documented in help bar; Shift+Tab is a silent bonus |
| Infinite loop if all tabs disabled | Guard with `tabCount` iteration cap |
| Tab key might be intercepted by terminal multiplexer (tmux) | Low risk with `tea.WithAltScreen()`; no action needed beyond documenting known limitation |
