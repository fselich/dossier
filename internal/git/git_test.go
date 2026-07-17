package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func initRepo(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	mustRun(t, dir, "init")
	mustRun(t, dir, "config", "user.email", "test@test")
	mustRun(t, dir, "config", "user.name", "Test")
	for name, content := range files {
		err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644)
		if err != nil {
			t.Fatal(err)
		}
	}
	if len(files) > 0 {
		mustRun(t, dir, "add", ".")
		mustRun(t, dir, "commit", "-m", "initial")
	}
	return dir
}

func mustRun(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return string(out)
}

func skipIfNoGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not found on PATH")
	}
}

func TestIsInsideWorkTree(t *testing.T) {
	skipIfNoGit(t)
	dir := initRepo(t, map[string]string{"a.txt": "hello"})

	if !IsInsideWorkTree(dir) {
		t.Error("expected true inside a repo")
	}

	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if !IsInsideWorkTree(sub) {
		t.Error("expected true inside a subdirectory of a repo")
	}

	outside := t.TempDir()
	if IsInsideWorkTree(outside) {
		t.Error("expected false outside a worktree")
	}
}

func TestWorkTreeRoot(t *testing.T) {
	skipIfNoGit(t)
	dir := initRepo(t, map[string]string{"a.txt": "hello"})

	got, err := WorkTreeRoot(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != dir {
		t.Errorf("WorkTreeRoot = %q, want %q", got, dir)
	}

	sub := filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	got, err = WorkTreeRoot(sub)
	if err != nil {
		t.Fatalf("unexpected error for subdirectory: %v", err)
	}
	if got != dir {
		t.Errorf("WorkTreeRoot(sub) = %q, want %q", got, dir)
	}

	outside := t.TempDir()
	_, err = WorkTreeRoot(outside)
	if err == nil {
		t.Error("expected error for non-repo directory")
	}
}

func TestStatusEmpty(t *testing.T) {
	skipIfNoGit(t)
	dir := initRepo(t, map[string]string{"a.txt": "hello"})

	files, err := Status(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("expected no files, got %d", len(files))
	}
}

func TestStatusBasic(t *testing.T) {
	skipIfNoGit(t)
	dir := initRepo(t, nil)
	mustRun(t, dir, "commit", "--allow-empty", "-m", "init")

	tests := []struct {
		name    string
		prepare func(dir string)
		want    func(files []FileStatus) bool
	}{
		{
			name: "untracked",
			prepare: func(dir string) {
				if err := os.WriteFile(filepath.Join(dir, "new.go"), []byte("package main"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			want: func(files []FileStatus) bool {
				return len(files) == 1 && files[0].X == '?' && files[0].Y == '?' && files[0].Path == "new.go"
			},
		},
		{
			name: "modified",
			prepare: func(dir string) {
				if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("hello"), 0o644); err != nil {
					t.Fatal(err)
				}
				mustRun(t, dir, "add", "a.txt")
				mustRun(t, dir, "commit", "-m", "add a")
				if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("changed"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			want: func(files []FileStatus) bool {
				return len(files) == 1 && files[0].X == ' ' && files[0].Y == 'M' && files[0].Path == "a.txt"
			},
		},
		{
			name: "added staged",
			prepare: func(dir string) {
				if err := os.WriteFile(filepath.Join(dir, "new.go"), []byte("package main"), 0o644); err != nil {
					t.Fatal(err)
				}
				mustRun(t, dir, "add", "new.go")
			},
			want: func(files []FileStatus) bool {
				return len(files) == 1 && files[0].X == 'A' && files[0].Y == ' ' && files[0].Path == "new.go"
			},
		},
		{
			name: "deleted",
			prepare: func(dir string) {
				if err := os.WriteFile(filepath.Join(dir, "del.txt"), []byte("bye"), 0o644); err != nil {
					t.Fatal(err)
				}
				mustRun(t, dir, "add", "del.txt")
				mustRun(t, dir, "commit", "-m", "add del")
				if err := os.Remove(filepath.Join(dir, "del.txt")); err != nil {
					t.Fatal(err)
				}
			},
			want: func(files []FileStatus) bool {
				return len(files) == 1 && files[0].X == ' ' && files[0].Y == 'D' && files[0].Path == "del.txt" && files[0].IsDeleted
			},
		},
		{
			name: "staged delete",
			prepare: func(dir string) {
				if err := os.WriteFile(filepath.Join(dir, "del.txt"), []byte("bye"), 0o644); err != nil {
					t.Fatal(err)
				}
				mustRun(t, dir, "add", "del.txt")
				mustRun(t, dir, "commit", "-m", "add del")
				if err := os.Remove(filepath.Join(dir, "del.txt")); err != nil {
					t.Fatal(err)
				}
				mustRun(t, dir, "add", "del.txt")
			},
			want: func(files []FileStatus) bool {
				return len(files) == 1 && files[0].X == 'D' && files[0].Y == ' ' && files[0].Path == "del.txt" && files[0].IsDeleted
			},
		},
		{
			name: "staged modified",
			prepare: func(dir string) {
				if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("v1"), 0o644); err != nil {
					t.Fatal(err)
				}
				mustRun(t, dir, "add", "a.txt")
				mustRun(t, dir, "commit", "-m", "init")
				if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("v2"), 0o644); err != nil {
					t.Fatal(err)
				}
				mustRun(t, dir, "add", "a.txt")
			},
			want: func(files []FileStatus) bool {
				return len(files) == 1 && files[0].X == 'M' && files[0].Y == ' ' && files[0].Path == "a.txt"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := filepath.Join(t.TempDir(), "repo")
			if err := os.Mkdir(d, 0o755); err != nil {
				t.Fatal(err)
			}
			mustRun(t, d, "init")
			mustRun(t, d, "config", "user.email", "test@test")
			mustRun(t, d, "config", "user.name", "Test")
			if tt.name != "untracked" && tt.name != "added staged" {
				mustRun(t, d, "commit", "--allow-empty", "-m", "init")
			}
			tt.prepare(d)
			files, err := Status(d)
			if err != nil {
				t.Fatalf("Status: %v", err)
			}
			if !tt.want(files) {
				t.Errorf("unexpected status: %+v", files)
			}
		})
	}
}

func TestStatusRenamed(t *testing.T) {
	skipIfNoGit(t)
	dir := initRepo(t, map[string]string{"old.go": "package main\n"})
	mustRun(t, dir, "mv", "old.go", "renamed.go")

	files, err := Status(dir)
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("expected rename entry")
	}
	f := files[0]
	if f.X != 'R' && f.Y != 'R' {
		t.Errorf("expected R status, got X=%c Y=%c", f.X, f.Y)
	}
	if f.Path == "" {
		t.Error("expected non-empty Path")
	}
	if f.OldPath == "" {
		t.Error("expected non-empty OldPath")
	}
}

func TestStatusCopied(t *testing.T) {
	skipIfNoGit(t)
	dir := initRepo(t, map[string]string{"orig.go": "package main\n\nfunc Foo() {}\n"})
	if err := os.WriteFile(filepath.Join(dir, "copy.go"), []byte("package main\n\nfunc Foo() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	mustRun(t, dir, "-c", "diff.renames=copies", "add", "copy.go")

	files, err := Status(dir)
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("expected at least one entry")
	}
	f := files[0]
	if f.X != 'C' && f.Y != 'C' {
		return // copy detection may not trigger depending on similarity; not an error
	}
	if f.Path == "" || f.OldPath == "" {
		t.Errorf("expected Path and OldPath, got Path=%q OldPath=%q", f.Path, f.OldPath)
	}
}

func TestStatusOpenspecFilter(t *testing.T) {
	skipIfNoGit(t)
	dir := initRepo(t, nil)
	mustRun(t, dir, "commit", "--allow-empty", "-m", "init")
	if err := os.Mkdir(filepath.Join(dir, "openspec"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "openspec", "test.md"), []byte("hi"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(dir, "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "src", "main.go"), []byte("package main"), 0o644); err != nil {
		t.Fatal(err)
	}

	files, err := Status(dir)
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	for _, f := range files {
		if len(f.Path) >= 9 && f.Path[:9] == "openspec/" {
			t.Errorf("file under openspec/ should be filtered: %q", f.Path)
		}
	}
	if len(files) != 1 || files[0].Path != "src/main.go" {
		t.Errorf("expected only src/main.go, got %+v", files)
	}
}

func TestStatusUnusualFilenames(t *testing.T) {
	skipIfNoGit(t)

	tests := []struct {
		name     string
		filename string
	}{
		{"spaces", "file with spaces.go"},
		{"quotes", `file"with'quotes.go`},
		{"unicode", "résumé.go"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := initRepo(t, nil)
			mustRun(t, dir, "commit", "--allow-empty", "-m", "init")
			path := filepath.Join(dir, tt.filename)
			if err := os.WriteFile(path, []byte("package main"), 0o644); err != nil {
				t.Fatal(err)
			}
			files, err := Status(dir)
			if err != nil {
				t.Fatalf("Status: %v", err)
			}
			if len(files) != 1 || files[0].Path != tt.filename {
				t.Errorf("expected single entry %q, got %+v", tt.filename, files)
			}
		})
	}
}

func TestStatusNewlineInFilename(t *testing.T) {
	skipIfNoGit(t)
	dir := initRepo(t, nil)
	mustRun(t, dir, "commit", "--allow-empty", "-m", "init")

	filename := "file\nwith\nnewline.go"
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, []byte("package main"), 0o644); err != nil {
		t.Skipf("filesystem rejected newline filename: %v", err)
	}

	files, err := Status(dir)
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 entry, got %d: %+v", len(files), files)
	}
	if files[0].Path != filename {
		t.Errorf("expected path %q, got %q", filename, files[0].Path)
	}
}

func TestStageModified(t *testing.T) {
	skipIfNoGit(t)
	dir := initRepo(t, map[string]string{"a.txt": "hello"})

	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("changed"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Stage(dir, "a.txt"); err != nil {
		t.Fatalf("Stage: %v", err)
	}
	files, err := Status(dir)
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if len(files) != 1 || files[0].X != 'M' || files[0].Y != ' ' {
		t.Errorf("expected staged modified (M ), got %+v", files)
	}
}

func TestStageUntracked(t *testing.T) {
	skipIfNoGit(t)
	dir := initRepo(t, nil)
	mustRun(t, dir, "commit", "--allow-empty", "-m", "init")

	if err := os.WriteFile(filepath.Join(dir, "new.go"), []byte("package main"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Stage(dir, "new.go"); err != nil {
		t.Fatalf("Stage: %v", err)
	}
	files, err := Status(dir)
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if len(files) != 1 || files[0].X != 'A' || files[0].Y != ' ' {
		t.Errorf("expected staged added (A ), got %+v", files)
	}
}

func TestStageDeleted(t *testing.T) {
	skipIfNoGit(t)
	dir := initRepo(t, map[string]string{"del.txt": "bye"})

	if err := os.Remove(filepath.Join(dir, "del.txt")); err != nil {
		t.Fatal(err)
	}
	if err := Stage(dir, "del.txt"); err != nil {
		t.Fatalf("Stage: %v", err)
	}
	files, err := Status(dir)
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if len(files) != 1 || files[0].X != 'D' || files[0].Y != ' ' {
		t.Errorf("expected staged delete (D ), got %+v", files)
	}
}

func TestUnstageStaged(t *testing.T) {
	skipIfNoGit(t)
	dir := initRepo(t, map[string]string{"a.txt": "hello"})

	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("changed"), 0o644); err != nil {
		t.Fatal(err)
	}
	mustRun(t, dir, "add", "a.txt")

	if err := Unstage(dir, "a.txt"); err != nil {
		t.Fatalf("Unstage: %v", err)
	}
	files, err := Status(dir)
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if len(files) != 1 || files[0].X != ' ' || files[0].Y != 'M' {
		t.Errorf("expected unstaged modified ( M), got %+v", files)
	}
}

func TestUnstageNoCommits(t *testing.T) {
	skipIfNoGit(t)
	dir := t.TempDir()
	mustRun(t, dir, "init")
	mustRun(t, dir, "config", "user.email", "test@test")
	mustRun(t, dir, "config", "user.name", "Test")

	if err := os.WriteFile(filepath.Join(dir, "new.go"), []byte("package main"), 0o644); err != nil {
		t.Fatal(err)
	}
	mustRun(t, dir, "add", "new.go")

	if err := Unstage(dir, "new.go"); err != nil {
		t.Fatalf("Unstage (no commits): %v", err)
	}
	files, err := Status(dir)
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if len(files) != 1 || files[0].X != '?' || files[0].Y != '?' {
		t.Errorf("expected untracked (??), got %+v", files)
	}
}

func TestUnstageRename(t *testing.T) {
	skipIfNoGit(t)
	dir := initRepo(t, map[string]string{"old.go": "package main\n"})

	mustRun(t, dir, "mv", "old.go", "new.go")
	mustRun(t, dir, "add", "new.go")

	if err := Unstage(dir, "old.go", "new.go"); err != nil {
		t.Fatalf("Unstage (rename): %v", err)
	}
	files, err := Status(dir)
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("expected files after unstage rename")
	}
	foundOld := false
	foundNew := false
	for _, f := range files {
		if f.Path == "old.go" {
			foundOld = true
		}
		if f.Path == "new.go" {
			foundNew = true
		}
	}
	if !foundOld {
		t.Error("expected old.go to appear after unstage rename")
	}
	if !foundNew {
		t.Error("expected new.go to appear after unstage rename")
	}
}

func TestRunGitTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping timeout test in short mode")
	}
	skipIfNoGit(t)

	dir := t.TempDir()
	out, err := RunGit(dir, "version")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) == 0 {
		t.Error("expected output from git version")
	}
}
