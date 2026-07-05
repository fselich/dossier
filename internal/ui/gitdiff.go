package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"charm.land/lipgloss/v2"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

type DiffLineType int

const (
	LineContext DiffLineType = iota
	LineAdded
	LineRemoved
	LineHunkHeader

	addBgColor    = "#1a3a1a"
	removeBgColor = "#3a1a1a"
)

type DiffLine struct {
	Type    DiffLineType
	Content string
	OldNum  int
	NewNum  int
}

var (
	lexerCache  sync.Map
	chromaStyle *chroma.Style
	chromaOnce  sync.Once
)

func initChromaStyle() {
	chromaOnce.Do(func() {
		chromaStyle = styles.Get("monokai")
		if chromaStyle == nil {
			chromaStyle = styles.Fallback
		}
	})
}

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

func highlightLine(content, filename, bgColor string) string {
	initChromaStyle()
	if content == "" {
		return ""
	}
	lexer := getLexer(filename)
	iterator, err := lexer.Tokenise(nil, content)
	if err != nil {
		return content
	}
	var b strings.Builder
	for _, token := range iterator.Tokens() {
		entry := chromaStyle.Get(token.Type)
		style := lipgloss.NewStyle()
		if entry.Colour.IsSet() {
			style = style.Foreground(lipgloss.Color(entry.Colour.String()))
		}
		if bgColor != "" {
			style = style.Background(lipgloss.Color(bgColor))
		}
		b.WriteString(style.Render(token.Value))
	}
	return b.String()
}

func parseDiff(raw string) []DiffLine {
	lines := strings.Split(raw, "\n")
	result := make([]DiffLine, 0, len(lines))
	oldNum, newNum := 0, 0
	for _, line := range lines {
		dl := parseDiffLine(line, &oldNum, &newNum)
		if dl != nil {
			result = append(result, *dl)
		}
	}
	return result
}

func parseHunkHeader(line string, oldNum, newNum *int) {
	// @@ -old,count +new,count @@
	parts := strings.SplitN(line, "@@", 3)
	if len(parts) < 2 {
		return
	}
	ranges := strings.TrimSpace(parts[1])
	for _, r := range strings.Fields(ranges) {
		if strings.HasPrefix(r, "-") {
			nums := strings.SplitN(r[1:], ",", 2)
			if n, err := strconv.Atoi(nums[0]); err == nil {
				*oldNum = n
			}
		} else if strings.HasPrefix(r, "+") {
			nums := strings.SplitN(r[1:], ",", 2)
			if n, err := strconv.Atoi(nums[0]); err == nil {
				*newNum = n
			}
		}
	}
}

func parseDiffLine(line string, oldNum, newNum *int) *DiffLine {
	switch {
	case strings.HasPrefix(line, "diff --git"),
		strings.HasPrefix(line, "index "),
		strings.HasPrefix(line, "new file"),
		strings.HasPrefix(line, "deleted file"),
		strings.HasPrefix(line, "similarity"),
		strings.HasPrefix(line, "rename"),
		strings.HasPrefix(line, "old mode"),
		strings.HasPrefix(line, "new mode"),
		strings.HasPrefix(line, "--- "),
		strings.HasPrefix(line, "+++ "):
		return nil
	case strings.HasPrefix(line, "@@"):
		parseHunkHeader(line, oldNum, newNum)
		return &DiffLine{Type: LineHunkHeader, Content: line, OldNum: -1, NewNum: -1}
	case strings.HasPrefix(line, "+"):
		dl := &DiffLine{Type: LineAdded, Content: line[1:], OldNum: -1, NewNum: *newNum}
		*newNum++
		return dl
	case strings.HasPrefix(line, "-"):
		dl := &DiffLine{Type: LineRemoved, Content: line[1:], OldNum: *oldNum, NewNum: -1}
		*oldNum++
		return dl
	case strings.HasPrefix(line, `\`):
		return nil
	case strings.HasPrefix(line, " "):
		dl := &DiffLine{Type: LineContext, Content: line[1:], OldNum: *oldNum, NewNum: *newNum}
		*oldNum++
		*newNum++
		return dl
	case line == "":
		return nil
	default:
		dl := &DiffLine{Type: LineContext, Content: line, OldNum: *oldNum, NewNum: *newNum}
		*oldNum++
		*newNum++
		return dl
	}
}

func fmtLineNum(n int) string {
	if n < 0 {
		return "    "
	}
	return fmt.Sprintf("%4d", n)
}

const lineNumWidth = 4

func renderDiff(lines []DiffLine, filename string, width int, scrollX int) string {
	if len(lines) == 0 {
		return helpStyle.Render("  (no diff available)")
	}
	codeWidth := width - 1 - lineNumWidth*2 - 3
	if codeWidth < 10 {
		codeWidth = 10
	}
	var sb strings.Builder
	for _, dl := range lines {
		switch dl.Type {
		case LineHunkHeader:
			sb.WriteString(" " + gitStatusRenamed.Render("    ···  "+dl.Content) + "\n")
		case LineAdded:
			oldNum := fmtLineNum(dl.OldNum)
			newNum := fmtLineNum(dl.NewNum)
			nums := gitStatusAdded.Render(oldNum + " " + newNum)
			content := scrollContent(dl.Content, scrollX)
			indicator := gitStatusAdded.Render("+ ")
			highlighted := highlightLine(content, filename, addBgColor)
			line := nums + " " + indicator + highlighted
			line = padLine(line, codeWidth+lineNumWidth*2+3, addBgColor)
			sb.WriteString(" " + line + "\n")
		case LineRemoved:
			oldNum := fmtLineNum(dl.OldNum)
			newNum := fmtLineNum(dl.NewNum)
			nums := diffRemoved.Render(oldNum + " " + newNum)
			content := scrollContent(dl.Content, scrollX)
			indicator := diffRemoved.Render("- ")
			highlighted := highlightLine(content, filename, removeBgColor)
			line := nums + " " + indicator + highlighted
			line = padLine(line, codeWidth+lineNumWidth*2+3, removeBgColor)
			sb.WriteString(" " + line + "\n")
		case LineContext:
			oldNum := fmtLineNum(dl.OldNum)
			newNum := fmtLineNum(dl.NewNum)
			nums := helpStyle.Render(oldNum + " " + newNum)
			content := scrollContent(dl.Content, scrollX)
			highlighted := highlightLine(content, filename, "")
			line := nums + "  " + highlighted
			sb.WriteString(" " + line + "\n")
		}
	}
	return sb.String()
}

func scrollContent(content string, scrollX int) string {
	runes := []rune(content)
	if scrollX >= len(runes) {
		return ""
	}
	return string(runes[scrollX:])
}

func padLine(line string, targetWidth int, bgColor string) string {
	w := lipgloss.Width(line)
	if w >= targetWidth {
		return line
	}
	if bgColor == "" {
		return line
	}
	return line + lipgloss.NewStyle().Background(lipgloss.Color(bgColor)).Render(strings.Repeat(" ", targetWidth-w))
}

func computeDiffLines(root string, x, y byte, rel, oldPath string) []DiffLine {
	if len(rel) == 0 {
		return nil
	}
	var cmds []*exec.Cmd
	switch {
	case x == ' ' && y != ' ':
		cmds = []*exec.Cmd{exec.Command("git", "-C", root, "diff", "--", rel)}
	case x != ' ' && y == ' ':
		cmds = []*exec.Cmd{exec.Command("git", "-C", root, "diff", "--cached", "--", rel)}
	default:
		cmds = []*exec.Cmd{
			exec.Command("git", "-C", root, "diff", "HEAD", "--", rel),
			exec.Command("git", "-C", root, "diff", "--cached", "--", rel),
			exec.Command("git", "-C", root, "diff", "--", rel),
		}
	}
	for _, cmd := range cmds {
		out, err := cmd.Output()
		if err == nil && len(out) > 0 {
			return parseDiff(string(out))
		}
	}
	return nil
}

func untrackedFileDiffLines(root, rel string) []DiffLine {
	path := rel
	if !filepath.IsAbs(path) {
		path = filepath.Join(root, path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	lines := strings.Split(string(data), "\n")
	result := make([]DiffLine, 0, len(lines))
	for i, line := range lines {
		n := i + 1
		result = append(result, DiffLine{Type: LineContext, Content: line, OldNum: n, NewNum: n})
	}
	return result
}

func (m *Model) toggleGitDiff() {
	if m.gitState.ShowingDiff {
		m.gitState.ShowingDiff = false
		m.gitState.DiffLines = nil
		m.gitState.DiffFile = ""
		m.gitState.ScrollX = 0
		m.refreshGitViewport()
		return
	}
	if m.gitState.Cursor >= len(m.gitState.Files) {
		return
	}
	f := m.gitState.Files[m.gitState.Cursor]
	if f.IsDeleted {
		return
	}
	var lines []DiffLine
	if f.X == '?' && f.Y == '?' {
		lines = untrackedFileDiffLines(m.gitRoot, f.Path)
	} else {
		lines = computeDiffLines(m.gitRoot, f.X, f.Y, f.Path, f.OldPath)
	}
	if lines == nil {
		m.gitState.DiffLines = nil
	}
	m.gitState.DiffLines = lines
	m.gitState.DiffFile = f.Path
	m.gitState.ShowingDiff = true
	m.gitState.ScrollX = 0
	m.refreshGitViewport()
}
