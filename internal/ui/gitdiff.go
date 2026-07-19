package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"charm.land/lipgloss/v2"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/fselich/dossier/internal/git"
)

type DiffLineType int

const (
	LineContext DiffLineType = iota
	LineAdded
	LineRemoved
	LineHunkHeader
)

type DiffLine struct {
	Type    DiffLineType
	Content string
	OldNum  int
	NewNum  int
}

var (
	lexerCache       sync.Map
	chromaStyleCache = map[string]*chroma.Style{}
)

func getChromaStyle(name string) *chroma.Style {
	if s, ok := chromaStyleCache[name]; ok {
		return s
	}
	s := styles.Get(name)
	if s == nil {
		s = styles.Fallback
	}
	chromaStyleCache[name] = s
	return s
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

func highlightLine(content, filename, bgColor, chromaStyleName string) string {
	if content == "" {
		return ""
	}
	cs := getChromaStyle(chromaStyleName)
	lexer := getLexer(filename)
	iterator, err := lexer.Tokenise(nil, content)
	if err != nil {
		return content
	}
	var b strings.Builder
	for _, token := range iterator.Tokens() {
		entry := cs.Get(token.Type)
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

func (m *Model) renderDiff(lines []DiffLine, filename string, width int, scrollX int, chromaStyleName string, addBgColor, removeBgColor string) string {
	if len(lines) == 0 {
		return m.theme.Styles.Help.Render("  (no diff available)")
	}
	codeWidth := width - 1 - lineNumWidth*2 - 3
	if codeWidth < 10 {
		codeWidth = 10
	}
	var sb strings.Builder
	for _, dl := range lines {
		switch dl.Type {
		case LineHunkHeader:
			sb.WriteString(" " + m.theme.Styles.GitRenamed.Render("    ···  "+dl.Content) + "\n")
		case LineAdded:
			oldNum := fmtLineNum(dl.OldNum)
			newNum := fmtLineNum(dl.NewNum)
			nums := m.theme.Styles.GitAdded.Render(oldNum + " " + newNum)
			content := scrollContent(dl.Content, scrollX)
			indicator := m.theme.Styles.GitAdded.Render("+ ")
			highlighted := highlightLine(content, filename, addBgColor, chromaStyleName)
			line := nums + " " + indicator + highlighted
			line = padLine(line, codeWidth+lineNumWidth*2+3, addBgColor)
			sb.WriteString(" " + line + "\n")
		case LineRemoved:
			oldNum := fmtLineNum(dl.OldNum)
			newNum := fmtLineNum(dl.NewNum)
			nums := m.theme.Styles.DiffRemoved.Render(oldNum + " " + newNum)
			content := scrollContent(dl.Content, scrollX)
			indicator := m.theme.Styles.DiffRemoved.Render("- ")
			highlighted := highlightLine(content, filename, removeBgColor, chromaStyleName)
			line := nums + " " + indicator + highlighted
			line = padLine(line, codeWidth+lineNumWidth*2+3, removeBgColor)
			sb.WriteString(" " + line + "\n")
		case LineContext:
			oldNum := fmtLineNum(dl.OldNum)
			newNum := fmtLineNum(dl.NewNum)
			nums := m.theme.Styles.Help.Render(oldNum + " " + newNum)
			content := scrollContent(dl.Content, scrollX)
			highlighted := highlightLine(content, filename, "", chromaStyleName)
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
	var argsList [][]string
	switch {
	case x == ' ' && y != ' ':
		argsList = [][]string{{"diff", "--", rel}}
	case x != ' ' && y == ' ':
		argsList = [][]string{{"diff", "--cached", "--", rel}}
	default:
		argsList = [][]string{
			{"diff", "HEAD", "--", rel},
			{"diff", "--cached", "--", rel},
			{"diff", "--", rel},
		}
	}
	for _, args := range argsList {
		out, err := git.RunGit(root, args...)
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
	m.loadDiffForFile(m.gitState.Cursor)
}

func (m *Model) loadDiffForFile(cursor int) {
	if cursor >= len(m.gitState.Files) {
		return
	}
	f := m.gitState.Files[cursor]
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
