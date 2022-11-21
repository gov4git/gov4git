package id

import (
	"context"
	"fmt"
	"testing"

	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/testutil"
)

type TestID struct {
	Home  testutil.LocalAddress
	Vault testutil.LocalAddress
}

func (x TestID) OwnerAddress() OwnerAddress {
	return OwnerAddress{Home: x.HomeAddress(), Vault: x.VaultAddress()}
}

func (x TestID) OwnerTree() OwnerTree {
	return OwnerTree{Home: x.Home.Tree, Vault: x.Vault.Tree}
}

func (x TestID) HomeAddress() HomeAddress {
	return HomeAddress(x.Home.Address)
}

func (x TestID) VaultAddress() VaultAddress {
	return VaultAddress(x.Vault.Address)
}

func (x TestID) String() string {
	return fmt.Sprintf("test identity, home_dir=%v vault_dir=%v\n", x.Home.Dir, x.Vault.Dir)
}

func NewTestID(ctx context.Context, t *testing.T, branch git.Branch, isBare bool) TestID {
	return TestID{
		Home:  testutil.NewLocalAddress(ctx, t, branch, isBare),
		Vault: testutil.NewLocalAddress(ctx, t, branch, isBare),
	}
}
