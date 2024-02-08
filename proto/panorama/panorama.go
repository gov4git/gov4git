package panorama

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotapi"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotio"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/id"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/gov4git/v2/proto/motion/motionapi"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
)

type Panoramic struct {
	RealBalance      float64                 `json:"real_balance"`
	ProjectedBalance float64                 `json:"projected_balance"`
	RealMotions      motionproto.MotionViews `json:"real_motions"`
	ProjectedMotions motionproto.MotionViews `json:"projected_motions"`
}

func Panorama(
	ctx context.Context,
	addr gov.Address,
	voterAddr id.OwnerAddress,

) *Panoramic {

	voterOwner := id.CloneOwner(ctx, voterAddr)
	return Panorama_Local(ctx, gov.Clone(ctx, addr), voterAddr, voterOwner)
}

func Panorama_Local(
	ctx context.Context,
	cloned gov.Cloned,
	voterAddr id.OwnerAddress,
	voterOwner id.OwnerCloned,

) *Panoramic {

	voterUser := member.FindClonedUser_Local(ctx, cloned, voterOwner)
	// voterProfile := member.GetUser_Local(ctx, cloned, voterUser)
	voterAccountID := member.UserAccountID(voterUser)

	realBalance := account.Get_Local(ctx, cloned, account.AccountID(voterAccountID)).Balance(account.PluralAsset).Quantity

	realMVS := motionapi.TrackMotionBatch_Local(ctx, cloned, voterAddr, voterOwner)

	// apply pending votes to governance
	for _, ad := range ballotapi.List_Local(ctx, cloned) {
		// TODO: simulates votes by directly processing the vote mail (rather than the vote log).
		// will require processing the mail without using the private repo for credentials (i.e. without verifying).
		//
		// ballotapi.TallyVoterCloned_StageOnly(
		// 	ctx,
		// 	gov.LiftCloned(ctx, cloned),
		// 	ad.ID,
		// 	voterUser,
		// 	voterProfile,
		// 	voterOwner.PublicClone(),
		// )
		if ballotio.TryLookupPolicy(ctx, ad.Policy) == nil { // only consider ballots with known policies
			continue
		}
		if ad.Closed {
			continue
		}
		vs := ballotapi.Track_StageOnly(
			ctx,
			voterAddr,
			voterOwner,
			cloned,
			ad.ID,
		)
		fetchedVote := ballotapi.FetchedVote{
			Voter:     voterUser,
			Address:   voterAddr.Public,
			Elections: vs.PendingVotes,
		}
		ballotapi.TallyFetchedVotes_StageOnly(
			ctx,
			cloned,
			ad.ID,
			ballotapi.FetchedVotes{fetchedVote},
		)
	}

	// rescore and update motions
	motionapi.Pipeline(ctx, gov.LiftCloned(ctx, cloned))

	projMVS := motionapi.TrackMotionBatch_Local(ctx, cloned, voterAddr, voterOwner)
	// TODO: fully simulate tallying by processing voter's mail (see above)
	for i := range projMVS {
		projMVS[i].Voter = nil
	}

	projBalance := account.Get_Local(ctx, cloned, account.AccountID(voterAccountID)).Balance(account.PluralAsset).Quantity

	return &Panoramic{
		RealBalance:      realBalance,
		ProjectedBalance: projBalance,
		RealMotions:      realMVS,
		ProjectedMotions: projMVS,
	}
}
