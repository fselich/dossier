package openspec

import (
	"os"
	"path/filepath"
	"testing"
)

func setupProjectDir(t *testing.T, changes []string) string {
	t.Helper()
	root := t.TempDir()
	openspecDir := filepath.Join(root, "openspec")
	changesDir := filepath.Join(openspecDir, "changes")
	if err := os.MkdirAll(changesDir, 0755); err != nil {
		t.Fatal(err)
	}
	for _, name := range changes {
		dir := filepath.Join(changesDir, name)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, ".openspec.yaml"), []byte("created: 2026-05-24"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestLoadConfigFrom(t *testing.T) {
	t.Run("valid file", func(t *testing.T) {
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, "openspec"), 0755); err != nil {
			t.Fatal(err)
		}
		yaml := `schema: spec-driven
context: |
  Tech stack: Go.
  Domain: TUI tool.
rules:
  proposal:
    - Keep it concise
  tasks:
    - Small steps
`
		if err := os.WriteFile(filepath.Join(root, "openspec", "config.yaml"), []byte(yaml), 0644); err != nil {
			t.Fatal(err)
		}
		cfg, err := LoadConfigFrom(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Context == "" {
			t.Error("expected non-empty Context")
		}
		if len(cfg.Rules) != 2 {
			t.Errorf("expected 2 rule groups, got %d", len(cfg.Rules))
		}
		if len(cfg.Rules["proposal"]) != 1 {
			t.Errorf("expected 1 proposal rule, got %d", len(cfg.Rules["proposal"]))
		}
	})

	t.Run("missing file", func(t *testing.T) {
		root := t.TempDir()
		cfg, err := LoadConfigFrom(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Context != "" {
			t.Error("expected empty Context for missing file")
		}
		if len(cfg.Rules) != 0 {
			t.Error("expected empty Rules for missing file")
		}
	})

	t.Run("malformed YAML", func(t *testing.T) {
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, "openspec"), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, "openspec", "config.yaml"), []byte("{bad yaml: ["), 0644); err != nil {
			t.Fatal(err)
		}
		_, err := LoadConfigFrom(root)
		if err == nil {
			t.Error("expected error for malformed YAML")
		}
	})

	t.Run("empty context", func(t *testing.T) {
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, "openspec"), 0755); err != nil {
			t.Fatal(err)
		}
		yaml := `schema: spec-driven
rules:
  proposal:
    - Keep it concise
`
		if err := os.WriteFile(filepath.Join(root, "openspec", "config.yaml"), []byte(yaml), 0644); err != nil {
			t.Fatal(err)
		}
		cfg, err := LoadConfigFrom(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Context != "" {
			t.Errorf("expected empty Context, got %q", cfg.Context)
		}
		if len(cfg.Rules["proposal"]) != 1 {
			t.Errorf("expected 1 proposal rule, got %d", len(cfg.Rules["proposal"]))
		}
	})

	t.Run("missing rules", func(t *testing.T) {
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, "openspec"), 0755); err != nil {
			t.Fatal(err)
		}
		yaml := `schema: spec-driven
context: Just context.
`
		if err := os.WriteFile(filepath.Join(root, "openspec", "config.yaml"), []byte(yaml), 0644); err != nil {
			t.Fatal(err)
		}
		cfg, err := LoadConfigFrom(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Context != "Just context." {
			t.Errorf("expected context, got %q", cfg.Context)
		}
		if len(cfg.Rules) != 0 {
			t.Errorf("expected empty Rules, got %v", cfg.Rules)
		}
	})
}

func TestLoadFrom(t *testing.T) {
	t.Run("valid project with changes", func(t *testing.T) {
		root := setupProjectDir(t, []string{"feat-a", "feat-b"})
		proj, err := LoadFrom(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(proj.Changes) != 2 {
			t.Errorf("expected 2 changes, got %d", len(proj.Changes))
		}
	})

	t.Run("missing openspec directory", func(t *testing.T) {
		root := t.TempDir()
		_, err := LoadFrom(root)
		if err == nil {
			t.Error("expected error for missing openspec/ directory")
		}
	})

	t.Run("empty changes directory", func(t *testing.T) {
		root := setupProjectDir(t, nil)
		proj, err := LoadFrom(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(proj.Changes) != 0 {
			t.Errorf("expected 0 changes, got %d", len(proj.Changes))
		}
	})
}

func TestLoadProjectSpecsFrom(t *testing.T) {
	t.Run("specs with subdirectories", func(t *testing.T) {
		root := t.TempDir()
		specsDir := filepath.Join(root, "openspec", "specs")
		for _, name := range []string{"auth", "profile"} {
			dir := filepath.Join(specsDir, name)
			if err := os.MkdirAll(dir, 0755); err != nil {
				t.Fatal(err)
			}
			content := "### Requirement: " + name + "-login\nDescription\n"
			if err := os.WriteFile(filepath.Join(dir, "spec.md"), []byte(content), 0644); err != nil {
				t.Fatal(err)
			}
		}
		specs, err := LoadProjectSpecsFrom(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(specs) != 2 {
			t.Errorf("expected 2 specs, got %d", len(specs))
		}
		if specs[0].Name != "auth" {
			t.Errorf("expected first spec 'auth', got %q", specs[0].Name)
		}
		if specs[1].Name != "profile" {
			t.Errorf("expected second spec 'profile', got %q", specs[1].Name)
		}
		if specs[0].RequirementCount != 1 {
			t.Errorf("expected 1 requirement in auth, got %d", specs[0].RequirementCount)
		}
	})

	t.Run("missing specs directory", func(t *testing.T) {
		root := t.TempDir()
		specs, err := LoadProjectSpecsFrom(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(specs) != 0 {
			t.Errorf("expected 0 specs, got %d", len(specs))
		}
	})

	t.Run("specs without spec.md", func(t *testing.T) {
		root := t.TempDir()
		dir := filepath.Join(root, "openspec", "specs", "empty-spec")
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		specs, err := LoadProjectSpecsFrom(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(specs) != 1 {
			t.Errorf("expected 1 spec entry, got %d", len(specs))
		}
		if specs[0].RequirementCount != 0 {
			t.Errorf("expected 0 requirements, got %d", specs[0].RequirementCount)
		}
	})
}

func TestListChangeNamesFrom(t *testing.T) {
	t.Run("with active changes", func(t *testing.T) {
		root := setupProjectDir(t, []string{"feat-a", "feat-b"})
		names, err := ListChangeNamesFrom(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(names) != 2 {
			t.Errorf("expected 2 names, got %d", len(names))
		}
	})

	t.Run("empty directory", func(t *testing.T) {
		root := setupProjectDir(t, nil)
		names, err := ListChangeNamesFrom(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(names) != 0 {
			t.Errorf("expected 0 names, got %d", len(names))
		}
	})

	t.Run("missing changes directory", func(t *testing.T) {
		root := t.TempDir()
		names, err := ListChangeNamesFrom(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(names) != 0 {
			t.Errorf("expected 0 names, got %d", len(names))
		}
	})
}

func TestListArchiveChangesFrom(t *testing.T) {
	t.Run("with archived changes", func(t *testing.T) {
		root := t.TempDir()
		archiveDir := filepath.Join(root, "openspec", "changes", "archive")
		for _, name := range []string{"2026-05-10-old-feat", "2026-05-05-ancient"} {
			dir := filepath.Join(archiveDir, name)
			if err := os.MkdirAll(dir, 0755); err != nil {
				t.Fatal(err)
			}
		}
		changes, err := ListArchiveChangesFrom(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(changes) != 2 {
			t.Errorf("expected 2 archived changes, got %d", len(changes))
		}
		if changes[0].Name != "old-feat" {
			t.Errorf("expected most recent first, got %q", changes[0].Name)
		}
	})

	t.Run("missing archive directory", func(t *testing.T) {
		root := t.TempDir()
		changes, err := ListArchiveChangesFrom(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(changes) != 0 {
			t.Errorf("expected 0 changes, got %d", len(changes))
		}
	})
}

func TestListArchiveNamesFrom(t *testing.T) {
	t.Run("with entries", func(t *testing.T) {
		root := t.TempDir()
		archiveDir := filepath.Join(root, "openspec", "changes", "archive")
		for _, name := range []string{"2026-05-10-old", "2026-05-05-older"} {
			if err := os.MkdirAll(filepath.Join(archiveDir, name), 0755); err != nil {
				t.Fatal(err)
			}
		}
		names, err := ListArchiveNamesFrom(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(names) != 2 {
			t.Errorf("expected 2 names, got %d", len(names))
		}
		if names[0] != "2026-05-10-old" {
			t.Errorf("expected most recent first, got %q", names[0])
		}
	})

	t.Run("missing archive directory", func(t *testing.T) {
		root := t.TempDir()
		names, err := ListArchiveNamesFrom(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(names) != 0 {
			t.Errorf("expected 0 names, got %d", len(names))
		}
	})
}

func TestListSpecNamesFrom(t *testing.T) {
	t.Run("with entries", func(t *testing.T) {
		root := t.TempDir()
		specsDir := filepath.Join(root, "openspec", "specs")
		for _, name := range []string{"auth", "profile"} {
			if err := os.MkdirAll(filepath.Join(specsDir, name), 0755); err != nil {
				t.Fatal(err)
			}
		}
		names, err := ListSpecNamesFrom(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(names) != 2 {
			t.Errorf("expected 2 names, got %d", len(names))
		}
		if names[0] != "auth" {
			t.Errorf("expected 'auth' first, got %q", names[0])
		}
	})

	t.Run("missing specs directory", func(t *testing.T) {
		root := t.TempDir()
		names, err := ListSpecNamesFrom(root)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(names) != 0 {
			t.Errorf("expected 0 names, got %d", len(names))
		}
	})
}

func TestLoadFromPath(t *testing.T) {
	t.Run("valid change path", func(t *testing.T) {
		root := t.TempDir()
		dir := filepath.Join(root, "openspec", "changes", "my-change")
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, ".openspec.yaml"), []byte("created: 2026-05-24"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "proposal.md"), []byte("# Proposal"), 0644); err != nil {
			t.Fatal(err)
		}
		proj, err := LoadFromPath(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(proj.Changes) != 1 {
			t.Errorf("expected 1 change, got %d", len(proj.Changes))
		}
		if proj.Changes[0].Name != "my-change" {
			t.Errorf("expected 'my-change', got %q", proj.Changes[0].Name)
		}
	})

	t.Run("nonexistent path", func(t *testing.T) {
		_, err := LoadFromPath("/nonexistent/path")
		if err == nil {
			t.Error("expected error for nonexistent path")
		}
	})

	t.Run("path without openspec yaml", func(t *testing.T) {
		root := t.TempDir()
		dir := filepath.Join(root, "not-a-change")
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		_, err := LoadFromPath(dir)
		if err == nil {
			t.Error("expected error for path without .openspec.yaml")
		}
	})
}

func TestReloadChange(t *testing.T) {
	t.Run("file modified on disk produces updated content", func(t *testing.T) {
		root := setupProjectDir(t, []string{"my-change"})
		proj, err := LoadFrom(root)
		if err != nil {
			t.Fatal(err)
		}
		ch := proj.Changes[0]

		tasksFile := filepath.Join(ch.Path, "tasks.md")
		if err := os.WriteFile(tasksFile, []byte("- [ ] updated task"), 0644); err != nil {
			t.Fatal(err)
		}

		reloaded := ReloadChange(ch)
		if reloaded.Tasks.Content != "- [ ] updated task" {
			t.Errorf("expected updated content, got %q", reloaded.Tasks.Content)
		}
		if !reloaded.Tasks.Present {
			t.Error("expected Tasks to be present")
		}
	})

	t.Run("file deleted produces absent artifact", func(t *testing.T) {
		root := setupProjectDir(t, []string{"my-change"})
		proj, err := LoadFrom(root)
		if err != nil {
			t.Fatal(err)
		}
		ch := proj.Changes[0]

		if err := os.WriteFile(filepath.Join(ch.Path, "proposal.md"), []byte("# Proposal"), 0644); err != nil {
			t.Fatal(err)
		}
		initial := ReloadChange(ch)
		if !initial.Proposal.Present {
			t.Fatal("expected Proposal to be present after write")
		}

		if err := os.Remove(filepath.Join(ch.Path, "proposal.md")); err != nil {
			t.Fatal(err)
		}
		afterDelete := ReloadChange(ch)
		if afterDelete.Proposal.Present {
			t.Error("expected Proposal to be absent after delete")
		}
	})
}

func TestExtractPurpose(t *testing.T) {
	t.Run("purpose present", func(t *testing.T) {
		content := "# spec\n\n## Purpose\nDefines the layout.\n\n## Requirements\n"
		got := ExtractPurpose(content)
		want := "Defines the layout."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("no purpose heading", func(t *testing.T) {
		content := "# spec\n\n## Requirements\n"
		got := ExtractPurpose(content)
		if got != "" {
			t.Errorf("expected empty, got %q", got)
		}
	})

	t.Run("purpose at EOF", func(t *testing.T) {
		content := "# spec\n\n## Purpose\nDefines the layout and main behavior."
		got := ExtractPurpose(content)
		want := "Defines the layout and main behavior."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("multi-line purpose", func(t *testing.T) {
		content := "# spec\n\n## Purpose\nDefines the layout.\nMain behavior.\n\n## Requirements\n"
		got := ExtractPurpose(content)
		want := "Defines the layout. Main behavior."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("markdown stripped from purpose", func(t *testing.T) {
		content := "# spec\n\n## Purpose\n**Bold** and *italic* and `code` and [link](url).\n\n## Requirements\n"
		got := ExtractPurpose(content)
		want := "Bold and italic and code and link."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
