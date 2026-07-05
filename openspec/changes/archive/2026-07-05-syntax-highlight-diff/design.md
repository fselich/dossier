## Context

The diff view currently runs `git diff --color=always` and displays the raw ANSI output. This works but only colors diff markers — the code inside the diff is plain monochrome text.

The `differ` project demonstrates a technique that gives full syntax highlighting inside diffs:

1. Parse the raw (uncolored) `git diff` into structured `DiffLine` records
2. For each code line, strip the `+`/`-`/` ` prefix
3. Tokenize the line with chroma's language lexer
4. Apply chroma foreground (syntax colors) + line-type background tint

This is efficient because chroma tokenizes per-line (each line is small), and the lexer is cached per file extension.

Chroma is already available as a direct dependency. The `lexers.Match(filename)` function auto-detects the language from the filename.

## Goals / Non-Goals

**Goals:**
- Show syntax-highlighted code in diff view for tracked files.
- Color-code diff indicators: green for `+`, red for `-`, plain for context.
- Hunk headers (`@@ ... @@`) rendered in cyan.
- Background tint: subtle green for added lines, subtle red for removed lines.
- Lexer cache per file extension (via `sync.Map`).
- Fall back to plain text if chroma lexer is unavailable.

**Non-Goals:**
- Line numbers (can be added later).
- Side-by-side or split diff views.
- Staged/unstaged indication within the diff (the XY code is already visible in the file list).
- New mode or UI flow changes — same toggle behavior with `d`.

## Decisions

### Per-line chroma tokenization

Chroma's `Tokenise(nil, content)` works on arbitrary strings, not just whole files. For a single line of code, it correctly identifies keywords, strings, comments, etc. The result is a stream of tokens with type information.

Each token's foreground color is extracted from the chosen chroma style (monokai). The background color is applied based on the diff line type:

```
highlightLine(line, filename, bgColor):
    lexer = getLexer(filename)     // cached by extension
    for token in lexer.Tokenise(line):
        fg = chromaStyle[token.Type].Colour  // chroma's color for this token type
        style = Foreground(fg) + Background(bgColor)
        emit styled(token.Value)
```

### Background colors for diff types

| Line type | Background | Indicator |
|-----------|-----------|-----------|
| Added (`+`) | `#1a3a1a` dark green | `+` in green |
| Removed (`-`) | `#3a1a1a` dark red | `-` in red |
| Context (` `) | None | ` ` in gray |
| Hunk header (`@@`) | None | Header in cyan |

Colors chosen to be subtle enough not to overpower the syntax highlighting, but visible enough to distinguish line types.

### Diff parser replaces git --color=always

The `computeDiff` function is refactored:
1. Run `git diff --no-color` (raw output) instead of `git diff --color=always`
2. Parse the raw output into `[]DiffLine`
3. For tracked files: `parseDiff` + `renderDiff`
4. For untracked files: keep existing `highlightUntrackedFile`

The diff parser strips git headers (`diff --git`, `index`, `---`, `+++`, `new file`, etc.) and keeps only code lines and hunk headers. The file name is already displayed as a header by `renderDiffContent`.

### Lexer cache

A `sync.Map` keyed by file extension avoids calling `lexers.Match()` for every line:

```go
var lexerCache sync.Map

func getLexer(filename string) chroma.Lexer {
    ext := filepath.Ext(filename)
    if ext == "" {
        ext = filepath.Base(filename)
    }
    if cached, ok := lexerCache.Load(ext); ok {
        return cached.(chroma.Lexer)
    }
    lexer := lexers.Match(filename)
    if lexer == nil {
        lexer = lexers.Fallback
    }
    lexer = chroma.Coalesce(lexer)
    lexerCache.Store(ext, lexer)
    return lexer
}
```

### Chroma style initialization

The monokai style is initialized lazily via `sync.Once` to avoid repeated style lookups.

## Architecture

```
internal/ui/gitdiff.go (REFACTORED):
  - parseDiff(raw string) ParsedDiff       — parse unified diff into []DiffLine
  - parseDiffLine(line) *DiffLine          — classify a single line
  - renderDiff(parsed, filename, width) string  — render with chroma highlighting
  - highlightLine(content, filename, bgColor) string — per-line chroma tokenization
  - getLexer(filename) chroma.Lexer        — cached lexer lookup
  - computeDiff(m *Model) string           — updated to use parseDiff + renderDiff
  - highlightUntrackedFile(root, rel) string — unchanged
  - toggleGitDiff(m *Model)                — unchanged

No changes to other files.
```

## Risks / Trade-offs

- **Per-line tokenization accuracy**: Chroma tokenizing individual lines may miss multi-line constructs (e.g., multi-line strings, block comments). In practice, this is rare within diff hunks and the result is still significantly better than no highlighting at all. The `chroma.Coalesce(lexer)` call improves handling.
- **Performance**: Tokenizing every line of a large diff could be slow. Mitigations: lexer cache, lazy style init, and diffs are typically limited in size (git status is a small working set). For a 1000-line diff, tokenization is O(lines) and still sub-second.
- **Color consistency**: The monokai chroma theme may look slightly different from the "dark" glamour theme. Since they're used in different contexts (markdown vs code), this is acceptable.
- **Empty diff lines**: Some diff lines may be empty (e.g., after stripping the prefix). The renderer handles this by emitting a plain newline.

## Migration Plan

1. Refactor `gitdiff.go`: add diff parser functions, chroma renderer, cache.
2. Update `computeDiff` to use raw `git diff` + parser + renderer.
3. Remove `git diff --color=always` usage (now `git diff --no-color`).
4. `make build && make lint && make test`.
