package id

import (
	"context"

	"github.com/gov4git/lib4git/git"
)

// PublicAddress points to the user's public repo.
type PublicAddress git.Address

func (x PublicAddress) Git() git.Address {
	return git.Address(x)
}

func (x PublicAddress) IsEmpty() bool {
	return x.IsEmpty()
}

// PrivateAdress points to the user's private repo.
type PrivateAddress git.Address

var ZeroPrivateAddress PrivateAddress

func (x PrivateAddress) Git() git.Address {
	return git.Address(x)
}

type OwnerAddress struct {
	Public  PublicAddress
	Private PrivateAddress
}

func LiftAddress(pub PublicAddress) OwnerAddress {
	return OwnerAddress{
		Public:  pub,
		Private: ZeroPrivateAddress,
	}
}

// Cloned
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

// OwnerCloned
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

func LiftCloned(ctx context.Context, cloned Cloned) OwnerCloned {
	return OwnerCloned{
		address: LiftAddress(cloned.Address()),
		Public:  cloned.Cloned,
		Private: nil,
	}
}
