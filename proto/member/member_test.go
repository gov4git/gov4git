package member

import (
	"context"
	"testing"

	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/testutil"
)

func TestMember(t *testing.T) {
	base.LogVerbosely()
	ctx := context.Background()
	repo := testutil.InitPlainRepo(t, ctx)
	wt := git.Worktree(ctx, repo.Repo)

	u1 := User("user1")
	r1 := Account{
		PublicAddress: id.PublicAddress{
			Repo:   git.URL("http://1"),
			Branch: git.MainBranch,
		},
	}
	AddUser_StageOnly(ctx, wt, u1, r1)
	r1Got := GetUser_Local(ctx, wt, u1)
	if r1 != r1Got {
		t.Fatalf("expecting %v, got %v", r1, r1Got)
	}

	if !IsMember_Local(ctx, wt, u1, Everybody) {
		t.Fatalf("expecting is member")
	}

	allUsers := ListGroupUsers_Local(ctx, wt, Everybody)
	if len(allUsers) != 1 || allUsers[0] != u1 {
		t.Fatalf("unexpected list of users in group everybody")
	}

	allGroups := ListUserGroups_Local(ctx, wt, u1)
	if len(allGroups) != 1 || allGroups[0] != Everybody {
		t.Fatalf("unexpected list of groups for user")
	}

	RemoveUser_StageOnly(ctx, wt, u1)
	err := must.Try(func() {
		GetUser_Local(ctx, wt, u1)
	})
	if err == nil {
		t.Fatalf("expecting error")
	}

	if IsMember_Local(ctx, wt, u1, Everybody) {
		t.Fatalf("expecting no membership")
	}
}
