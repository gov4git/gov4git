package ballot

import (
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/lib/ns"
	"github.com/gov4git/gov4git/lib/util"
	"github.com/gov4git/gov4git/mod"
	"github.com/gov4git/gov4git/mod/gov"
	"github.com/gov4git/gov4git/mod/id"
	"github.com/gov4git/gov4git/mod/member"
)

var (
	ballotNS      = mod.RootNS.Sub("ballots")
	adFilebase    = "ballot_ad.json"
	tallyFilebase = "ballot_tally.json"
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
	Ad             AdForm             `json:"ballot_ad"`
	Participations ParticipationForms `json:"ballot_participations"`
	ChoiceTallies  ChoiceTallyForms   `json:"ballot_choice_tallies"`
}

type ParticipationForm struct {
	Name      member.User      `json:"voter_user_name"`
	Address   id.PublicAddress `json:"voter_address"`
	Envelopes []VoteEnvelope   `json:"voter_vote_envelopes"`
}

type ParticipationForms []ParticipationForm

func (x ParticipationForms) Len() int           { return len(x) }
func (x ParticipationForms) Less(i, j int) bool { return x[i].Name < x[j].Name }
func (x ParticipationForms) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type ChoiceTallyForm struct {
	Choice string  `json:"choice"`
	Score  float64 `json:"score"`
}

type ChoiceTallyForms []ChoiceTallyForm

func (x ChoiceTallyForms) Len() int           { return len(x) }
func (x ChoiceTallyForms) Less(i, j int) bool { return x[i].Score > x[j].Score }
func (x ChoiceTallyForms) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
