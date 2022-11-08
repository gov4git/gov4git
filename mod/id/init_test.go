package id

import (
	"testing"

	"github.com/gov4git/gov4git/lib/must"
	"github.com/gov4git/gov4git/lib/testutil"
)

func TestInit(t *testing.T) {
	ctx := testutil.NewCtx()
	testID := InitTestID(ctx, t, true)
	Init(ctx, testID.Public.Address, testID.Private.Address)
	if err := must.Try(func() { Init(ctx, testID.Public.Address, testID.Private.Address) }); err == nil {
		t.Fatal("second init must fail")
	}
}
