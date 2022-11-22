package member

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func SetGroup(ctx context.Context, addr gov.GovAddress, name Group) {
	r, t := gov.Clone(ctx, addr)
	chg := SetGroupStageOnly(ctx, t, name)
	git.Commit(ctx, t, chg.Msg)
	git.Push(ctx, r)
}

func SetGroupStageOnly(ctx context.Context, t *git.Tree, name Group) git.ChangeNoResult {
	return groupsKV.Set(ctx, groupsNS, t, name, form.None{})
}

func IsGroup(ctx context.Context, addr gov.GovAddress, name Group) bool {
	_, t := gov.Clone(ctx, addr)
	x := IsGroupLocal(ctx, t, name)
	return x
}

func IsGroupLocal(ctx context.Context, t *git.Tree, name Group) bool {
	err := must.Try(func() { groupsKV.Get(ctx, groupsNS, t, name) })
	return err == nil
}

func AddGroup(ctx context.Context, addr gov.GovAddress, name Group) {
	r, t := gov.Clone(ctx, addr)
	chg := AddGroupStageOnly(ctx, t, name)
	git.Commit(ctx, t, chg.Msg)
	git.Push(ctx, r)
}

func AddGroupStageOnly(ctx context.Context, t *git.Tree, name Group) git.ChangeNoResult {
	if IsGroupLocal(ctx, t, name) {
		must.Panic(ctx, fmt.Errorf("group already exists"))
	}
	return SetGroupStageOnly(ctx, t, name)
}

func RemoveGroup(ctx context.Context, addr gov.GovAddress, name Group) {
	r, t := gov.Clone(ctx, addr)
	chg := RemoveGroupStageOnly(ctx, t, name)
	git.Commit(ctx, t, chg.Msg)
	git.Push(ctx, r)
}

func RemoveGroupStageOnly(ctx context.Context, t *git.Tree, name Group) git.ChangeNoResult {
	groupsKV.Remove(ctx, groupsNS, t, name)
	groupUsersKKV.RemovePrimary(ctx, groupUsersNS, t, name) // remove memberships
	return git.ChangeNoResult{
		Msg: fmt.Sprintf("Remove group %v", name),
	}
}
