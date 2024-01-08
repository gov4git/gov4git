package ballotapi

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotio"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func Change(
	ctx context.Context,
	addr gov.OwnerAddress,
	id ballotproto.BallotID,
	title string,
	description string,

) git.Change[form.Map, ballotproto.Ad] {

	cloned := gov.CloneOwner(ctx, addr)

	chg := Change_StageOnly(ctx, cloned, id, title, description)
	proto.Commit(ctx, cloned.Public.Tree(), chg)
	cloned.Public.Push(ctx)
	return chg
}

func Change_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	id ballotproto.BallotID,
	title string,
	description string,

) git.Change[form.Map, ballotproto.Ad] {

	ad, _ := ballotio.LoadPolicy(ctx, cloned.Public.Tree(), id)
	ad.Title = title
	ad.Description = description
	git.ToFileStage(ctx, cloned.Public.Tree(), id.AdNS(), ad)

	return git.NewChange(
		fmt.Sprintf("Change ballot %v info", id),
		"ballot_change",
		form.Map{"name": id, "title": title, "description": description},
		ad,
		nil,
	)
}
