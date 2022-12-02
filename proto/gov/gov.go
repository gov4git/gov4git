package gov

import (
	"context"

	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/git"
)

type GovAddress id.PublicAddress

func Clone(ctx context.Context, addr GovAddress) git.Cloned {
	return git.Clone(ctx, git.Address(addr))
}

type OrganizerAddress id.OwnerAddress

func CloneOrganizer(ctx context.Context, addr OrganizerAddress) id.OwnerCloned {
	return id.CloneOwner(ctx, id.OwnerAddress(addr))
}
