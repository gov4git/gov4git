package kv

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/lib/must"
	"github.com/gov4git/gov4git/mod"
)

const (
	keyFilebase   = "key.json"
	valueFilebase = "value.json"
)

type Key = git.URL
type Value = form.Form

func KeyNS(ns mod.NS, key Key) mod.NS {
	return ns.Sub(form.StringHashForFilename(string(key)))
}

func Set[V Value](ctx context.Context, ns mod.NS, t *git.Tree, key Key, value V) mod.Change[form.None] {
	keyNS := KeyNS(ns, key)
	err := t.Filesystem.MkdirAll(keyNS.Path(), 0755)
	must.NoError(ctx, err)
	form.ToFile(ctx, t.Filesystem, filepath.Join(keyNS.Path(), keyFilebase), key)
	form.ToFile(ctx, t.Filesystem, filepath.Join(keyNS.Path(), valueFilebase), value)
	git.Add(ctx, t, keyNS.Path())
	return mod.Change[form.None]{
		Msg: fmt.Sprintf("Change value of %v in namespace %v", key, ns),
	}
}

func Get[V Value](ctx context.Context, ns mod.NS, t *git.Tree, key Key) V {
	return form.FromFile[V](ctx, t.Filesystem, filepath.Join(KeyNS(ns, key).Path(), valueFilebase))
}

func GetMany[V Value](ctx context.Context, ns mod.NS, t *git.Tree, keys []Key) []V {
	r := make([]V, len(keys))
	for i, k := range keys {
		r[i] = Get[V](ctx, ns, t, k)
	}
	return r
}

func Remove(ctx context.Context, ns mod.NS, t *git.Tree, key Key) mod.Change[form.None] {
	_, err := t.Remove(KeyNS(ns, key).Path())
	must.NoError(ctx, err)
	return mod.Change[form.None]{
		Msg: fmt.Sprintf("Remove value for %v in namespace %v", key, ns),
	}
}

func ListKeys(ctx context.Context, ns mod.NS, t *git.Tree) []Key {
	infos, err := t.Filesystem.ReadDir(ns.Path())
	must.NoError(ctx, err)
	r := make([]Key, len(infos))
	for i, info := range infos {
		r[i] = form.FromFile[Key](ctx, t.Filesystem, filepath.Join(ns.Path(), info.Name(), keyFilebase))
	}
	return r
}
