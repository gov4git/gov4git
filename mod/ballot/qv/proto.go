package qv

import (
	"github.com/gov4git/gov4git/mod/balance"
	"github.com/gov4git/gov4git/mod/ballot/proto"
)

const (
	VotingCredits       balance.Balance = "voting_credits"
	VotingCreditsOnHold balance.Balance = "voting_credits_on_hold"
)

const (
	SummaryAdopted   proto.Summary = "adopted"
	SummaryAbandoned proto.Summary = "abandoned"
)
