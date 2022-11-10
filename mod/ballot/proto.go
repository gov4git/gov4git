package ballot

import (
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/lib/ns"
	"github.com/gov4git/gov4git/mod"
	"github.com/gov4git/gov4git/mod/gov"
	"github.com/gov4git/gov4git/mod/member"
)

var (
	ballotNS   = mod.RootNS.Sub("ballots")
	adFilebase = "ballot_ad.json"
)

func OpenBallotNS[S Strategy](name ns.NS) ns.NS {
	var s S
	return ballotNS.Sub("open").Sub(s.StrategyName()).Join(name)
}

func ClosedBallotNS[S Strategy](name ns.NS) ns.NS {
	var s S
	return ballotNS.Sub("closed").Sub(s.StrategyName()).Join(name)
}

func BallotTopic[S Strategy](name ns.NS) string {
	var s S
	return s.StrategyName() + ":" + name.Path()
}

type AdForm struct {
	Community    gov.CommunityAddress `json:"community"`
	Name         ns.NS                `json:"path"`
	Title        string               `json:"title"`
	Description  string               `json:"description"`
	Choices      []string             `json:"choices"`
	Strategy     string               `json:"strategy"`
	Participants member.Group         `json:"participants_group"`
	ParentCommit git.CommitHash       `json:"parent_commit"`
}

type VoteForm struct {
	VoteChoice         string  `json:"vote_choice"`
	VoteStrengthChange float64 `json:"vote_strength_change"`
}

type VoteEnvelope struct {
	AdCommit  git.CommitHash `json:"ballot_ad_commit"`
	Ad        AdForm         `json:"ballot_ad"`
	Elections []VoteForm     `json:"ballot_elections"`
}

type TallyForm struct{}
