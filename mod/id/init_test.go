package id

import (
	"testing"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/lib/must"
	"github.com/gov4git/gov4git/lib/testutil"
)

func TestInit(t *testing.T) {
	ctx := testutil.NewCtx()
	testID := NewTestID(ctx, t, git.MainBranch, true)
	Init(ctx, testID.OwnerAddress())
	if err := must.Try(func() { Init(ctx, testID.OwnerAddress()) }); err == nil {
		t.Fatal("second init must fail")
	}
}
