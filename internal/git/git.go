package git

import (
	"os/exec"
	"strings"
)

type FileStatus struct {
	X, Y      byte
	Path      string
	OldPath   string
	IsDeleted bool
}

func IsInsideWorkTree(root string) bool {
	cmd := exec.Command("git", "-C", root, "rev-parse", "--is-inside-work-tree")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "true"
}

func WorkTreeRoot(root string) string {
	cmd := exec.Command("git", "-C", root, "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return root
	}
	return strings.TrimSpace(string(out))
}

func Status(root string) ([]FileStatus, error) {
	cmd := exec.Command("git", "-C", root, "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	raw := strings.TrimRight(string(out), "\n\r")
	if raw == "" {
		return nil, nil
	}
	lines := strings.Split(raw, "\n")
	files := make([]FileStatus, 0, len(lines))
	for _, line := range lines {
		if len(line) < 4 {
			continue
		}
		fs := FileStatus{X: line[0], Y: line[1]}
		path := line[3:]

		if fs.X == 'R' || fs.Y == 'R' || fs.X == 'C' || fs.Y == 'C' {
			if idx := strings.Index(path, " -> "); idx >= 0 {
				fs.OldPath = path[:idx]
				fs.Path = path[idx+4:]
			} else {
				fs.Path = path
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
