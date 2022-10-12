package strategy

import (
	"context"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto"
)

type PriorityPoll proto.PriorityPollStrategy

func (x PriorityPoll) Tally(ctx context.Context, community git.Local, ad proto.GovBallotTally) error {
	return nil
}
