package member

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history/trace"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func SetGroup(ctx context.Context, addr gov.Address, name Group) {
	cloned := gov.Clone(ctx, addr)
	chg := SetGroup_StageOnly(ctx, cloned, name)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
}

func SetGroup_StageOnly(ctx context.Context, cloned gov.Cloned, name Group) git.ChangeNoResult {
	return groupsKV.Set(ctx, groupsNS, cloned.Tree(), name, form.None{})
}

func IsGroup(ctx context.Context, addr gov.Address, name Group) bool {
	return IsGroup_Local(ctx, gov.Clone(ctx, addr), name)
}

func IsGroup_Local(ctx context.Context, cloned gov.Cloned, name Group) bool {
	err := must.Try(func() { groupsKV.Get(ctx, groupsNS, cloned.Tree(), name) })
	return err == nil
}

func AddGroup(ctx context.Context, addr gov.Address, name Group) {
	cloned := gov.Clone(ctx, addr)
	chg := AddGroup_StageOnly(ctx, cloned, name)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
}

func AddGroup_StageOnly(ctx context.Context, cloned gov.Cloned, name Group) git.ChangeNoResult {
	if IsGroup_Local(ctx, cloned, name) {
		must.Panic(ctx, fmt.Errorf("group already exists"))
	}
	chg := SetGroup_StageOnly(ctx, cloned, name)

	// log
	trace.Log_StageOnly(ctx, cloned, &trace.Event{
		Op:     "group_add",
		Args:   trace.M{"name": name},
		Result: nil,
	})

	return chg
}

func RemoveGroup(ctx context.Context, addr gov.Address, name Group) {
	cloned := gov.Clone(ctx, addr)
	chg := RemoveGroup_StageOnly(ctx, cloned, name)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
}

func RemoveGroup_StageOnly(ctx context.Context, cloned gov.Cloned, name Group) git.ChangeNoResult {
	groupsKV.Remove(ctx, groupsNS, cloned.Tree(), name)
	groupUsersKKV.RemovePrimary(ctx, groupUsersNS, cloned.Tree(), name) // remove memberships

	// log
	trace.Log_StageOnly(ctx, cloned, &trace.Event{
		Op:     "group_remove",
		Args:   trace.M{"name": name},
		Result: nil,
	})

	return git.NewChangeNoResult(fmt.Sprintf("Remove group %v", name), "member_remove_group")
}
