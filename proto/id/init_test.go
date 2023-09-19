package id

import (
	"testing"

	"github.com/gov4git/gov4git/runtime"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/testutil"
)

func TestInit(t *testing.T) {
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	testID := NewTestID(ctx, t, git.MainBranch, true)
	Init(ctx, testID.OwnerAddress())
	if err := must.Try(func() { Init(ctx, testID.OwnerAddress()) }); err == nil {
		t.Fatal("second init must fail")
	}
}
