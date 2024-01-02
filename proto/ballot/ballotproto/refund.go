package ballotproto

import (
	"sort"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/member"
)

func FlattenRefunds(m map[member.User]account.Holding) Refunds {
	r := Refunds{}
	for k, v := range m {
		r = append(r, Refund{User: k, Amount: v})
	}
	r.Sort()
	return r
}

type Refund struct {
	User   member.User     `json:"user"`
	Amount account.Holding `json:"amount"`
}

type Refunds []Refund

func (x Refunds) Len() int {
	return len(x)
}

func (x Refunds) Less(i, j int) bool {
	return x[i].User < x[j].User
}

func (x Refunds) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x Refunds) Sort() {
	sort.Sort(x)
}
