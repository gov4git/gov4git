// Package member implements governance member management services
package member

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

const (
	Everybody = Group("everybody")
)

type User string
type Group string

func AddMember(ctx context.Context, addr gov.GovPublicAddress, user User, group Group) {
	cloned := gov.Clone(ctx, addr)
	chg := AddMember_StageOnly(ctx, cloned.Tree(), user, group)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
}

func AddMember_StageOnly(ctx context.Context, t *git.Tree, user User, group Group) git.ChangeNoResult {
	userGroupsKKV.Set(ctx, userGroupsNS, t, user, group, true)
	groupUsersKKV.Set(ctx, groupUsersNS, t, group, user, true)
	return git.NewChangeNoResult(fmt.Sprintf("Added user %v to group %v", user, group), "member_add_member")
}

func IsMember(ctx context.Context, addr gov.GovPublicAddress, user User, group Group) bool {
	return IsMember_Local(ctx, gov.Clone(ctx, addr).Tree(), user, group)
}

func IsMember_Local(ctx context.Context, t *git.Tree, user User, group Group) bool {
	var userHasGroup, groupHasUser bool
	must.Try(
		func() { userHasGroup = userGroupsKKV.Get(ctx, userGroupsNS, t, user, group) },
	)
	must.Try(
		func() { groupHasUser = groupUsersKKV.Get(ctx, groupUsersNS, t, group, user) },
	)
	return userHasGroup && groupHasUser
}

func RemoveMember(ctx context.Context, addr gov.GovPublicAddress, user User, group Group) {
	cloned := gov.Clone(ctx, addr)
	chg := RemoveMember_StageOnly(ctx, cloned.Tree(), user, group)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
}

func RemoveMember_StageOnly(ctx context.Context, t *git.Tree, user User, group Group) git.ChangeNoResult {
	userGroupsKKV.Remove(ctx, userGroupsNS, t, user, group)
	groupUsersKKV.Remove(ctx, groupUsersNS, t, group, user)
	return git.NewChangeNoResult(fmt.Sprintf("Removed user %v from group %v", user, group), "member_remove_member")
}

func ListUserGroups(ctx context.Context, addr gov.GovPublicAddress, user User) []Group {
	return ListUserGroups_Local(ctx, gov.Clone(ctx, addr).Tree(), user)
}

func ListUserGroups_Local(ctx context.Context, t *git.Tree, user User) []Group {
	return userGroupsKKV.ListSecondaryKeys(ctx, userGroupsNS, t, user)
}

func ListGroupUsers(ctx context.Context, addr gov.GovPublicAddress, group Group) []User {
	return ListGroupUsers_Local(ctx, gov.Clone(ctx, addr).Tree(), group)
}

func ListGroupUsers_Local(ctx context.Context, t *git.Tree, group Group) []User {
	return groupUsersKKV.ListSecondaryKeys(ctx, groupUsersNS, t, group)
}
