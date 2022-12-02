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

func AddMember(ctx context.Context, addr gov.GovAddress, user User, group Group) {
	cloned := gov.Clone(ctx, addr)
	chg := AddMemberStageOnly(ctx, cloned.Tree(), user, group)
	git.Commit(ctx, cloned.Tree(), chg.Msg)
	cloned.Push(ctx)
}

func AddMemberStageOnly(ctx context.Context, t *git.Tree, user User, group Group) git.ChangeNoResult {
	userGroupsKKV.Set(ctx, userGroupsNS, t, user, group, true)
	groupUsersKKV.Set(ctx, groupUsersNS, t, group, user, true)
	return git.ChangeNoResult{
		Msg: fmt.Sprintf("Added user %v to group %v", user, group),
	}
}

func IsMember(ctx context.Context, addr gov.GovAddress, user User, group Group) bool {
	return IsMemberLocal(ctx, gov.Clone(ctx, addr).Tree(), user, group)
}

func IsMemberLocal(ctx context.Context, t *git.Tree, user User, group Group) bool {
	var userHasGroup, groupHasUser bool
	must.Try(
		func() { userHasGroup = userGroupsKKV.Get(ctx, userGroupsNS, t, user, group) },
	)
	must.Try(
		func() { groupHasUser = groupUsersKKV.Get(ctx, groupUsersNS, t, group, user) },
	)
	return userHasGroup && groupHasUser
}

func RemoveMember(ctx context.Context, addr gov.GovAddress, user User, group Group) {
	cloned := gov.Clone(ctx, addr)
	chg := RemoveMemberStageOnly(ctx, cloned.Tree(), user, group)
	proto.Commit(ctx, cloned.Tree(), chg.Msg)
	cloned.Push(ctx)
}

func RemoveMemberStageOnly(ctx context.Context, t *git.Tree, user User, group Group) git.ChangeNoResult {
	userGroupsKKV.Remove(ctx, userGroupsNS, t, user, group)
	groupUsersKKV.Remove(ctx, groupUsersNS, t, group, user)
	return git.ChangeNoResult{
		Msg: fmt.Sprintf("Removed user %v from group %v", user, group),
	}
}

func ListUserGroups(ctx context.Context, addr gov.GovAddress, user User) []Group {
	return ListUserGroupsLocal(ctx, gov.Clone(ctx, addr).Tree(), user)
}

func ListUserGroupsLocal(ctx context.Context, t *git.Tree, user User) []Group {
	return userGroupsKKV.ListSecondaryKeys(ctx, userGroupsNS, t, user)
}

func ListGroupUsers(ctx context.Context, addr gov.GovAddress, group Group) []User {
	return ListGroupUsersLocal(ctx, gov.Clone(ctx, addr).Tree(), group)
}

func ListGroupUsersLocal(ctx context.Context, t *git.Tree, group Group) []User {
	return groupUsersKKV.ListSecondaryKeys(ctx, groupUsersNS, t, group)
}
