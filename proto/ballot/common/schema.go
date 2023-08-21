package common

import (
	"path/filepath"
	"time"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/form"
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

	VoteLogNS = proto.RootNS.Sub("votes") // namespace in voter's repo for recording votes
)

// VoteLog records the votes of a user to a ballot within a given governance.
type VoteLog struct {
	GovID         id.ID          `json:"governance_id"`
	GovAddress    gov.GovAddress `json:"governance_address"`
	Ballot        ns.NS          `json:"ballot_name"`
	VoteEnvelopes VoteEnvelopes  `json:"vote_envelopes"` // in the order in which they were sent
}

// VoterStatus reflects the state of an individual user's votes within a ballot.
type VoterStatus struct {
	GovID         id.ID             `json:"governance_id"`
	GovAddress    gov.GovAddress    `json:"governance_address"`
	BallotName    ns.NS             `json:"ballot_name"`
	AcceptedVotes AcceptedElections `json:"accepted_votes"`
	RejectedVotes RejectedElections `json:"rejected_votes"`
	PendingVotes  Elections         `json:"pending_votes"`
}

func VoteLogPath(govID id.ID, ballotName ns.NS) ns.NS {
	return VoteLogNS.Sub(
		filepath.Join(
			form.StringHashForFilename(string(govID)),
			form.StringHashForFilename(BallotTopic(ballotName)),
		),
	)
}

func BallotTopic(ballotName ns.NS) string {
	return "ballot:" + ballotName.Path()
}

func BallotPath(path ns.NS) ns.NS {
	return BallotNS.Join(path)
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
	Closed       bool           `json:"closed"` // closed ballots cannot be re-opened
	Cancelled    bool           `json:"cancelled"`
	ParentCommit git.CommitHash `json:"parent_commit"`
}

type BallotAddress struct {
	Gov  gov.GovAddress
	Name ns.NS
}

type Election struct {
	VoteID             id.ID     `json:"vote_id"`
	VoteTime           time.Time `json:"vote_time"`
	VoteChoice         string    `json:"vote_choice"`
	VoteStrengthChange float64   `json:"vote_strength_change"`
}

func NewElection(choice string, strength float64) Election {
	return Election{
		VoteID:             id.GenerateUniqueID(), //XXX
		VoteTime:           time.Now(),
		VoteChoice:         choice,
		VoteStrengthChange: strength,
	}
}

type Elections []Election

type VoteEnvelope struct {
	AdCommit  git.CommitHash `json:"ballot_ad_commit"`
	Ad        Advertisement  `json:"ballot_ad"`
	Elections Elections      `json:"ballot_elections"`
}

type VoteEnvelopes []VoteEnvelope

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
	Ad            Advertisement                               `json:"ballot_advertisement"`
	Scores        map[string]float64                          `json:"ballot_scores"`        // choice -> score
	VotesByUser   map[member.User]map[string]StrengthAndScore `json:"ballot_votes_by_user"` // user -> choice -> signed voting credits spent on the choice by the user
	AcceptedVotes map[member.User]AcceptedElections           `json:"ballot_accepted_votes"`
	RejectedVotes map[member.User]RejectedElections           `json:"ballot_rejected_votes"`
	Charges       map[member.User]float64                     `json:"ballot_charges"`
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

type AdStrategyTally struct {
	Ad       Advertisement `json:"ballot_advertisement"`
	Strategy Strategy      `json:"ballot_strategy"`
	Tally    Tally         `json:"ballot_tally"`
}

type Summary string

type Outcome struct {
	Summary string             `json:"ballot_summary"`
	Scores  map[string]float64 `json:"ballot_scores"`
}
