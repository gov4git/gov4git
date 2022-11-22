package id

import (
	"context"

	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func FetchOwnerCredentials(ctx context.Context, addr OwnerAddress) PrivateCredentials {
	_, tree := CloneOwner(ctx, addr)
	return GetOwnerCredentials(ctx, tree)
}

func GetOwnerCredentials(ctx context.Context, owner OwnerTree) PrivateCredentials {
	return GetPrivateCredentials(ctx, owner.Private)
}

func GetPrivateCredentials(ctx context.Context, vault *git.Tree) PrivateCredentials {
	return form.FromFile[PrivateCredentials](ctx, vault.Filesystem, PrivateCredentialsNS.Path())
}
