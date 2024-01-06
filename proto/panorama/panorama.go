package panorama

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/id"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/gov4git/v2/proto/motion/motionapi"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
)

type Panoramic struct {
	RealBalance      float64                 `json:"real_balance"`
	EffectiveBalance float64                 `json:"effective_balance"`
	Motions          motionproto.MotionViews `json:"motions"`
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

	u := member.FindClonedUser_Local(ctx, cloned, voterOwner)
	accountID := member.UserAccountID(u)
	a := account.Get_Local(
		ctx,
		cloned,
		account.AccountID(accountID),
	)
	real := a.Balance(account.PluralAsset).Quantity

	mvs := motionapi.TrackMotionBatch_Local(ctx, cloned, voterAddr, voterOwner)

	// eff := XXX

	return &Panoramic{
		RealBalance: real,
		// EffectiveBalance: eff,
		Motions: mvs,
	}
}
