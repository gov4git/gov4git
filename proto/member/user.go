package member

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history/metric"
	"github.com/gov4git/gov4git/v2/proto/history/trace"
	"github.com/gov4git/gov4git/v2/proto/id"
	"github.com/gov4git/gov4git/v2/proto/kv"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func setUser_StageOnly(ctx context.Context, cloned gov.Cloned, name User, user UserProfile) git.ChangeNoResult {
	SetGroup_StageOnly(ctx, cloned, Everybody)                  // create everybody group, if it doesn't exist
	AddMember_StageOnly(ctx, cloned, name, Everybody)           // add membership of user to everybody
	return usersKV.Set(ctx, usersNS, cloned.Tree(), name, user) // create the user record
}

func GetUser(ctx context.Context, addr gov.Address, name User) UserProfile {
	return GetUser_Local(ctx, gov.Clone(ctx, addr), name)
}

func IsUser_Local(ctx context.Context, cloned gov.Cloned, name User) bool {
	err := must.Try(
		func() {
			GetUser_Local(ctx, cloned, name)
		},
	)
	if err == nil {
		return true
	}
	if git.IsNotExist(err) {
		return false
	}
	must.Panic(ctx, err)
	return false
}

func GetUser_Local(ctx context.Context, cloned gov.Cloned, name User) UserProfile {
	return usersKV.Get(ctx, usersNS, cloned.Tree(), name)
}

func AddUserByPublicAddress(ctx context.Context, govAddr gov.Address, name User, userAddr id.PublicAddress) {
	cred := id.FetchPublicCredentials(ctx, userAddr)
	AddUser(ctx, govAddr, name, UserProfile{ID: cred.ID, PublicAddress: userAddr})
}

func AddUserByPublicAddress_StageOnly(ctx context.Context, cloned gov.Cloned, name User, userAddr id.PublicAddress) {
	cred := id.FetchPublicCredentials(ctx, userAddr)
	AddUser_StageOnly(ctx, cloned, name, UserProfile{ID: cred.ID, PublicAddress: userAddr})
}

func AddUser(ctx context.Context, addr gov.Address, name User, acct UserProfile) {
	cloned := gov.Clone(ctx, addr)
	chg := AddUser_StageOnly(ctx, cloned, name, acct)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
}

func AddUser_StageOnly(ctx context.Context, cloned gov.Cloned, name User, profile UserProfile) git.ChangeNoResult {
	if err := must.Try(func() { GetUser_Local(ctx, cloned, name) }); err == nil {
		must.Panic(ctx, fmt.Errorf("user already exists"))
	}
	account.Create_StageOnly(ctx, cloned, UserAccountID(name), account.NobodyAccountID, fmt.Sprintf("account for user %v", name))
	chg := setUser_StageOnly(ctx, cloned, name, profile)

	// log
	metric.Log_StageOnly(ctx, cloned, &metric.Event{
		Join: &metric.JoinEvent{
			User: metric.User(name),
		},
	})

	return chg
}

func RemoveUser(ctx context.Context, addr gov.Address, name User) {
	cloned := gov.Clone(ctx, addr)
	chg := RemoveUser_StageOnly(ctx, cloned, name)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
}

func RemoveUser_StageOnly(ctx context.Context, cloned gov.Cloned, name User) git.ChangeNoResult {
	must.Assertf(ctx, IsUser_Local(ctx, cloned, name), "%v is not a name", name)

	// burn user's balance and remove user account
	account.Remove_StageOnly(ctx, cloned, UserAccountID(name), "removing user")

	// remove all group memberships of the user
	for _, g := range ListUserGroups_Local(ctx, cloned, name) {
		RemoveMember_StageOnly(ctx, cloned, name, g)
	}

	// remove user record
	usersKV.Remove(ctx, usersNS, cloned.Tree(), name)
	chg := git.NewChangeNoResult(fmt.Sprintf("Remove user %v", name), "member_remove_user")

	// log
	trace.Log_StageOnly(ctx, cloned, &trace.Event{
		Op:     "user_remove",
		Args:   trace.M{"name": name},
		Result: nil,
	})

	return chg
}

// set prop

func SetUserProp[V form.Form](ctx context.Context, addr gov.Address, user User, key string, value V) {
	cloned := gov.Clone(ctx, addr)
	chg := SetUserProp_StageOnly(ctx, cloned, user, key, value)
	proto.Commit(ctx, cloned.Tree(), chg)
	cloned.Push(ctx)
}

func SetUserProp_StageOnly[V form.Form](
	ctx context.Context,
	cloned gov.Cloned,
	user User,
	key string,
	value V,
) git.ChangeNoResult {
	must.Assertf(ctx, IsUser_Local(ctx, cloned, user), "%v is not a user", user)
	propKV := kv.KV[string, V]{}
	return propKV.Set(ctx, usersKV.KeyNS(usersNS, user), cloned.Tree(), key, value)
}

// get prop

func GetUserProp[V form.Form](ctx context.Context, addr gov.Address, user User, key string) V {
	return GetUserProp_Local[V](ctx, gov.Clone(ctx, addr), user, key)
}

func GetUserProp_Local[V form.Form](ctx context.Context, cloned gov.Cloned, user User, key string) V {
	must.Assertf(ctx, IsUser_Local(ctx, cloned, user), "%v is not a user", user)
	propKV := kv.KV[string, V]{}
	return propKV.Get(ctx, usersKV.KeyNS(usersNS, user), cloned.Tree(), key)
}

func GetUserPropOrDefault[V form.Form](ctx context.Context, addr gov.Address, user User, key string, default_ V) V {
	return GetUserPropOrDefault_Local[V](ctx, gov.Clone(ctx, addr), user, key, default_)
}

func GetUserPropOrDefault_Local[V form.Form](
	ctx context.Context,
	cloned gov.Cloned,
	user User,
	key string,
	default_ V,
) V {
	must.Assertf(ctx, IsUser_Local(ctx, cloned, user), "%v is not a user", user)
	v, err := must.Try1(func() V { return GetUserProp_Local[V](ctx, cloned, user, key) })
	if err != nil {
		return default_
	}
	return v
}

// lookup

func LookupUserByAddress(ctx context.Context, govAddr gov.Address, userAddr id.PublicAddress) []User {
	return LookupUserByAddress_Local(ctx, gov.Clone(ctx, govAddr), userAddr)
}

func LookupUserByAddress_Local(ctx context.Context, cloned gov.Cloned, userAddr id.PublicAddress) []User {
	us := usersKV.ListKeys(ctx, usersNS, cloned.Tree())
	r := []User{}
	for _, u := range us {
		acct := GetUser_Local(ctx, cloned, u)
		if acct.PublicAddress == userAddr {
			r = append(r, u)
		}
	}
	return r
}

func LookupUserByID(ctx context.Context, govAddr gov.Address, userID id.ID) []User {
	return LookupUserByID_Local(ctx, gov.Clone(ctx, govAddr), userID)
}

func LookupUserByID_Local(ctx context.Context, cloned gov.Cloned, userID id.ID) []User {
	us := usersKV.ListKeys(ctx, usersNS, cloned.Tree())
	r := []User{}
	for _, u := range us {
		acct := GetUser_Local(ctx, cloned, u)
		if acct.ID == userID {
			r = append(r, u)
		}
	}
	return r
}
