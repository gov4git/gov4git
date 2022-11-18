package core

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/mod"
	"github.com/gov4git/gov4git/mod/ballot/load"
	"github.com/gov4git/gov4git/mod/ballot/proto"
	"github.com/gov4git/gov4git/mod/gov"
	"github.com/gov4git/gov4git/mod/id"
	"github.com/gov4git/gov4git/mod/mail"
	"github.com/gov4git/gov4git/mod/member"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

func Tally(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	ballotName ns.NS,
) git.Change[proto.Tally] {

	govRepo, govTree := id.CloneOwner(ctx, id.OwnerAddress(govAddr))
	chg := TallyStageOnly(ctx, govAddr, govRepo, govTree, ballotName)
	mod.Commit(ctx, git.Worktree(ctx, govRepo.Public), chg.Msg)
	git.Push(ctx, govRepo.Public)
	return chg
}

func TallyStageOnly(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	govRepo id.OwnerRepo,
	govTree id.OwnerTree,
	ballotName ns.NS,
) git.Change[proto.Tally] {

	communityTree := govTree.Public

	ad, strat := load.LoadStrategy(ctx, communityTree, ballotName)

	// list participating users
	users := member.ListGroupUsersLocal(ctx, communityTree, ad.Participants)

	// get user accounts
	accounts := make([]member.Account, len(users))
	for i, user := range users {
		accounts[i] = member.GetUserLocal(ctx, communityTree, user)
	}

	// fetch votes from users
	var fetchedVotes proto.FetchedVotes
	for i, account := range accounts {
		fetchedVotes = append(fetchedVotes,
			fetchVotes(ctx, govAddr, govRepo, govTree, ballotName, users[i], account).Result...)
	}

	// read current tally
	var currentTally *proto.Tally
	if tryCurrentTally, err := must.Try1(func() proto.Tally { return LoadTally(ctx, communityTree, ballotName) }); err == nil {
		currentTally = &tryCurrentTally
	}

	updatedTally := strat.Tally(ctx, govRepo, govTree, &ad, currentTally, fetchedVotes).Result

	// write updated tally
	openTallyNS := proto.OpenBallotNS(ballotName).Sub(proto.TallyFilebase)
	git.ToFileStage(ctx, communityTree, openTallyNS.Path(), updatedTally)

	return git.Change[proto.Tally]{
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
) git.Change[proto.FetchedVotes] {

	fetched := proto.FetchedVotes{}
	respond := func(ctx context.Context, req proto.VoteEnvelope, _ id.SignedPlaintext) (resp proto.VoteEnvelope, err error) {

		if !req.VerifyConsistency() {
			return proto.VoteEnvelope{}, fmt.Errorf("vote envelope is not valid")
		}
		fetched = append(fetched,
			proto.FetchedVote{
				Voter:     user,
				Address:   account.Home,
				Elections: req.Elections,
			})
		return req, nil
	}

	_, voterPublicTree := git.Clone(ctx, git.Address(account.Home))
	mail.ReceiveSignedStageOnly(
		ctx,
		govTree,
		account.Home,
		voterPublicTree,
		proto.BallotTopic(ballotName),
		respond,
	)

	return git.Change[proto.FetchedVotes]{
		Result: fetched,
		Msg:    fmt.Sprintf("Fetched votes from user %v on ballot %v", user, ballotName),
	}
}

func LoadTally(
	ctx context.Context,
	communityTree *git.Tree,
	ballotName ns.NS,
) proto.Tally {
	openTallyNS := proto.OpenBallotNS(ballotName).Sub(proto.TallyFilebase)
	return git.FromFile[proto.Tally](ctx, communityTree, openTallyNS.Path())
}
