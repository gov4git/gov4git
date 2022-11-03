package git

import (
	"path/filepath"
	"strings"
)

func MakeNonAbs(p string) string {
	return strings.TrimLeft(filepath.Clean(p), "/")
}

func MakeNonAbsPaths(paths []string) []string {
	r := make([]string, len(paths))
	for i, p := range paths {
		r[i] = MakeNonAbs(p)
	}
	return r
}
