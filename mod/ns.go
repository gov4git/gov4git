package mod

import "path/filepath"

type NS string

func (ns NS) Path() string {
	return filepath.Join(Root, filepath.Clean(string(ns)))
}

func (ns NS) Subpath(p ...string) string {
	return filepath.Join(append([]string{ns.Path()}, p...)...)
}

const Root = ".gov"
