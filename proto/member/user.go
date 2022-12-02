package member

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/kv"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func SetUser(ctx context.Context, addr gov.GovAddress, name User, acct Account) {
	cloned := gov.Clone(ctx, addr)
	chg := SetUserStageOnly(ctx, cloned.Tree(), name, acct)
	proto.Commit(ctx, cloned.Tree(), chg.Msg)
	cloned.Push(ctx)
}

func SetUserStageOnly(ctx context.Context, t *git.Tree, name User, user Account) git.ChangeNoResult {
	SetGroupStageOnly(ctx, t, Everybody)
	AddMemberStageOnly(ctx, t, name, Everybody)
	return usersKV.Set(ctx, usersNS, t, name, user)
}

func GetUser(ctx context.Context, addr gov.GovAddress, name User) Account {
	return GetUserLocal(ctx, gov.Clone(ctx, addr).Tree(), name)
}

func GetUserLocal(ctx context.Context, t *git.Tree, name User) Account {
	return usersKV.Get(ctx, usersNS, t, name)
}

func AddUserByPublicAddress(ctx context.Context, govAddr gov.GovAddress, name User, userAddr id.PublicAddress) {
	cred := id.FetchPublicCredentials(ctx, userAddr)
	AddUser(ctx, govAddr, name, Account{ID: cred.ID, PublicAddress: userAddr})
}

func AddUserByPublicAddressStageOnly(ctx context.Context, t *git.Tree, name User, userAddr id.PublicAddress) {
	cred := id.FetchPublicCredentials(ctx, userAddr)
	AddUserStageOnly(ctx, t, name, Account{ID: cred.ID, PublicAddress: userAddr})
}

func AddUser(ctx context.Context, addr gov.GovAddress, name User, acct Account) {
	cloned := gov.Clone(ctx, addr)
	chg := AddUserStageOnly(ctx, cloned.Tree(), name, acct)
	proto.Commit(ctx, cloned.Tree(), chg.Msg)
	cloned.Push(ctx)
}

func AddUserStageOnly(ctx context.Context, t *git.Tree, name User, user Account) git.ChangeNoResult {
	if err := must.Try(func() { GetUserLocal(ctx, t, name) }); err == nil {
		must.Panic(ctx, fmt.Errorf("user already exists"))
	}
	return SetUserStageOnly(ctx, t, name, user)
}

func RemoveUser(ctx context.Context, addr gov.GovAddress, name User) {
	cloned := gov.Clone(ctx, addr)
	chg := RemoveUserStageOnly(ctx, cloned.Tree(), name)
	proto.Commit(ctx, cloned.Tree(), chg.Msg)
	cloned.Push(ctx)
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
	cloned := gov.Clone(ctx, addr)
	chg := SetUserPropStageOnly(ctx, cloned.Tree(), user, key, value)
	proto.Commit(ctx, cloned.Tree(), chg.Msg)
	cloned.Push(ctx)
}

func SetUserPropStageOnly[V form.Form](ctx context.Context, t *git.Tree, user User, key string, value V) git.ChangeNoResult {
	propKV := kv.KV[string, V]{}
	return propKV.Set(ctx, usersKV.KeyNS(usersNS, user), t, key, value)
}

// get prop

func GetUserProp[V form.Form](ctx context.Context, addr gov.GovAddress, user User, key string) V {
	return GetUserPropLocal[V](ctx, gov.Clone(ctx, addr).Tree(), user, key)
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
	return LookupUserByAddressLocal(ctx, gov.Clone(ctx, govAddr).Tree(), userAddr)
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

func LookupUserByID(ctx context.Context, govAddr gov.GovAddress, userID id.ID) []User {
	return LookupUserByIDLocal(ctx, gov.Clone(ctx, govAddr).Tree(), userID)
}

func LookupUserByIDLocal(ctx context.Context, t *git.Tree, userID id.ID) []User {
	us := usersKV.ListKeys(ctx, usersNS, t)
	r := []User{}
	for _, u := range us {
		acct := GetUserLocal(ctx, t, u)
		if acct.ID == userID {
			r = append(r, u)
		}
	}
	return r
}
