package gov

import (
	"context"

	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/git"
)

type GovPublicAddress id.PublicAddress

func Clone(ctx context.Context, addr GovPublicAddress) git.Cloned {
	return git.CloneOne(ctx, git.Address(addr))
}

type GovPrivateAddress id.OwnerAddress

func CloneOrganizer(ctx context.Context, addr GovPrivateAddress) id.OwnerCloned {
	return id.CloneOwner(ctx, id.OwnerAddress(addr))
}
