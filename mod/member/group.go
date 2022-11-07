package member

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/lib/must"
	"github.com/gov4git/gov4git/mod"
)

func SetGroup(ctx context.Context, t *git.Tree, name Group) mod.Change[form.None] {
	return groupsKV.Set(ctx, groupsNS, t, name, form.None{})
}

func GetGroup(ctx context.Context, t *git.Tree, name Group) {
	groupsKV.Get(ctx, groupsNS, t, name)
}

func AddGroup(ctx context.Context, t *git.Tree, name Group) mod.Change[form.None] {
	if err := must.Try0(func() { GetGroup(ctx, t, name) }); err == nil {
		must.Panic(ctx, fmt.Errorf("group already exists"))
	}
	return SetGroup(ctx, t, name)
}

func RemoveGroup(ctx context.Context, t *git.Tree, name Group) mod.Change[form.None] {
	groupsKV.Remove(ctx, groupsNS, t, name)
	//XXX: remove memberships
	return mod.Change[form.None]{
		Msg: fmt.Sprintf("Remove group %v", name),
	}
}
