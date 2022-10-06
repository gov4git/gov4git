package files

import (
	"context"
	"os"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/base"
	"github.com/gov4git/gov4git/lib/form"
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

func WriteByteFile(path string, bytes []byte) error {
	dir, _ := filepath.Split(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(path, bytes, 0644); err != nil {
		return err
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

func ReadFormFile(ctx context.Context, path string, f any) (FormFile, error) {
	if err := form.DecodeFormFromFile(ctx, path, f); err != nil {
		return FormFile{}, err
	}
	return FormFile{Path: path, Form: f}, nil
}

func WriteFormFile(ctx context.Context, path string, f form.Form) error {
	base.Infof("writing form file %v", path)
	dir, _ := filepath.Split(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if err := form.EncodeFormToFile(ctx, f, path); err != nil {
		return err
	}
	base.Infof("wrote form file %v", path)
	return nil
}
