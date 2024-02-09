package sv

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
)

type SV struct {
	Kernel ScoreKernel
}

type ScoreKernel interface {
	Score(
		ctx context.Context,
		cloned gov.Cloned,
		ad *ballotproto.Ad,
		el ballotproto.AcceptedElections,

	) ScoredVotes

	CalcJS(
		ctx context.Context,
		cloned gov.Cloned,
		ad *ballotproto.Ad,
		tally *ballotproto.Tally,

	) *ballotproto.Margin
}

type ScoredVotes struct {
	Votes ballotproto.AcceptedElections
	Score map[string]ballotproto.StrengthAndScore // choice -> voting strength and resulting score
	Cost  float64
}

func (x SV) GetScorer(ctx context.Context) ScoreKernel {
	if x.Kernel == nil {
		return MakeQVScoreKernel(ctx, 1.0)
	}
	return x.Kernel
}
