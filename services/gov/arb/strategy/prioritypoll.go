package strategy

import (
	"context"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/govproto"
)

type PriorityPoll govproto.PriorityPollStrategy

func (x PriorityPoll) Tally(ctx context.Context, community git.Local, ad govproto.GovBallotTally) error {
	return nil
}
