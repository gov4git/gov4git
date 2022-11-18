package bureau

import (
	"testing"

	"github.com/gov4git/gov4git/proto/balance"
	"github.com/gov4git/gov4git/proto/bureau"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/gov4git/test"
	"github.com/gov4git/lib4git/testutil"
)

func TestBureau(t *testing.T) {
	ctx := testutil.NewCtx()
	cty := test.NewTestCommunity(t, ctx, 2)

	usd := balance.Balance("usd")

	// credit user 0 with some cash
	balance.Set(ctx, cty.Community(), cty.MemberUser(0), usd, 3.0)

	// user 0 requests transfer to user 1
	bureau.Transfer(ctx, cty.MemberOwner(0), cty.Community(), member.User(""), usd, cty.MemberUser(1), usd, 1.0)

	// process request
	bureau.Process(ctx, cty.Organizer(), member.Everybody)

	// get resulting balances
	u0 := balance.Get(ctx, cty.Community(), cty.MemberUser(0), usd)
	u1 := balance.Get(ctx, cty.Community(), cty.MemberUser(1), usd)

	if u0 != 2.0 {
		t.Errorf("expecting 2, got %v", u0)
	}
	if u1 != 1.0 {
		t.Errorf("expecting 1, got %v", u1)
	}
}
