package id

import (
	"context"

	"github.com/gov4git/lib4git/git"
)

type HomeAddress git.Address

type VaultAddress git.Address

type OwnerAddress struct {
	Home  HomeAddress
	Vault VaultAddress
}

type OwnerRepo struct {
	Home  *git.Repository
	Vault *git.Repository
}

type OwnerTree struct {
	Home  *git.Tree
	Vault *git.Tree
}

func CloneTree(ctx context.Context, addr HomeAddress) *git.Tree {
	_, homeTree := git.Clone(ctx, git.Address(addr))
	return homeTree
}

func CloneOwner(ctx context.Context, ownerAddr OwnerAddress) (OwnerRepo, OwnerTree) {
	homeRepo, homeTree := git.CloneOrInit(ctx, git.Address(ownerAddr.Home))
	vaultRepo, vaultTree := git.CloneOrInit(ctx, git.Address(ownerAddr.Vault))
	return OwnerRepo{Home: homeRepo, Vault: vaultRepo}, OwnerTree{Home: homeTree, Vault: vaultTree}
}
