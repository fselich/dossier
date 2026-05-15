package openspec

import (
	"os"
	"regexp"
	"strings"
)

type ItemKind int

const (
	KindSection ItemKind = iota
	KindTask
)

type TaskItem struct {
	Kind    ItemKind
	Text    string
	Done    bool
	LineNum int
}

var (
	rxSection = regexp.MustCompile(`^## (.+)$`)
	rxPending = regexp.MustCompile(`^- \[ \] (.+)$`)
	rxDone    = regexp.MustCompile(`^- \[x\] (.+)$`)
)

func ParseTasks(content string) []TaskItem {
	lines := strings.Split(content, "\n")
	var items []TaskItem
	for i, line := range lines {
		switch {
		case rxSection.MatchString(line):
			m := rxSection.FindStringSubmatch(line)
			items = append(items, TaskItem{Kind: KindSection, Text: m[1], LineNum: i})
		case rxPending.MatchString(line):
			m := rxPending.FindStringSubmatch(line)
			items = append(items, TaskItem{Kind: KindTask, Text: m[1], Done: false, LineNum: i})
		case rxDone.MatchString(line):
			m := rxDone.FindStringSubmatch(line)
			items = append(items, TaskItem{Kind: KindTask, Text: m[1], Done: true, LineNum: i})
		}
	}
	return items
}

// FindCursorByText returns the index of the first KindTask item with the given
// text, or the index of the first KindTask item if the text is not found.
func FindCursorByText(items []TaskItem, text string) int {
	first := -1
	for i, item := range items {
		if item.Kind != KindTask {
			continue
		}
		if first == -1 {
			first = i
		}
		if item.Text == text {
			return i
		}
	}
	if first == -1 {
		return 0
	}
	return first
}

// ToggleTask flips the done state of items[idx] in memory and on disk.
func ToggleTask(path string, items []TaskItem, idx int) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	if idx >= len(items) || items[idx].LineNum >= len(lines) {
		return nil
	}
	if items[idx].Done {
		lines[items[idx].LineNum] = strings.Replace(lines[items[idx].LineNum], "- [x] ", "- [ ] ", 1)
		items[idx].Done = false
	} else {
		lines[items[idx].LineNum] = strings.Replace(lines[items[idx].LineNum], "- [ ] ", "- [x] ", 1)
		items[idx].Done = true
	}
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)
}
