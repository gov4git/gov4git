package proposal

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotapi"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotio"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotstrategies/sv"
	"github.com/gov4git/gov4git/v2/proto/gov"
)

func init() {
	ctx := context.Background()
	ballotio.Install(
		ctx,
		ProposalApprovalPollStrategyName,
		sv.SV{
			Kernel: ScoreKernel{},
		},
	)
}

type ScoreKernel struct{}

type ScoreKernelState struct {
	Bounty float64 `json:"bounty"`
}

func (sk ScoreKernel) Score(
	ctx context.Context,
	cloned gov.Cloned,
	ad *ballotproto.Advertisement,
	el ballotproto.AcceptedElections,

) sv.ScoredVotes {

	qvSK := sv.QVScoreKernel{}
	return qvSK.Score(ctx, cloned, ad, el)
}

func (sk ScoreKernel) CalcJS(
	ctx context.Context,
	cloned gov.Cloned,
	ad *ballotproto.Advertisement,

) ballotproto.MarginCalcJS {

	state := ballotapi.LoadStrategyState_Local[ScoreKernelState](ctx, cloned, ad.ID)
	js := fmt.Sprintf(scoreKernelMarginJS, state.Bounty)
	return ballotproto.MarginCalcJS(js)
}

const scoreKernelMarginJS = `

function calcMargin(currentTally, voteUser, voteChoice, targetVote) {

	// compute currentVote and currentCharge (for logged in user)

	var currentVote = 0.0;
	var currentCharge = 0.0;

	var currentScoresByUser = currentTally.scores_by_user[voteUser];
	if (currentScoresByUser !== undefined) {
		var currentChoiceByUser = currentScoresByUser[voteChoice];
		if (currentChoiceByUser !== undefined) {
			currentVote = currentChoiceByUser.score;
			currentCharge = currentChoiceByUser.strength;
		}
	}

	// compute targetCharge and cost difference

     var targetCharge = targetVote * targetVote;
     var cost = targetCharge - currentCharge;

	// compute reward for reviewers

	var reward = 0.0;
	if (targetVote > 0.0) {

		var totalCharges = cost;
		var totalChargesForPositiveVotes = cost;

		for (const [user, charge] of Object.entries(currentTally.charges)) {
			var userVote = 0.0;
			var scoresByUser = currentTally.scores_by_user[user];
			if (scoresByUser !== undefined) {
				var choiceByUser = scoresByUser[voteChoice];
				if (choiceByUser !== undefined) {
					userVote = choiceByUser.score;
				}
			}
			if (userVote > 0.0) {
				totalCharges += charge;
			}
		}

		reward = totalCharges * ((currentCharge + cost) / totalChargesForPositiveVotes);
	}

     return {
          "currentVote": {
               "label": "Current vote",
               "description": "Your current vote",
               "value": currentVote,
          },
          "targetVote": {
               "label": "Target vote",
               "description": "Your target vote",
               "value": targetVote,
          },
          "cost": {
               "label": "Cost",
               "description": "Cost of changing your vote",
               "value": cost,
          },
		"bounty": {
               "label": "Bounty",
               "description": "Bounty rewarded to the author of the pull request, if merged.",
               "value": %v,
		},
		"reward": {
               "label": "Reward",
               "description": "Potential reward to you as reviewer, if pull request is merged.",
               "value": reward,
		},
     }
}
`
