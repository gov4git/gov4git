package id

import (
	"context"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/mod"
)

type PrivateMod struct {
	mod.NS
	Public  git.Address
	Private git.Address
}

func (m PrivateMod) GetPrivateCredentials(ctx context.Context, wt *git.Tree) PrivateCredentials {
	return form.FromFile[PrivateCredentials](ctx, wt.Filesystem, m.Sub(PrivateCredentialsFilebase).Path())
}
