package id

import (
	"context"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/lib/ns"
)

type PrivateMod struct {
	ns.NS
	Public  git.Address
	Private git.Address
}

func Private(publicAddr git.Address, privateAddr git.Address) PrivateMod {
	return PrivateMod{NS: PrivateNS, Public: publicAddr, Private: privateAddr}
}

func GetPrivateCredentials(ctx context.Context, m PrivateMod, priv *git.Tree) PrivateCredentials {
	return form.FromFile[PrivateCredentials](ctx, priv.Filesystem, m.Sub(PrivateCredentialsFilebase).Path())
}
