package gov

import (
	"context"

	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/git"
)

type GovAddress id.PublicAddress

func Clone(ctx context.Context, addr GovAddress) (*git.Repository, *git.Tree) {
	r, t := git.Clone(ctx, git.Address(addr))
	return r, t
}

type OrganizerAddress id.OwnerAddress

func CloneOrganizer(ctx context.Context, addr OrganizerAddress) (id.OwnerRepo, id.OwnerTree) {
	return id.CloneOwner(ctx, id.OwnerAddress(addr))
}