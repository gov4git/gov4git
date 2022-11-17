package ballot

import (
	"github.com/gov4git/gov4git/mod"
	"github.com/gov4git/gov4git/mod/gov"
	"github.com/gov4git/gov4git/mod/id"
	"github.com/gov4git/gov4git/mod/member"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/ns"
	"github.com/gov4git/lib4git/util"
)

var (
	ballotNS         = mod.RootNS.Sub("ballots")
	adFilebase       = "ballot_ad.json"
	strategyFilebase = "ballot_strategy.json"
	tallyFilebase    = "ballot_tally.json"
)

func OpenBallotNS(name ns.NS) ns.NS {
	return ballotNS.Sub("open").Join(name)
}

func ClosedBallotNS(name ns.NS) ns.NS {
	return ballotNS.Sub("closed").Join(name)
}

func BallotTopic(name ns.NS) string {
	return "ballot:" + name.Path()
}

type Advertisement struct {
	Community    gov.CommunityAddress `json:"community"`
	Name         ns.NS                `json:"path"`
	Title        string               `json:"title"`
	Description  string               `json:"description"`
	Choices      []string             `json:"choices"`
	Strategy     string               `json:"strategy"`
	Participants member.Group         `json:"participants_group"`
	ParentCommit git.CommitHash       `json:"parent_commit"`
}

type BallotAddress struct {
	Gov  gov.CommunityAddress
	Name ns.NS
}

type Election struct {
	VoteChoice         string  `json:"vote_choice"`
	VoteStrengthChange float64 `json:"vote_strength_change"`
}

type Elections []Election

type VoteEnvelope struct {
	AdCommit  git.CommitHash `json:"ballot_ad_commit"`
	Ad        Advertisement  `json:"ballot_ad"`
	Elections []Election     `json:"ballot_elections"`
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

type TallyForm struct {
	Ad           Advertisement `json:"ballot_ad"`
	FetchedVotes FetchedVotes  `json:"ballot_fetched_votes"`
	ChoiceScores ChoiceScores  `json:"ballot_choice_scores"`
}

type FetchedVote struct {
	Voter     member.User      `json:"voter_user"`
	Address   id.PublicAddress `json:"voter_address"`
	Elections Elections        `json:"voter_elections"`
}

type FetchedVotes []FetchedVote

func (x FetchedVotes) Len() int           { return len(x) }
func (x FetchedVotes) Less(i, j int) bool { return x[i].Voter < x[j].Voter }
func (x FetchedVotes) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type ChoiceScore struct {
	Choice string  `json:"choice"`
	Score  float64 `json:"score"`
}

type ChoiceScores []ChoiceScore

func (x ChoiceScores) Len() int           { return len(x) }
func (x ChoiceScores) Less(i, j int) bool { return x[i].Score > x[j].Score }
func (x ChoiceScores) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
