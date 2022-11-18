package common

import (
	"context"

	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

type Parameters interface{}

type Strategy interface {
	form.Form

	Name() string

	VerifyElections(
		ctx context.Context,
		voterAddr id.OwnerAddress,
		govAddr gov.CommunityAddress,
		voterTree id.OwnerTree,
		govTree *git.Tree,
		ad Advertisement,
		elections Elections,
	)

	Tally(
		ctx context.Context,
		govRepo id.OwnerRepo,
		govTree id.OwnerTree,
		ad *Advertisement,
		current *Tally,
		fetched []FetchedVote,
	) git.Change[Tally] // tallying can change other aspects of the repo, like user balances

	Close(
		ctx context.Context,
		govRepo id.OwnerRepo,
		govTree id.OwnerTree,
		ad *Advertisement,
		tally *Tally,
		summary Summary,
	) git.Change[Outcome]
}
