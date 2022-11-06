package kkv

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/mod"
	"github.com/gov4git/gov4git/mod/kv"
)

func Set[V kv.Value](ctx context.Context, ns mod.NS, t *git.Tree, k1 kv.Key, k2 kv.Key, v V) mod.Change[form.None] {
	kv.Set(ctx, ns, t, k1, form.None{})
	kv.Set(ctx, kv.KeyNS(ns, k1), t, k2, v)
	return mod.Change[form.None]{
		Msg: fmt.Sprintf("Change value of (%v, %v) in namespace %v", k1, k2, ns),
	}
}

func Get[V kv.Value](ctx context.Context, ns mod.NS, t *git.Tree, k1 kv.Key, k2 kv.Key) V {
	return kv.Get[V](ctx, kv.KeyNS(ns, k1), t, k2)
}

func Remove(ctx context.Context, ns mod.NS, t *git.Tree, k1 kv.Key, k2 kv.Key) mod.Change[form.None] {
	//TODO: garbage-collect empty primary keys
	return kv.Remove(ctx, kv.KeyNS(ns, k1), t, k2)
}

func ListPrimaryKeys(ctx context.Context, ns mod.NS, t *git.Tree) []kv.Key {
	return kv.ListKeys(ctx, ns, t)
}

func ListSecondaryKeys(ctx context.Context, ns mod.NS, t *git.Tree, k1 kv.Key) []kv.Key {
	return kv.ListKeys(ctx, kv.KeyNS(ns, k1), t)
}
