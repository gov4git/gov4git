package proposal

import (
	"context"

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
	MotionID              motionproto.MotionID `json:"motion_id"`
	InverseCostMultiplier float64              `json:"inverse_cost_multiplier"`
	Bounty                float64              `json:"bounty"`
}

func (sk ScoreKernel) Score(
	ctx context.Context,
	cloned gov.Cloned,
	ad *ballotproto.Ad,
	el ballotproto.AcceptedElections,

) sv.ScoredVotes {

	state := ballotapi.LoadPolicyState_Local[ScoreKernelState](ctx, cloned, ad.ID)
	qvSK := sv.MakeQVScoreKernel(ctx, state.InverseCostMultiplier)
	return qvSK.Score(ctx, cloned, ad, el)
}

func (sk ScoreKernel) CalcJS(
	ctx context.Context,
	cloned gov.Cloned,
	ad *ballotproto.Ad,
	tally *ballotproto.Tally,

) *ballotproto.Margin {

	state := ballotapi.LoadPolicyState_Local[ScoreKernelState](ctx, cloned, ad.ID)
	qvSK := sv.MakeQVScoreKernel(ctx, state.InverseCostMultiplier)
	return qvSK.CalcJS(ctx, cloned, ad, tally)
}
