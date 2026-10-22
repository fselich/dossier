package openspec

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fakeArchiveClient struct {
	result ArchiveResult
	err    error
	run    func()
}

func (f *fakeArchiveClient) Archive(_ string, _ string) (ArchiveResult, error) {
	if f.run != nil {
		f.run()
	}
	return f.result, f.err
}

func writeLifecycleFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLifecycleReactivateAndUndo(t *testing.T) {
	root := t.TempDir()
	archive := filepath.Join(root, "openspec", "changes", "archive", "2026-08-08-fix-cache")
	writeLifecycleFile(t, filepath.Join(archive, "proposal.md"), "proposal")
	service := NewLifecycleService(root, nil)

	record, err := service.Reactivate(archive)
	if err != nil {
		t.Fatal(err)
	}
	active := filepath.Join(root, "openspec", "changes", "fix-cache")
	if _, err := os.Stat(filepath.Join(active, "proposal.md")); err != nil {
		t.Fatalf("active change not restored: %v", err)
	}
	if err := service.Undo(record); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(archive, "proposal.md")); err != nil {
		t.Fatalf("archive not restored: %v", err)
	}
}

func TestLifecycleReactivateFallsBackToCopyAndRemove(t *testing.T) {
	root := t.TempDir()
	archive := filepath.Join(root, "openspec", "changes", "archive", "2026-08-08-fix-cache")
	writeLifecycleFile(t, filepath.Join(archive, "proposal.md"), "proposal")
	service := NewLifecycleService(root, nil)
	service.rename = func(string, string) error { return errors.New("cross-device link") }
	if _, err := service.Reactivate(archive); err != nil {
		t.Fatal(err)
	}
	active := filepath.Join(root, "openspec", "changes", "fix-cache")
	if got, err := os.ReadFile(filepath.Join(active, "proposal.md")); err != nil || string(got) != "proposal" {
		t.Fatalf("fallback did not copy active change: %q, %v", got, err)
	}
	if _, err := os.Stat(archive); !os.IsNotExist(err) {
		t.Fatalf("fallback did not remove source: %v", err)
	}
}

func TestLifecycleReactivateRestoresSourceWhenFallbackCopyFails(t *testing.T) {
	root := t.TempDir()
	archive := filepath.Join(root, "openspec", "changes", "archive", "2026-08-08-fix-cache")
	writeLifecycleFile(t, filepath.Join(archive, "proposal.md"), "proposal")
	service := NewLifecycleService(root, nil)
	service.rename = func(string, string) error { return errors.New("cross-device link") }
	service.copy = func(_, target string) error {
		writeLifecycleFile(t, filepath.Join(target, "partial.md"), "partial")
		return errors.New("copy failed")
	}
	_, err := service.Reactivate(archive)
	if err == nil || !strings.Contains(err.Error(), "copy failed") {
		t.Fatalf("expected copy failure, got %v", err)
	}
	if got, readErr := os.ReadFile(filepath.Join(archive, "proposal.md")); readErr != nil || string(got) != "proposal" {
		t.Fatalf("source not retained after move failure: %q, %v", got, readErr)
	}
	target := filepath.Join(root, "openspec", "changes", "fix-cache")
	if _, statErr := os.Stat(target); !os.IsNotExist(statErr) {
		t.Fatalf("partial target not removed: %v", statErr)
	}
}

func TestLifecycleReactivateSupportsLegacyUndatedArchiveName(t *testing.T) {
	root := t.TempDir()
	archive := filepath.Join(root, "openspec", "changes", "archive", "legacy-change")
	writeLifecycleFile(t, filepath.Join(archive, "proposal.md"), "proposal")
	record, err := NewLifecycleService(root, nil).Reactivate(archive)
	if err != nil {
		t.Fatal(err)
	}
	if record.ChangeName != "legacy-change" {
		t.Fatalf("clean legacy name = %q", record.ChangeName)
	}
	if _, err := os.Stat(filepath.Join(root, "openspec", "changes", "legacy-change", "proposal.md")); err != nil {
		t.Fatalf("legacy archive not reactivated: %v", err)
	}
}

func TestLifecycleReactivateRefusesCollision(t *testing.T) {
	root := t.TempDir()
	archive := filepath.Join(root, "openspec", "changes", "archive", "2026-08-08-fix-cache")
	writeLifecycleFile(t, filepath.Join(archive, "proposal.md"), "archived")
	writeLifecycleFile(t, filepath.Join(root, "openspec", "changes", "fix-cache", "proposal.md"), "active")
	_, err := NewLifecycleService(root, nil).Reactivate(archive)
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("expected collision, got %v", err)
	}
}

func TestLifecycleReactivateRejectsSymlinkEscape(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "2026-08-08-fix-cache")
	writeLifecycleFile(t, filepath.Join(outside, "proposal.md"), "outside")
	archiveRoot := filepath.Join(root, "openspec", "changes", "archive")
	if err := os.MkdirAll(archiveRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(archiveRoot, "2026-08-08-fix-cache")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	_, err := NewLifecycleService(root, nil).Reactivate(link)
	if err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("expected symlink rejection, got %v", err)
	}
}

func TestLifecycleReactivateRejectsOutsideArchiveRoot(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "2026-08-08-fix-cache")
	writeLifecycleFile(t, filepath.Join(outside, "proposal.md"), "outside")
	_, err := NewLifecycleService(root, nil).Reactivate(outside)
	if err == nil || !strings.Contains(err.Error(), "outside") {
		t.Fatalf("expected containment error, got %v", err)
	}
}

func TestLifecycleReactivationUndoRefusesChangedFiles(t *testing.T) {
	root := t.TempDir()
	archive := filepath.Join(root, "openspec", "changes", "archive", "2026-08-08-fix-cache")
	writeLifecycleFile(t, filepath.Join(archive, "proposal.md"), "proposal")
	service := NewLifecycleService(root, nil)
	record, err := service.Reactivate(archive)
	if err != nil {
		t.Fatal(err)
	}
	activeFile := filepath.Join(root, "openspec", "changes", "fix-cache", "proposal.md")
	writeLifecycleFile(t, activeFile, "externally changed")
	if err := service.Undo(record); err == nil || !strings.Contains(err.Error(), filepath.Dir(activeFile)) {
		t.Fatalf("expected changed path conflict, got %v", err)
	}
	if got, _ := os.ReadFile(activeFile); string(got) != "externally changed" {
		t.Fatalf("undo partially changed file: %q", got)
	}
}

func TestLifecycleArchiveAndUndoRestoresSpecs(t *testing.T) {
	root := t.TempDir()
	active := filepath.Join(root, "openspec", "changes", "fix-cache")
	writeLifecycleFile(t, filepath.Join(active, "proposal.md"), "proposal")
	writeLifecycleFile(t, filepath.Join(active, "specs", "cache", "spec.md"), "delta")
	spec := filepath.Join(root, "openspec", "specs", "cache", "spec.md")
	writeLifecycleFile(t, spec, "before")
	archive := filepath.Join(root, "openspec", "changes", "archive", "2026-08-08-fix-cache")
	client := &fakeArchiveClient{result: ArchiveResult{ChangeName: "fix-cache", ArchivePath: archive}}
	client.run = func() {
		writeLifecycleFile(t, spec, "after")
		if err := os.MkdirAll(filepath.Dir(archive), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Rename(active, archive); err != nil {
			t.Fatal(err)
		}
	}
	service := NewLifecycleService(root, client)
	record, _, err := service.Archive(Change{Name: "fix-cache", Path: active})
	if err != nil {
		t.Fatal(err)
	}
	if err := service.Undo(record); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.ReadFile(spec); string(got) != "before" {
		t.Fatalf("spec not restored: %q", got)
	}
	if got, _ := os.ReadFile(filepath.Join(active, "proposal.md")); string(got) != "proposal" {
		t.Fatalf("change not restored: %q", got)
	}
}

func TestLifecycleUndoArchiveRemovesNewCanonicalSpecDirectory(t *testing.T) {
	root := t.TempDir()
	active := filepath.Join(root, "openspec", "changes", "fix-cache")
	writeLifecycleFile(t, filepath.Join(active, "specs", "new-capability", "spec.md"), "delta")
	spec := filepath.Join(root, "openspec", "specs", "new-capability", "spec.md")
	archive := filepath.Join(root, "openspec", "changes", "archive", "2026-08-08-fix-cache")
	client := &fakeArchiveClient{result: ArchiveResult{ArchivePath: archive}}
	client.run = func() {
		writeLifecycleFile(t, spec, "created")
		if err := os.MkdirAll(filepath.Dir(archive), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Rename(active, archive); err != nil {
			t.Fatal(err)
		}
	}
	service := NewLifecycleService(root, client)
	record, _, err := service.Archive(Change{Name: "fix-cache", Path: active})
	if err != nil {
		t.Fatal(err)
	}
	if err := service.Undo(record); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Dir(spec)); !os.IsNotExist(err) {
		t.Fatalf("new canonical spec directory remains after undo: %v", err)
	}
}

func TestLifecycleArchiveBackupContainsManualSpecPreimages(t *testing.T) {
	root := t.TempDir()
	active := filepath.Join(root, "openspec", "changes", "fix-cache")
	writeLifecycleFile(t, filepath.Join(active, "specs", "cache", "spec.md"), "delta")
	spec := filepath.Join(root, "openspec", "specs", "cache", "spec.md")
	writeLifecycleFile(t, spec, "before")
	archive := filepath.Join(root, "openspec", "changes", "archive", "2026-08-08-fix-cache")
	client := &fakeArchiveClient{result: ArchiveResult{ArchivePath: archive}}
	client.run = func() {
		writeLifecycleFile(t, spec, "after")
		if err := os.MkdirAll(filepath.Dir(archive), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Rename(active, archive); err != nil {
			t.Fatal(err)
		}
	}
	service := NewLifecycleService(root, client)
	record, _, err := service.Archive(Change{Name: "fix-cache", Path: active})
	if err != nil {
		t.Fatal(err)
	}
	defer service.Cleanup(record)
	backupSpec := filepath.Join(record.backupDir, "project-specs", "cache", "spec.md")
	if got, readErr := os.ReadFile(backupSpec); readErr != nil || string(got) != "before" {
		t.Fatalf("manual spec preimage missing: %q, %v", got, readErr)
	}
	manifest, readErr := os.ReadFile(filepath.Join(record.backupDir, "project-specs", "MANIFEST.txt"))
	if readErr != nil || !strings.Contains(string(manifest), spec) {
		t.Fatalf("manual recovery manifest missing target: %q, %v", manifest, readErr)
	}
}

func TestLifecycleArchiveFailureRestoresPartialMutation(t *testing.T) {
	root := t.TempDir()
	active := filepath.Join(root, "openspec", "changes", "fix-cache")
	writeLifecycleFile(t, filepath.Join(active, "specs", "cache", "spec.md"), "delta")
	spec := filepath.Join(root, "openspec", "specs", "cache", "spec.md")
	writeLifecycleFile(t, spec, "before")
	client := &fakeArchiveClient{err: errors.New("archive failed")}
	client.run = func() { writeLifecycleFile(t, spec, "partial") }
	_, _, err := NewLifecycleService(root, client).Archive(Change{Name: "fix-cache", Path: active})
	if err == nil || !strings.Contains(err.Error(), "archive failed") {
		t.Fatalf("expected operation failure, got %v", err)
	}
	if got, _ := os.ReadFile(spec); string(got) != "before" {
		t.Fatalf("partial mutation not restored: %q", got)
	}
}

func TestLifecycleArchiveReportsRecoveryFailureAndPreservesRecoveryCopies(t *testing.T) {
	root := t.TempDir()
	active := filepath.Join(root, "openspec", "changes", "fix-cache")
	writeLifecycleFile(t, filepath.Join(active, "proposal.md"), "proposal")
	archive := filepath.Join(root, "openspec", "changes", "archive", "2026-08-08-fix-cache")
	client := &fakeArchiveClient{err: errors.New("archive failed")}
	service := NewLifecycleService(root, client)
	client.run = func() {
		if err := os.MkdirAll(filepath.Dir(archive), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Rename(active, archive); err != nil {
			t.Fatal(err)
		}
		service.copy = func(string, string) error { return errors.New("restore copy failed") }
	}
	_, _, err := service.Archive(Change{Name: "fix-cache", Path: active})
	if err == nil || !strings.Contains(err.Error(), "archive failed") || !strings.Contains(err.Error(), "restore copy failed") || !strings.Contains(err.Error(), "dossier-lifecycle-") {
		t.Fatalf("expected operation, restoration, and backup diagnostics, got %v", err)
	}
	if got, readErr := os.ReadFile(filepath.Join(archive, "proposal.md")); readErr != nil || string(got) != "proposal" {
		t.Fatalf("archive recovery copy was not preserved: %q, %v", got, readErr)
	}
}

func TestLifecycleArchiveRejectsMismatchedDestinationAndRecovers(t *testing.T) {
	root := t.TempDir()
	active := filepath.Join(root, "openspec", "changes", "fix-cache")
	writeLifecycleFile(t, filepath.Join(active, "proposal.md"), "proposal")
	wrong := filepath.Join(root, "openspec", "changes", "archive", "2026-08-08-other-change")
	client := &fakeArchiveClient{result: ArchiveResult{ArchivePath: wrong}}
	client.run = func() {
		if err := os.MkdirAll(filepath.Dir(wrong), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Rename(active, wrong); err != nil {
			t.Fatal(err)
		}
	}
	_, _, err := NewLifecycleService(root, client).Archive(Change{Name: "fix-cache", Path: active})
	if err == nil || !strings.Contains(err.Error(), "does not match") {
		t.Fatalf("expected mismatched destination error, got %v", err)
	}
	if got, readErr := os.ReadFile(filepath.Join(active, "proposal.md")); readErr != nil || string(got) != "proposal" {
		t.Fatalf("active change not recovered: %q, %v", got, readErr)
	}
	if _, statErr := os.Stat(wrong); !os.IsNotExist(statErr) {
		t.Fatalf("wrong archive destination remains: %v", statErr)
	}
}

func TestLifecycleArchiveFailureDoesNotRemoveUnrelatedNewArchive(t *testing.T) {
	root := t.TempDir()
	active := filepath.Join(root, "openspec", "changes", "fix-cache")
	writeLifecycleFile(t, filepath.Join(active, "proposal.md"), "proposal")
	unrelated := filepath.Join(root, "openspec", "changes", "archive", "2026-08-08-other-change")
	client := &fakeArchiveClient{err: errors.New("archive failed")}
	client.run = func() { writeLifecycleFile(t, filepath.Join(unrelated, "proposal.md"), "unrelated") }
	_, _, err := NewLifecycleService(root, client).Archive(Change{Name: "fix-cache", Path: active})
	if err == nil {
		t.Fatal("expected archive failure")
	}
	if got, readErr := os.ReadFile(filepath.Join(unrelated, "proposal.md")); readErr != nil || string(got) != "unrelated" {
		t.Fatalf("recovery removed unrelated archive: %q, %v", got, readErr)
	}
}

func TestLifecycleUndoArchiveRefusesOccupiedActiveTarget(t *testing.T) {
	root := t.TempDir()
	active := filepath.Join(root, "openspec", "changes", "fix-cache")
	writeLifecycleFile(t, filepath.Join(active, "proposal.md"), "proposal")
	archive := filepath.Join(root, "openspec", "changes", "archive", "2026-08-08-fix-cache")
	client := &fakeArchiveClient{result: ArchiveResult{ArchivePath: archive}}
	client.run = func() {
		if err := os.MkdirAll(filepath.Dir(archive), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Rename(active, archive); err != nil {
			t.Fatal(err)
		}
	}
	service := NewLifecycleService(root, client)
	record, _, err := service.Archive(Change{Name: "fix-cache", Path: active})
	if err != nil {
		t.Fatal(err)
	}
	writeLifecycleFile(t, filepath.Join(active, "external.md"), "external")
	if err := service.Undo(record); err == nil || !strings.Contains(err.Error(), active) {
		t.Fatalf("expected occupied target conflict, got %v", err)
	}
	if got, readErr := os.ReadFile(filepath.Join(active, "external.md")); readErr != nil || string(got) != "external" {
		t.Fatalf("preflight partially changed occupied target: %q, %v", got, readErr)
	}
	if _, statErr := os.Stat(archive); statErr != nil {
		t.Fatalf("archive changed despite preflight refusal: %v", statErr)
	}
}

func TestLifecycleCleanupDisposesArchiveBackup(t *testing.T) {
	root := t.TempDir()
	active := filepath.Join(root, "openspec", "changes", "fix-cache")
	writeLifecycleFile(t, filepath.Join(active, "proposal.md"), "proposal")
	archive := filepath.Join(root, "openspec", "changes", "archive", "2026-08-08-fix-cache")
	client := &fakeArchiveClient{result: ArchiveResult{ArchivePath: archive}}
	client.run = func() {
		if err := os.MkdirAll(filepath.Dir(archive), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Rename(active, archive); err != nil {
			t.Fatal(err)
		}
	}
	service := NewLifecycleService(root, client)
	record, _, err := service.Archive(Change{Name: "fix-cache", Path: active})
	if err != nil {
		t.Fatal(err)
	}
	backup := record.backupDir
	service.Cleanup(record)
	if _, statErr := os.Stat(backup); !os.IsNotExist(statErr) {
		t.Fatalf("backup not disposed: %v", statErr)
	}
	if err := service.Undo(record); err == nil || !strings.Contains(err.Error(), "backup unavailable") {
		t.Fatalf("disposed record remained usable: %v", err)
	}
}

func TestLifecycleUndoArchiveRefusesChangedSpecWithoutPartialWrites(t *testing.T) {
	root := t.TempDir()
	active := filepath.Join(root, "openspec", "changes", "fix-cache")
	writeLifecycleFile(t, filepath.Join(active, "specs", "cache", "spec.md"), "delta")
	spec := filepath.Join(root, "openspec", "specs", "cache", "spec.md")
	writeLifecycleFile(t, spec, "before")
	archive := filepath.Join(root, "openspec", "changes", "archive", "2026-08-08-fix-cache")
	client := &fakeArchiveClient{result: ArchiveResult{ArchivePath: archive}}
	client.run = func() {
		writeLifecycleFile(t, spec, "after")
		if err := os.MkdirAll(filepath.Dir(archive), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Rename(active, archive); err != nil {
			t.Fatal(err)
		}
	}
	service := NewLifecycleService(root, client)
	record, _, err := service.Archive(Change{Name: "fix-cache", Path: active})
	if err != nil {
		t.Fatal(err)
	}
	writeLifecycleFile(t, spec, "external")
	if err := service.Undo(record); err == nil || !strings.Contains(err.Error(), spec) {
		t.Fatalf("expected conflict, got %v", err)
	}
	if _, err := os.Stat(active); !os.IsNotExist(err) {
		t.Fatalf("undo wrote active path despite failed preflight: %v", err)
	}
}
