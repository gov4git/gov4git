package ballotproto

import (
	"math"
	"time"

	"github.com/gov4git/gov4git/v2/proto/member"
)

type Tally struct {
	Ad            Ad                                          `json:"advertisement"`
	Scores        map[string]float64                          `json:"scores"`         // choice -> score
	ScoresByUser  map[member.User]map[string]StrengthAndScore `json:"scores_by_user"` // user -> choice -> signed voting credits spent on the choice by the user
	AcceptedVotes map[member.User]AcceptedElections           `json:"accepted_votes"`
	RejectedVotes map[member.User]RejectedElections           `json:"rejected_votes"`
	Charges       map[member.User]float64                     `json:"charges"`
}

func (x Tally) NumVoters() int {
	return len(x.AcceptedVotes)
}

func (x Tally) Capitalization() float64 {
	cap := 0.0
	for _, spent := range x.Charges {
		cap += spent
	}
	return cap
}

func (x Tally) Attention() float64 {
	score := 0.0
	for _, choices := range x.ScoresByUser {
		for _, ss := range choices {
			score += math.Abs(ss.Score)
		}
	}
	return score
}

type StrengthAndScore struct {
	Strength float64 `json:"strength"` // signed number of voting credits spent by the user
	Score    float64 `json:"score"`    // vote impact, based on the voting strength (above)
}

func (ss StrengthAndScore) Vote() float64 {
	return ss.Score
}

type AcceptedElection struct {
	Time time.Time `json:"accepted_time"`
	Vote Election  `json:"accepted_vote"`
}

type AcceptedElections []AcceptedElection

type RejectedElection struct {
	Time   time.Time `json:"rejected_time"`
	Vote   Election  `json:"rejected_vote"`
	Reason string    `json:"rejected_reason"`
}

type RejectedElections []RejectedElection
