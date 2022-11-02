package strategy

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/govproto"
)

type BallotStrategy interface {
	Tally(ctx context.Context, community git.Local, ad govproto.GovBallotTally) error
}

func ParseStrategy(s govproto.GovBallotStrategy) (BallotStrategy, error) {
	switch {
	case s.PriorityPoll != nil:
		return PriorityPoll(*s.PriorityPoll), nil
	}
	return nil, fmt.Errorf("cannot parse ballot strategy")
}
