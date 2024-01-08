package sv

import (
	"context"
	"math"

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
	) ballotproto.MarginCalcJS
}

type ScoredVotes struct {
	Votes ballotproto.AcceptedElections
	Score map[string]ballotproto.StrengthAndScore // choice -> voting strength and resulting score
	Cost  float64
}

func (x SV) GetScorer() ScoreKernel {
	if x.Kernel == nil {
		return QVScoreKernel{}
	}
	return x.Kernel
}

type QVScoreKernel struct{}

func (QVScoreKernel) Score(
	ctx context.Context,
	cloned gov.Cloned,
	ad *ballotproto.Ad,
	el ballotproto.AcceptedElections,

) ScoredVotes {

	// aggregate voting strength on each choice
	score := map[string]ballotproto.StrengthAndScore{}
	for _, el := range el {
		x := score[el.Vote.VoteChoice]
		x.Strength += el.Vote.VoteStrengthChange
		score[el.Vote.VoteChoice] = x
	}
	// compute score per choice
	for choice, ss := range score {
		score[choice] = ballotproto.StrengthAndScore{
			Strength: ss.Strength,
			Score:    qvScoreFromStrength(ss.Strength),
		}
	}
	// compute aggregate cost
	cost := 0.0
	for _, x := range score {
		cost += math.Abs(x.Strength)
	}
	return ScoredVotes{Votes: el, Score: score, Cost: cost}
}

func qvScoreFromStrength(strength float64) float64 {
	sign := 1.0
	if strength < 0 {
		sign = -1.0
	}
	return sign * math.Sqrt(math.Abs(strength))
}

func (QVScoreKernel) CalcJS(
	context.Context,
	gov.Cloned,
	*ballotproto.Ad,

) ballotproto.MarginCalcJS {

	return ballotproto.MarginCalcJS(qvMarginCalcJS)
}
