package ballot

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/load"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/ns"
)

func Change(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	name ns.NS,
	title string,
	description string,
) git.Change[form.Map, common.Advertisement] {

	govCloned := id.CloneOwner(ctx, id.OwnerAddress(govAddr))

	chg := Change_StageOnly(ctx, govAddr, govCloned, name, title, description)
	proto.Commit(ctx, govCloned.Public.Tree(), chg)
	govCloned.Public.Push(ctx)
	return chg
}

func Change_StageOnly(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	govCloned id.OwnerCloned,
	name ns.NS,
	title string,
	description string,
) git.Change[form.Map, common.Advertisement] {

	adNS := common.BallotPath(name).Append(common.AdFilebase)

	ad, _ := load.LoadStrategy(ctx, govCloned.Public.Tree(), name)
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
