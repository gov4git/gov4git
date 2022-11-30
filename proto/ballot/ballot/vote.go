package ballot

import (
	"context"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/load"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/mail"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

func Vote(
	ctx context.Context,
	voterAddr id.OwnerAddress,
	govAddr gov.GovAddress,
	ballotName ns.NS,
	elections common.Elections,
) git.Change[mail.SeqNo] {

	govRepo := git.CloneRepo(ctx, git.Address(govAddr))
	voterRepo, voterTree := id.CloneOwner(ctx, voterAddr)
	chg := VoteStageOnly(ctx, voterAddr, govAddr, voterTree, govRepo, ballotName, elections)
	proto.Commit(ctx, voterTree.Public, chg.Msg)
	git.Push(ctx, voterRepo.Public)

	return chg
}

func VoteStageOnly(
	ctx context.Context,
	voterAddr id.OwnerAddress,
	govAddr gov.GovAddress,
	voterTree id.OwnerTree,
	govRepo *git.Repository,
	ballotName ns.NS,
	elections common.Elections,
) git.Change[mail.SeqNo] {

	govTree := git.Worktree(ctx, govRepo)

	ad, strat := load.LoadStrategy(ctx, govTree, ballotName, false)

	verifyElections(ctx, strat, voterAddr, govAddr, voterTree, govTree, ad, elections)
	envelope := common.VoteEnvelope{
		AdCommit:  git.Head(ctx, govRepo),
		Ad:        ad,
		Elections: elections,
	}

	return mail.SendSignedStageOnly(ctx, voterTree, govTree, common.BallotTopic(ballotName), envelope)
}

func verifyElections(
	ctx context.Context,
	strat common.Strategy,
	voterAddr id.OwnerAddress,
	govAddr gov.GovAddress,
	voterTree id.OwnerTree,
	govTree *git.Tree,
	ad common.Advertisement,
	elections common.Elections,
) {
	// check elections use available choices
	if len(ad.Choices) > 0 {
		for _, e := range elections {
			if !stringIsIn(e.VoteChoice, ad.Choices) {
				must.Errorf(ctx, "election %v is not an available choice", e.VoteChoice)
			}
		}
	}

	strat.VerifyElections(ctx, voterAddr, govAddr, voterTree, govTree, ad, elections)
}

func stringIsIn(s string, in []string) bool {
	for _, in := range in {
		if s == in {
			return true
		}
	}
	return false
}
