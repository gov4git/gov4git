package files

import (
	"context"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func WithDir(ctx context.Context, dir Dir) context.Context {
	return context.WithValue(ctx, ctxDirKey{}, dir)
}

type ctxDirKey struct{}

func DirOf(ctx context.Context) Dir {
	return ctx.Value(ctxDirKey{}).(Dir)
}

type Dir struct {
	Path string
}

func (d Dir) Abs(path string) string {
	return filepath.Join(d.Path, path)
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

func EphemeralPath(topic string) string {
	t := time.Now()
	return filepath.Join(
		"ephemeral",
		strconv.Itoa(t.Year()),
		strconv.Itoa(int(t.Month())),
		strconv.Itoa(t.Day()),
		strconv.Itoa(t.Hour()),
		strings.Join([]string{topic, strconv.FormatUint(uint64(rand.Int63()), 64)}, "."),
	)
}

func (d Dir) MakeEphemeralDir(topic string) (string, error) {
	eph := EphemeralPath(topic)
	if err := os.MkdirAll(d.Abs(eph), 0755); err != nil {
		return "", err
	}
	return eph, nil
}
