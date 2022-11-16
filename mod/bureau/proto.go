package bureau

import (
	"github.com/gov4git/gov4git/mod/balance"
	"github.com/gov4git/gov4git/mod/id"
	"github.com/gov4git/gov4git/mod/member"
)

const BureauTopic = "bureau"

type Request struct {
	Transfer *TransferRequest `json:"transfer"`
}

type Requests []Request

type TransferRequest struct {
	FromUser    member.User     `json:"from_user"`
	FromBalance balance.Balance `json:"from_balance"`
	ToUser      member.User     `json:"to_user"`
	ToBalance   balance.Balance `json:"to_balance"`
	Amount      float64         `json:"amount"`
}

type FetchedRequest struct {
	User     member.User      `json:"requesting_user"`
	Address  id.PublicAddress `json:"requesting_address"`
	Requests Requests         `json:"requests"`
}

type FetchedRequests []FetchedRequest
