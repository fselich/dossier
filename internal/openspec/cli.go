package openspec

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const openSpecTimeout = 10 * time.Second

type OpenSpecRunner interface {
	Run(dir string, args ...string) (stdout, stderr []byte, err error)
}

type OSOpenSpecRunner struct{}

func (OSOpenSpecRunner) Run(dir string, args ...string) ([]byte, []byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), openSpecTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "openspec", args...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if ctx.Err() != nil {
		err = fmt.Errorf("openspec command timed out after %s: %w", openSpecTimeout, ctx.Err())
	}
	return stdout.Bytes(), stderr.Bytes(), err
}

type CLI struct {
	runner OpenSpecRunner
}

func NewCLI(runner OpenSpecRunner) *CLI {
	return &CLI{runner: runner}
}

func NewOSCLI() *CLI {
	return NewCLI(OSOpenSpecRunner{})
}

type ChangeListItem struct {
	Name           string `json:"name"`
	CompletedTasks int    `json:"completedTasks"`
	TotalTasks     int    `json:"totalTasks"`
	LastModified   string `json:"lastModified"`
	Status         string `json:"status"`
}

type SchemaInfo struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Artifacts   []string `json:"artifacts"`
	Source      string   `json:"source"`
}

type ArtifactPathStatus struct {
	OutputPath          string   `json:"outputPath"`
	ResolvedOutputPath  string   `json:"resolvedOutputPath"`
	ExistingOutputPaths []string `json:"existingOutputPaths"`
}

type StatusArtifact struct {
	ID         string   `json:"id"`
	OutputPath string   `json:"outputPath"`
	Status     string   `json:"status"`
	Requires   []string `json:"requires"`
}

type CreateChangeResult struct {
	Name   string
	Schema string
	Path   string
}

type ValidationResult struct {
	Valid   bool
	Summary string
	Errors  []string
}

type ArchiveResult struct {
	ChangeName   string
	ArchivedAs   string
	ArchivePath  string
	SpecsUpdated bool
}

type ChangeStatus struct {
	ChangeName    string                        `json:"changeName"`
	SchemaName    string                        `json:"schemaName"`
	ChangeRoot    string                        `json:"changeRoot"`
	ArtifactPaths map[string]ArtifactPathStatus `json:"artifactPaths"`
	IsComplete    bool                          `json:"isComplete"`
	ApplyRequires []string                      `json:"applyRequires"`
	Artifacts     []StatusArtifact              `json:"artifacts"`
}

type listChangesResponse struct {
	Changes []ChangeListItem `json:"changes"`
}

func (c *CLI) ListChanges(root string) ([]ChangeListItem, error) {
	args := []string{"list", "--json"}
	stdout, err := c.run(root, args)
	if err != nil {
		return nil, err
	}
	var response listChangesResponse
	if err := json.Unmarshal(stdout, &response); err != nil {
		return nil, fmt.Errorf("decode openspec %s: %w", strings.Join(args, " "), err)
	}
	return response.Changes, nil
}

func (c *CLI) Schemas(root string) ([]SchemaInfo, error) {
	args := []string{"schemas", "--json"}
	stdout, err := c.run(root, args)
	if err != nil {
		return nil, err
	}
	var schemas []SchemaInfo
	if err := json.Unmarshal(stdout, &schemas); err != nil {
		return nil, fmt.Errorf("decode openspec %s: %w", strings.Join(args, " "), err)
	}
	return schemas, nil
}

func (c *CLI) CreateChange(root, name, schema string) (CreateChangeResult, error) {
	args := []string{"new", "change", name, "--schema", schema, "--json"}
	var response struct {
		Change struct {
			ID     string `json:"id"`
			Path   string `json:"path"`
			Schema string `json:"schema"`
		} `json:"change"`
	}
	if err := c.runJSON(root, args, &response); err != nil {
		return CreateChangeResult{}, err
	}
	if response.Change.ID == "" {
		return CreateChangeResult{}, fmt.Errorf("decode openspec %s: missing change.id", strings.Join(args, " "))
	}
	return CreateChangeResult{Name: response.Change.ID, Schema: response.Change.Schema, Path: response.Change.Path}, nil
}

func (c *CLI) Validate(root, name, kind string) (ValidationResult, error) {
	if kind != "change" && kind != "spec" {
		return ValidationResult{}, fmt.Errorf("invalid OpenSpec validation type %q", kind)
	}
	args := []string{"validate", name, "--type", kind, "--json", "--no-interactive"}
	var response struct {
		Items []struct {
			ID     string `json:"id"`
			Type   string `json:"type"`
			Valid  bool   `json:"valid"`
			Issues []struct {
				Level   string `json:"level"`
				Path    string `json:"path"`
				Message string `json:"message"`
			} `json:"issues"`
		} `json:"items"`
	}
	if err := c.runJSON(root, args, &response); err != nil {
		return ValidationResult{}, err
	}
	if len(response.Items) != 1 {
		return ValidationResult{}, fmt.Errorf("decode openspec %s: expected one validation item, got %d", strings.Join(args, " "), len(response.Items))
	}
	item := response.Items[0]
	result := ValidationResult{Valid: item.Valid}
	for _, issue := range item.Issues {
		message := issue.Message
		if issue.Path != "" {
			message = issue.Path + ": " + message
		}
		if issue.Level != "" {
			message = issue.Level + " " + message
		}
		result.Errors = append(result.Errors, message)
	}
	if result.Valid {
		result.Summary = fmt.Sprintf("%s %s is valid", kind, name)
	} else {
		result.Summary = fmt.Sprintf("%s %s has %d validation issue(s)", kind, name, len(result.Errors))
	}
	return result, nil
}

func (c *CLI) Archive(root, changeName string) (ArchiveResult, error) {
	args := []string{"archive", changeName, "--yes", "--json"}
	var response struct {
		Archive struct {
			Change       string `json:"change"`
			ArchivedAs   string `json:"archivedAs"`
			Path         string `json:"path"`
			SpecsUpdated bool   `json:"specsUpdated"`
		} `json:"archive"`
	}
	if err := c.runJSON(root, args, &response); err != nil {
		return ArchiveResult{}, err
	}
	if response.Archive.Path == "" {
		return ArchiveResult{}, fmt.Errorf("decode openspec %s: missing archive.path", strings.Join(args, " "))
	}
	return ArchiveResult{ChangeName: response.Archive.Change, ArchivedAs: response.Archive.ArchivedAs, ArchivePath: response.Archive.Path, SpecsUpdated: response.Archive.SpecsUpdated}, nil
}

func (c *CLI) Status(root, changeName string) (ChangeStatus, error) {
	args := []string{"status", "--change", changeName, "--json"}
	stdout, err := c.run(root, args)
	if err != nil {
		return ChangeStatus{}, err
	}
	var status ChangeStatus
	if err := json.Unmarshal(stdout, &status); err != nil {
		return ChangeStatus{}, fmt.Errorf("decode openspec %s: %w", strings.Join(args, " "), err)
	}
	return status, nil
}

// StatusAll fetches status for every active change in a single CLI
// invocation (`openspec status --all --json`) — one process spawn instead
// of one per change, which dominates startup time for multi-change repos.
func (c *CLI) StatusAll(root string) ([]ChangeStatus, error) {
	args := []string{"status", "--all", "--json"}
	stdout, err := c.run(root, args)
	if err != nil {
		return nil, err
	}
	var response struct {
		Changes []ChangeStatus `json:"changes"`
	}
	if err := json.Unmarshal(stdout, &response); err != nil {
		return nil, fmt.Errorf("decode openspec %s: %w", strings.Join(args, " "), err)
	}
	return response.Changes, nil
}

func (c *CLI) runJSON(root string, args []string, target any) error {
	stdout, err := c.run(root, args)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(stdout, target); err != nil {
		return fmt.Errorf("decode openspec %s: %w", strings.Join(args, " "), err)
	}
	return nil
}

func (c *CLI) run(root string, args []string) ([]byte, error) {
	stdout, stderr, err := c.runner.Run(root, args...)
	if err == nil {
		return stdout, nil
	}
	message := strings.TrimSpace(string(stderr))
	if message == "" {
		message = err.Error()
	}
	return nil, fmt.Errorf("openspec %s: %s: %w", strings.Join(args, " "), message, err)
}

func (l *Loader) EnrichChange(ch Change, status ChangeStatus) Change {
	fallback := append([]ChangeArtifact(nil), ch.Artifacts...)
	ch.StatusEnriched = true
	fallbackByID := make(map[string]ChangeArtifact, len(fallback))
	for _, artifact := range fallback {
		fallbackByID[artifact.ID] = artifact
	}

	if status.ChangeName != "" && status.ChangeName != ch.Name {
		appendChangeDiagnostic(&ch, fmt.Sprintf("status change name %q does not match %q", status.ChangeName, ch.Name))
	}
	if status.ChangeRoot != "" && !sameCleanPath(status.ChangeRoot, ch.Path) {
		appendChangeDiagnostic(&ch, fmt.Sprintf("status change root %q does not match %q", status.ChangeRoot, ch.Path))
	}
	if status.SchemaName != "" {
		ch.Schema = status.SchemaName
	}

	ch.Artifacts = make([]ChangeArtifact, 0, len(status.Artifacts)+len(fallback))
	statusIndexByID := make(map[string]int, len(status.Artifacts))
	claimed := make(map[string]bool)
	for _, raw := range status.Artifacts {
		pathStatus := status.ArtifactPaths[raw.ID]
		outputPattern := raw.OutputPath
		if outputPattern == "" {
			outputPattern = pathStatus.OutputPath
		}
		artifact := ChangeArtifact{
			ID:            raw.ID,
			OutputPattern: outputPattern,
			Status:        parseArtifactStatus(raw.Status, &ch),
			Requires:      append([]string(nil), raw.Requires...),
			Source:        ArtifactSourceStatus,
		}
		for _, outputPath := range pathStatus.ExistingOutputPaths {
			relativePath, err := safeRelativeOutputPath(ch.Path, outputPath)
			if err != nil {
				appendChangeDiagnostic(&ch, fmt.Sprintf("artifact %s output %q: %v", raw.ID, outputPath, err))
				continue
			}
			content, err := l.fs.ReadFile(filepath.Join(ch.Path, filepath.FromSlash(relativePath)))
			if err != nil {
				if fallbackOutput, ok := findFallbackOutput(fallbackByID[raw.ID], relativePath); ok {
					artifact.Outputs = append(artifact.Outputs, fallbackOutput)
					claimed[relativePath] = true
					continue
				}
				appendChangeDiagnostic(&ch, fmt.Sprintf("read status output %s: %v", relativePath, err))
				continue
			}
			artifact.Outputs = append(artifact.Outputs, ArtifactOutput{
				RelativePath: relativePath,
				DisplayName:  discoveredOutputName(relativePath),
				Content:      string(content),
				Present:      true,
			})
			claimed[relativePath] = true
		}
		statusIndexByID[raw.ID] = len(ch.Artifacts)
		ch.Artifacts = append(ch.Artifacts, artifact)
	}
	// Custom schemas generate outputs whose filename stems differ from the
	// schema artifact IDs (e.g. artifact "idea" generating "1-idea.md"), so
	// discovery-derived artifacts must never re-show a file a status artifact
	// already carries: filter fallback outputs to unclaimed paths, donate
	// same-ID leftovers to their status artifact, and drop emptied artifacts.
	for _, artifact := range fallback {
		unclaimed := unclaimedOutputs(artifact.Outputs, claimed)
		if len(unclaimed) == 0 {
			continue
		}
		for _, output := range unclaimed {
			claimed[output.RelativePath] = true
		}
		if index, ok := statusIndexByID[artifact.ID]; ok {
			ch.Artifacts[index].Outputs = append(ch.Artifacts[index].Outputs, unclaimed...)
			continue
		}
		artifact.Outputs = unclaimed
		ch.Artifacts = append(ch.Artifacts, artifact)
	}
	return ch
}

func unclaimedOutputs(outputs []ArtifactOutput, claimed map[string]bool) []ArtifactOutput {
	var remaining []ArtifactOutput
	for _, output := range outputs {
		if claimed[output.RelativePath] {
			continue
		}
		remaining = append(remaining, output)
	}
	return remaining
}

func parseArtifactStatus(status string, ch *Change) ArtifactStatus {
	switch ArtifactStatus(status) {
	case ArtifactStatusBlocked, ArtifactStatusReady, ArtifactStatusDone:
		return ArtifactStatus(status)
	case ArtifactStatusUnknown:
		return ArtifactStatusUnknown
	default:
		appendChangeDiagnostic(ch, fmt.Sprintf("unknown artifact status %q", status))
		return ArtifactStatusUnknown
	}
}

func safeRelativeOutputPath(root, path string) (string, error) {
	absoluteRoot, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("resolve change root: %w", err)
	}
	absolutePath := path
	if !filepath.IsAbs(absolutePath) {
		absolutePath = filepath.Join(absoluteRoot, absolutePath)
	}
	absolutePath, err = filepath.Abs(absolutePath)
	if err != nil {
		return "", fmt.Errorf("resolve output: %w", err)
	}
	relativePath, err := filepath.Rel(absoluteRoot, absolutePath)
	if err != nil {
		return "", fmt.Errorf("resolve output relative path: %w", err)
	}
	if !isContainedRelativePath(relativePath) {
		return "", fmt.Errorf("outside change root")
	}

	evaluatedRoot, err := filepath.EvalSymlinks(absoluteRoot)
	if err != nil {
		return "", fmt.Errorf("resolve change root symlinks: %w", err)
	}
	evaluatedPath, err := filepath.EvalSymlinks(absolutePath)
	if err != nil {
		return "", fmt.Errorf("resolve output symlinks: %w", err)
	}
	evaluatedRelative, err := filepath.Rel(evaluatedRoot, evaluatedPath)
	if err != nil {
		return "", fmt.Errorf("resolve evaluated output relative path: %w", err)
	}
	if !isContainedRelativePath(evaluatedRelative) {
		return "", fmt.Errorf("outside change root after resolving symlinks")
	}
	return filepath.ToSlash(relativePath), nil
}

func isContainedRelativePath(path string) bool {
	return path != ".." && !strings.HasPrefix(path, ".."+string(filepath.Separator)) && !filepath.IsAbs(path)
}

func sameCleanPath(a, b string) bool {
	absoluteA, errA := filepath.Abs(a)
	absoluteB, errB := filepath.Abs(b)
	if errA != nil || errB != nil {
		return filepath.Clean(a) == filepath.Clean(b)
	}
	return filepath.Clean(absoluteA) == filepath.Clean(absoluteB)
}

func findFallbackOutput(artifact ChangeArtifact, relativePath string) (ArtifactOutput, bool) {
	for _, output := range artifact.Outputs {
		if output.RelativePath == relativePath {
			return output, true
		}
	}
	return ArtifactOutput{}, false
}
