package member

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/mod/gov"
	"github.com/gov4git/gov4git/mod/kv"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func SetUser(ctx context.Context, addr gov.CommunityAddress, name User, acct Account) {
	r, t := gov.CloneCommunity(ctx, addr)
	chg := SetUserStageOnly(ctx, t, name, acct)
	git.Commit(ctx, t, chg.Msg)
	git.Push(ctx, r)
}

func SetUserStageOnly(ctx context.Context, t *git.Tree, name User, user Account) git.ChangeNoResult {
	SetGroupStageOnly(ctx, t, Everybody)
	AddMemberStageOnly(ctx, t, name, Everybody)
	return usersKV.Set(ctx, usersNS, t, name, user)
}

func GetUser(ctx context.Context, addr gov.CommunityAddress, name User) Account {
	_, t := gov.CloneCommunity(ctx, addr)
	x := GetUserLocal(ctx, t, name)
	return x
}

func GetUserLocal(ctx context.Context, t *git.Tree, name User) Account {
	return usersKV.Get(ctx, usersNS, t, name)
}

func AddUser(ctx context.Context, addr gov.CommunityAddress, name User, acct Account) {
	r, t := gov.CloneCommunity(ctx, addr)
	chg := AddUserStageOnly(ctx, t, name, acct)
	git.Commit(ctx, t, chg.Msg)
	git.Push(ctx, r)
}

func AddUserStageOnly(ctx context.Context, t *git.Tree, name User, user Account) git.ChangeNoResult {
	if err := must.Try(func() { GetUserLocal(ctx, t, name) }); err == nil {
		must.Panic(ctx, fmt.Errorf("user already exists"))
	}
	return SetUserStageOnly(ctx, t, name, user)
}

func RemoveUser(ctx context.Context, addr gov.CommunityAddress, name User) {
	r, t := gov.CloneCommunity(ctx, addr)
	chg := RemoveUserStageOnly(ctx, t, name)
	git.Commit(ctx, t, chg.Msg)
	git.Push(ctx, r)
}

func RemoveUserStageOnly(ctx context.Context, t *git.Tree, name User) git.ChangeNoResult {
	usersKV.Remove(ctx, usersNS, t, name)
	userGroupsKKV.RemovePrimary(ctx, userGroupsNS, t, name) // remove memberships
	return git.ChangeNoResult{
		Msg: fmt.Sprintf("Remove user %v", name),
	}
}

// props

func SetUserProp[V form.Form](ctx context.Context, addr gov.CommunityAddress, user User, key string, value V) {
	r, t := gov.CloneCommunity(ctx, addr)
	chg := SetUserPropStageOnly(ctx, t, user, key, value)
	git.Commit(ctx, t, chg.Msg)
	git.Push(ctx, r)
}

func SetUserPropStageOnly[V form.Form](ctx context.Context, t *git.Tree, user User, key string, value V) git.ChangeNoResult {
	propKV := kv.KV[string, V]{}
	return propKV.Set(ctx, usersKV.KeyNS(usersNS, user), t, key, value)
}

func GetUserProp[V form.Form](ctx context.Context, addr gov.CommunityAddress, user User, key string) V {
	_, t := gov.CloneCommunity(ctx, addr)
	x := GetUserPropLocal[V](ctx, t, user, key)
	return x
}

func GetUserPropLocal[V form.Form](ctx context.Context, t *git.Tree, user User, key string) V {
	propKV := kv.KV[string, V]{}
	return propKV.Get(ctx, usersKV.KeyNS(usersNS, user), t, key)
}
