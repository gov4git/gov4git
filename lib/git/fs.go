package git

import (
	"context"

	"github.com/gov4git/gov4git/lib/must"
	"github.com/gov4git/gov4git/lib/ns"
)

func TreeMkdirAll(ctx context.Context, t *Tree, ns ns.NS) {
	err := t.Filesystem.MkdirAll(ns.Path(), 0755)
	must.NoError(ctx, err)
}
