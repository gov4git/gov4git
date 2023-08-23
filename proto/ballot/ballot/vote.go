package ballot

import (
	"context"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/load"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/mail"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Vote(
	ctx context.Context,
	voterAddr id.OwnerAddress,
	govAddr gov.GovAddress,
	ballotName common.BallotName,
	elections common.Elections,
) git.Change[form.Map, mail.RequestEnvelope[common.VoteEnvelope]] {

	govCloned := git.CloneOne(ctx, git.Address(govAddr))
	voterOwner := id.CloneOwner(ctx, voterAddr)
	chg := Vote_StageOnly(ctx, voterAddr, govAddr, voterOwner, govCloned, ballotName, elections)
	proto.Commit(ctx, voterOwner.Public.Tree(), chg)
	voterOwner.Public.Push(ctx)

	return chg
}

func Vote_StageOnly(
	ctx context.Context,
	voterAddr id.OwnerAddress,
	govAddr gov.GovAddress,
	voterOwner id.OwnerCloned,
	govCloned git.Cloned,
	ballotName common.BallotName,
	elections common.Elections,
) git.Change[form.Map, mail.RequestEnvelope[common.VoteEnvelope]] {

	ad, strat := load.LoadStrategy(ctx, govCloned.Tree(), ballotName)

	must.Assertf(ctx, !ad.Closed, "ballot is closed")
	must.Assertf(ctx, !ad.Frozen, "ballot is frozen")

	verifyElections(ctx, strat, voterAddr, govAddr, voterOwner, govCloned, ad, elections)
	envelope := common.VoteEnvelope{
		AdCommit:  git.Head(ctx, govCloned.Repo()),
		Ad:        ad,
		Elections: elections,
	}

	// record vote in voter's repo
	voterTree := voterOwner.Public.Tree()
	govCred := id.GetPublicCredentials(ctx, govCloned.Tree())
	voteLogNS := common.VoteLogPath(govCred.ID, ballotName)
	// read current vote log
	voteLog, err := git.TryFromFile[common.VoteLog](ctx, voterTree, voteLogNS)
	if git.IsNotExist(err) {
		voteLog = common.VoteLog{
			GovID:         govCred.ID,
			GovAddress:    govAddr,
			Ballot:        ballotName,
			VoteEnvelopes: nil,
		}
	} else {
		must.NoError(ctx, err)
	}
	// append new vote
	voteLog.VoteEnvelopes = append(voteLog.VoteEnvelopes, envelope)
	git.ToFileStage(ctx, voterTree, voteLogNS, voteLog)

	// send vote to community by mail
	sendChg := mail.Request_StageOnly(ctx, voterOwner, govCloned.Tree(), common.BallotTopic(ballotName), envelope)
	return git.NewChange(
		"Cast vote",
		"ballot_vote",
		form.Map{"ballot_name": ballotName, "elections": elections},
		sendChg.Result,
		form.Forms{sendChg},
	)
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

	tally := LoadTally(ctx, govCloned.Tree(), ad.Name)
	strat.VerifyElections(ctx, voterAddr, govAddr, voterOwner, govCloned, &ad, &tally, elections)
}

func stringIsIn(s string, in []string) bool {
	for _, in := range in {
		if s == in {
			return true
		}
	}
	return false
}
