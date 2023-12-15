package sv

import (
	"context"
	"time"

	"github.com/gov4git/gov4git/v2/proto/ballot/common"
	"github.com/gov4git/gov4git/v2/proto/member"
)

func augmentAndScoreUserVotes(
	ctx context.Context,
	kernel ScoreKernel,
	oldVotes common.AcceptedElections,
	newVotes common.Elections,
) (oldScore, augmentedScore ScoredVotes) {

	oldScore = kernel.Score(ctx, oldVotes)
	augmentedVotes := append(append(common.AcceptedElections{}, oldVotes...), acceptVotes(newVotes)...)
	augmentedScore = kernel.Score(ctx, augmentedVotes)
	return
}

func acceptVotes(votes common.Elections) common.AcceptedElections {
	r := make(common.AcceptedElections, len(votes))
	for i, v := range votes {
		r[i] = common.AcceptedElection{Time: time.Now(), Vote: v}
	}
	return r
}

func rejectVotes(votes common.Elections, reason error) common.RejectedElections {
	r := make(common.RejectedElections, len(votes))
	for i, v := range votes {
		r[i] = common.RejectedElection{Time: time.Now(), Vote: v, Reason: reason.Error()}
	}
	return r
}

func totalScore(choices []string, votesByUser map[member.User]map[string]common.StrengthAndScore) map[string]float64 {
	scores := map[string]float64{}
	for _, choice := range choices {
		scores[choice] = 0.0
	}
	for _, userVotes := range votesByUser {
		for choice, votes := range userVotes {
			scores[choice] += votes.Score
		}
	}
	return scores
}
