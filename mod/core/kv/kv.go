package kv

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/lib/must"
	"github.com/gov4git/gov4git/mod"
	"github.com/gov4git/gov4git/mod/runtime"
)

const (
	keyFilebase   = "key.json"
	valueFilebase = "value.json"
)

func keyDirpath(m mod.Mod, key git.URL) string {
	return m.Subpath(form.StringHashForFilename(string(key)))
}

func Set[V form.Form](ctx context.Context, m mod.Mod, t *git.Tree, key git.URL, value V) runtime.Change[struct{}] {
	dirpath := keyDirpath(m, key)
	err := t.Filesystem.MkdirAll(dirpath, 0755)
	must.NoError(ctx, err)
	form.MustEncodeToFile(ctx, t.Filesystem, filepath.Join(dirpath, keyFilebase), key)
	form.MustEncodeToFile(ctx, t.Filesystem, filepath.Join(dirpath, valueFilebase), value)
	git.Add(ctx, t, dirpath)
	return runtime.Change[struct{}]{Msg: fmt.Sprintf("Change value of %v in namespace %v", key, m.Namespace)}
}

func Get[V form.Form](ctx context.Context, m mod.Mod, t *git.Tree, key git.URL) V {
	dirpath := keyDirpath(m, key)
	return form.MustDecodeFromFile[V](ctx, t.Filesystem, filepath.Join(dirpath, valueFilebase))
}

func GetMany[V form.Form](ctx context.Context, m mod.Mod, t *git.Tree, keys []git.URL) []V {
	r := make([]V, len(keys))
	for i, k := range keys {
		r[i] = Get[V](ctx, m, t, k)
	}
	return r
}

func Remove(ctx context.Context, m mod.Mod, t *git.Tree, key git.URL) runtime.Change[struct{}] {
	_, err := t.Remove(keyDirpath(m, key))
	must.NoError(ctx, err)
	return runtime.Change[struct{}]{Msg: fmt.Sprintf("Remove value for %v in namespace %v", key, m.Namespace)}
}

func ListKeys(ctx context.Context, m mod.Mod, t *git.Tree) []git.URL {
	infos, err := t.Filesystem.ReadDir(m.Namespace)
	must.NoError(ctx, err)
	r := make([]git.URL, len(infos))
	for i, info := range infos {
		r[i] = form.MustDecodeFromFile[git.URL](ctx, t.Filesystem, filepath.Join(m.Namespace, info.Name(), keyFilebase))
	}
	return r
}

// XXX: Add, Multiply
