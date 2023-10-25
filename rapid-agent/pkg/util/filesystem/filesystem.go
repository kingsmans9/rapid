package filesystem

import (
	"io"
	"os"
	"path/filepath"
)

// Interface defines methods for interacting with an underlying file system.
type Interface interface {
	TempDir(dir, prefix string) (string, error)
	MkdirAll(path string, perm os.FileMode) error
	Create(name string) (io.WriteCloser, error)
	OpenFile(name string, flag int, perm os.FileMode) (io.WriteCloser, error)
	RemoveAll(path string) error
	ReadDir(dirname string) ([]os.FileInfo, error)
	ReadFile(filename string) ([]byte, error)
	DirExists(path string) (bool, error)
	TempFile(dir, prefix string) (NameWriteCloser, error)
	Stat(path string) (os.FileInfo, error)
	Glob(path string) ([]string, error)
}

type NameWriteCloser interface {
	io.WriteCloser

	Name() string
}

func NewFileSystem() Interface {
	return &osFileSystem{}
}

type osFileSystem struct{}

func (fs *osFileSystem) Glob(path string) ([]string, error) {
	return filepath.Glob(path)
}

func (fs *osFileSystem) TempDir(dir, prefix string) (string, error) {
	return os.MkdirTemp(dir, prefix)
}

func (fs *osFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (fs *osFileSystem) Create(name string) (io.WriteCloser, error) {
	return os.Create(name)
}

func (fs *osFileSystem) OpenFile(name string, flag int, perm os.FileMode) (io.WriteCloser, error) {
	return os.OpenFile(name, flag, perm)
}

func (fs *osFileSystem) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (fs *osFileSystem) ReadDir(dirname string) ([]os.FileInfo, error) {
	var fileInfos []os.FileInfo
	dirInfos, ise := os.ReadDir(dirname)
	if ise != nil {
		return fileInfos, ise
	}
	for _, dirInfo := range dirInfos {
		fileInfo, ise := dirInfo.Info()
		if ise == nil {
			fileInfos = append(fileInfos, fileInfo)
		}
	}
	return fileInfos, nil
}

func (fs *osFileSystem) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func (fs *osFileSystem) DirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (fs *osFileSystem) TempFile(dir, prefix string) (NameWriteCloser, error) {
	return os.CreateTemp(dir, prefix)
}

func (fs *osFileSystem) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}
