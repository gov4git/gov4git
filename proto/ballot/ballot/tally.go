package ballot

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/load"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/mail"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

func Tally(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	ballotName ns.NS,
) git.Change[common.Tally] {

	govRepo, govTree := id.CloneOwner(ctx, id.OwnerAddress(govAddr))
	chg := TallyStageOnly(ctx, govAddr, govRepo, govTree, ballotName)
	proto.Commit(ctx, git.Worktree(ctx, govRepo.Home), chg.Msg)
	git.Push(ctx, govRepo.Home)
	return chg
}

func TallyStageOnly(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	govRepo id.OwnerRepo,
	govTree id.OwnerTree,
	ballotName ns.NS,
) git.Change[common.Tally] {

	communityTree := govTree.Home

	ad, strat := load.LoadStrategy(ctx, communityTree, ballotName)

	// list participating users
	users := member.ListGroupUsersLocal(ctx, communityTree, ad.Participants)

	// get user accounts
	accounts := make([]member.Account, len(users))
	for i, user := range users {
		accounts[i] = member.GetUserLocal(ctx, communityTree, user)
	}

	// fetch votes from users
	var fetchedVotes common.FetchedVotes
	for i, account := range accounts {
		fetchedVotes = append(fetchedVotes,
			fetchVotes(ctx, govAddr, govRepo, govTree, ballotName, users[i], account).Result...)
	}

	// read current tally
	var currentTally *common.Tally
	if tryCurrentTally, err := must.Try1(func() common.Tally { return LoadTally(ctx, communityTree, ballotName) }); err == nil {
		currentTally = &tryCurrentTally
	}

	updatedTally := strat.Tally(ctx, govRepo, govTree, &ad, currentTally, fetchedVotes).Result

	// write updated tally
	openTallyNS := common.OpenBallotNS(ballotName).Sub(common.TallyFilebase)
	git.ToFileStage(ctx, communityTree, openTallyNS.Path(), updatedTally)

	return git.Change[common.Tally]{
		Result: updatedTally,
		Msg:    fmt.Sprintf("Tally votes on ballot %v", ballotName),
	}
}

func fetchVotes(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	govRepo id.OwnerRepo,
	govTree id.OwnerTree,
	ballotName ns.NS,
	user member.User,
	account member.Account,
) git.Change[common.FetchedVotes] {

	fetched := common.FetchedVotes{}
	respond := func(ctx context.Context, req common.VoteEnvelope, _ id.SignedPlaintext) (resp common.VoteEnvelope, err error) {

		if !req.VerifyConsistency() {
			return common.VoteEnvelope{}, fmt.Errorf("vote envelope is not valid")
		}
		fetched = append(fetched,
			common.FetchedVote{
				Voter:     user,
				Address:   account.Home,
				Elections: req.Elections,
			})
		return req, nil
	}

	_, voterHomeTree := git.Clone(ctx, git.Address(account.Home))
	mail.ReceiveSignedStageOnly(
		ctx,
		govTree,
		account.Home,
		voterHomeTree,
		common.BallotTopic(ballotName),
		respond,
	)

	return git.Change[common.FetchedVotes]{
		Result: fetched,
		Msg:    fmt.Sprintf("Fetched votes from user %v on ballot %v", user, ballotName),
	}
}

func LoadTally(
	ctx context.Context,
	communityTree *git.Tree,
	ballotName ns.NS,
) common.Tally {
	openTallyNS := common.OpenBallotNS(ballotName).Sub(common.TallyFilebase)
	return git.FromFile[common.Tally](ctx, communityTree, openTallyNS.Path())
}
