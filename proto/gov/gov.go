package gov

import (
	"context"

	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/git"
)

type GovAddress id.PublicAddress

func Clone(ctx context.Context, addr GovAddress) git.Cloned {
	return git.CloneOne(ctx, git.Address(addr))
}

type GovOwnerAddress id.OwnerAddress

func CloneOwner(ctx context.Context, addr GovOwnerAddress) GovOwnerCloned {
	return GovOwnerCloned(id.CloneOwner(ctx, id.OwnerAddress(addr)))
}

type GovOwnerCloned id.OwnerCloned

func (x GovOwnerCloned) IDOwnerCloned() id.OwnerCloned {
	return id.OwnerCloned(x)
}
