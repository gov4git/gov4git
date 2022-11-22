// Package member implements governance member management services
package member

import (
	"context"
	"fmt"

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
	r, t := gov.Clone(ctx, addr)
	chg := AddMemberStageOnly(ctx, t, user, group)
	git.Commit(ctx, t, chg.Msg)
	git.Push(ctx, r)
}

func AddMemberStageOnly(ctx context.Context, t *git.Tree, user User, group Group) git.ChangeNoResult {
	userGroupsKKV.Set(ctx, userGroupsNS, t, user, group, true)
	groupUsersKKV.Set(ctx, groupUsersNS, t, group, user, true)
	return git.ChangeNoResult{
		Msg: fmt.Sprintf("Added user %v to group %v", user, group),
	}
}

func IsMember(ctx context.Context, addr gov.GovAddress, user User, group Group) bool {
	_, t := gov.Clone(ctx, addr)
	x := IsMemberLocal(ctx, t, user, group)
	return x
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
	r, t := gov.Clone(ctx, addr)
	chg := RemoveMemberStageOnly(ctx, t, user, group)
	git.Commit(ctx, t, chg.Msg)
	git.Push(ctx, r)
}

func RemoveMemberStageOnly(ctx context.Context, t *git.Tree, user User, group Group) git.ChangeNoResult {
	userGroupsKKV.Remove(ctx, userGroupsNS, t, user, group)
	groupUsersKKV.Remove(ctx, groupUsersNS, t, group, user)
	return git.ChangeNoResult{
		Msg: fmt.Sprintf("Removed user %v from group %v", user, group),
	}
}

func ListUserGroups(ctx context.Context, addr gov.GovAddress, user User) []Group {
	_, t := gov.Clone(ctx, addr)
	x := ListUserGroupsLocal(ctx, t, user)
	return x
}

func ListUserGroupsLocal(ctx context.Context, t *git.Tree, user User) []Group {
	return userGroupsKKV.ListSecondaryKeys(ctx, userGroupsNS, t, user)
}

func ListGroupUsers(ctx context.Context, addr gov.GovAddress, group Group) []User {
	_, t := gov.Clone(ctx, addr)
	x := ListGroupUsersLocal(ctx, t, group)
	return x
}

func ListGroupUsersLocal(ctx context.Context, t *git.Tree, group Group) []User {
	return groupUsersKKV.ListSecondaryKeys(ctx, groupUsersNS, t, group)
}
