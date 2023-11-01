package id

import (
	"context"

	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func FetchPublicCredentials(ctx context.Context, addr PublicAddress) PublicCredentials {
	return GetPublicCredentials(ctx, git.CloneOne(ctx, git.Address(addr)).Tree())
}

func GetPublicCredentials(ctx context.Context, t *git.Tree) PublicCredentials {
	cred := form.FromFile[PublicCredentials](ctx, t.Filesystem, PublicCredentialsNS)
	must.Assertf(ctx, cred.IsValid(), "credentials are not valid")
	return cred
}
