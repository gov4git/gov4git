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

type Key interface {
	~string
}

type Value = form.Form

type KV[K Key, V Value] struct{}

func (KV[K, V]) KeyNS(ns mod.NS, key K) mod.NS {
	return ns.Sub(form.StringHashForFilename(string(key)))
}

func (x KV[K, V]) Set(ctx context.Context, ns mod.NS, t *git.Tree, key K, value V) mod.Change[form.None] {
	keyNS := x.KeyNS(ns, key)
	err := t.Filesystem.MkdirAll(keyNS.Path(), 0755)
	must.NoError(ctx, err)
	form.ToFile(ctx, t.Filesystem, filepath.Join(keyNS.Path(), keyFilebase), key)
	form.ToFile(ctx, t.Filesystem, filepath.Join(keyNS.Path(), valueFilebase), value)
	git.Add(ctx, t, keyNS.Path())
	return mod.Change[form.None]{
		Msg: fmt.Sprintf("Change value of %v in namespace %v", key, ns),
	}
}

func (x KV[K, V]) Get(ctx context.Context, ns mod.NS, t *git.Tree, key K) V {
	return form.FromFile[V](ctx, t.Filesystem, filepath.Join(x.KeyNS(ns, key).Path(), valueFilebase))
}

func (x KV[K, V]) GetMany(ctx context.Context, ns mod.NS, t *git.Tree, keys []K) []V {
	r := make([]V, len(keys))
	for i, k := range keys {
		r[i] = x.Get(ctx, ns, t, k)
	}
	return r
}

func (x KV[K, V]) Remove(ctx context.Context, ns mod.NS, t *git.Tree, key K) mod.Change[form.None] {
	_, err := t.Remove(x.KeyNS(ns, key).Path())
	must.NoError(ctx, err)
	return mod.Change[form.None]{
		Msg: fmt.Sprintf("Remove value for %v in namespace %v", key, ns),
	}
}

func (x KV[K, V]) ListKeys(ctx context.Context, ns mod.NS, t *git.Tree) []K {
	infos, err := t.Filesystem.ReadDir(ns.Path())
	must.NoError(ctx, err)
	r := []K{}
	for _, info := range infos {
		if !info.IsDir() { // TODO: filter dirs with key hashes?
			continue
		}
		k := form.FromFile[K](ctx, t.Filesystem, filepath.Join(ns.Path(), info.Name(), keyFilebase))
		r = append(r, k)
	}
	return r
}
