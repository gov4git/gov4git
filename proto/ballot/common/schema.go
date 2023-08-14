package common

import (
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/ns"
	"github.com/gov4git/lib4git/util"
)

var (
	BallotNS         = proto.RootNS.Sub("ballots")
	AdFilebase       = "ballot_ad.json"
	StrategyFilebase = "ballot_strategy.json"
	TallyFilebase    = "ballot_tally.json"
	OutcomeFilebase  = "ballot_outcome.json"
)

func OpenBallotNS(name ns.NS) ns.NS {
	return BallotNS.Sub("open").Join(name)
}

func ClosedBallotNS(name ns.NS) ns.NS {
	return BallotNS.Sub("closed").Join(name)
}

func BallotTopic(name ns.NS) string {
	return "ballot:" + name.Path()
}

type Advertisement struct {
	Gov          gov.GovAddress `json:"community"`
	Name         ns.NS          `json:"path"`
	Title        string         `json:"title"`
	Description  string         `json:"description"`
	Choices      []string       `json:"choices"`
	Strategy     string         `json:"strategy"`
	Participants member.Group   `json:"participants_group"`
	Frozen       bool           `json:"frozen"` // if frozen, the ballot is not accepting votes
	ParentCommit git.CommitHash `json:"parent_commit"`
}

type BallotAddress struct {
	Gov  gov.GovAddress
	Name ns.NS
}

type Election struct {
	// XXX: add timestamp
	VoteChoice         string  `json:"vote_choice"`
	VoteStrengthChange float64 `json:"vote_strength_change"`
}

type Elections []Election

type VoteEnvelope struct {
	AdCommit  git.CommitHash `json:"ballot_ad_commit"`
	Ad        Advertisement  `json:"ballot_ad"`
	Elections Elections      `json:"ballot_elections"`
}

// Verify verifies that elections are consistent with the ballot ad.
func (x VoteEnvelope) VerifyConsistency() bool {
	for _, v := range x.Elections {
		if !util.IsIn(v.VoteChoice, x.Ad.Choices...) {
			return false
		}
	}
	return true
}

type Tally struct {
	Ad     Advertisement `json:"ballot_advertisement"`
	Votes  FetchedVotes  `json:"ballot_fetched_votes"`
	Scores ChoiceScores  `json:"ballot_choice_scores"`
}

type Tally2 struct {
	Ad            Advertisement                               `json:"ballot_advertisement"`
	Scores        map[string]float64                          `json:"ballot_scores"`        // choice -> score
	VotesByUser   map[member.User]map[string]StrengthAndScore `json:"ballot_votes_by_user"` // user -> choice -> signed voting credits spent on the choice by the user
	AcceptedVotes map[member.User]Elections                   `json:"ballot_accepted_votes"`
	RejectedVotes map[member.User]RejectedElections           `json:"ballot_rejected_votes"`
	Charges       map[member.User]float64                     `json:"ballot_charges"`
}

type StrengthAndScore struct {
	Strength float64 `json:"strength"` // signed number of voting credits spent by the user
	Score    float64 `json:"score"`    // qv score, based on the voting strength (above)
}

type RejectedElection struct {
	// XXX: add timestamp
	Vote   Election `json:"rejected_vote"`
	Reason string   `json:"rejection_reason"`
}

type RejectedElections []RejectedElection

type FetchedVote struct {
	Voter     member.User      `json:"voter_user"`
	Address   id.PublicAddress `json:"voter_address"`
	Elections Elections        `json:"voter_elections"`
}

type FetchedVotes []FetchedVote

func (x FetchedVotes) Len() int           { return len(x) }
func (x FetchedVotes) Less(i, j int) bool { return x[i].Voter < x[j].Voter }
func (x FetchedVotes) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// XXX: shouldn't be necessary (below)

type ChoiceScore struct {
	Choice string  `json:"choice"`
	Score  float64 `json:"score"`
}

type ChoiceScores []ChoiceScore

func (x ChoiceScores) Len() int           { return len(x) }
func (x ChoiceScores) Less(i, j int) bool { return x[i].Score > x[j].Score }
func (x ChoiceScores) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type AdStrategyTally struct {
	Ad       Advertisement `json:"ballot_advertisement"`
	Strategy Strategy      `json:"ballot_strategy"`
	Tally    Tally         `json:"ballot_tally"`
}

type Summary string

type Outcome struct {
	Summary Summary      `json:"ballot_summary"`
	Scores  ChoiceScores `json:"ballot_choice_scores"`
}
