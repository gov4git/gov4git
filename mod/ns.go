package mod

import (
	"path/filepath"
)

const Root = ".gov"

type NS string

func (ns NS) Path() string {
	return filepath.Join(Root, filepath.Clean(string(ns)))
}

func (ns NS) Sub(path string) NS {
	return NS(filepath.Join(string(ns), path))
}
