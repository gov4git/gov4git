package files

import (
	"context"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gov4git/gov4git/lib/form"
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

func (d Dir) Remove(path string) error {
	return os.Remove(d.Abs(path))
}

func (d Dir) RemoveAll(path string) error {
	return os.RemoveAll(d.Abs(path))
}

func TrimSlashes(p string) string {
	return strings.Trim(filepath.Clean(p), string(filepath.Separator))
}

func (d Dir) Glob(pattern string) ([]string, error) {
	m, err := filepath.Glob(filepath.Join(d.Path, pattern))
	if err != nil {
		return nil, err
	}
	for i := range m {
		m[i] = TrimSlashes(m[i][len(d.Path):]) // remove dir prefix
	}
	return m, nil
}

func (d Dir) WriteByteFile(path string, bytes []byte) error {
	return WriteByteFile(d.Abs(path), bytes)
}

func (d Dir) WriteByteFiles(files ByteFiles) error {
	for _, f := range files {
		if err := d.WriteByteFile(f.Path, f.Bytes); err != nil {
			return err
		}
	}
	return nil
}

func (d Dir) ReadByteFile(path string) (ByteFile, error) {
	return ReadByteFile(d.Abs(path))
}

func (d Dir) WriteFormFile(ctx context.Context, path string, f form.Form) error {
	return WriteFormFile(ctx, d.Abs(path), f)
}

func (d Dir) WriteFormFiles(ctx context.Context, files FormFiles) error {
	for _, f := range files {
		if err := d.WriteFormFile(ctx, f.Path, f.Form); err != nil {
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
	return prefix + "-" + strconv.FormatUint(uint64(rand.Int63()), 36) + "-" + suffix
	// t := time.Now()
	// return filepath.Join(
	// 	"ephemeral",
	// 	t.Format("2006-01-02"),
	// 	t.Format("15:04:05"),
	// 	prefix+"-"+strconv.FormatUint(uint64(rand.Int63()), 36)+"-"+suffix,
	// )
}

func (d Dir) MkEphemeralDir(prefix, suffix string) (Dir, error) {
	eph := Dir{Path: d.Abs(EphemeralPath(prefix, suffix))}
	if err := os.MkdirAll(eph.Path, 0755); err != nil {
		return Dir{}, err
	}
	return eph, nil
}

func Rename(from Dir, to Dir) error {
	return os.Rename(from.Path, to.Path)
}
