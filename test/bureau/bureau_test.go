package bureau

import (
	"testing"

	"github.com/gov4git/gov4git/proto/account"
	"github.com/gov4git/gov4git/proto/bureau"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/gov4git/runtime"
	"github.com/gov4git/gov4git/test"
	"github.com/gov4git/lib4git/testutil"
)

func TestBureau(t *testing.T) {
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	// credit user 0 with some cash
	account.Deposit(ctx, cty.Gov(), cty.MemberAccountID(0), account.H(account.PluralAsset, 3.0), "test")

	// user 0 requests transfer to user 1
	bureau.Transfer(ctx, cty.MemberOwner(0), cty.Gov(), member.User(""), cty.MemberUser(1), 1.0)

	// process request
	bureau.Process(ctx, cty.Organizer(), member.Everybody)

	// get resulting balances
	u0 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(0)).Balance(account.PluralAsset).Quantity
	u1 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(1)).Balance(account.PluralAsset).Quantity

	if u0 != 2.0 {
		t.Errorf("expecting 2, got %v", u0)
	}
	if u1 != 1.0 {
		t.Errorf("expecting 1, got %v", u1)
	}
}
