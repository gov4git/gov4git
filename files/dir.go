package files

import (
	"os"
	"path/filepath"
)

type Dir struct {
	Path string
}

func (d Dir) Abs(path string) string {
	return filepath.Join(d.Path, path)
}

func (d Dir) Mkdir(path string) error {
	return os.MkdirAll(d.Abs(path), 0755)
}

func (d Dir) WriteByteFile(file ByteFile) error {
	return WriteByteFiles(d.Path, ByteFiles{file})
}

func (d Dir) ReadByteFile(path string) (ByteFile, error) {
	return ReadByteFile(d.Abs(path))
}

func (d Dir) WriteFormFile(file FormFile) error {
	return WriteFormFiles(d.Path, FormFiles{file})
}

func (d Dir) ReadFormFile(path string, f any) (FormFile, error) {
	return ReadFormFile(d.Abs(path), f)
}
