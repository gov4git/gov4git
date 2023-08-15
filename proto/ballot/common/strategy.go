package common

import (
	"context"

	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/member"
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
		govAddr gov.GovAddress,
		voterCloned id.OwnerCloned,
		govCloned git.Cloned,
		ad *Advertisement,
		prior *Tally,
		elections Elections,
	)

	Tally(
		ctx context.Context,
		gov id.OwnerCloned,
		ad *Advertisement,
		current *Tally,
		fetched map[member.User]Elections,
	) git.Change[form.Map, Tally] // tallying can change other aspects of the repo, like user balances

	Close(
		ctx context.Context,
		gov id.OwnerCloned,
		ad *Advertisement,
		tally *Tally,
	) git.Change[form.Map, Outcome]

	Cancel(
		ctx context.Context,
		gov id.OwnerCloned,
		ad *Advertisement,
		tally *Tally,
	) git.Change[form.Map, Outcome]
}
