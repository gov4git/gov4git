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
	chg := SetGroup_StageOnly(ctx, cloned.Tree(), name)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
}

func SetGroup_StageOnly(ctx context.Context, t *git.Tree, name Group) git.ChangeNoResult {
	return groupsKV.Set(ctx, groupsNS, t, name, form.None{})
}

func IsGroup(ctx context.Context, addr gov.GovAddress, name Group) bool {
	return IsGroup_Local(ctx, gov.Clone(ctx, addr).Tree(), name)
}

func IsGroup_Local(ctx context.Context, t *git.Tree, name Group) bool {
	err := must.Try(func() { groupsKV.Get(ctx, groupsNS, t, name) })
	return err == nil
}

func AddGroup(ctx context.Context, addr gov.GovAddress, name Group) {
	cloned := gov.Clone(ctx, addr)
	chg := AddGroup_StageOnly(ctx, cloned.Tree(), name)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
}

func AddGroup_StageOnly(ctx context.Context, t *git.Tree, name Group) git.ChangeNoResult {
	if IsGroup_Local(ctx, t, name) {
		must.Panic(ctx, fmt.Errorf("group already exists"))
	}
	return SetGroup_StageOnly(ctx, t, name)
}

func RemoveGroup(ctx context.Context, addr gov.GovAddress, name Group) {
	cloned := gov.Clone(ctx, addr)
	chg := RemoveGroup_StageOnly(ctx, cloned.Tree(), name)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
}

func RemoveGroup_StageOnly(ctx context.Context, t *git.Tree, name Group) git.ChangeNoResult {
	groupsKV.Remove(ctx, groupsNS, t, name)
	groupUsersKKV.RemovePrimary(ctx, groupUsersNS, t, name) // remove memberships
	return git.NewChangeNoResult(fmt.Sprintf("Remove group %v", name), "member_remove_group")
}
