## 1. Diff parser

- [x] 1.1 Add `DiffLine` struct with `Type` (Context, Added, Removed, HunkHeader) and `Content string`
- [x] 1.2 Add `parseDiff(raw string) []DiffLine` that parses unified diff output, stripping git headers
- [x] 1.3 Add `parseDiffLine(line string) *DiffLine` that classifies a single line by prefix

## 2. Chroma rendering

- [x] 2.1 Add `getLexer(filename string) chroma.Lexer` with `sync.Map` cache per file extension
- [x] 2.2 Add `initChromaStyle()` with `sync.Once` to load monokai style
- [x] 2.3 Add `highlightLine(content, filename, bgColor string) string` that tokenizes a single line with chroma

## 3. Diff renderer

- [x] 3.1 Add `renderDiff(lines []DiffLine, filename string) string` that renders parsed diff with chroma syntax highlighting
- [x] 3.2 Apply correct colors: green + green bg for added, red + red bg for removed, cyan for hunk headers
- [x] 3.3 Handle edge cases: empty lines, binary diffs, git headers stripped

## 4. Refactor computeDiff

- [x] 4.1 Update `computeDiff` to run `git diff` (no --color) and pass through `parseDiff` + `renderDiff`
- [x] 4.2 Keep `highlightUntrackedFile` unchanged for untracked files
- [x] 4.3 Keeping `os/exec` for git commands (still needed)

## 5. Verify

- [x] 5.1 `make build` succeeds
- [x] 5.6 `make lint` passes
- [x] 5.7 `make test` passes
- [x] 5.2 `d` on modified Go file shows syntax-highlighted diff
- [x] 5.3 `d` on modified YAML/Markdown file uses correct lexer
- [x] 5.4 Added lines have green tint, removed lines have red tint
- [x] 5.5 `d` on untracked file still shows syntax-highlighted content
