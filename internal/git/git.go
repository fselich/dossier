package git

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const gitTimeout = 2 * time.Second

type FileStatus struct {
	X, Y      byte
	Path      string
	OldPath   string
	IsDeleted bool
}

func RunGit(dir string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gitTimeout)
	defer cancel()
	gitArgs := append([]string{"-C", dir}, args...)
	cmd := exec.CommandContext(ctx, "git", gitArgs...)
	return cmd.Output()
}

func Stage(root string, paths ...string) error {
	args := append([]string{"add", "--"}, paths...)
	_, err := RunGit(root, args...)
	return err
}

func Unstage(root string, paths ...string) error {
	args := append([]string{"reset", "-q", "HEAD", "--"}, paths...)
	_, err := RunGit(root, args...)
	if err != nil {
		args2 := append([]string{"rm", "--cached", "-q", "--"}, paths...)
		_, err2 := RunGit(root, args2...)
		return err2
	}
	return nil
}

func IsInsideWorkTree(root string) bool {
	out, err := RunGit(root, "rev-parse", "--is-inside-work-tree")
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "true"
}

func WorkTreeRoot(root string) (string, error) {
	out, err := RunGit(root, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("git rev-parse --show-toplevel: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func Status(root string) ([]FileStatus, error) {
	out, err := RunGit(root, "status", "--porcelain=v1", "-z", "-u")
	if err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, nil
	}
	parts := strings.Split(string(out), "\x00")
	files := make([]FileStatus, 0, len(parts))
	for i := 0; i < len(parts); i++ {
		p := parts[i]
		if len(p) < 4 {
			continue
		}
		fs := FileStatus{X: p[0], Y: p[1]}
		path := p[3:]

		if fs.X == 'R' || fs.Y == 'R' || fs.X == 'C' || fs.Y == 'C' {
			fs.OldPath = path
			i++
			if i < len(parts) {
				fs.Path = parts[i]
			}
		} else {
			fs.Path = path
		}

		if fs.X == 'D' || fs.Y == 'D' {
			fs.IsDeleted = true
		}

		if strings.HasPrefix(fs.Path, "openspec/") {
			continue
		}

		files = append(files, fs)
	}
	return files, nil
}
