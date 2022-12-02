package id

import (
	"context"

	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func FetchPublicCredentials(ctx context.Context, addr PublicAddress) PublicCredentials {
	return GetPublicCredentials(ctx, git.Clone(ctx, git.Address(addr)).Tree())
}

func GetPublicCredentials(ctx context.Context, t *git.Tree) PublicCredentials {
	return form.FromFile[PublicCredentials](ctx, t.Filesystem, PublicCredentialsNS.Path())
}
