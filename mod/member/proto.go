package member

import (
	"github.com/gov4git/gov4git/mod"
	"github.com/gov4git/gov4git/mod/id"
	"github.com/gov4git/gov4git/mod/kv"
	"github.com/gov4git/lib4git/form"
)

var (
	membersNS = mod.RootNS.Sub("members")

	usersNS = membersNS.Sub("users")
	usersKV = kv.KV[User, Account]{}

	groupsNS = membersNS.Sub("groups")
	groupsKV = kv.KV[Group, form.None]{}

	userGroupsNS  = membersNS.Sub("user_groups")
	userGroupsKKV = kv.KKV[User, Group, bool]{}

	groupUsersNS  = membersNS.Sub("group_users")
	groupUsersKKV = kv.KKV[Group, User, bool]{}
)

type Account struct {
	Home id.PublicAddress `json:"home"`
}
