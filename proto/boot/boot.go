package boot

import (
	"context"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func Boot(
	ctx context.Context,
	ownerAddr id.OwnerAddress,
) git.Change[form.None, id.PrivateCredentials] {

	ownerCloned := id.CloneOwner(ctx, ownerAddr)
	privChg := BootLocal(ctx, ownerAddr, ownerCloned)
	ownerCloned.Public.Push(ctx)
	ownerCloned.Private.Push(ctx)
	return privChg
}

func BootLocal(
	ctx context.Context,
	ownerAddr id.OwnerAddress,
	ownerCloned id.OwnerCloned,
) git.Change[form.None, id.PrivateCredentials] {

	chg := id.InitLocal(ctx, ownerAddr, ownerCloned)
	chg2 := member.SetGroupStageOnly(ctx, ownerCloned.Public.Tree(), member.Everybody)
	proto.Commit(ctx, ownerCloned.Public.Tree(), chg2)

	return chg
}
