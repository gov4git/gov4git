package id

import (
	"context"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/mod"
)

type PublicMod struct {
	mod.NS
	Public git.Address
}

func (m PublicMod) GetPublicCredentials(ctx context.Context, wt *git.Tree) PublicCredentials {
	return form.MustDecodeFromFile[PublicCredentials](ctx, wt.Filesystem, m.Subpath(PublicCredentialsFilebase))
}
