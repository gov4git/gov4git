package qv

import (
	"github.com/gov4git/gov4git/proto/balance"
	"github.com/gov4git/gov4git/proto/ballot/common"
)

var (
	VotingCredits = balance.Balance{"voting_credits"}
)

const (
	SummaryAdopted   common.Summary = "adopted"
	SummaryAbandoned common.Summary = "abandoned"
)
