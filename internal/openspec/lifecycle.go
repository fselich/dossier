package openspec

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type ArchiveClient interface {
	Archive(root, changeName string) (ArchiveResult, error)
}

type LifecycleKind string

const (
	LifecycleArchive    LifecycleKind = "archive"
	LifecycleReactivate LifecycleKind = "reactivate"
)

type filePreimage struct {
	path   string
	exists bool
	data   []byte
	mode   os.FileMode
}

// LifecycleUndo is an opaque, session-local record for one successful lifecycle operation.
type LifecycleUndo struct {
	Kind       LifecycleKind
	ChangeName string
	source     string
	dest       string
	backupDir  string
	preimages  []filePreimage
	post       map[string]string
}

type LifecycleService struct {
	root   string
	cli    ArchiveClient
	rename func(string, string) error
	copy   func(string, string) error
}

func NewLifecycleService(root string, cli ArchiveClient) *LifecycleService {
	return &LifecycleService{root: filepath.Clean(root), cli: cli, rename: os.Rename, copy: copyDir}
}

func (s *LifecycleService) Reactivate(archivePath string) (*LifecycleUndo, error) {
	archivePath, err := s.containedPath(filepath.Join(s.root, "openspec", "changes", "archive"), archivePath)
	if err != nil {
		return nil, fmt.Errorf("reactivate source outside archive root: %w", err)
	}
	info, err := os.Lstat(archivePath)
	if err != nil {
		return nil, fmt.Errorf("check archived change %s: %w", archivePath, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("archived change path must not be a symlink: %s", archivePath)
	}
	cleanName, err := activeNameFromArchive(filepath.Base(archivePath))
	if err != nil {
		return nil, err
	}
	target := filepath.Join(s.root, "openspec", "changes", cleanName)
	if _, err := os.Lstat(target); err == nil {
		return nil, fmt.Errorf("active target already exists: %s", target)
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("check active target %s: %w", target, err)
	}
	if _, err := os.Stat(archivePath); err != nil {
		return nil, fmt.Errorf("check archived change %s: %w", archivePath, err)
	}
	if err := s.moveDir(archivePath, target); err != nil {
		return nil, fmt.Errorf("reactivate %s: %w", cleanName, err)
	}
	fingerprint, err := fingerprintPath(target)
	if err != nil {
		_ = s.moveDir(target, archivePath)
		return nil, fmt.Errorf("fingerprint reactivated change: %w", err)
	}
	return &LifecycleUndo{
		Kind: LifecycleReactivate, ChangeName: cleanName, source: archivePath, dest: target,
		post: map[string]string{target: fingerprint, archivePath: "missing"},
	}, nil
}

func (s *LifecycleService) Archive(change Change) (*LifecycleUndo, ArchiveResult, error) {
	if s.cli == nil {
		return nil, ArchiveResult{}, fmt.Errorf("OpenSpec CLI is required to archive")
	}
	activeRoot := filepath.Join(s.root, "openspec", "changes")
	activePath, err := s.containedPath(activeRoot, change.Path)
	if err != nil || filepath.Dir(activePath) != filepath.Clean(activeRoot) {
		return nil, ArchiveResult{}, fmt.Errorf("archive source outside active changes root: %s", change.Path)
	}
	backupParent, err := os.MkdirTemp("", "dossier-lifecycle-*")
	if err != nil {
		return nil, ArchiveResult{}, fmt.Errorf("create lifecycle backup: %w", err)
	}
	backup := filepath.Join(backupParent, "change")
	if err := copyDir(activePath, backup); err != nil {
		_ = os.RemoveAll(backupParent)
		return nil, ArchiveResult{}, fmt.Errorf("backup change %s: %w", change.Name, err)
	}
	backupFingerprint, err := fingerprintPath(backup)
	if err != nil {
		_ = os.RemoveAll(backupParent)
		return nil, ArchiveResult{}, fmt.Errorf("fingerprint change backup: %w", err)
	}
	preimages, err := captureAffectedSpecs(s.root, activePath)
	if err != nil {
		_ = os.RemoveAll(backupParent)
		return nil, ArchiveResult{}, err
	}
	if err := writePreimageBackup(backupParent, preimages); err != nil {
		_ = os.RemoveAll(backupParent)
		return nil, ArchiveResult{}, fmt.Errorf("persist project-spec preimages: %w", err)
	}
	beforeArchives, err := archiveEntries(s.root)
	if err != nil {
		_ = os.RemoveAll(backupParent)
		return nil, ArchiveResult{}, err
	}

	result, operationErr := s.cli.Archive(s.root, change.Name)
	if operationErr != nil {
		recoveryErr := s.restoreArchivePrestate(change.Name, "", activePath, backup, backupFingerprint, preimages, beforeArchives)
		if recoveryErr != nil {
			return nil, ArchiveResult{}, fmt.Errorf("archive %s failed: %v; automatic restoration failed: %v; backup retained at %s", change.Name, operationErr, recoveryErr, backupParent)
		}
		_ = os.RemoveAll(backupParent)
		return nil, ArchiveResult{}, operationErr
	}

	archivePath, err := s.resolveArchiveResult(change.Name, result, beforeArchives)
	if err != nil {
		recoveryErr := s.restoreArchivePrestate(change.Name, result.ArchivePath, activePath, backup, backupFingerprint, preimages, beforeArchives)
		if recoveryErr != nil {
			return nil, ArchiveResult{}, fmt.Errorf("archive result invalid: %v; automatic restoration failed: %v; backup retained at %s", err, recoveryErr, backupParent)
		}
		_ = os.RemoveAll(backupParent)
		return nil, ArchiveResult{}, err
	}
	result.ArchivePath = archivePath
	post := make(map[string]string, len(preimages)+2)
	post[activePath] = "missing"
	for _, image := range preimages {
		fingerprint, fingerprintErr := fingerprintPath(image.path)
		if fingerprintErr != nil {
			_ = os.RemoveAll(backupParent)
			return nil, ArchiveResult{}, fmt.Errorf("capture archive post-state %s: %w", image.path, fingerprintErr)
		}
		post[image.path] = fingerprint
	}
	archiveFingerprint, err := fingerprintPath(archivePath)
	if err != nil {
		_ = os.RemoveAll(backupParent)
		return nil, ArchiveResult{}, fmt.Errorf("capture archived change post-state: %w", err)
	}
	post[archivePath] = archiveFingerprint
	return &LifecycleUndo{
		Kind: LifecycleArchive, ChangeName: change.Name, source: activePath, dest: archivePath,
		backupDir: backupParent, preimages: preimages, post: post,
	}, result, nil
}

func (s *LifecycleService) Undo(record *LifecycleUndo) error {
	if record == nil {
		return fmt.Errorf("no lifecycle action to undo")
	}
	paths := make([]string, 0, len(record.post))
	for path := range record.post {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		got, err := fingerprintPath(path)
		if err != nil {
			return fmt.Errorf("check undo path %s: %w", path, err)
		}
		if got != record.post[path] {
			return fmt.Errorf("cannot undo: affected path changed: %s", path)
		}
	}

	switch record.Kind {
	case LifecycleReactivate:
		if err := s.moveDir(record.dest, record.source); err != nil {
			return fmt.Errorf("undo reactivation %s: %w", record.ChangeName, err)
		}
	case LifecycleArchive:
		if record.backupDir == "" {
			return fmt.Errorf("undo archive %s: backup unavailable", record.ChangeName)
		}
		if err := s.copy(filepath.Join(record.backupDir, "change"), record.source); err != nil {
			return fmt.Errorf("restore active change %s: %w", record.source, err)
		}
		if err := restorePreimages(record.preimages); err != nil {
			return fmt.Errorf("restore project specifications: %w", err)
		}
		if err := os.RemoveAll(record.dest); err != nil {
			return fmt.Errorf("remove archived result %s: %w", record.dest, err)
		}
		s.Cleanup(record)
	default:
		return fmt.Errorf("unknown lifecycle undo kind %q", record.Kind)
	}
	return nil
}

func (s *LifecycleService) Cleanup(record *LifecycleUndo) {
	if record != nil && record.backupDir != "" {
		_ = os.RemoveAll(record.backupDir)
		record.backupDir = ""
	}
}

func (s *LifecycleService) moveDir(source, target string) error {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	if err := s.rename(source, target); err == nil {
		return nil
	}
	if err := s.copy(source, target); err != nil {
		_ = os.RemoveAll(target)
		return err
	}
	if err := os.RemoveAll(source); err != nil {
		_ = os.RemoveAll(target)
		return err
	}
	return nil
}

func (s *LifecycleService) restoreArchivePrestate(changeName, explicitArchivePath, activePath, backup, backupFingerprint string, preimages []filePreimage, before map[string]bool) error {
	var failures []string
	activeRestored := true
	activeFingerprint, fingerprintErr := fingerprintPath(activePath)
	if fingerprintErr != nil || activeFingerprint != backupFingerprint {
		if err := os.RemoveAll(activePath); err != nil {
			activeRestored = false
			failures = append(failures, fmt.Sprintf("%s: %v", activePath, err))
		} else if err := s.copy(backup, activePath); err != nil {
			activeRestored = false
			failures = append(failures, fmt.Sprintf("%s: %v", activePath, err))
		}
	}
	if err := restoreChangedPreimages(preimages); err != nil {
		failures = append(failures, err.Error())
	}
	after, err := archiveEntries(s.root)
	if err != nil {
		failures = append(failures, err.Error())
	} else if activeRestored {
		explicit := explicitArchivePath
		if explicit != "" && !filepath.IsAbs(explicit) {
			explicit = filepath.Join(s.root, explicit)
		}
		explicit = filepath.Clean(explicit)
		for path := range after {
			if before[path] {
				continue
			}
			cleanName, parseErr := activeNameFromArchive(filepath.Base(path))
			selectedDestination := parseErr == nil && cleanName == changeName
			if !selectedDestination && explicit != "" && filepath.Clean(path) == explicit {
				fingerprint, fingerprintErr := fingerprintPath(path)
				selectedDestination = fingerprintErr == nil && fingerprint == backupFingerprint
			}
			if !selectedDestination {
				continue
			}
			if err := os.RemoveAll(path); err != nil {
				failures = append(failures, fmt.Sprintf("%s: %v", path, err))
			}
		}
	}
	if len(failures) > 0 {
		return fmt.Errorf("manual recovery may be required for: %s", strings.Join(failures, "; "))
	}
	return nil
}

func (s *LifecycleService) resolveArchiveResult(name string, result ArchiveResult, before map[string]bool) (string, error) {
	if result.ArchivePath != "" {
		path := result.ArchivePath
		if !filepath.IsAbs(path) {
			path = filepath.Join(s.root, path)
		}
		contained, err := s.containedPath(filepath.Join(s.root, "openspec", "changes", "archive"), path)
		if err != nil {
			return "", fmt.Errorf("archive destination outside archive root: %w", err)
		}
		if _, err := os.Stat(contained); err != nil {
			return "", fmt.Errorf("archive destination %s: %w", contained, err)
		}
		cleanName, parseErr := activeNameFromArchive(filepath.Base(contained))
		if parseErr != nil || cleanName != name {
			return "", fmt.Errorf("archive destination %s does not match change %s", contained, name)
		}
		return contained, nil
	}
	after, err := archiveEntries(s.root)
	if err != nil {
		return "", err
	}
	var candidates []string
	for path := range after {
		if !before[path] {
			clean, parseErr := activeNameFromArchive(filepath.Base(path))
			if parseErr == nil && clean == name {
				candidates = append(candidates, path)
			}
		}
	}
	if len(candidates) != 1 {
		return "", fmt.Errorf("cannot identify archive destination for %s", name)
	}
	return candidates[0], nil
}

func (s *LifecycleService) containedPath(root, path string) (string, error) {
	absoluteRoot, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	relative, err := filepath.Rel(absoluteRoot, absolutePath)
	if err != nil || !isContainedRelativePath(relative) {
		return "", fmt.Errorf("%s is outside %s", path, root)
	}
	return filepath.Clean(absolutePath), nil
}

func activeNameFromArchive(base string) (string, error) {
	name := base
	if len(base) > 11 && base[10] == '-' {
		if _, err := time.Parse("2006-01-02", base[:10]); err == nil {
			name = base[11:]
		}
	}
	if name == "" || name == "." || name == ".." || strings.ContainsAny(name, `/\\`) {
		return "", fmt.Errorf("invalid archived change name %q", name)
	}
	return name, nil
}

func captureAffectedSpecs(root, changePath string) ([]filePreimage, error) {
	deltaRoot := filepath.Join(changePath, "specs")
	entries, err := os.ReadDir(deltaRoot)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("discover affected specifications: %w", err)
	}
	var images []filePreimage
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		path := filepath.Join(root, "openspec", "specs", entry.Name(), "spec.md")
		image := filePreimage{path: path, mode: 0o644}
		info, statErr := os.Stat(path)
		switch {
		case statErr == nil:
			image.exists = true
			image.mode = info.Mode().Perm()
			image.data, err = os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("snapshot project specification %s: %w", path, err)
			}
		case os.IsNotExist(statErr):
		default:
			return nil, fmt.Errorf("snapshot project specification %s: %w", path, statErr)
		}
		images = append(images, image)
	}
	sort.Slice(images, func(i, j int) bool { return images[i].path < images[j].path })
	return images, nil
}

func writePreimageBackup(backupParent string, images []filePreimage) error {
	root := filepath.Join(backupParent, "project-specs")
	if err := os.MkdirAll(root, 0o700); err != nil {
		return err
	}
	var manifest strings.Builder
	for _, image := range images {
		state := "absent"
		if image.exists {
			state = "present"
			name := filepath.Base(filepath.Dir(image.path))
			path := filepath.Join(root, name, "spec.md")
			if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
				return err
			}
			if err := os.WriteFile(path, image.data, 0o600); err != nil {
				return err
			}
		}
		_, _ = fmt.Fprintf(&manifest, "%s\t%s\n", state, image.path)
	}
	return os.WriteFile(filepath.Join(root, "MANIFEST.txt"), []byte(manifest.String()), 0o600)
}

func restoreChangedPreimages(images []filePreimage) error {
	var changed []filePreimage
	for _, image := range images {
		matches, err := matchesPreimage(image)
		if err != nil || !matches {
			changed = append(changed, image)
		}
	}
	return restorePreimages(changed)
}

func matchesPreimage(image filePreimage) (bool, error) {
	info, err := os.Stat(image.path)
	if os.IsNotExist(err) {
		return !image.exists, nil
	}
	if err != nil {
		return false, err
	}
	if !image.exists || info.Mode().Perm() != image.mode {
		return false, nil
	}
	data, err := os.ReadFile(image.path)
	if err != nil {
		return false, err
	}
	return string(data) == string(image.data), nil
}

func restorePreimages(images []filePreimage) error {
	for _, image := range images {
		if !image.exists {
			if err := os.Remove(image.path); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("remove newly created %s: %w", image.path, err)
			}
			// Archive creates one capability directory for a previously absent spec.
			// Remove it only when empty so unrelated external files are preserved.
			_ = os.Remove(filepath.Dir(image.path))
			continue
		}
		if err := os.MkdirAll(filepath.Dir(image.path), 0o755); err != nil {
			return fmt.Errorf("prepare %s: %w", image.path, err)
		}
		if err := os.WriteFile(image.path, image.data, image.mode); err != nil {
			return fmt.Errorf("restore %s: %w", image.path, err)
		}
	}
	return nil
}

func archiveEntries(root string) (map[string]bool, error) {
	dir := filepath.Join(root, "openspec", "changes", "archive")
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return map[string]bool{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read archive directory: %w", err)
	}
	result := make(map[string]bool)
	for _, entry := range entries {
		if entry.IsDir() {
			result[filepath.Join(dir, entry.Name())] = true
		}
	}
	return result, nil
}

func fingerprintPath(path string) (string, error) {
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return "missing", nil
	}
	if err != nil {
		return "", err
	}
	hash := sha256.New()
	if !info.IsDir() {
		if err := fingerprintEntry(hash, path, filepath.Base(path), info); err != nil {
			return "", err
		}
		return hex.EncodeToString(hash.Sum(nil)), nil
	}
	var paths []string
	err = filepath.Walk(path, func(current string, currentInfo os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		paths = append(paths, current)
		return nil
	})
	if err != nil {
		return "", err
	}
	sort.Strings(paths)
	for _, current := range paths {
		currentInfo, err := os.Lstat(current)
		if err != nil {
			return "", err
		}
		relative, _ := filepath.Rel(path, current)
		if err := fingerprintEntry(hash, current, filepath.ToSlash(relative), currentInfo); err != nil {
			return "", err
		}
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func fingerprintEntry(hash io.Writer, path, relative string, info os.FileInfo) error {
	_, _ = fmt.Fprintf(hash, "%s\x00%d\x00", relative, info.Mode())
	if info.Mode().IsRegular() {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		_, copyErr := io.Copy(hash, file)
		closeErr := file.Close()
		if copyErr != nil {
			return copyErr
		}
		return closeErr
	}
	if info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(path)
		if err != nil {
			return err
		}
		_, _ = io.WriteString(hash, target)
	}
	return nil
}

func copyDir(source, target string) error {
	if _, err := os.Lstat(target); err == nil {
		return fmt.Errorf("destination already exists: %s", target)
	} else if !os.IsNotExist(err) {
		return err
	}
	return filepath.Walk(source, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		destination := filepath.Join(target, relative)
		switch {
		case info.IsDir():
			return os.MkdirAll(destination, info.Mode().Perm())
		case info.Mode()&os.ModeSymlink != 0:
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			return os.Symlink(link, destination)
		case info.Mode().IsRegular():
			input, err := os.Open(path)
			if err != nil {
				return err
			}
			output, err := os.OpenFile(destination, os.O_CREATE|os.O_EXCL|os.O_WRONLY, info.Mode().Perm())
			if err != nil {
				_ = input.Close()
				return err
			}
			_, copyErr := io.Copy(output, input)
			inputErr := input.Close()
			outputErr := output.Close()
			if copyErr != nil {
				return copyErr
			}
			if inputErr != nil {
				return inputErr
			}
			return outputErr
		default:
			return fmt.Errorf("unsupported file type: %s", path)
		}
	})
}
