// Package member implements community member management services
package member

import (
	"context"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/mod"
)

const (
	everybody = "everybody"
)

var (
	userNS  = mod.NS("users")
	groupNS = mod.NS("groups")
)

func AddUser(ctx context.Context, ns mod.NS, t *git.Tree, name string, url git.URL) {
	XXX
}
