package common

import (
	"time"

	"github.com/gov4git/gov4git/proto/member"
)

type Tally struct {
	Ad            Advertisement                               `json:"ballot_advertisement"`
	Scores        map[string]float64                          `json:"ballot_scores"`        // choice -> score
	VotesByUser   map[member.User]map[string]StrengthAndScore `json:"ballot_votes_by_user"` // user -> choice -> signed voting credits spent on the choice by the user
	AcceptedVotes map[member.User]AcceptedElections           `json:"ballot_accepted_votes"`
	RejectedVotes map[member.User]RejectedElections           `json:"ballot_rejected_votes"`
	Charges       map[member.User]float64                     `json:"ballot_charges"`
}

func (x Tally) Capitalization() float64 {
	cap := 0.0
	for _, spent := range x.Charges {
		cap += spent
	}
	return cap
}

type StrengthAndScore struct {
	Strength float64 `json:"strength"` // signed number of voting credits spent by the user
	Score    float64 `json:"score"`    // qv score, based on the voting strength (above)
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
