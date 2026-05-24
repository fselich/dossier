package openspec

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Artifact struct {
	Content string
	Present bool
}

type NamedSpec struct {
	Name    string
	Content string
}

type Change struct {
	Name        string
	Path        string
	Created     string
	DisplayDate string
	Proposal    Artifact
	Design      Artifact
	Tasks       Artifact
	Specs       Artifact
	SpecFiles   []NamedSpec
}

type Project struct {
	Name    string
	Changes []Change
}

type ProjectSpec struct {
	Name             string
	RequirementCount int
	RequirementNames []string
	Content          string
}

type openspecMeta struct {
	Schema  string `yaml:"schema"`
	Created string `yaml:"created"`
}

type ProjectConfig struct {
	Context string
	Rules   map[string][]string
}

type projectConfigYAML struct {
	Context string              `yaml:"context"`
	Rules   map[string][]string `yaml:"rules"`
}

// ── *From(root) variants ──────────────────────────────────────────────────────

func LoadFrom(root string) (*Project, error) {
	openspecDir := filepath.Join(root, "openspec")
	if _, err := os.Stat(openspecDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("no openspec/ directory found in %s", root)
	}

	project := &Project{Name: filepath.Base(root)}

	changesDir := filepath.Join(openspecDir, "changes")
	entries, err := os.ReadDir(changesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return project, nil
		}
		return nil, err
	}

	for _, e := range entries {
		if !e.IsDir() || e.Name() == "archive" {
			continue
		}
		cp := filepath.Join(changesDir, e.Name())
		ch := Change{Name: e.Name(), Path: cp}

		if raw, err := os.ReadFile(filepath.Join(cp, ".openspec.yaml")); err == nil {
			var m openspecMeta
			_ = yaml.Unmarshal(raw, &m)
			ch.Created = m.Created
		}

		ch.Proposal = loadFile(filepath.Join(cp, "proposal.md"))
		ch.Design = loadFile(filepath.Join(cp, "design.md"))
		ch.Tasks = loadFile(filepath.Join(cp, "tasks.md"))
		ch.Specs, ch.SpecFiles = loadSpecs(filepath.Join(cp, "specs"))

		project.Changes = append(project.Changes, ch)
	}
	return project, nil
}

func LoadConfigFrom(root string) (ProjectConfig, error) {
	data, err := os.ReadFile(filepath.Join(root, "openspec", "config.yaml"))
	if err != nil {
		if os.IsNotExist(err) {
			return ProjectConfig{}, nil
		}
		return ProjectConfig{}, nil
	}
	var raw projectConfigYAML
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return ProjectConfig{}, fmt.Errorf("openspec/config.yaml: %w", err)
	}
	return ProjectConfig{Context: strings.TrimSpace(raw.Context), Rules: raw.Rules}, nil
}

func LoadProjectSpecsFrom(root string) ([]ProjectSpec, error) {
	specsDir := filepath.Join(root, "openspec", "specs")
	entries, err := os.ReadDir(specsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var specs []ProjectSpec
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		ps := ProjectSpec{Name: e.Name()}
		if data, err := os.ReadFile(filepath.Join(specsDir, e.Name(), "spec.md")); err == nil {
			ps.Content = string(data)
			for _, line := range strings.Split(ps.Content, "\n") {
				if strings.HasPrefix(line, "### Requirement: ") {
					ps.RequirementCount++
					ps.RequirementNames = append(ps.RequirementNames, strings.TrimPrefix(line, "### Requirement: "))
				}
			}
		}
		specs = append(specs, ps)
	}
	sort.Slice(specs, func(i, j int) bool { return specs[i].Name < specs[j].Name })
	return specs, nil
}

func ListChangeNamesFrom(root string) ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(root, "openspec", "changes"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() && e.Name() != "archive" {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

func ListArchiveChangesFrom(root string) ([]Change, error) {
	archiveDir := filepath.Join(root, "openspec", "changes", "archive")
	entries, err := os.ReadDir(archiveDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	dirs := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e.Name())
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(dirs)))

	var changes []Change
	for _, dir := range dirs {
		cp := filepath.Join(archiveDir, dir)
		cleanName, dispDate := parseArchiveName(dir)
		ch := Change{Name: cleanName, Path: cp, DisplayDate: dispDate}
		if raw, err := os.ReadFile(filepath.Join(cp, ".openspec.yaml")); err == nil {
			var m openspecMeta
			_ = yaml.Unmarshal(raw, &m)
			ch.Created = m.Created
		}
		ch.Proposal = loadFile(filepath.Join(cp, "proposal.md"))
		ch.Design = loadFile(filepath.Join(cp, "design.md"))
		ch.Tasks = loadFile(filepath.Join(cp, "tasks.md"))
		ch.Specs, ch.SpecFiles = loadSpecs(filepath.Join(cp, "specs"))
		changes = append(changes, ch)
	}
	return changes, nil
}

func ListArchiveNamesFrom(root string) ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(root, "openspec", "changes", "archive"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(names)))
	return names, nil
}

func ListSpecNamesFrom(root string) ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(root, "openspec", "specs"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}

// ── Zero-argument wrappers (delegate to *From with os.Getwd()) ─────────────────

func Load() (*Project, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return LoadFrom(cwd)
}

func LoadConfig() (ProjectConfig, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return ProjectConfig{}, err
	}
	return LoadConfigFrom(cwd)
}

func LoadProjectSpecs() ([]ProjectSpec, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return LoadProjectSpecsFrom(cwd)
}

func ListArchiveChanges() ([]Change, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return ListArchiveChangesFrom(cwd)
}

func ListArchiveNames() ([]string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return ListArchiveNamesFrom(cwd)
}

func ListSpecNames() ([]string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return ListSpecNamesFrom(cwd)
}

func ListChangeNames() ([]string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return ListChangeNamesFrom(cwd)
}

// ── Path-based loader (unchanged) ──────────────────────────────────────────────

func LoadFromPath(path string) (*Project, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("path not found: %s", path)
	}
	if _, err := os.Stat(filepath.Join(path, ".openspec.yaml")); os.IsNotExist(err) {
		return nil, fmt.Errorf("not a valid change directory (missing .openspec.yaml): %s", path)
	}

	ch := Change{Name: filepath.Base(path), Path: path}

	if raw, err := os.ReadFile(filepath.Join(path, ".openspec.yaml")); err == nil {
		var m openspecMeta
		_ = yaml.Unmarshal(raw, &m)
		ch.Created = m.Created
	}

	ch.Proposal = loadFile(filepath.Join(path, "proposal.md"))
	ch.Design = loadFile(filepath.Join(path, "design.md"))
	ch.Tasks = loadFile(filepath.Join(path, "tasks.md"))
	ch.Specs, ch.SpecFiles = loadSpecs(filepath.Join(path, "specs"))

	project := &Project{
		Name:    filepath.Base(filepath.Dir(path)),
		Changes: []Change{ch},
	}
	return project, nil
}

// ── Helpers ────────────────────────────────────────────────────────────────────

func parseArchiveName(dir string) (name, date string) {
	if len(dir) > 11 && dir[4] == '-' && dir[7] == '-' && dir[10] == '-' {
		t, err := time.Parse("2006-01-02", dir[:10])
		if err == nil {
			return dir[11:], t.Format("02/01/2006")
		}
	}
	return dir, ""
}

func ReloadChange(ch Change) Change {
	ch.Proposal = loadFile(filepath.Join(ch.Path, "proposal.md"))
	ch.Design = loadFile(filepath.Join(ch.Path, "design.md"))
	ch.Tasks = loadFile(filepath.Join(ch.Path, "tasks.md"))
	ch.Specs, ch.SpecFiles = loadSpecs(filepath.Join(ch.Path, "specs"))
	return ch
}

func loadFile(path string) Artifact {
	data, err := os.ReadFile(path)
	if err != nil {
		return Artifact{}
	}
	return Artifact{Content: string(data), Present: true}
}

func loadSpecs(dir string) (Artifact, []NamedSpec) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return Artifact{}, nil
	}
	var parts []string
	var files []NamedSpec
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name(), "spec.md"))
		if err != nil {
			continue
		}
		content := string(data)
		files = append(files, NamedSpec{Name: e.Name(), Content: content})
		parts = append(parts, "# "+e.Name()+"\n\n"+content)
	}
	if len(parts) == 0 {
		return Artifact{}, nil
	}
	return Artifact{Content: strings.Join(parts, "\n\n---\n\n"), Present: true}, files
}
