package id

import (
	"context"

	"github.com/gov4git/lib4git/git"
)

type PublicAddress git.Address

type PrivateAddress git.Address

type OwnerAddress struct {
	Public  PublicAddress
	Private PrivateAddress
}

type OwnerRepo struct {
	Public  *git.Repository
	Private *git.Repository
}

type OwnerTree struct {
	Public  *git.Tree
	Private *git.Tree
}

func CloneTree(ctx context.Context, addr PublicAddress) *git.Tree {
	_, publicTree := git.Clone(ctx, git.Address(addr))
	return publicTree
}

func CloneOwner(ctx context.Context, ownerAddr OwnerAddress) (OwnerRepo, OwnerTree) {
	publicRepo, publicTree := git.CloneOrInit(ctx, git.Address(ownerAddr.Public))
	privateRepo, privateTree := git.CloneOrInit(ctx, git.Address(ownerAddr.Private))
	return OwnerRepo{Public: publicRepo, Private: privateRepo}, OwnerTree{Public: publicTree, Private: privateTree}
}
