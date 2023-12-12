package sv

import (
	"context"
	_ "embed"
	"math"

	"github.com/gov4git/gov4git/proto/ballot/common"
)

type SV struct {
	Kernel ScoreKernel
}

func (x SV) GetScorer() ScoreKernel {
	if x.Kernel == nil {
		return QVScore{}
	}
	return x.Kernel
}

type ScoreKernel interface {
	Score(ctx context.Context, el common.AcceptedElections) ScoredVotes
	CalcJS(ctx context.Context) string
}

type ScoredVotes struct {
	Votes common.AcceptedElections
	Score map[string]common.StrengthAndScore // choice -> voting strength and resulting score
	Cost  float64
}

type QVScore struct{}

func (QVScore) Score(ctx context.Context, el common.AcceptedElections) ScoredVotes {
	// aggregate voting strength on each choice
	score := map[string]common.StrengthAndScore{}
	for _, el := range el {
		x := score[el.Vote.VoteChoice]
		x.Strength += el.Vote.VoteStrengthChange
		score[el.Vote.VoteChoice] = x
	}
	// compute score per choice
	for choice, ss := range score {
		score[choice] = common.StrengthAndScore{
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

var (
	//go:embed calc.js
	calcJS string
)

func (QVScore) CalcJS(context.Context) string {
	return calcJS
}
