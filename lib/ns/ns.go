package ns

import (
	"path/filepath"
)

type NS string

func (ns NS) Path() string {
	return filepath.Clean(string(ns))
}

func (ns NS) Sub(path string) NS {
	return NS(filepath.Join(string(ns), path))
}

func (ns NS) Join(sub NS) NS {
	return ns.Sub(string(sub))
}
