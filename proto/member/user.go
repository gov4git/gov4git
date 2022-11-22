package member

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/kv"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func SetUser(ctx context.Context, addr gov.GovAddress, name User, acct Account) {
	r, t := gov.Clone(ctx, addr)
	chg := SetUserStageOnly(ctx, t, name, acct)
	git.Commit(ctx, t, chg.Msg)
	git.Push(ctx, r)
}

func SetUserStageOnly(ctx context.Context, t *git.Tree, name User, user Account) git.ChangeNoResult {
	SetGroupStageOnly(ctx, t, Everybody)
	AddMemberStageOnly(ctx, t, name, Everybody)
	return usersKV.Set(ctx, usersNS, t, name, user)
}

func GetUser(ctx context.Context, addr gov.GovAddress, name User) Account {
	_, t := gov.Clone(ctx, addr)
	x := GetUserLocal(ctx, t, name)
	return x
}

func GetUserLocal(ctx context.Context, t *git.Tree, name User) Account {
	return usersKV.Get(ctx, usersNS, t, name)
}

func AddUser(ctx context.Context, addr gov.GovAddress, name User, acct Account) {
	r, t := gov.Clone(ctx, addr)
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

func RemoveUser(ctx context.Context, addr gov.GovAddress, name User) {
	r, t := gov.Clone(ctx, addr)
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

// set prop

func SetUserProp[V form.Form](ctx context.Context, addr gov.GovAddress, user User, key string, value V) {
	r, t := gov.Clone(ctx, addr)
	chg := SetUserPropStageOnly(ctx, t, user, key, value)
	git.Commit(ctx, t, chg.Msg)
	git.Push(ctx, r)
}

func SetUserPropStageOnly[V form.Form](ctx context.Context, t *git.Tree, user User, key string, value V) git.ChangeNoResult {
	propKV := kv.KV[string, V]{}
	return propKV.Set(ctx, usersKV.KeyNS(usersNS, user), t, key, value)
}

// get prop

func GetUserProp[V form.Form](ctx context.Context, addr gov.GovAddress, user User, key string) V {
	_, t := gov.Clone(ctx, addr)
	x := GetUserPropLocal[V](ctx, t, user, key)
	return x
}

func GetUserPropLocal[V form.Form](ctx context.Context, t *git.Tree, user User, key string) V {
	propKV := kv.KV[string, V]{}
	return propKV.Get(ctx, usersKV.KeyNS(usersNS, user), t, key)
}

func GetUserPropOrDefault[V form.Form](ctx context.Context, addr gov.GovAddress, user User, key string, default_ V) V {
	r := default_
	r, _ = must.Try1(func() V { return GetUserProp[V](ctx, addr, user, key) })
	return r
}

func GetUserPropLocalOrDefault[V form.Form](ctx context.Context, t *git.Tree, user User, key string, default_ V) V {
	r := default_
	r, _ = must.Try1(func() V { return GetUserPropLocal[V](ctx, t, user, key) })
	return r
}

// lookup

func LookupUserByAddress(ctx context.Context, govAddr gov.GovAddress, userAddr id.PublicAddress) []User {
	_, t := gov.Clone(ctx, govAddr)
	return LookupUserByAddressLocal(ctx, t, userAddr)
}

func LookupUserByAddressLocal(ctx context.Context, t *git.Tree, userAddr id.PublicAddress) []User {
	us := usersKV.ListKeys(ctx, usersNS, t)
	r := []User{}
	for _, u := range us {
		acct := GetUserLocal(ctx, t, u)
		if acct.PublicAddress == userAddr {
			r = append(r, u)
		}
	}
	return r
}
