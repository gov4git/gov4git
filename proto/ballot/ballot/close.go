package ballot

import (
	"context"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/load"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/ns"
)

func Close(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	ballotName ns.NS,
	summary common.Summary,
) git.Change[common.Outcome] {

	govRepo, govTree := id.CloneOwner(ctx, id.OwnerAddress(govAddr))
	chg := CloseStageOnly(ctx, govAddr, govRepo, govTree, ballotName, summary)
	proto.Commit(ctx, govTree.Home, chg.Msg)
	git.Push(ctx, govRepo.Home)
	return chg
}

func CloseStageOnly(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	govRepo id.OwnerRepo,
	govTree id.OwnerTree,
	ballotName ns.NS,
	summary common.Summary,
) git.Change[common.Outcome] {

	// verify ad and strategy are present
	ad, strat := load.LoadStrategy(ctx, govTree.Home, ballotName)
	tally := LoadTally(ctx, govTree.Home, ballotName)
	chg := strat.Close(ctx, govRepo, govTree, &ad, &tally, summary)

	// write outcome
	openOutcomeNS := common.OpenBallotNS(ballotName).Sub(common.OutcomeFilebase)
	git.ToFileStage(ctx, govTree.Home, openOutcomeNS.Path(), chg.Result)

	openNS := common.OpenBallotNS(ballotName)
	closedNS := common.ClosedBallotNS(ballotName)
	git.RenameStage(ctx, govTree.Home, openNS.Path(), closedNS.Path())

	return chg
}
