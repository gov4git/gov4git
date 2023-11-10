package gov

import (
	"context"

	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/git"
)

type Address id.PublicAddress

func Clone(ctx context.Context, addr Address) git.Cloned {
	return git.CloneOne(ctx, git.Address(addr))
}

type OwnerAddress id.OwnerAddress

func CloneOwner(ctx context.Context, addr OwnerAddress) OwnerCloned {
	return OwnerCloned(id.CloneOwner(ctx, id.OwnerAddress(addr)))
}

type OwnerCloned id.OwnerCloned

func (x OwnerCloned) GovOwnerAddress() OwnerAddress {
	return OwnerAddress(id.OwnerCloned(x).Address())
}

func (x OwnerCloned) IDOwnerCloned() id.OwnerCloned {
	return id.OwnerCloned(x)
}
