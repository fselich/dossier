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
	DisplayDate string // "DD Mon" for archive entries, empty for active changes
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

// LoadProjectSpecs reads openspec/specs/ and returns one ProjectSpec per subdirectory,
// sorted alphabetically by name.
func LoadProjectSpecs() []ProjectSpec {
	cwd, err := os.Getwd()
	if err != nil {
		return nil
	}
	specsDir := filepath.Join(cwd, "openspec", "specs")
	entries, err := os.ReadDir(specsDir)
	if err != nil {
		return nil
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
	return specs
}

type openspecMeta struct {
	Schema  string `yaml:"schema"`
	Created string `yaml:"created"`
}

func Load() (*Project, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	openspecDir := filepath.Join(cwd, "openspec")
	if _, err := os.Stat(openspecDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("no openspec/ directory found in current directory")
	}

	project := &Project{Name: filepath.Base(cwd)}

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

// LoadFromPath loads a single change from an explicit directory path.
// The directory must exist and contain a .openspec.yaml file.
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

// parseArchiveName splits a directory name of the form "YYYY-MM-DD-<name>"
// into a clean name and a formatted date string ("dd/mm/yyyy").
// If the format does not match, it returns the full dir name and an empty date.
func parseArchiveName(dir string) (name, date string) {
	if len(dir) > 11 && dir[4] == '-' && dir[7] == '-' && dir[10] == '-' {
		t, err := time.Parse("2006-01-02", dir[:10])
		if err == nil {
			return dir[11:], t.Format("02/01/2006")
		}
	}
	return dir, ""
}

// ListArchiveChanges loads all changes from openspec/changes/archive/, most recent first.
func ListArchiveChanges() []Change {
	cwd, err := os.Getwd()
	if err != nil {
		return nil
	}
	archiveDir := filepath.Join(cwd, "openspec", "changes", "archive")
	entries, err := os.ReadDir(archiveDir)
	if err != nil {
		return nil
	}

	// Collect dir names and sort descending (most recent first by name prefix).
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
	return changes
}

// ListArchiveNames returns the names of archived change directories on disk,
// sorted descending to match ListArchiveChanges() order.
// It only reads the directory listing — no artifact files — so it is cheap
// to call on every poll tick.
func ListArchiveNames() []string {
	cwd, err := os.Getwd()
	if err != nil {
		return nil
	}
	entries, err := os.ReadDir(filepath.Join(cwd, "openspec", "changes", "archive"))
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(names)))
	return names
}

// ListSpecNames returns the names of project spec directories on disk,
// sorted alphabetically to match LoadProjectSpecs() order.
// It only reads the directory listing — no artifact files — so it is cheap
// to call on every poll tick.
func ListSpecNames() []string {
	cwd, err := os.Getwd()
	if err != nil {
		return nil
	}
	entries, err := os.ReadDir(filepath.Join(cwd, "openspec", "specs"))
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	return names
}

// ListChangeNames returns the names of active change directories on disk.
// It only reads the directory listing — no artifact files — so it is cheap
// to call on every poll tick.
func ListChangeNames() []string {
	cwd, err := os.Getwd()
	if err != nil {
		return nil
	}
	entries, err := os.ReadDir(filepath.Join(cwd, "openspec", "changes"))
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() && e.Name() != "archive" {
			names = append(names, e.Name())
		}
	}
	return names
}

// ReloadChange rereads all artifact files for an existing Change from disk.
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
