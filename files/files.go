package files

import (
	"os"
	"path/filepath"

	"github.com/petar/gitsoc/form"
)

type ByteFile struct {
	Path  string
	Bytes []byte
}

type ByteFiles []ByteFile

func (x ByteFiles) Paths() []string {
	p := make([]string, len(x))
	for i, f := range x {
		p[i] = f.Path
	}
	return p
}

func ReadByteFile(path string) (ByteFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ByteFile{}, err
	}
	return ByteFile{Path: path, Bytes: data}, nil
}

func ReadByteFiles(paths []string) (ByteFiles, error) {
	files := make(ByteFiles, len(paths))
	for i, p := range paths {
		file, err := ReadByteFile(p)
		if err != nil {
			return nil, err
		}
		files[i] = file
	}
	return files, nil
}

func WriteByteFiles(dir string, files ByteFiles) error {
	for _, file := range files {
		fdir, _ := filepath.Split(file.Path)
		if err := os.MkdirAll(filepath.Join(dir, fdir), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dir, file.Path), file.Bytes, 0644); err != nil {
			return err
		}
	}
	return nil
}

type FormFile struct {
	Path string
	Form any
}

type FormFiles []FormFile

func (x FormFiles) Paths() []string {
	p := make([]string, len(x))
	for i, f := range x {
		p[i] = f.Path
	}
	return p
}

func ReadFormFile(path string, f any) (FormFile, error) {
	if err := form.DecodeFormFromFile(path, f); err != nil {
		return FormFile{}, err
	}
	return FormFile{Path: path, Form: f}, nil
}

func WriteFormFiles(root string, formFiles FormFiles) error {
	for _, ff := range formFiles {
		fdir, _ := filepath.Split(ff.Path)
		if err := os.MkdirAll(filepath.Join(root, fdir), 0755); err != nil {
			return err
		}
		if err := form.EncodeFormToFile(ff.Form, filepath.Join(root, ff.Path)); err != nil {
			return err
		}
	}
	return nil
}