package gov

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/id"
	"github.com/gov4git/lib4git/git"
)

// Non-owner

type Address id.PublicAddress

func (a Address) Git() git.Address {
	return id.PublicAddress(a).Git()
}

func Clone(ctx context.Context, addr Address) Cloned {
	return Cloned(id.Clone(ctx, id.PublicAddress(addr)))
}

type Cloned id.Cloned

func (x Cloned) Address() Address {
	return Address(id.Cloned(x).Address())
}

// Owner

type OwnerAddress id.OwnerAddress

func CloneOwner(ctx context.Context, addr OwnerAddress) OwnerCloned {
	cloned := OwnerCloned(id.CloneOwner(ctx, id.OwnerAddress(addr)))
	invokePostCloners(ctx, cloned)
	return cloned
}

func LiftCloned(ctx context.Context, cloned Cloned) OwnerCloned {
	return OwnerCloned(id.LiftCloned(ctx, id.Cloned(cloned)))
}

// OwnerCloned

type OwnerCloned id.OwnerCloned

func (x OwnerCloned) PublicClone() Cloned {
	return Cloned(id.OwnerCloned(x).PublicClone())
}

func (x OwnerCloned) GovAddress() Address {
	return Address(id.OwnerCloned(x).Address().Public)
}

func (x OwnerCloned) GovOwnerAddress() OwnerAddress {
	return OwnerAddress(id.OwnerCloned(x).Address())
}

func (x OwnerCloned) IDOwnerCloned() id.OwnerCloned {
	return id.OwnerCloned(x)
}
