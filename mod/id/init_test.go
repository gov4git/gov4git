package id

import (
	"testing"

	"github.com/gov4git/gov4git/lib/must"
	"github.com/gov4git/gov4git/lib/testutil"
)

func TestInit(t *testing.T) {
	ctx := testutil.NewCtx()
	testID := InitTestID(ctx, t, true)
	ownerAddr := OwnerAddress{
		Public:  PublicAddress(testID.Public.Address),
		Private: PrivateAddress(testID.Private.Address),
	}
	Init(ctx, ownerAddr)
	if err := must.Try(func() { Init(ctx, ownerAddr) }); err == nil {
		t.Fatal("second init must fail")
	}
}
