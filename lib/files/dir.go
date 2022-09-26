package files

import (
	"context"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// WithWorkDir attaches a working directory to the given context and returns the updated context.
func WithWorkDir(ctx context.Context, dir Dir) context.Context {
	return context.WithValue(ctx, ctxDirKey{}, dir)
}

type ctxDirKey struct{}

func WorkDir(ctx context.Context) *Dir {
	d, ok := ctx.Value(ctxDirKey{}).(Dir)
	if !ok {
		return nil
	}
	return &d
}

type Dir struct {
	Path string
}

func PathDir(path string) Dir {
	return Dir{Path: path}
}

func TempDir() Dir {
	return Dir{Path: os.TempDir()}
}

func (d Dir) Abs(path string) string {
	return filepath.Join(d.Path, path)
}

func (d Dir) Stat(path string) (os.FileInfo, error) {
	return os.Stat(d.Abs(path))
}

func (d Dir) Subdir(path string) Dir {
	return Dir{Path: filepath.Join(d.Path, path)}
}

func (d Dir) Mk() error {
	return os.MkdirAll(d.Path, 0755)
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

func (d Dir) WriteFormFile(ctx context.Context, file FormFile) error {
	return WriteFormFile(ctx, d.Path, file)
}

func (d Dir) WriteFormFiles(ctx context.Context, files FormFiles) error {
	for _, f := range files {
		if err := WriteFormFile(ctx, d.Path, f); err != nil {
			return err
		}
	}
	return nil
}

func (d Dir) ReadFormFile(ctx context.Context, path string, f any) (FormFile, error) {
	return ReadFormFile(ctx, d.Abs(path), f)
}

// EphemeralPath returns /prefix/YYYY-MM-DD/HH:MM:SS/suffix/nonce
func EphemeralPath(prefix, suffix string) string {
	t := time.Now()
	return filepath.Join(
		"ephemeral",
		t.Format("2006-01-02"),
		t.Format("15:04:05"),
		strconv.FormatUint(uint64(rand.Int63()), 36),
	)
}

func (d Dir) MkEphemeralDir(prefix, suffix string) (Dir, error) {
	eph := Dir{Path: d.Abs(EphemeralPath(prefix, suffix))}
	if err := os.MkdirAll(eph.Path, 0755); err != nil {
		return Dir{}, err
	}
	return eph, nil
}
