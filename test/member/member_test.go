package member

import (
	"testing"

	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/gov4git/test"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/testutil"
	"github.com/gov4git/lib4git/util"
)

func TestUserAddRemove(t *testing.T) {
	ctx := testutil.NewCtx(t, true)
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

func TestGroupAddRemove(t *testing.T) {
	ctx := testutil.NewCtx(t, true)
	cty := test.NewTestCommunity(t, ctx, 2)

	u1 := member.User("testuser1")
	g1 := member.Group("testgroup1")

	// add user to group, check user is a member
	member.AddUserByPublicAddress(ctx, cty.Gov(), u1, cty.MemberOwner(0).Public)
	member.AddGroup(ctx, cty.Gov(), g1)
	member.AddMember(ctx, cty.Gov(), u1, g1)
	users1 := member.ListGroupUsers(ctx, cty.Gov(), g1)
	if len(users1) != 1 || users1[0] != u1 {
		t.Fatalf("expecting %v, got %v", []member.User{u1}, users1)
	}
	if !member.IsMember(ctx, cty.Gov(), u1, g1) {
		t.Errorf("expecting user to be a member")
	}

	// check user's groups are `everybody` and `testgroup1`
	groups1 := member.ListUserGroups(ctx, cty.Gov(), u1)
	if !util.IsIn(g1, groups1...) {
		t.Errorf("expecting group to be in user's memberships")
	}

	// remove user from group, check group has no members
	member.RemoveMember(ctx, cty.Gov(), u1, g1)
	users2 := member.ListGroupUsers(ctx, cty.Gov(), g1)
	if len(users2) != 0 {
		t.Fatalf("expecting no members, got %v", users2)
	}

	// verify user is in `everybody`group
	users3 := member.ListGroupUsers(ctx, cty.Gov(), member.Everybody)
	if !util.IsIn(u1, users3...) {
		t.Fatalf("expecting %v, got %v", []member.User{u1}, users3)
	}
}
