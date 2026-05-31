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

type Loader struct {
	fs fileSystem
}

func NewLoader(fs fileSystem) *Loader {
	return &Loader{fs: fs}
}

var defaultLoader = NewLoader(OSFS{})

// ── *From(root) variants ──────────────────────────────────────────────────────

func (l *Loader) LoadFrom(root string) (*Project, error) {
	openspecDir := filepath.Join(root, "openspec")
	if _, err := l.fs.Stat(openspecDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("no openspec/ directory found in %s", root)
	}

	project := &Project{Name: filepath.Base(root)}

	changesDir := filepath.Join(openspecDir, "changes")
	entries, err := l.fs.ReadDir(changesDir)
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
		ch := l.loadChangeFromDir(filepath.Join(changesDir, e.Name()), e.Name(), "")
		project.Changes = append(project.Changes, ch)
	}

	sort.SliceStable(project.Changes, func(i, j int) bool {
		a, b := project.Changes[i].Created, project.Changes[j].Created
		switch {
		case a == "" && b == "":
			return project.Changes[i].Name < project.Changes[j].Name
		case a == "":
			return false
		case b == "":
			return true
		default:
			return a > b
		}
	})

	return project, nil
}

func (l *Loader) LoadConfigFrom(root string) (ProjectConfig, error) {
	data, err := l.fs.ReadFile(filepath.Join(root, "openspec", "config.yaml"))
	if err != nil {
		if os.IsNotExist(err) {
			return ProjectConfig{}, nil
		}
		return ProjectConfig{}, err
	}
	var raw projectConfigYAML
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return ProjectConfig{}, fmt.Errorf("openspec/config.yaml: %w", err)
	}
	return ProjectConfig{Context: strings.TrimSpace(raw.Context), Rules: raw.Rules}, nil
}

func (l *Loader) LoadProjectSpecsFrom(root string) ([]ProjectSpec, error) {
	specsDir := filepath.Join(root, "openspec", "specs")
	entries, err := l.fs.ReadDir(specsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	specs := make([]ProjectSpec, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		ps := ProjectSpec{Name: e.Name()}
		if data, err := l.fs.ReadFile(filepath.Join(specsDir, e.Name(), "spec.md")); err == nil {
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

func (l *Loader) ListChangeNamesFrom(root string) ([]string, error) {
	entries, err := l.fs.ReadDir(filepath.Join(root, "openspec", "changes"))
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

func (l *Loader) ListArchiveChangesFrom(root string) ([]Change, error) {
	archiveDir := filepath.Join(root, "openspec", "changes", "archive")
	entries, err := l.fs.ReadDir(archiveDir)
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
		cleanName, dispDate := parseArchiveName(dir)
		ch := l.loadChangeFromDir(filepath.Join(archiveDir, dir), cleanName, dispDate)
		changes = append(changes, ch)
	}
	return changes, nil
}

func (l *Loader) ListArchiveNamesFrom(root string) ([]string, error) {
	entries, err := l.fs.ReadDir(filepath.Join(root, "openspec", "changes", "archive"))
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

func (l *Loader) ListSpecNamesFrom(root string) ([]string, error) {
	entries, err := l.fs.ReadDir(filepath.Join(root, "openspec", "specs"))
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

// ── Path-based loader ─────────────────────────────────────────────────────────

func (l *Loader) LoadFromPath(path string) (*Project, error) {
	if _, err := l.fs.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("path not found: %s", path)
	}
	if _, err := l.fs.Stat(filepath.Join(path, ".openspec.yaml")); os.IsNotExist(err) {
		return nil, fmt.Errorf("not a valid change directory (missing .openspec.yaml): %s", path)
	}

	ch := l.loadChangeFromDir(path, filepath.Base(path), "")

	project := &Project{
		Name:    filepath.Base(filepath.Dir(path)),
		Changes: []Change{ch},
	}
	return project, nil
}

func (l *Loader) ReloadChange(ch Change) Change {
	ch.Proposal = l.loadFile(filepath.Join(ch.Path, "proposal.md"))
	ch.Design = l.loadFile(filepath.Join(ch.Path, "design.md"))
	ch.Tasks = l.loadFile(filepath.Join(ch.Path, "tasks.md"))
	ch.Specs, ch.SpecFiles = l.loadSpecs(filepath.Join(ch.Path, "specs"))
	return ch
}

// ── Helpers ────────────────────────────────────────────────────────────────────

func (l *Loader) loadChangeFromDir(dir, name, displayDate string) Change {
	ch := Change{Name: name, Path: dir, DisplayDate: displayDate}
	if raw, err := l.fs.ReadFile(filepath.Join(dir, ".openspec.yaml")); err == nil {
		var m openspecMeta
		// Ignore unmarshal errors: .openspec.yaml is optional metadata,
		// missing or malformed fields are non-fatal.
		_ = yaml.Unmarshal(raw, &m)
		ch.Created = m.Created
	}
	ch.Proposal = l.loadFile(filepath.Join(dir, "proposal.md"))
	ch.Design = l.loadFile(filepath.Join(dir, "design.md"))
	ch.Tasks = l.loadFile(filepath.Join(dir, "tasks.md"))
	ch.Specs, ch.SpecFiles = l.loadSpecs(filepath.Join(dir, "specs"))
	return ch
}

func (l *Loader) loadFile(path string) Artifact {
	data, err := l.fs.ReadFile(path)
	if err != nil {
		return Artifact{}
	}
	return Artifact{Content: string(data), Present: true}
}

func (l *Loader) loadSpecs(dir string) (Artifact, []NamedSpec) {
	entries, err := l.fs.ReadDir(dir)
	if err != nil {
		return Artifact{}, nil
	}
	var parts []string
	var files []NamedSpec
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		data, err := l.fs.ReadFile(filepath.Join(dir, e.Name(), "spec.md"))
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

func parseArchiveName(dir string) (name, date string) {
	if len(dir) > 11 && dir[4] == '-' && dir[7] == '-' && dir[10] == '-' {
		t, err := time.Parse("2006-01-02", dir[:10])
		if err == nil {
			return dir[11:], t.Format("02/01/2006")
		}
	}
	return dir, ""
}

func ExtractRequirement(raw, name string) string {
	target := "### Requirement: " + name
	lines := strings.Split(raw, "\n")
	start := -1
	for i, l := range lines {
		if l == target {
			start = i
			break
		}
	}
	if start < 0 {
		return ""
	}
	block := []string{lines[start]}
	for _, l := range lines[start+1:] {
		if strings.HasPrefix(l, "### Requirement: ") {
			break
		}
		block = append(block, l)
	}
	return strings.Join(block, "\n")
}

func ConfigToMarkdown(cfg ProjectConfig) string {
	var sb strings.Builder
	if cfg.Context != "" {
		sb.WriteString("## Context\n\n")
		sb.WriteString(cfg.Context)
		sb.WriteString("\n")
	}
	if len(cfg.Rules) > 0 {
		if cfg.Context != "" {
			sb.WriteString("\n")
		}
		sb.WriteString("## Rules\n")
		keys := make([]string, 0, len(cfg.Rules))
		for k := range cfg.Rules {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			sb.WriteString("\n### ")
			sb.WriteString(k)
			sb.WriteString("\n\n")
			for _, item := range cfg.Rules[k] {
				sb.WriteString("- ")
				sb.WriteString(item)
				sb.WriteString("\n")
			}
		}
	}
	return sb.String()
}

// ── Backward-compatible package-level wrappers ─────────────────────────────────

func LoadFrom(root string) (*Project, error) {
	return defaultLoader.LoadFrom(root)
}

func LoadConfigFrom(root string) (ProjectConfig, error) {
	return defaultLoader.LoadConfigFrom(root)
}

func LoadProjectSpecsFrom(root string) ([]ProjectSpec, error) {
	return defaultLoader.LoadProjectSpecsFrom(root)
}

func ListChangeNamesFrom(root string) ([]string, error) {
	return defaultLoader.ListChangeNamesFrom(root)
}

func ListArchiveChangesFrom(root string) ([]Change, error) {
	return defaultLoader.ListArchiveChangesFrom(root)
}

func ListArchiveNamesFrom(root string) ([]string, error) {
	return defaultLoader.ListArchiveNamesFrom(root)
}

func ListSpecNamesFrom(root string) ([]string, error) {
	return defaultLoader.ListSpecNamesFrom(root)
}

func LoadFromPath(path string) (*Project, error) {
	return defaultLoader.LoadFromPath(path)
}

func ReloadChange(ch Change) Change {
	return defaultLoader.ReloadChange(ch)
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
