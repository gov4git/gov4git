package metrics

import (
	"context"
	"sort"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/member"
)

type CapTable []UserCap

func (x CapTable) Len() int {
	return len(x)
}

func (x CapTable) Less(i, j int) bool {
	return x[i].Cap < x[j].Cap
}

func (x CapTable) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

type UserCap struct {
	Name member.User `json:"name"`
	Cap  float64     `json:"cap"`
}

func GetCapTable_Local(
	ctx context.Context,
	cloned gov.Cloned,

) CapTable {

	users := member.ListGroupUsers_Local(ctx, cloned, member.Everybody)
	table := make(CapTable, len(users))
	for i := range users {
		table[i] = UserCap{
			Name: users[i],
			Cap:  account.Get_Local(ctx, cloned, member.UserAccountID(users[i])).Balance(account.PluralAsset).Quantity,
		}
	}
	sort.Sort(table)
	return table
}
