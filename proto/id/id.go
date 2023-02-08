package id

import (
	"context"

	"github.com/gov4git/lib4git/git"
)

// PublicAddress points to the user's public repo.
type PublicAddress git.Address

// PrivateAdress points to the user's private repo.
type PrivateAddress git.Address

type OwnerAddress struct {
	Public  PublicAddress
	Private PrivateAddress
}

func Clone(ctx context.Context, addr PublicAddress) git.Cloned {
	return git.CloneOne(ctx, git.Address(addr))
}

type OwnerCloned struct {
	Public  git.Cloned
	Private git.Cloned
}

func CloneOwner(ctx context.Context, ownerAddr OwnerAddress) OwnerCloned {
	return OwnerCloned{
		Public:  git.CloneOne(ctx, git.Address(ownerAddr.Public)),
		Private: git.CloneOne(ctx, git.Address(ownerAddr.Private)),
	}
}
