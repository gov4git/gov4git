package member

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/lib/must"
	"github.com/gov4git/gov4git/mod"
	"github.com/gov4git/gov4git/mod/kv"
)

func SetUser(ctx context.Context, t *git.Tree, name User, url git.URL) mod.Change[form.None] {
	AddGroup(ctx, t, Everybody)
	AddMember(ctx, t, name, Everybody)
	return usersKV.Set(ctx, usersNS, t, name, url)
}

func GetUser(ctx context.Context, t *git.Tree, name User) git.URL {
	return usersKV.Get(ctx, usersNS, t, name)
}

func AddUser(ctx context.Context, t *git.Tree, name User, url git.URL) mod.Change[form.None] {
	if err := must.Try0(func() { GetUser(ctx, t, name) }); err == nil {
		must.Panic(ctx, fmt.Errorf("user already exists"))
	}
	return SetUser(ctx, t, name, url)
}

func RemoveUser(ctx context.Context, t *git.Tree, name User) mod.Change[form.None] {
	usersKV.Remove(ctx, usersNS, t, name)
	userGroupsKKV.RemovePrimary(ctx, userGroupsNS, t, name) // remove memberships
	return mod.Change[form.None]{
		Msg: fmt.Sprintf("Remove user %v", name),
	}
}

// props

func SetUserProp[V form.Form](ctx context.Context, t *git.Tree, user User, key string, value V) mod.Change[form.None] {
	propKV := kv.KV[string, V]{}
	return propKV.Set(ctx, usersKV.KeyNS(usersNS, user), t, key, value)
}

func GetUserProp[V form.Form](ctx context.Context, t *git.Tree, user User, key string, value V) V {
	propKV := kv.KV[string, V]{}
	return propKV.Get(ctx, usersKV.KeyNS(usersNS, user), t, key)
}
