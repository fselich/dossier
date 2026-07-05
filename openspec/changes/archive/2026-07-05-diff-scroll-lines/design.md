## Context

The diff view currently pre-renders the entire diff into a single string stored in `gitState.DiffContent`. Once rendered with chroma highlighting, the string contains ANSI codes that make post-hoc truncation fragile (splitting ANSI sequences corrupts display).

The fix: store parsed `DiffLine` structs instead of a pre-rendered string. Re-render with a horizontal scroll offset each time the view updates. Truncate the raw content *before* chroma tokenization, so ANSI codes are never split.

## Goals / Non-Goals

**Goals:**
- Horizontal scroll in diff view via `→`/`←` (10 runes per step).
- Line numbers (`OldNum`/`NewNum`) displayed as 4-char columns.
- Scroll offset trims raw content before highlighting — no ANSI corruption.
- Regenerate content on every scroll step (cached lexer keeps it fast).
- `d`/`Esc` reset scroll to 0 when returning to file list.

**Non-Goals:**
- Word wrapping (replaced by horizontal scroll approach).
- Synchronized vertical + horizontal scroll (vertical is already handled by viewport).
- Scroll beyond content bounds (clamped to 0).

## Decisions

### Store parsed diff, not rendered string

`gitState` changes from:

```
DiffContent string   // pre-rendered ANSI string
```

to:

```
DiffLines  []DiffLine   // structured diff data
DiffFile   string        // which file is being viewed
ScrollX    int           // horizontal scroll offset
```

On every render, `renderDiffContent` calls `renderDiff(gitState.DiffLines, filename, m.width, gitState.ScrollX)` to regenerate the ANSI content.

This adds a small cost per scroll step (tokenizing all diff lines), but:
- Lexer is cached per file extension
- Each line tokenization is sub-millisecond
- Diffs are typically small (< 500 lines)

### Truncate before highlighting

For each diff line, the rendering process is:

1. Take raw `Content` string
2. Skip `ScrollX` runes (or 0 if content is shorter)
3. Tokenize the remaining portion with chroma
4. Apply foreground (chroma) + background (diff type) styles
5. Prepend line number + indicator

This ensures ANSI codes are never in the truncation path.

### Line number tracking

`DiffLine` gets `OldNum` and `NewNum` fields (like differ). The `parseDiffLine` function maintains these via pointers passed to it. Line numbers are reset per hunk header (`@@ -13,6 +14,8 @@`):

```
@@ -13,6 +14,8 @@ func main() {
```

- `oldNum` set to 13, `newNum` set to 14
- Context lines increment both
- Added lines increment only newNum
- Removed lines increment only oldNum
- Hunk headers don't increment either

Line numbers are displayed as 4-character right-aligned columns:

```
       1  package main
  +    2  import "fmt"
```

### Scroll keys

In diff view, `j`/`k` scroll vertically (viewport), new keys scroll horizontally:

| Key | Action |
|-----|--------|
| `→` / `Shift+L` | Scroll right +10 |
| `←` / `Shift+H` | Scroll left -10 |
| `d` / `Esc` | Reset scroll to 0, return to list |

Regular `h`/`l` still navigate between changes even in diff view (existing behavior).

## Architecture

```
internal/ui/model.go:
  gitState:
    - DiffContent string       → REMOVED
    + DiffLines  []DiffLine    → stored parsed diff lines
    + ScrollX    int           → horizontal offset in runes

internal/ui/gitdiff.go:
  + DiffLine.{OldNum, NewNum}         → line number tracking
  + parseDiffLine(...) with num tracking → oldNum/newNum increment logic
  + renderDiff(lines, filename, width, scrollX) → scroll-aware rendering
  + highlightLine → unchanged
  + toggleGitDiff → sets ScrollX=0, populates DiffLines

internal/ui/git.go:
  + renderDiffContent → calls renderDiff with current scrollX, viewport GotoTop() on scroll
  + pollGitStatus → clears DiffLines, ScrollX

internal/ui/viewer.go:
  + "right" / "L" (shift+l) → ScrollX += 10 in diff view
  + "left" / "H" (shift+h) → ScrollX -= 10 in diff view
  + "d" / "esc" → reset ScrollX=0

internal/ui/view.go:
  + renderHelpBar → show scroll keys in diff mode
```

## Migration Plan

1. Add `OldNum`, `NewNum` to `DiffLine` struct.
2. Update `parseDiffLine` to track line numbers.
3. Replace `DiffContent` with `DiffLines` + `ScrollX` in `gitState`.
4. Update `renderDiff` to accept width + scrollX, truncate raw content before highlighting.
5. Add scroll key handlers in viewer.
6. Update help bar.
7. `make build && make lint && make test`.
