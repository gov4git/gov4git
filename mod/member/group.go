package member

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/lib/must"
)

func SetGroup(ctx context.Context, t *git.Tree, name Group) git.ChangeNoResult {
	return groupsKV.Set(ctx, groupsNS, t, name, form.None{})
}

func IsGroup(ctx context.Context, t *git.Tree, name Group) bool {
	err := must.Try(func() { groupsKV.Get(ctx, groupsNS, t, name) })
	return err == nil
}

func AddGroup(ctx context.Context, t *git.Tree, name Group) git.ChangeNoResult {
	if IsGroup(ctx, t, name) {
		must.Panic(ctx, fmt.Errorf("group already exists"))
	}
	return SetGroup(ctx, t, name)
}

func RemoveGroup(ctx context.Context, t *git.Tree, name Group) git.ChangeNoResult {
	groupsKV.Remove(ctx, groupsNS, t, name)
	groupUsersKKV.RemovePrimary(ctx, groupUsersNS, t, name) // remove memberships
	return git.ChangeNoResult{
		Msg: fmt.Sprintf("Remove group %v", name),
	}
}
