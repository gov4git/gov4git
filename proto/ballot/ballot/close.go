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

	govCloned := id.CloneOwner(ctx, id.OwnerAddress(govAddr))
	chg := CloseStageOnly(ctx, govAddr, govCloned, ballotName, summary)
	proto.Commit(ctx, govCloned.Public.Tree(), chg.Msg)
	govCloned.Public.Push(ctx)
	return chg
}

func CloseStageOnly(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	govCloned id.OwnerCloned,
	ballotName ns.NS,
	summary common.Summary,
) git.Change[common.Outcome] {

	govTree := govCloned.Public.Tree()

	// verify ad and strategy are present
	ad, strat := load.LoadStrategy(ctx, govTree, ballotName, false)
	tally := LoadTally(ctx, govTree, ballotName, false)
	chg := strat.Close(ctx, govCloned, &ad, &tally, summary)

	// write outcome
	openOutcomeNS := common.OpenBallotNS(ballotName).Sub(common.OutcomeFilebase)
	git.ToFileStage(ctx, govTree, openOutcomeNS.Path(), chg.Result)

	openNS := common.OpenBallotNS(ballotName)
	closedNS := common.ClosedBallotNS(ballotName)
	git.RenameStage(ctx, govTree, openNS.Path(), closedNS.Path())

	return chg
}
