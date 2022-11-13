package id

import (
	"context"

	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func GetPrivateCredentials(ctx context.Context, priv *git.Tree) PrivateCredentials {
	return form.FromFile[PrivateCredentials](ctx, priv.Filesystem, PrivateCredentialsNS.Path())
}

func GetOwnerCredentials(ctx context.Context, owner OwnerTree) PrivateCredentials {
	return GetPrivateCredentials(ctx, owner.Private)
}
