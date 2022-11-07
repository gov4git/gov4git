package member

import (
	"context"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/mod"
)

func AddGroup(ctx context.Context, ns mod.NS, t *git.Tree, name string, url git.URL) {
	XXX
}

func SetGroup(ctx context.Context, ns mod.NS, t *git.Tree, name string, url git.URL) {
	XXX
}

func GetGroup(ctx context.Context, ns mod.NS, t *git.Tree, name string) git.URL {
	XXX
}

func RemoveGroup(ctx context.Context, ns mod.NS, t *git.Tree, name string) {
	XXX
}
