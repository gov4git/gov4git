package git

import "path/filepath"

func ListFilesRecursively(t *Tree, dir string) ([]string, error) {
	fs := t.Filesystem
	infos, err := fs.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	list := []string{}
	for _, info := range infos {
		if info.IsDir() {
			sublist, err := ListFilesRecursively(t, filepath.Join(dir, info.Name()))
			if err != nil {
				return nil, err
			}
			list = append(list, prefixPaths(dir, sublist)...)
		} else {
			list = append(list, filepath.Join(dir, info.Name()))
		}
	}
	return list, nil
}

func prefixPaths(prefix string, paths []string) []string {
	r := make([]string, len(paths))
	for i, p := range paths {
		r[i] = filepath.Join(prefix, p)
	}
	return r
}
