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
		Home: id.HomeAddress{
			Repo:   git.URL("http://1"),
			Branch: git.MainBranch,
		},
	}
	AddUserStageOnly(ctx, wt, u1, r1)
	r1Got := GetUserLocal(ctx, wt, u1)
	if r1 != r1Got {
		t.Fatalf("expecting %v, got %v", r1, r1Got)
	}

	if !IsMemberLocal(ctx, wt, u1, Everybody) {
		t.Fatalf("expecting is member")
	}

	allUsers := ListGroupUsersLocal(ctx, wt, Everybody)
	if len(allUsers) != 1 || allUsers[0] != u1 {
		t.Fatalf("unexpected list of users in group everybody")
	}

	allGroups := ListUserGroupsLocal(ctx, wt, u1)
	if len(allGroups) != 1 || allGroups[0] != Everybody {
		t.Fatalf("unexpected list of groups for user")
	}

	RemoveUserStageOnly(ctx, wt, u1)
	err := must.Try(func() {
		GetUserLocal(ctx, wt, u1)
	})
	if err == nil {
		t.Fatalf("expecting error")
	}

	if IsMemberLocal(ctx, wt, u1, Everybody) {
		t.Fatalf("expecting no membership")
	}
}
