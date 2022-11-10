package member

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/lib/must"
	"github.com/gov4git/gov4git/mod/kv"
)

func SetUser(ctx context.Context, t *git.Tree, name User, user Account) git.ChangeNoResult {
	AddGroup(ctx, t, Everybody)
	AddMember(ctx, t, name, Everybody)
	return usersKV.Set(ctx, usersNS, t, name, user)
}

func GetUser(ctx context.Context, t *git.Tree, name User) Account {
	return usersKV.Get(ctx, usersNS, t, name)
}

func AddUser(ctx context.Context, t *git.Tree, name User, user Account) git.ChangeNoResult {
	if err := must.Try(func() { GetUser(ctx, t, name) }); err == nil {
		must.Panic(ctx, fmt.Errorf("user already exists"))
	}
	return SetUser(ctx, t, name, user)
}

func RemoveUser(ctx context.Context, t *git.Tree, name User) git.ChangeNoResult {
	usersKV.Remove(ctx, usersNS, t, name)
	userGroupsKKV.RemovePrimary(ctx, userGroupsNS, t, name) // remove memberships
	return git.ChangeNoResult{
		Msg: fmt.Sprintf("Remove user %v", name),
	}
}

// props

func SetUserProp[V form.Form](ctx context.Context, t *git.Tree, user User, key string, value V) git.ChangeNoResult {
	propKV := kv.KV[string, V]{}
	return propKV.Set(ctx, usersKV.KeyNS(usersNS, user), t, key, value)
}

func GetUserProp[V form.Form](ctx context.Context, t *git.Tree, user User, key string, value V) V {
	propKV := kv.KV[string, V]{}
	return propKV.Get(ctx, usersKV.KeyNS(usersNS, user), t, key)
}
