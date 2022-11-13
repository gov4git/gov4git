package mod

import (
	"context"

	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/ns"
)

var RootNS = ns.NS("")

func Commit(ctx context.Context, t *git.Tree, msg string) {
	git.Commit(ctx, t, "gov4git: "+msg)
}
