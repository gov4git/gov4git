package kv

import (
	"context"
	"fmt"

	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/ns"
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
	return git.NewChangeNoResult(
		fmt.Sprintf("Change value of (%v, %v) in namespace %v of key-key-value.", k1, k2, ns),
		"kkv_set",
	)
}

func (x KKV[K1, K2, V]) Get(ctx context.Context, ns ns.NS, t *git.Tree, k1 K1, k2 K2) V {
	return x.Secondary().Get(ctx, x.Primary().KeyNS(ns, k1), t, k2)
}

func (x KKV[K1, K2, V]) Remove(ctx context.Context, ns ns.NS, t *git.Tree, k1 K1, k2 K2) git.ChangeNoResult {
	//TODO: garbage-collect empty primary keys
	kvChg := x.Secondary().Remove(ctx, x.Primary().KeyNS(ns, k1), t, k2)
	return git.NewChange(
		"Remove key-key-value.",
		"kkv_remove",
		form.None{},
		kvChg.Result,
		form.Forms{kvChg},
	)
}

func (x KKV[K1, K2, V]) RemovePrimary(ctx context.Context, ns ns.NS, t *git.Tree, k1 K1) git.ChangeNoResult {
	kvChg := x.Primary().Remove(ctx, ns, t, k1)
	return git.NewChange(
		"Remove primary key-key-value.",
		"kkv_remove",
		form.None{},
		kvChg.Result,
		form.Forms{kvChg},
	)
}

func (x KKV[K1, K2, V]) ListPrimaryKeys(ctx context.Context, ns ns.NS, t *git.Tree) []K1 {
	return x.Primary().ListKeys(ctx, ns, t)
}

func (x KKV[K1, K2, V]) ListSecondaryKeys(ctx context.Context, ns ns.NS, t *git.Tree, k1 K1) []K2 {
	return x.Secondary().ListKeys(ctx, x.Primary().KeyNS(ns, k1), t)
}
