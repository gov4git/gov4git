package ballotapi

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotio"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/id"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Erase(
	ctx context.Context,
	govAddr gov.OwnerAddress,
	ballotName ballotproto.BallotName,
) git.Change[form.Map, bool] {

	govCloned := id.CloneOwner(ctx, id.OwnerAddress(govAddr))
	chg := Erase_StageOnly(ctx, govCloned, ballotName)
	proto.Commit(ctx, govCloned.Public.Tree(), chg)
	govCloned.Public.Push(ctx)
	return chg
}

func Erase_StageOnly(
	ctx context.Context,
	govCloned id.OwnerCloned,
	ballotName ballotproto.BallotName,
) git.Change[form.Map, bool] {

	govTree := govCloned.Public.Tree()

	// verify ad and strategy are present
	ballotio.LoadStrategy(ctx, govTree, ballotName)

	// write outcome
	ballotNS := ballotproto.BallotPath(ballotName)
	_, err := git.TreeRemove(ctx, govTree, ballotNS)
	must.NoError(ctx, err)

	return git.NewChange(
		fmt.Sprintf("Erased ballot %v", ballotName),
		"ballot_erase",
		form.Map{"name": ballotName},
		true,
		nil,
	)
}
