package core

import (
	"context"

	"github.com/gov4git/gov4git/mod"
	"github.com/gov4git/gov4git/mod/ballot/load"
	"github.com/gov4git/gov4git/mod/ballot/proto"
	"github.com/gov4git/gov4git/mod/gov"
	"github.com/gov4git/gov4git/mod/id"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/ns"
)

func Close(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	ballotName ns.NS,
	summary proto.Summary,
) git.Change[proto.Outcome] {

	govRepo, govTree := id.CloneOwner(ctx, id.OwnerAddress(govAddr))
	chg := CloseStageOnly(ctx, govAddr, govRepo, govTree, ballotName, summary)
	mod.Commit(ctx, govTree.Public, chg.Msg)
	git.Push(ctx, govRepo.Public)
	return chg
}

func CloseStageOnly(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	govRepo id.OwnerRepo,
	govTree id.OwnerTree,
	ballotName ns.NS,
	summary proto.Summary,
) git.Change[proto.Outcome] {

	// verify ad and strategy are present
	ad, strat := load.LoadStrategy(ctx, govTree.Public, ballotName)
	tally := LoadTally(ctx, govTree.Public, ballotName)
	chg := strat.Close(ctx, govRepo, govTree, &ad, &tally, summary)

	// write outcome
	openOutcomeNS := proto.OpenBallotNS(ballotName).Sub(proto.OutcomeFilebase)
	git.ToFileStage(ctx, govTree.Public, openOutcomeNS.Path(), chg.Result)

	openNS := proto.OpenBallotNS(ballotName)
	closedNS := proto.ClosedBallotNS(ballotName)
	git.RenameStage(ctx, govTree.Public, openNS.Path(), closedNS.Path())

	return chg
}
