package ballotproto

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/id"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

type StrategyName string

func (x StrategyName) String() string {
	return string(x)
}

type Strategy interface {
	VerifyElections(
		ctx context.Context,
		voterAddr id.OwnerAddress,
		govAddr gov.Address,
		voterCloned id.OwnerCloned,
		govCloned gov.Cloned,
		ad *Advertisement,
		prior *Tally,
		elections Elections,
	)

	Tally(
		ctx context.Context,
		cloned gov.OwnerCloned,
		ad *Advertisement,
		current *Tally,
		fetched map[member.User]Elections,

	) git.Change[form.Map, Tally] // tallying can change other aspects of the repo, like user balances

	Margin(
		ctx context.Context,
		cloned gov.Cloned,
		ad *Advertisement,
		current *Tally,

	) *Margin

	Open(
		ctx context.Context,
		cloned gov.OwnerCloned,
		ad *Advertisement,

	) *Tally

	Close(
		ctx context.Context,
		cloned gov.OwnerCloned,
		ad *Advertisement,
		tally *Tally,

	) git.Change[form.Map, Outcome]

	Cancel(
		ctx context.Context,
		cloned gov.OwnerCloned,
		ad *Advertisement,
		tally *Tally,

	) git.Change[form.Map, Outcome]

	Reopen(
		ctx context.Context,
		cloned gov.OwnerCloned,
		ad *Advertisement,
		tally *Tally,

	) git.Change[form.Map, form.None]
}
