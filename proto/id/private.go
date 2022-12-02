package id

import (
	"context"

	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func FetchOwnerCredentials(ctx context.Context, addr OwnerAddress) PrivateCredentials {
	return GetOwnerCredentials(ctx, CloneOwner(ctx, addr))
}

func GetOwnerCredentials(ctx context.Context, owner OwnerCloned) PrivateCredentials {
	return GetPrivateCredentials(ctx, owner.Private.Tree())
}

func GetPrivateCredentials(ctx context.Context, privateTree *git.Tree) PrivateCredentials {
	return form.FromFile[PrivateCredentials](ctx, privateTree.Filesystem, PrivateCredentialsNS.Path())
}
