package member

import (
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/kv"
	"github.com/gov4git/lib4git/form"
)

var (
	membersNS = proto.RootNS.Sub("members")

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
	ID            id.ID            `json:"id"`
	PublicAddress id.PublicAddress `json:"public_address"`
}
