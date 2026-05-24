package openspec

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseTasks(t *testing.T) {
	t.Run("sections and tasks", func(t *testing.T) {
		content := "## Section A\n- [ ] pending task\n- [x] done task\n## Section B\n- [ ] another pending\n"
		items := ParseTasks(content)
		if len(items) != 5 {
			t.Fatalf("expected 5 items, got %d", len(items))
		}
		if items[0].Kind != KindSection || items[0].Text != "Section A" {
			t.Errorf("item 0: expected section 'Section A', got kind=%d text=%q", items[0].Kind, items[0].Text)
		}
		if items[1].Kind != KindTask || items[1].Done || items[1].Text != "pending task" {
			t.Errorf("item 1: expected pending task, got done=%v text=%q", items[1].Done, items[1].Text)
		}
		if items[2].Kind != KindTask || !items[2].Done || items[2].Text != "done task" {
			t.Errorf("item 2: expected done task, got done=%v text=%q", items[2].Done, items[2].Text)
		}
		if items[3].Kind != KindSection || items[3].Text != "Section B" {
			t.Errorf("item 3: expected section 'Section B', got kind=%d text=%q", items[3].Kind, items[3].Text)
		}
		if items[4].Kind != KindTask || items[4].Done || items[4].Text != "another pending" {
			t.Errorf("item 4: expected pending task, got done=%v text=%q", items[4].Done, items[4].Text)
		}
	})

	t.Run("only sections no tasks", func(t *testing.T) {
		items := ParseTasks("## Only Section\n")
		if len(items) != 1 {
			t.Fatalf("expected 1 item, got %d", len(items))
		}
		if items[0].Kind != KindSection {
			t.Error("expected section")
		}
	})

	t.Run("empty content", func(t *testing.T) {
		items := ParseTasks("")
		if len(items) != 0 {
			t.Errorf("expected 0 items, got %d", len(items))
		}
	})

	t.Run("mixed checkbox formats", func(t *testing.T) {
		content := "- [ ] pending\n- [x] done\n- [X] uppercase done\nsome random line\n"
		items := ParseTasks(content)
		if len(items) != 2 {
			t.Fatalf("expected 2 items ('- [X]' not matched, random line ignored), got %d", len(items))
		}
	})
}

func TestFindCursorByText(t *testing.T) {
	items := ParseTasks("## Section\n- [ ] task A\n- [x] task B\n- [ ] task C\n")

	t.Run("text found", func(t *testing.T) {
		idx := FindCursorByText(items, "task B")
		if idx != 2 {
			t.Errorf("expected idx 2, got %d", idx)
		}
	})

	t.Run("text not found returns first task", func(t *testing.T) {
		idx := FindCursorByText(items, "nonexistent")
		if idx != 1 {
			t.Errorf("expected idx 1 (first task), got %d", idx)
		}
	})

	t.Run("only sections returns 0", func(t *testing.T) {
		sectionOnly := ParseTasks("## Solo section\n")
		idx := FindCursorByText(sectionOnly, "anything")
		if idx != 0 {
			t.Errorf("expected 0 for section-only items, got %d", idx)
		}
	})
}

func TestToggleTask(t *testing.T) {
	t.Run("pending to done writes to disk", func(t *testing.T) {
		root := t.TempDir()
		tasksFile := filepath.Join(root, "tasks.md")
		content := "- [ ] pending task\n- [x] done task\n"
		if err := os.WriteFile(tasksFile, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		items := ParseTasks(content)
		taskIdx := 0 // the pending task
		if err := ToggleTask(tasksFile, items, taskIdx); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !items[taskIdx].Done {
			t.Error("expected item to be marked done in memory")
		}

		data, err := os.ReadFile(tasksFile)
		if err != nil {
			t.Fatal(err)
		}
		expected := "- [x] pending task\n- [x] done task\n"
		if string(data) != expected {
			t.Errorf("expected file content %q, got %q", expected, string(data))
		}
	})

	t.Run("done to pending writes to disk", func(t *testing.T) {
		root := t.TempDir()
		tasksFile := filepath.Join(root, "tasks.md")
		content := "- [x] done task\n"
		if err := os.WriteFile(tasksFile, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		items := ParseTasks(content)
		taskIdx := 0
		if err := ToggleTask(tasksFile, items, taskIdx); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if items[taskIdx].Done {
			t.Error("expected item to be marked pending in memory")
		}

		data, err := os.ReadFile(tasksFile)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != "- [ ] done task\n" {
			t.Errorf("expected '- [ ] done task', got %q", string(data))
		}
	})

	t.Run("idx out of range returns nil", func(t *testing.T) {
		root := t.TempDir()
		tasksFile := filepath.Join(root, "tasks.md")
		if err := os.WriteFile(tasksFile, []byte("- [ ] task\n"), 0644); err != nil {
			t.Fatal(err)
		}
		items := ParseTasks("- [ ] task\n")
		if err := ToggleTask(tasksFile, items, 99); err != nil {
			t.Errorf("expected nil for out-of-range idx, got %v", err)
		}
	})

	t.Run("file does not exist returns error", func(t *testing.T) {
		items := ParseTasks("- [ ] task\n")
		err := ToggleTask("/nonexistent/tasks.md", items, 0)
		if err == nil {
			t.Error("expected error for nonexistent file")
		}
	})
}
