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
	govAddr gov.OwnerAddress,
	name ballotproto.BallotName,
	title string,
	description string,
) git.Change[form.Map, ballotproto.Advertisement] {

	govCloned := gov.CloneOwner(ctx, govAddr)

	chg := Change_StageOnly(ctx, govCloned, name, title, description)
	proto.Commit(ctx, govCloned.Public.Tree(), chg)
	govCloned.Public.Push(ctx)
	return chg
}

func Change_StageOnly(
	ctx context.Context,
	govCloned gov.OwnerCloned,
	name ballotproto.BallotName,
	title string,
	description string,
) git.Change[form.Map, ballotproto.Advertisement] {

	adNS := ballotproto.BallotPath(name).Append(ballotproto.AdFilebase)

	ad, _ := ballotio.LoadStrategy(ctx, govCloned.Public.Tree(), name)
	ad.Title = title
	ad.Description = description
	git.ToFileStage(ctx, govCloned.Public.Tree(), adNS, ad)

	return git.NewChange(
		fmt.Sprintf("Change ballot %v info", name),
		"ballot_change",
		form.Map{"name": name, "title": title, "description": description},
		ad,
		nil,
	)
}
