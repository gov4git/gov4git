package id

import (
	"context"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
)

func GetPrivateCredentials(ctx context.Context, priv *git.Tree) PrivateCredentials {
	return form.FromFile[PrivateCredentials](ctx, priv.Filesystem, PrivateCredentialsNS.Path())
}
