package member

import (
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/account"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/kv"
	"github.com/gov4git/lib4git/form"
)

var (
	membersNS = proto.RootNS.Append("members")

	usersNS = membersNS.Append("users")
	usersKV = kv.KV[User, Account]{}

	groupsNS = membersNS.Append("groups")
	groupsKV = kv.KV[Group, form.None]{}

	userGroupsNS  = membersNS.Append("user_groups")
	userGroupsKKV = kv.KKV[User, Group, bool]{}

	groupUsersNS  = membersNS.Append("group_users")
	groupUsersKKV = kv.KKV[Group, User, bool]{}
)

type Account struct {
	ID            id.ID            `json:"id"`
	PublicAddress id.PublicAddress `json:"public_address"`
}

func UserAccountID(user User) account.AccountID {
	return account.AccountID("user:" + user)
}

func UserOwnerID(user User) account.OwnerID {
	return account.OwnerID("user:" + user)
}
