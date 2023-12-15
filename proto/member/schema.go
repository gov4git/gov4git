package member

import (
	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/id"
	"github.com/gov4git/gov4git/v2/proto/kv"
	"github.com/gov4git/lib4git/form"
)

var (
	membersNS = proto.RootNS.Append("members")

	usersNS = membersNS.Append("users")
	usersKV = kv.KV[User, UserProfile]{}

	groupsNS = membersNS.Append("groups")
	groupsKV = kv.KV[Group, form.None]{}

	userGroupsNS  = membersNS.Append("user_groups")
	userGroupsKKV = kv.KKV[User, Group, bool]{}

	groupUsersNS  = membersNS.Append("group_users")
	groupUsersKKV = kv.KKV[Group, User, bool]{}
)

type UserProfile struct {
	ID            id.ID            `json:"id"`
	PublicAddress id.PublicAddress `json:"public_address"`
}

func UserAccountID(user User) account.AccountID {
	return account.AccountIDFromLine(account.Pair("user", string(user)))
}

func UserOwnerID(user User) account.OwnerID {
	return account.OwnerIDFromLine(account.Pair("user", string(user)))
}
