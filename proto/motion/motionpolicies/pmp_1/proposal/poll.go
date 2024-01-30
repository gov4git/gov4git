package proposal

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotapi"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotio"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotpolicies/sv"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
)

func init() {
	ctx := context.Background()
	ballotio.Install(
		ctx,
		ProposalApprovalPollPolicyName,
		sv.SV{
			Kernel: ScoreKernel{},
		},
	)
}

type ScoreKernel struct{}

type ScoreKernelState struct {
	MotionID       motionproto.MotionID `json:"motion_id"`
	CostMultiplier float64              `json:"cost_multiplier"`
	Bounty         float64              `json:"bounty"`
}

func (sk ScoreKernel) Score(
	ctx context.Context,
	cloned gov.Cloned,
	ad *ballotproto.Ad,
	el ballotproto.AcceptedElections,

) sv.ScoredVotes {

	state := ballotapi.LoadPolicyState_Local[ScoreKernelState](ctx, cloned, ad.ID)
	qvSK := sv.MakeQVScoreKernel(ctx, state.CostMultiplier)
	return qvSK.Score(ctx, cloned, ad, el)
}

func (sk ScoreKernel) CalcJS(
	ctx context.Context,
	cloned gov.Cloned,
	ad *ballotproto.Ad,

) ballotproto.MarginCalcJS {

	//XXX: needs re-implementing, not reflective of Score at the moment
	state := ballotapi.LoadPolicyState_Local[ScoreKernelState](ctx, cloned, ad.ID)
	js := fmt.Sprintf(scoreKernelMarginJS, state.Bounty)
	return ballotproto.MarginCalcJS(js)
}

const scoreKernelMarginJS = `

// XXX: PLACEHOLDER. DO NOT USE.
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
          "help": {
               "label": "Help",
               "description": "This ballot uses Quadratic Voting to determine the approval score of a PR. " +
				"Accepted PRs reward the authors with the bounties accumulated in the resolved issues,
				and reviewers with the proceeds of votes on the approval score. Rejected PRs refund voters.",
               "value": null,
          },
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
