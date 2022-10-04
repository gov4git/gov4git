package files

import (
	"path/filepath"
	"strings"
)

func MakeNonAbs(p string) string {
	return strings.TrimLeft(filepath.Clean(p), "/") //XXX: unix specific?
}

func MakeNonAbsPaths(paths []string) []string {
	r := make([]string, len(paths))
	for i, p := range paths {
		r[i] = MakeNonAbs(p)
	}
	return r
}
