package member

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func SetGroup(ctx context.Context, addr gov.GovAddress, name Group) {
	cloned := gov.Clone(ctx, addr)
	chg := SetGroupStageOnly(ctx, cloned.Tree(), name)
	proto.Commit(ctx, cloned.Tree(), chg.Msg)
	cloned.Push(ctx)
}

func SetGroupStageOnly(ctx context.Context, t *git.Tree, name Group) git.ChangeNoResult {
	return groupsKV.Set(ctx, groupsNS, t, name, form.None{})
}

func IsGroup(ctx context.Context, addr gov.GovAddress, name Group) bool {
	return IsGroupLocal(ctx, gov.Clone(ctx, addr).Tree(), name)
}

func IsGroupLocal(ctx context.Context, t *git.Tree, name Group) bool {
	err := must.Try(func() { groupsKV.Get(ctx, groupsNS, t, name) })
	return err == nil
}

func AddGroup(ctx context.Context, addr gov.GovAddress, name Group) {
	cloned := gov.Clone(ctx, addr)
	chg := AddGroupStageOnly(ctx, cloned.Tree(), name)
	proto.Commit(ctx, cloned.Tree(), chg.Msg)
	cloned.Push(ctx)
}

func AddGroupStageOnly(ctx context.Context, t *git.Tree, name Group) git.ChangeNoResult {
	if IsGroupLocal(ctx, t, name) {
		must.Panic(ctx, fmt.Errorf("group already exists"))
	}
	return SetGroupStageOnly(ctx, t, name)
}

func RemoveGroup(ctx context.Context, addr gov.GovAddress, name Group) {
	cloned := gov.Clone(ctx, addr)
	chg := RemoveGroupStageOnly(ctx, cloned.Tree(), name)
	proto.Commit(ctx, cloned.Tree(), chg.Msg)
	cloned.Push(ctx)
}

func RemoveGroupStageOnly(ctx context.Context, t *git.Tree, name Group) git.ChangeNoResult {
	groupsKV.Remove(ctx, groupsNS, t, name)
	groupUsersKKV.RemovePrimary(ctx, groupUsersNS, t, name) // remove memberships
	return git.ChangeNoResult{
		Msg: fmt.Sprintf("Remove group %v", name),
	}
}
