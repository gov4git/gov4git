package kv

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/lib/ns"
)

type KKV[K1 Key, K2 Key, V Value] struct{}

func (x KKV[K1, K2, V]) Primary() KV[K1, form.None] {
	return KV[K1, form.None]{}
}

func (x KKV[K1, K2, V]) Secondary() KV[K2, V] {
	return KV[K2, V]{}
}

func (x KKV[K1, K2, V]) Set(ctx context.Context, ns ns.NS, t *git.Tree, k1 K1, k2 K2, v V) git.ChangeNoResult {
	x.Primary().Set(ctx, ns, t, k1, form.None{})
	x.Secondary().Set(ctx, x.Primary().KeyNS(ns, k1), t, k2, v)
	return git.ChangeNoResult{
		Msg: fmt.Sprintf("Change value of (%v, %v) in namespace %v", k1, k2, ns),
	}
}

func (x KKV[K1, K2, V]) Get(ctx context.Context, ns ns.NS, t *git.Tree, k1 K1, k2 K2) V {
	return x.Secondary().Get(ctx, x.Primary().KeyNS(ns, k1), t, k2)
}

func (x KKV[K1, K2, V]) Remove(ctx context.Context, ns ns.NS, t *git.Tree, k1 K1, k2 K2) git.ChangeNoResult {
	//TODO: garbage-collect empty primary keys
	return x.Secondary().Remove(ctx, x.Primary().KeyNS(ns, k1), t, k2)
}

func (x KKV[K1, K2, V]) RemovePrimary(ctx context.Context, ns ns.NS, t *git.Tree, k1 K1) git.ChangeNoResult {
	return x.Primary().Remove(ctx, ns, t, k1)
}

func (x KKV[K1, K2, V]) ListPrimaryKeys(ctx context.Context, ns ns.NS, t *git.Tree) []K1 {
	return x.Primary().ListKeys(ctx, ns, t)
}

func (x KKV[K1, K2, V]) ListSecondaryKeys(ctx context.Context, ns ns.NS, t *git.Tree, k1 K1) []K2 {
	return x.Secondary().ListKeys(ctx, x.Primary().KeyNS(ns, k1), t)
}
