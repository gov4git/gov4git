package sv

import (
	"context"
	"fmt"
	"math"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/lib4git/form"
)

type QVScoreKernel struct {
	InverseCostMultiplier float64 `json:"inverse_cost_multiplier"`
}

func MakeQVScoreKernel(ctx context.Context, inverseCostMultiplier float64) QVScoreKernel {
	return QVScoreKernel{
		InverseCostMultiplier: inverseCostMultiplier,
	}
}

func (k QVScoreKernel) Score(
	ctx context.Context,
	cloned gov.Cloned,
	ad *ballotproto.Ad,
	el ballotproto.AcceptedElections,

) ScoredVotes {

	// aggregate voting strength on each choice
	score := map[string]ballotproto.StrengthAndScore{}
	for _, el := range el {
		x := score[el.Vote.VoteChoice]
		x.Strength += el.Vote.VoteStrengthChange
		score[el.Vote.VoteChoice] = x
	}
	// compute score per choice
	for choice, ss := range score {
		score[choice] = ballotproto.StrengthAndScore{
			Strength: ss.Strength,
			Score:    qvScoreFromStrength(ss.Strength, k.InverseCostMultiplier),
		}
	}
	// compute aggregate cost
	cost := 0.0
	for _, x := range score {
		cost += math.Abs(x.Strength)
	}
	return ScoredVotes{Votes: el, Score: score, Cost: cost}
}

// score = SIGN(strength) * SQRT( |strength| * max(1, inverseCostMultiplier)) )
func qvScoreFromStrength(strength float64, inverseCostMultiplier float64) float64 {
	sign := 1.0
	if strength < 0 {
		sign = -1.0
	}
	return sign * math.Sqrt(math.Abs(strength*max(1, inverseCostMultiplier)))
}

// strength = SIGN(score) * score^2 / max(1, inverseCostMultiplier)
func qvStrengthFromScore(score float64, inverseCostMultiplier float64) float64 {
	sign := 1.0
	if score < 0 {
		sign = -1.0
	}
	return sign * score * score / max(1, inverseCostMultiplier)
}

func (k QVScoreKernel) CalcJS(
	ctx context.Context,
	cloned gov.Cloned,
	ad *ballotproto.Ad,
	tally *ballotproto.Tally,

) *ballotproto.Margin {

	return &ballotproto.Margin{
		Help: &ballotproto.MarginCalculator{
			Label:       "Help",
			Description: "Description of ballot",
			FnJS: fmt.Sprintf(
				`function() { return %q }`,
				fmt.Sprintf("The vote impact of `P` credits is `SQRT(%0.6f*P)`.", k.InverseCostMultiplier),
			),
		},
		Cost: &ballotproto.MarginCalculator{
			Label:       "Cost",
			Description: "Additional cost to reach a desired impact",
			FnJS:        fmt.Sprintf(qvCostJSFmt, k.InverseCostMultiplier, form.SprintJSON(tally)),
		},
		Impact: &ballotproto.MarginCalculator{
			Label:       "Impact",
			Description: "Additional impact to reach a desired cost",
			FnJS:        fmt.Sprintf(qvImpactJSFmt, k.InverseCostMultiplier, form.SprintJSON(tally)),
		},
	}
}

const (
	// Additional cost to reach a desired impact
	qvCostJSFmt = `
	function(voteUser, voteChoice, voteImpact) {
		let inverseCostMultiplier = %f;
		let tally = %s;
		var currentVoteImpact = 0.0;
		var currentVoteCost = 0.0;
		var currentScoresByUser = currentTally.scores_by_user[voteUser];
		if (currentScoresByUser !== undefined) {
			var currentChoiceByUser = currentScoresByUser[voteChoice];
			if (currentChoiceByUser !== undefined) {
				currentVoteImpact = currentChoiceByUser.score;
				currentVoteCost = Math.abs(currentChoiceByUser.strength);
			}
		}

		var voteCost = voteImpact * voteImpact / Math.max(1, inverseCostMultiplier);
		var costDiff = voteCost - currentVoteCost;
		return costDiff;
	}
	`

	// Additional impact to reach a desired cost
	qvImpactJSFmt = `
	function(voteUser, voteChoice, voteCost) {
		let inverseCostMultiplier = %f;
		let tally = %s;
		var currentVoteImpact = 0.0;
		var currentVoteCost = 0.0;
		var currentScoresByUser = currentTally.scores_by_user[voteUser];
		if (currentScoresByUser !== undefined) {
			var currentChoiceByUser = currentScoresByUser[voteChoice];
			if (currentChoiceByUser !== undefined) {
				currentVoteImpact = currentChoiceByUser.score;
				currentVoteCost = Math.abs(currentChoiceByUser.strength);
			}
		}

		XXX
		return costDiff;
	}
	`
)
