package qv

import (
	"context"
	"math"

	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/must"
)

func GroupVotesByUser(ctx context.Context, fv common.FetchedVotes) map[member.User]common.FetchedVote {
	group := map[member.User]common.FetchedVote{}
	for _, fetchedVote := range fv {
		_, ok := group[fetchedVote.Voter]
		must.Assertf(ctx, !ok, "multiple votes by user unexpected")
		group[fetchedVote.Voter] = fetchedVote
	}
	return group
}

type ScoredVotes struct {
	Votes common.FetchedVote
	Score map[string]float64 // choice -> score
	Cost  float64
}

func ScoreAndCostOfUserVotes(ctx context.Context, fv common.FetchedVote) ScoredVotes {
	score := map[string]float64{}
	for _, el := range fv.Elections {
		score[el.VoteChoice] += el.VoteStrengthChange
	}
	cost := 0.0
	for choice, x := range score {
		s := math.Sqrt(x)
		score[choice] = s
		cost += math.Abs(s)
	}
	return ScoredVotes{Votes: fv, Score: score, Cost: cost}
}
