package kv

import (
	"context"
	"fmt"

	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
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

func (KV[K, V]) KeyNS(ns ns.NS, key K) ns.NS {
	return ns.Append(form.StringHashForFilename(string(key)))
}

func (x KV[K, V]) Set(ctx context.Context, ns ns.NS, t *git.Tree, key K, value V) git.ChangeNoResult {
	keyNS := x.KeyNS(ns, key)
	git.TreeMkdirAll(ctx, t, keyNS)
	git.ToFileStage(ctx, t, keyNS.Append(keyFilebase), key)
	git.ToFileStage(ctx, t, keyNS.Append(valueFilebase), value)
	return git.NewChangeNoResult(
		fmt.Sprintf("Change value of %v in namespace %v", key, ns),
		"kv_set",
	)
}

func (x KV[K, V]) Contains(ctx context.Context, ns ns.NS, t *git.Tree, key K) bool {
	err := must.Try(
		func() {
			x.Get(ctx, ns, t, key)
		},
	)
	if err == nil {
		return true
	}
	return false
}

func (x KV[K, V]) Get(ctx context.Context, ns ns.NS, t *git.Tree, key K) V {
	return form.FromFile[V](ctx, t.Filesystem, x.KeyNS(ns, key).Append(valueFilebase))
}

func (x KV[K, V]) GetMany(ctx context.Context, ns ns.NS, t *git.Tree, keys []K) []V {
	r := make([]V, len(keys))
	for i, k := range keys {
		r[i] = x.Get(ctx, ns, t, k)
	}
	return r
}

func (x KV[K, V]) Remove(ctx context.Context, ns ns.NS, t *git.Tree, key K) git.ChangeNoResult {
	_, err := git.TreeRemove(ctx, t, x.KeyNS(ns, key))
	must.NoError(ctx, err)
	return git.NewChangeNoResult(
		fmt.Sprintf("Remove value for %v in namespace %v", key, ns),
		"kv_remove",
	)
}

func (x KV[K, V]) ListKeys(ctx context.Context, ns ns.NS, t *git.Tree) []K {
	infos, err := git.TreeReadDir(ctx, t, ns)
	must.NoError(ctx, err)
	r := []K{}
	for _, info := range infos {
		if !info.IsDir() {
			continue
		}
		keyFileNS := ns.Append(info.Name(), keyFilebase)
		k, err := must.Try1(
			func() K {
				return form.FromFile[K](ctx, t.Filesystem, keyFileNS)
			},
		)
		if err != nil {
			base.Errorf("unrecognizable kv dir %v", keyFileNS.Dir().GitPath())
			continue
		}
		r = append(r, k)
	}
	return r
}

func (x KV[K, V]) ListKeyValues(ctx context.Context, ns ns.NS, t *git.Tree) ([]K, []V) {
	keys := x.ListKeys(ctx, ns, t)
	values := make([]V, len(keys))
	for i, key := range keys {
		values[i] = x.Get(ctx, ns, t, key)
	}
	return keys, values
}
