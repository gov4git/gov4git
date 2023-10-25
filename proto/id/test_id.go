package id

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/testutil"
)

type TestID struct {
	Public  testutil.LocalAddress
	Private testutil.LocalAddress
}

func (x TestID) OwnerAddress() OwnerAddress {
	return OwnerAddress{Public: x.PublicAddress(), Private: x.PrivateAddress()}
}

func (x TestID) OwnerCloned() OwnerCloned {
	return OwnerCloned{Public: x.Public, Private: x.Private}
}

func (x TestID) PublicAddress() PublicAddress {
	return PublicAddress(x.Public.Address())
}

func (x TestID) PrivateAddress() PrivateAddress {
	return PrivateAddress(x.Private.Address())
}

func (x TestID) String() string {
	return fmt.Sprintf("test identity, public_dir=%v private_dir=%v\n", x.Public.Dir(), x.Private.Dir())
}

func NewTestID(ctx context.Context, t *testing.T, branch git.Branch, isBare bool) TestID {
	return TestID{
		Public:  testutil.NewLocalAddress(ctx, t, branch, isBare),
		Private: testutil.NewLocalAddress(ctx, t, branch, isBare),
	}
}

func (x TestID) Erase(ctx context.Context) {
	must.NoError(ctx, os.RemoveAll(x.Public.Dir()))
	must.NoError(ctx, os.RemoveAll(x.Private.Dir()))
}
