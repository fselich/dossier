package openspec

import "os"

type fileSystem interface {
	ReadFile(name string) ([]byte, error)
	WriteFile(name string, data []byte, perm os.FileMode) error
	ReadDir(name string) ([]os.DirEntry, error)
	Stat(name string) (os.FileInfo, error)
}

type OSFS struct{}

func (OSFS) ReadFile(name string) ([]byte, error) { return os.ReadFile(name) }
func (OSFS) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}
func (OSFS) ReadDir(name string) ([]os.DirEntry, error) { return os.ReadDir(name) }
func (OSFS) Stat(name string) (os.FileInfo, error)      { return os.Stat(name) }
