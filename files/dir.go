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
	return WriteByteFile(d.Path, file)
}

func (d Dir) WriteByteFiles(files ByteFiles) error {
	for _, f := range files {
		if err := WriteByteFile(d.Path, f); err != nil {
			return err
		}
	}
	return nil
}

func (d Dir) ReadByteFile(path string) (ByteFile, error) {
	return ReadByteFile(d.Abs(path))
}

func (d Dir) WriteFormFile(file FormFile) error {
	return WriteFormFile(d.Path, file)
}

func (d Dir) WriteFormFiles(files FormFiles) error {
	for _, f := range files {
		if err := WriteFormFile(d.Path, f); err != nil {
			return err
		}
	}
	return nil
}

func (d Dir) ReadFormFile(path string, f any) (FormFile, error) {
	return ReadFormFile(d.Abs(path), f)
}
