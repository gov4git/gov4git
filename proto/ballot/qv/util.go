package qv

import (
	"context"
	"math"
	"time"

	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/member"
)

type scoredVotes struct {
	Votes common.AcceptedElections
	Score map[string]common.StrengthAndScore // choice -> voting strength and resulting score
	Cost  float64
}

func scoreUserVotes(ctx context.Context, el common.AcceptedElections) scoredVotes {
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
	return scoredVotes{Votes: el, Score: score, Cost: cost}
}

func qvScoreFromStrength(strength float64) float64 {
	sign := 1.0
	if strength < 0 {
		sign = -1.0
	}
	return sign * math.Sqrt(math.Abs(strength))
}

func augmentAndScoreUserVotes(ctx context.Context, oldVotes common.AcceptedElections, newVotes common.Elections) (oldScore, augmentedScore scoredVotes) {
	oldScore = scoreUserVotes(ctx, oldVotes)
	augmentedVotes := append(append(common.AcceptedElections{}, oldVotes...), acceptVotes(newVotes)...)
	augmentedScore = scoreUserVotes(ctx, augmentedVotes)
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
