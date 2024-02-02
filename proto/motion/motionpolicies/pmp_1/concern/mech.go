package concern

import (
	"math"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_1"
)

func priorityScore(costOfPriority, idealDeficit, matchDeficit, matchFunds float64) float64 {

	return costOfPriority + matchRatio(matchFunds, matchDeficit)*idealDeficit
}

func idealFunding(tally *ballotproto.Tally) float64 {

	voteSum := 0.0
	for _, userSS := range tally.ScoresByUser {
		ss := userSS[pmp_1.ConcernBallotChoice]
		voteSum += ss.Score
	}
	return voteSum * voteSum
}

func matchRatio(matchFunds float64, matchDeficit float64) float64 {

	if matchDeficit <= 0 {
		return 1
	}
	if matchFunds <= 0 {
		return 0
	}
	return math.Min(1, matchFunds/matchDeficit)
}
