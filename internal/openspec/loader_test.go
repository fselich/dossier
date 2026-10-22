package openspec

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type failingReadFS struct {
	OSFS
	path string
}

func (fs failingReadFS) ReadFile(name string) ([]byte, error) {
	if name == fs.path {
		return nil, fmt.Errorf("injected read failure")
	}
	return fs.OSFS.ReadFile(name)
}

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

func TestSchemaAwareChangeModel(t *testing.T) {
	t.Run("loads schema metadata and conventional artifacts dynamically", func(t *testing.T) {
		root := setupProjectDir(t, []string{"my-change"})
		dir := filepath.Join(root, "openspec", "changes", "my-change")
		if err := os.WriteFile(filepath.Join(dir, ".openspec.yaml"), []byte("schema: spec-driven\ncreated: 2026-05-24\n"), 0644); err != nil {
			t.Fatal(err)
		}
		files := map[string]string{
			"proposal.md":            "# Proposal",
			"design.md":              "# Design",
			"tasks.md":               "- [ ] task",
			"specs/auth/spec.md":     "# Auth",
			"specs/payments/spec.md": "# Payments",
		}
		for name, content := range files {
			path := filepath.Join(dir, name)
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(content), 0644); err != nil {
				t.Fatal(err)
			}
		}

		proj, err := LoadFrom(root)
		if err != nil {
			t.Fatal(err)
		}
		ch := proj.Changes[0]
		if ch.Schema != "spec-driven" {
			t.Fatalf("expected schema spec-driven, got %q", ch.Schema)
		}
		wantIDs := []string{"proposal", "specs", "design", "tasks"}
		if len(ch.Artifacts) != len(wantIDs) {
			t.Fatalf("expected %d artifacts, got %d", len(wantIDs), len(ch.Artifacts))
		}
		for i, want := range wantIDs {
			if ch.Artifacts[i].ID != want {
				t.Errorf("artifact %d: expected %q, got %q", i, want, ch.Artifacts[i].ID)
			}
			if ch.Artifacts[i].Source != ArtifactSourceDiscovered {
				t.Errorf("artifact %q: expected discovered source, got %q", want, ch.Artifacts[i].Source)
			}
		}
		specs, ok := ch.ArtifactByID("specs")
		if !ok {
			t.Fatal("expected dynamic specs artifact")
		}
		if len(specs.Outputs) != 2 {
			t.Fatalf("expected 2 spec outputs, got %d", len(specs.Outputs))
		}
		if specs.Outputs[0].RelativePath != "specs/auth/spec.md" || specs.Outputs[1].RelativePath != "specs/payments/spec.md" {
			t.Fatalf("unexpected spec output order: %#v", specs.Outputs)
		}
	})

	t.Run("retains arbitrary artifact identifiers and dependency metadata", func(t *testing.T) {
		ch := Change{Artifacts: []ChangeArtifact{
			{ID: "proposal", Status: ArtifactStatusDone},
			{
				ID:       "diagnosis",
				Status:   ArtifactStatusReady,
				Requires: []string{"proposal", "reproduction"},
				Outputs: []ArtifactOutput{{
					RelativePath: "diagnosis.md",
					DisplayName:  "diagnosis",
					Content:      "# Root cause",
					Present:      true,
				}},
			},
		}}

		artifact, ok := ch.ArtifactByID("diagnosis")
		if !ok {
			t.Fatal("expected diagnosis artifact")
		}
		if artifact.Status != ArtifactStatusReady {
			t.Errorf("expected ready status, got %q", artifact.Status)
		}
		if len(artifact.Requires) != 2 || artifact.Requires[1] != "reproduction" {
			t.Errorf("unexpected requirements: %#v", artifact.Requires)
		}
		if len(artifact.Outputs) != 1 || artifact.Outputs[0].Content != "# Root cause" {
			t.Errorf("unexpected outputs: %#v", artifact.Outputs)
		}
	})

	t.Run("malformed metadata retains a diagnostic", func(t *testing.T) {
		root := setupProjectDir(t, []string{"my-change"})
		dir := filepath.Join(root, "openspec", "changes", "my-change")
		if err := os.WriteFile(filepath.Join(dir, ".openspec.yaml"), []byte("schema: [broken\n"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "proposal.md"), []byte("# Proposal"), 0644); err != nil {
			t.Fatal(err)
		}

		proj, err := LoadFrom(root)
		if err != nil {
			t.Fatal(err)
		}
		ch := proj.Changes[0]
		if ch.Schema != "" {
			t.Errorf("expected unknown schema, got %q", ch.Schema)
		}
		if ch.Diagnostic == "" {
			t.Fatal("expected malformed metadata diagnostic")
		}
		if _, ok := ch.ArtifactByID("proposal"); !ok {
			t.Fatal("expected readable proposal despite malformed metadata")
		}
	})
}

func TestDiscoverChangeArtifacts(t *testing.T) {
	t.Run("discovers arbitrary markdown artifacts deterministically", func(t *testing.T) {
		root := setupProjectDir(t, []string{"fix-cache"})
		dir := filepath.Join(root, "openspec", "changes", "fix-cache")
		if err := os.WriteFile(filepath.Join(dir, ".openspec.yaml"), []byte("schema: bugfix\ncreated: 2026-05-24\n"), 0644); err != nil {
			t.Fatal(err)
		}
		for name, content := range map[string]string{
			"proposal.md":     "# Proposal",
			"reproduction.md": "# Reproduction",
			"diagnosis.md":    "# Diagnosis",
			"tasks.md":        "- [ ] fix",
			"README.md":       "generated description",
			"notes.txt":       "not markdown",
			".hidden.md":      "hidden",
		} {
			if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
				t.Fatal(err)
			}
		}

		proj, err := LoadFrom(root)
		if err != nil {
			t.Fatal(err)
		}
		ch := proj.Changes[0]
		for _, id := range []string{"proposal", "reproduction", "diagnosis", "tasks"} {
			if _, ok := ch.ArtifactByID(id); !ok {
				t.Errorf("expected discovered artifact %q", id)
			}
		}
		for _, id := range []string{"README", "notes", ".hidden"} {
			if _, ok := ch.ArtifactByID(id); ok {
				t.Errorf("did not expect artifact %q", id)
			}
		}

		projAgain, err := LoadFrom(root)
		if err != nil {
			t.Fatal(err)
		}
		if len(ch.Artifacts) != len(projAgain.Changes[0].Artifacts) {
			t.Fatalf("artifact count changed between loads")
		}
		for i := range ch.Artifacts {
			if ch.Artifacts[i].ID != projAgain.Changes[0].Artifacts[i].ID {
				t.Errorf("artifact order changed at %d: %q != %q", i, ch.Artifacts[i].ID, projAgain.Changes[0].Artifacts[i].ID)
			}
		}
	})

	t.Run("groups nested markdown outputs under their top-level artifact", func(t *testing.T) {
		root := setupProjectDir(t, []string{"research-change"})
		dir := filepath.Join(root, "openspec", "changes", "research-change")
		files := map[string]string{
			"evidence/benchmarks/latency.md": "# Latency",
			"evidence/compatibility.md":      "# Compatibility",
			"specs/auth/spec.md":             "# Auth",
			"specs/payments/spec.md":         "# Payments",
		}
		for name, content := range files {
			path := filepath.Join(dir, name)
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, []byte(content), 0644); err != nil {
				t.Fatal(err)
			}
		}

		proj, err := LoadFrom(root)
		if err != nil {
			t.Fatal(err)
		}
		ch := proj.Changes[0]
		evidence, ok := ch.ArtifactByID("evidence")
		if !ok {
			t.Fatal("expected evidence artifact")
		}
		if len(evidence.Outputs) != 2 {
			t.Fatalf("expected 2 evidence outputs, got %d", len(evidence.Outputs))
		}
		if evidence.Outputs[0].RelativePath != "evidence/benchmarks/latency.md" || evidence.Outputs[1].RelativePath != "evidence/compatibility.md" {
			t.Fatalf("unexpected evidence output order: %#v", evidence.Outputs)
		}
		specs, ok := ch.ArtifactByID("specs")
		if !ok || len(specs.Outputs) != 2 {
			t.Fatalf("expected two specs outputs, got %#v", specs)
		}
	})

	t.Run("ignores symlinks that could escape the change root", func(t *testing.T) {
		root := setupProjectDir(t, []string{"safe-change"})
		dir := filepath.Join(root, "openspec", "changes", "safe-change")
		outside := filepath.Join(root, "outside.md")
		if err := os.WriteFile(outside, []byte("secret outside content"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(outside, filepath.Join(dir, "escaped.md")); err != nil {
			t.Skipf("symlink unavailable: %v", err)
		}
		if err := os.WriteFile(filepath.Join(dir, "proposal.md"), []byte("# Safe"), 0644); err != nil {
			t.Fatal(err)
		}

		proj, err := LoadFrom(root)
		if err != nil {
			t.Fatal(err)
		}
		ch := proj.Changes[0]
		if _, ok := ch.ArtifactByID("escaped"); ok {
			t.Fatal("symlinked output must not be discovered")
		}
		if ch.Diagnostic == "" {
			t.Fatal("expected diagnostic for skipped symlink")
		}
	})

	t.Run("caps discovery and reports the limit once", func(t *testing.T) {
		root := setupProjectDir(t, []string{"large-change"})
		dir := filepath.Join(root, "openspec", "changes", "large-change")
		if err := os.MkdirAll(filepath.Join(dir, "artifacts"), 0755); err != nil {
			t.Fatal(err)
		}
		for i := range maxDiscoveredArtifactFiles + 1 {
			name := filepath.Join(dir, "artifacts", fmt.Sprintf("artifact-%03d.md", i))
			if err := os.WriteFile(name, []byte("# Artifact"), 0644); err != nil {
				t.Fatal(err)
			}
		}
		if err := os.WriteFile(filepath.Join(dir, "z.md"), []byte("# Extra"), 0644); err != nil {
			t.Fatal(err)
		}

		proj, err := LoadFrom(root)
		if err != nil {
			t.Fatal(err)
		}
		ch := proj.Changes[0]
		artifact, ok := ch.ArtifactByID("artifacts")
		if !ok {
			t.Fatal("expected grouped artifacts output")
		}
		if len(artifact.Outputs) != maxDiscoveredArtifactFiles {
			t.Fatalf("expected %d discovered outputs, got %d", maxDiscoveredArtifactFiles, len(artifact.Outputs))
		}
		if count := strings.Count(ch.Diagnostic, "artifact discovery limited"); count != 1 {
			t.Fatalf("expected one limit diagnostic, got %d: %q", count, ch.Diagnostic)
		}
	})

	t.Run("reports unreadable markdown and preserves readable artifacts", func(t *testing.T) {
		root := setupProjectDir(t, []string{"partial-change"})
		dir := filepath.Join(root, "openspec", "changes", "partial-change")
		brokenPath := filepath.Join(dir, "broken.md")
		if err := os.WriteFile(brokenPath, []byte("# Broken"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "proposal.md"), []byte("# Proposal"), 0644); err != nil {
			t.Fatal(err)
		}
		loader := NewLoader(failingReadFS{OSFS: OSFS{}, path: brokenPath})

		proj, err := loader.LoadFrom(root)
		if err != nil {
			t.Fatal(err)
		}
		ch := proj.Changes[0]
		if _, ok := ch.ArtifactByID("proposal"); !ok {
			t.Fatal("expected readable proposal")
		}
		if _, ok := ch.ArtifactByID("broken"); ok {
			t.Fatal("unreadable output must not be present")
		}
		if !strings.Contains(ch.Diagnostic, "read broken.md: injected read failure") {
			t.Fatalf("expected read diagnostic, got %q", ch.Diagnostic)
		}
	})

	t.Run("discovers custom artifacts in archived changes", func(t *testing.T) {
		root := t.TempDir()
		archiveRoot := filepath.Join(root, "openspec", "changes", "archive")
		fixtures := []struct {
			dir       string
			schema    string
			artifacts map[string]string
		}{
			{
				dir:    "2026-05-24-investigate-cache",
				schema: "spike",
				artifacts: map[string]string{
					"proposal.md": "# Proposal", "research.md": "# Research", "decision.md": "# Decision",
				},
			},
			{
				dir:    "2026-05-23-add-export",
				schema: "feature",
				artifacts: map[string]string{
					"proposal.md": "# Proposal", "design.md": "# Design", "tasks.md": "- [x] done",
				},
			},
			{
				dir: "legacy-unknown",
				artifacts: map[string]string{
					"notes.md": "# Notes",
				},
			},
		}
		for _, fixture := range fixtures {
			dir := filepath.Join(archiveRoot, fixture.dir)
			if err := os.MkdirAll(dir, 0755); err != nil {
				t.Fatal(err)
			}
			if fixture.schema != "" {
				metadata := "schema: " + fixture.schema + "\ncreated: 2026-05-24\n"
				if err := os.WriteFile(filepath.Join(dir, ".openspec.yaml"), []byte(metadata), 0644); err != nil {
					t.Fatal(err)
				}
			}
			for name, content := range fixture.artifacts {
				if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
					t.Fatal(err)
				}
			}
		}

		changes, err := ListArchiveChangesFrom(root)
		if err != nil {
			t.Fatal(err)
		}
		if len(changes) != 3 {
			t.Fatalf("expected three archives, got %d", len(changes))
		}
		byName := make(map[string]Change, len(changes))
		for _, change := range changes {
			byName[change.Name] = change
		}
		spike := byName["investigate-cache"]
		if spike.Schema != "spike" || spike.DisplayDate != "24/05/2026" {
			t.Errorf("unexpected spike metadata: schema=%q date=%q", spike.Schema, spike.DisplayDate)
		}
		for _, id := range []string{"proposal", "research", "decision"} {
			if _, ok := spike.ArtifactByID(id); !ok {
				t.Errorf("expected archived spike artifact %q", id)
			}
		}
		feature := byName["add-export"]
		if feature.Schema != "feature" {
			t.Errorf("expected feature schema, got %q", feature.Schema)
		}
		unknown := byName["legacy-unknown"]
		if unknown.Schema != "" {
			t.Errorf("expected unknown schema, got %q", unknown.Schema)
		}
		if _, ok := unknown.ArtifactByID("notes"); !ok {
			t.Fatal("expected unknown-schema archive notes")
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
		tasks, ok := reloaded.ArtifactByID("tasks")
		if !ok || len(tasks.Outputs) != 1 {
			t.Fatalf("expected one dynamic tasks output, got %#v", tasks)
		}
		if tasks.Outputs[0].Content != "- [ ] updated task" || !tasks.Outputs[0].Present {
			t.Errorf("expected updated dynamic content, got %#v", tasks.Outputs[0])
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
		proposal, ok := initial.ArtifactByID("proposal")
		if !ok || len(proposal.Outputs) != 1 || !proposal.Outputs[0].Present {
			t.Fatal("expected dynamic proposal output after write")
		}

		if err := os.Remove(filepath.Join(ch.Path, "proposal.md")); err != nil {
			t.Fatal(err)
		}
		afterDelete := ReloadChange(ch)
		if _, ok := afterDelete.ArtifactByID("proposal"); ok {
			t.Error("expected dynamic proposal artifact to be absent after delete")
		}
	})
}
