package member

import (
	"testing"

	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/gov4git/test"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/testutil"
)

func TestUserAddRemove(t *testing.T) {
	ctx := testutil.NewCtx()
	cty := test.NewTestCommunity(t, ctx, 2)

	name := member.User("testuser")

	member.AddUserByPublicAddress(ctx, cty.Gov(), name, cty.MemberOwner(0).Public)

	acct := member.GetUser(ctx, cty.Gov(), name)
	if acct.PublicAddress != cty.MemberOwner(0).Public {
		t.Errorf("expecting %v, got %v", cty.MemberOwner(0).Public, acct.PublicAddress)
	}

	member.RemoveUser(ctx, cty.Gov(), name)

	if must.Try(func() { member.GetUser(ctx, cty.Gov(), name) }) == nil {
		t.Errorf("expecting user to be missing")
	}
}
