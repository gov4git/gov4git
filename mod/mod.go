package mod

import "path/filepath"

type Mod struct {
	Namespace string
}

func (m Mod) Path() string {
	return filepath.Join(Root, filepath.Clean(m.Namespace))
}

func (m Mod) Subpath(p ...string) string {
	return filepath.Join(append([]string{m.Path()}, p...)...)
}

const Root = ".gov"
