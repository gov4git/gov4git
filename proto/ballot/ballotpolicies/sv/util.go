package sv

import (
	"context"
	"time"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/member"
)

func augmentAndScoreUserVotes(
	ctx context.Context,
	cloned gov.Cloned,
	ad *ballotproto.Ad,
	kernel ScoreKernel,
	oldVotes ballotproto.AcceptedElections,
	newVotes ballotproto.Elections,
) (oldScore, augmentedScore ScoredVotes) {

	oldScore = kernel.Score(ctx, cloned, ad, oldVotes)
	augmentedVotes := append(append(ballotproto.AcceptedElections{}, oldVotes...), acceptVotes(newVotes)...)
	augmentedScore = kernel.Score(ctx, cloned, ad, augmentedVotes)
	return
}

func acceptVotes(votes ballotproto.Elections) ballotproto.AcceptedElections {
	r := make(ballotproto.AcceptedElections, len(votes))
	for i, v := range votes {
		r[i] = ballotproto.AcceptedElection{Time: time.Now(), Vote: v}
	}
	return r
}

func rejectVotes(votes ballotproto.Elections, reason error) ballotproto.RejectedElections {
	r := make(ballotproto.RejectedElections, len(votes))
	for i, v := range votes {
		r[i] = ballotproto.RejectedElection{Time: time.Now(), Vote: v, Reason: reason.Error()}
	}
	return r
}

func totalScore(choices []string, votesByUser map[member.User]map[string]ballotproto.StrengthAndScore) map[string]float64 {
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
