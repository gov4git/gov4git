package id

import (
	"context"

	"github.com/gov4git/lib4git/git"
)

// PublicAddress points to the user's public repo.
type PublicAddress git.Address

func (x PublicAddress) IsEmpty() bool {
	return x.IsEmpty()
}

// PrivateAdress points to the user's private repo.
type PrivateAddress git.Address

type OwnerAddress struct {
	Public  PublicAddress
	Private PrivateAddress
}

type Cloned struct {
	address PublicAddress
	git.Cloned
}

func (x Cloned) Address() PublicAddress {
	return x.address
}

func Clone(ctx context.Context, addr PublicAddress) Cloned {
	return Cloned{
		address: addr,
		Cloned:  git.CloneOne(ctx, git.Address(addr)),
	}
}

type OwnerCloned struct {
	address OwnerAddress
	Public  git.Cloned
	Private git.Cloned
}

func (x OwnerCloned) PublicClone() Cloned {
	return Cloned{
		address: x.address.Public,
		Cloned:  x.Public,
	}
}

func (x OwnerCloned) Address() OwnerAddress {
	return x.address
}

func CloneOwner(ctx context.Context, ownerAddr OwnerAddress) OwnerCloned {
	return OwnerCloned{
		address: ownerAddr,
		Public:  git.CloneOne(ctx, git.Address(ownerAddr.Public)),
		Private: git.CloneOne(ctx, git.Address(ownerAddr.Private)),
	}
}
