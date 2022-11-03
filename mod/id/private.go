package id

import (
	"context"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/mod"
)

type PrivateMod struct {
	mod.Mod
	Public  git.Address
	Private git.Address
}

func (m PrivateMod) GetPrivateCredentials(ctx context.Context, wt *git.Worktree) (cred PrivateCredentials, err error) {
	return form.DecodeFromFile[PrivateCredentials](ctx, wt.Filesystem, m.Subpath(PrivateCredentialsFilebase))
}
