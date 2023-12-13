// Package member implements governance member management services
package member

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/history"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

const (
	Everybody = Group("everybody")
)

type User string

func (u User) IsNone() bool {
	return u == ""
}

type Group string

func AddMember(ctx context.Context, addr gov.Address, user User, group Group) {
	cloned := gov.Clone(ctx, addr)
	chg := AddMember_StageOnly(ctx, cloned, user, group)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
}

func AddMember_StageOnly(ctx context.Context, cloned gov.Cloned, user User, group Group) git.ChangeNoResult {
	userGroupsKKV.Set(ctx, userGroupsNS, cloned.Tree(), user, group, true)
	groupUsersKKV.Set(ctx, groupUsersNS, cloned.Tree(), group, user, true)

	// log
	history.Log_StageOnly(ctx, cloned, &history.Event{
		Op: &history.Op{
			Op:     "add_user_to_group",
			Args:   history.M{"user": user, "group": group},
			Result: nil,
		},
	})

	return git.NewChangeNoResult(fmt.Sprintf("Added user %v to group %v", user, group), "member_add_member")
}

func IsMember(ctx context.Context, addr gov.Address, user User, group Group) bool {
	return IsMember_Local(ctx, gov.Clone(ctx, addr), user, group)
}

func IsMember_Local(ctx context.Context, cloned gov.Cloned, user User, group Group) bool {
	var userHasGroup, groupHasUser bool
	must.Try(
		func() { userHasGroup = userGroupsKKV.Get(ctx, userGroupsNS, cloned.Tree(), user, group) },
	)
	must.Try(
		func() { groupHasUser = groupUsersKKV.Get(ctx, groupUsersNS, cloned.Tree(), group, user) },
	)
	return userHasGroup && groupHasUser
}

func RemoveMember(ctx context.Context, addr gov.Address, user User, group Group) {
	cloned := gov.Clone(ctx, addr)
	chg := RemoveMember_StageOnly(ctx, cloned, user, group)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
}

func RemoveMember_StageOnly(ctx context.Context, cloned gov.Cloned, user User, group Group) git.ChangeNoResult {
	userGroupsKKV.Remove(ctx, userGroupsNS, cloned.Tree(), user, group)
	groupUsersKKV.Remove(ctx, groupUsersNS, cloned.Tree(), group, user)

	// log
	history.Log_StageOnly(ctx, cloned, &history.Event{
		Op: &history.Op{
			Op:     "remove_user_from_group",
			Args:   history.M{"user": user, "group": group},
			Result: nil,
		},
	})

	return git.NewChangeNoResult(fmt.Sprintf("Removed user %v from group %v", user, group), "member_remove_member")
}

func ListUserGroups(ctx context.Context, addr gov.Address, user User) []Group {
	return ListUserGroups_Local(ctx, gov.Clone(ctx, addr), user)
}

func ListUserGroups_Local(ctx context.Context, cloned gov.Cloned, user User) []Group {
	return userGroupsKKV.ListSecondaryKeys(ctx, userGroupsNS, cloned.Tree(), user)
}

func ListGroupUsers(ctx context.Context, addr gov.Address, group Group) []User {
	return ListGroupUsers_Local(ctx, gov.Clone(ctx, addr), group)
}

func ListGroupUsers_Local(ctx context.Context, cloned gov.Cloned, group Group) []User {
	return groupUsersKKV.ListSecondaryKeys(ctx, groupUsersNS, cloned.Tree(), group)
}
