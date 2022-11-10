package id

import (
	"context"
	"fmt"
	"testing"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/lib/testutil"
)

type TestID struct {
	Public  testutil.PlainAddress
	Private testutil.PlainAddress
}

func (x TestID) PublicAddress() PublicAddress {
	return PublicAddress(x.Public.Address)
}

func (x TestID) String() string {
	return fmt.Sprintf("test identity, public_dir=%v private_dir=%v\n", x.Public.Dir, x.Private.Dir)
}

func InitTestID(ctx context.Context, t *testing.T, isBare bool) TestID {
	return TestID{
		Public:  testutil.InitPlainAddress(ctx, t, git.MainBranch, isBare),
		Private: testutil.InitPlainAddress(ctx, t, git.MainBranch, isBare),
	}
}
