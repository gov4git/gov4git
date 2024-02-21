package concern

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotio"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotpolicies/sv"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
)

func init() {
	ctx := context.Background()
	ballotio.Install(
		ctx,
		ConcernPriorityPollPolicyName,
		sv.SV{
			Kernel: ScoreKernel{},
		},
	)
}

const ConcernPriorityPollPolicyName ballotproto.PolicyName = "pmp-concern-priority-v1"

type ScoreKernel struct{}

func (sk ScoreKernel) Score(
	ctx context.Context,
	cloned gov.Cloned,
	ad *ballotproto.Ad,
	el ballotproto.AcceptedElections,

) sv.ScoredVotes {

	qvSK := sv.MakeQVScoreKernel(ctx, 1.0)
	return qvSK.Score(ctx, cloned, ad, el)
}

func (sk ScoreKernel) CalcJS(
	ctx context.Context,
	cloned gov.Cloned,
	ad *ballotproto.Ad,
	tally *ballotproto.Tally,

) *ballotproto.Margin {

	qvSK := sv.MakeQVScoreKernel(ctx, 1.0)
	margin := qvSK.CalcJS(ctx, cloned, ad, tally)
	margin.Reward = &ballotproto.MarginCalculator{
		Label:       "Reward",
		Description: "Potential reward to the voter, assuming a favorable outcome",
		FnJS:        rewardJSFmt,
	}
	return margin
}

const (
	rewardJSFmt = `
	function(voteUser, voteChoice, voteImpact) {
		return 0;
	}
	`
)
