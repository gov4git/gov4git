package member

import (
	"context"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/mod"
)

func AddUser(ctx context.Context, ns mod.NS, t *git.Tree, name string, url git.URL) {
	XXX
}

func SetUser(ctx context.Context, ns mod.NS, t *git.Tree, name string, url git.URL) {
	XXX
}

func GetUser(ctx context.Context, ns mod.NS, t *git.Tree, name string) git.URL {
	XXX
}

func RemoveUser(ctx context.Context, ns mod.NS, t *git.Tree, name string) {
	XXX
}
