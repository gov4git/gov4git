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

	govCloned := git.Clone(ctx, git.Address(govAddr))
	voterOwner := id.CloneOwner(ctx, voterAddr)
	chg := VoteStageOnly(ctx, voterAddr, govAddr, voterOwner, govCloned, ballotName, elections)
	proto.Commit(ctx, voterOwner.Public.Tree(), chg.Msg)
	voterOwner.Public.Push(ctx)

	return chg
}

func VoteStageOnly(
	ctx context.Context,
	voterAddr id.OwnerAddress,
	govAddr gov.GovAddress,
	voterOwner id.OwnerCloned,
	govCloned git.Cloned,
	ballotName ns.NS,
	elections common.Elections,
) git.Change[mail.SeqNo] {

	ad, strat := load.LoadStrategy(ctx, govCloned.Tree(), ballotName, false)

	verifyElections(ctx, strat, voterAddr, govAddr, voterOwner, govCloned, ad, elections)
	envelope := common.VoteEnvelope{
		AdCommit:  git.Head(ctx, govCloned.Repo()),
		Ad:        ad,
		Elections: elections,
	}

	return mail.SendSignedStageOnly(ctx, voterOwner, govCloned.Tree(), common.BallotTopic(ballotName), envelope)
}

func verifyElections(
	ctx context.Context,
	strat common.Strategy,
	voterAddr id.OwnerAddress,
	govAddr gov.GovAddress,
	voterOwner id.OwnerCloned,
	govCloned git.Cloned,
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

	strat.VerifyElections(ctx, voterAddr, govAddr, voterOwner, govCloned, ad, elections)
}

func stringIsIn(s string, in []string) bool {
	for _, in := range in {
		if s == in {
			return true
		}
	}
	return false
}
