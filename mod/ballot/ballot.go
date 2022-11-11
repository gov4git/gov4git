package ballot

import (
	"context"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/lib/ns"
	"github.com/gov4git/gov4git/mod/gov"
	"github.com/gov4git/gov4git/mod/id"
)

type Strategy interface {
	StrategyName() string
	Tally(
		ctx context.Context,
		govRepo id.OwnerRepo,
		govTree id.OwnerTree,
		//
		ad *Advertisement,
		current *TallyForm,
		fetched []FetchedVote,
	) git.Change[TallyForm] // tallying can change other aspects of the repo, like user balances
}

type BallotAddress[S Strategy] struct {
	Gov  gov.CommunityAddress
	Name ns.NS
}
