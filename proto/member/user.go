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
	chg := SetUser_StageOnly(ctx, cloned.Tree(), name, acct)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
}

func SetUser_StageOnly(ctx context.Context, t *git.Tree, name User, user Account) git.ChangeNoResult {
	SetGroup_StageOnly(ctx, t, Everybody)           // create everybody group, if it doesn't exist
	AddMember_StageOnly(ctx, t, name, Everybody)    // add membership of user to everybody
	return usersKV.Set(ctx, usersNS, t, name, user) // create the user record
}

func GetUser(ctx context.Context, addr gov.GovAddress, name User) Account {
	return GetUser_Local(ctx, gov.Clone(ctx, addr).Tree(), name)
}

func GetUser_Local(ctx context.Context, t *git.Tree, name User) Account {
	return usersKV.Get(ctx, usersNS, t, name)
}

func AddUserByPublicAddress(ctx context.Context, govAddr gov.GovAddress, name User, userAddr id.PublicAddress) {
	cred := id.FetchPublicCredentials(ctx, userAddr)
	AddUser(ctx, govAddr, name, Account{ID: cred.ID, PublicAddress: userAddr})
}

func AddUserByPublicAddress_StageOnly(ctx context.Context, t *git.Tree, name User, userAddr id.PublicAddress) {
	cred := id.FetchPublicCredentials(ctx, userAddr)
	AddUser_StageOnly(ctx, t, name, Account{ID: cred.ID, PublicAddress: userAddr})
}

func AddUser(ctx context.Context, addr gov.GovAddress, name User, acct Account) {
	cloned := gov.Clone(ctx, addr)
	chg := AddUser_StageOnly(ctx, cloned.Tree(), name, acct)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
}

func AddUser_StageOnly(ctx context.Context, t *git.Tree, name User, user Account) git.ChangeNoResult {
	if err := must.Try(func() { GetUser_Local(ctx, t, name) }); err == nil {
		must.Panic(ctx, fmt.Errorf("user already exists"))
	}
	return SetUser_StageOnly(ctx, t, name, user)
}

func RemoveUser(ctx context.Context, addr gov.GovAddress, name User) {
	cloned := gov.Clone(ctx, addr)
	chg := RemoveUser_StageOnly(ctx, cloned.Tree(), name)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
}

func RemoveUser_StageOnly(ctx context.Context, t *git.Tree, name User) git.ChangeNoResult {
	// remove all group memberships of the user
	for _, g := range ListUserGroups_Local(ctx, t, name) {
		RemoveMember_StageOnly(ctx, t, name, g)
	}
	// remove user record
	usersKV.Remove(ctx, usersNS, t, name)
	return git.NewChangeNoResult(fmt.Sprintf("Remove user %v", name), "member_remove_user")
}

// set prop

func SetUserProp[V form.Form](ctx context.Context, addr gov.GovAddress, user User, key string, value V) {
	cloned := gov.Clone(ctx, addr)
	chg := SetUserProp_StageOnly(ctx, cloned.Tree(), user, key, value)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
}

func SetUserProp_StageOnly[V form.Form](ctx context.Context, t *git.Tree, user User, key string, value V) git.ChangeNoResult {
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
	return LookupUserByAddress_Local(ctx, gov.Clone(ctx, govAddr).Tree(), userAddr)
}

func LookupUserByAddress_Local(ctx context.Context, t *git.Tree, userAddr id.PublicAddress) []User {
	us := usersKV.ListKeys(ctx, usersNS, t)
	r := []User{}
	for _, u := range us {
		acct := GetUser_Local(ctx, t, u)
		if acct.PublicAddress == userAddr {
			r = append(r, u)
		}
	}
	return r
}

func LookupUserByID(ctx context.Context, govAddr gov.GovAddress, userID id.ID) []User {
	return LookupUserByID_Local(ctx, gov.Clone(ctx, govAddr).Tree(), userID)
}

func LookupUserByID_Local(ctx context.Context, t *git.Tree, userID id.ID) []User {
	us := usersKV.ListKeys(ctx, usersNS, t)
	r := []User{}
	for _, u := range us {
		acct := GetUser_Local(ctx, t, u)
		if acct.ID == userID {
			r = append(r, u)
		}
	}
	return r
}
