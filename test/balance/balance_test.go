package balance

import (
	"testing"

	"github.com/gov4git/gov4git/proto/balance"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/test"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/testutil"
)

func TestBalance(t *testing.T) {
	ctx := testutil.NewCtx(t, true)
	cty := test.NewTestCommunity(t, ctx, 2)

	bal := balance.Balance("test_balance")

	// test set/get roundtrip
	balance.Set(ctx, cty.Gov(), cty.MemberUser(0), bal, 30.0)
	actual1 := balance.Get(ctx, cty.Gov(), cty.MemberUser(0), bal)
	if actual1 != 30.0 {
		t.Errorf("expecting %v, got %v", 30.0, actual1)
	}

	// test balance transfer
	cloned := gov.Clone(ctx, cty.Gov())
	balance.TransferStageOnly(ctx, cloned.Tree(), cty.MemberUser(0), bal, cty.MemberUser(1), bal, 10.0)
	git.Commit(ctx, cloned.Tree(), "test commit")
	cloned.Push(ctx)
	actual2 := balance.Get(ctx, cty.Gov(), cty.MemberUser(1), bal)
	if actual2 != 10.0 {
		t.Errorf("expecting %v, got %v", 10.0, actual2)
	}

}
