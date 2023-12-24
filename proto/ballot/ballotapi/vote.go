package ballotapi

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotio"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/id"
	"github.com/gov4git/gov4git/v2/proto/mail"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Vote(
	ctx context.Context,
	voterAddr id.OwnerAddress,
	govAddr gov.Address,
	ballotName ballotproto.BallotName,
	elections ballotproto.Elections,
) git.Change[form.Map, mail.RequestEnvelope[ballotproto.VoteEnvelope]] {

	govCloned := gov.Clone(ctx, govAddr)
	voterOwner := id.CloneOwner(ctx, voterAddr)
	chg := Vote_StageOnly(ctx, voterAddr, voterOwner, govCloned, ballotName, elections)
	proto.Commit(ctx, voterOwner.Public.Tree(), chg)
	voterOwner.Public.Push(ctx)

	return chg
}

func Vote_StageOnly(
	ctx context.Context,
	voterAddr id.OwnerAddress,
	voterOwner id.OwnerCloned,
	govCloned gov.Cloned,
	ballotName ballotproto.BallotName,
	elections ballotproto.Elections,
) git.Change[form.Map, mail.RequestEnvelope[ballotproto.VoteEnvelope]] {

	ad, strat := ballotio.LoadStrategy(ctx, govCloned.Tree(), ballotName)

	must.Assertf(ctx, !ad.Closed, "ballot is closed")
	must.Assertf(ctx, !ad.Frozen, "ballot is frozen")

	verifyElections(ctx, strat, voterAddr, govCloned.Address(), voterOwner, govCloned, ad, elections)
	envelope := ballotproto.VoteEnvelope{
		AdCommit:  git.Head(ctx, govCloned.Repo()),
		Ad:        ad,
		Elections: elections,
	}

	// record vote in voter's repo
	voterTree := voterOwner.Public.Tree()
	govCred := id.GetPublicCredentials(ctx, govCloned.Tree())
	voteLogNS := ballotproto.VoteLogPath(govCred.ID, ballotName)
	// read current vote log
	voteLog, err := git.TryFromFile[ballotproto.VoteLog](ctx, voterTree, voteLogNS)
	if git.IsNotExist(err) {
		voteLog = ballotproto.VoteLog{
			GovID:         govCred.ID,
			GovAddress:    govCloned.Address(),
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
	sendChg := mail.Request_StageOnly(ctx, voterOwner, govCloned.Tree(), ballotproto.BallotTopic(ballotName), envelope)
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
	strat ballotproto.Strategy,
	voterAddr id.OwnerAddress,
	govAddr gov.Address,
	voterOwner id.OwnerCloned,
	govCloned gov.Cloned,
	ad ballotproto.Advertisement,
	elections ballotproto.Elections,
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
