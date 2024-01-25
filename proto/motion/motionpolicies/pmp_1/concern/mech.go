package concern

import (
	"math"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_1"
)

func priorityOfConcern(tally *ballotproto.Tally) float64 {

	panic("XXX")
}

func idealizedQuadraticFundingDeficitOfConcern(tally *ballotproto.Tally) float64 {
	return idealizedQuadraticFundingOfConcern(tally) - capitalistFundingOfConcern(tally)
}

func idealizedQuadraticFundingOfConcern(tally *ballotproto.Tally) float64 {

	voteSum := 0.0
	for _, userSS := range tally.ScoresByUser {
		ss := userSS[pmp_1.ConcernBallotChoice]
		voteSum += ss.Score
	}
	return voteSum * voteSum
}

func capitalistFundingOfConcern(tally *ballotproto.Tally) float64 {

	f := 0.0
	for _, userSS := range tally.ScoresByUser {
		ss := userSS[pmp_1.ConcernBallotChoice]
		f += math.Abs(ss.Strength)
	}
	return f
}
