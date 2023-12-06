package boot

import (
	"context"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/docket/policies/pmp"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/gov4git/proto/treasury"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func Boot(
	ctx context.Context,
	ownerAddr gov.OwnerAddress,
) git.Change[form.None, id.PrivateCredentials] {

	ownerCloned := gov.CloneOwner(ctx, ownerAddr)
	privChg := Boot_Local(ctx, ownerCloned)
	ownerCloned.Public.Push(ctx)
	ownerCloned.Private.Push(ctx)
	return privChg
}

func Boot_Local(
	ctx context.Context,
	ownerCloned gov.OwnerCloned,
) git.Change[form.None, id.PrivateCredentials] {

	// initialize project identity
	chg := id.Init_Local(ctx, ownerCloned.IDOwnerCloned())

	// create group everybody
	chg2 := member.SetGroup_StageOnly(ctx, ownerCloned.PublicClone(), member.Everybody)

	// create treasury accounts
	treasury.Boot_StageOnly(ctx, ownerCloned.PublicClone())

	// create PMP accounts
	pmp.Boot_StageOnly(ctx, ownerCloned.PublicClone())

	proto.Commit(ctx, ownerCloned.Public.Tree(), chg2)
	return chg
}
