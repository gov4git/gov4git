package ballot

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/lib/ns"
	"github.com/gov4git/gov4git/mod/gov"
	"github.com/gov4git/gov4git/mod/id"
	"github.com/gov4git/gov4git/mod/mail"
	"github.com/gov4git/gov4git/mod/member"
)

func Tally[S Strategy](
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	ballotName ns.NS,
) git.Change[TallyForm] {

	govRepo, govTree := id.CloneOwner(ctx, id.OwnerAddress(govAddr))
	chg := TallyStageOnly[S](ctx, govAddr, govRepo, govTree, ballotName)
	git.Commit(ctx, git.Worktree(ctx, govRepo.Public), chg.Msg)
	git.Push(ctx, govRepo.Public)
	return chg
}

func TallyStageOnly[S Strategy](
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	govRepo id.OwnerRepo,
	govTree id.OwnerTree,
	ballotName ns.NS,
) git.Change[TallyForm] {

	communityTree := govTree.Public

	// read ad
	openAdNS := OpenBallotNS[S](ballotName).Sub(adFilebase)
	ad := git.FromFile[AdForm](ctx, communityTree, openAdNS.Path())

	// list participating users
	users := member.ListGroupUsers(ctx, communityTree, ad.Participants)

	// get user accounts
	accounts := make([]member.Account, len(users))
	for i, user := range users {
		accounts[i] = member.GetUser(ctx, communityTree, user)
	}

	// fetch votes from users
	var fetchedVotes []FetchedVote
	for i, account := range accounts {
		fetchedVotes = append(fetchedVotes,
			fetchVotes[S](ctx, govAddr, govRepo, govTree, ballotName, users[i], account).Result...)
	}

	// read current tally
	openTallyNS := OpenBallotNS[S](ballotName).Sub(tallyFilebase)
	var currentTally *TallyForm
	if tryCurrentTally, err := git.TryFromFile[TallyForm](ctx, communityTree, openTallyNS.Path()); err == nil {
		currentTally = &tryCurrentTally
	}

	var s S
	updatedTally := s.Tally(ctx, govRepo, govTree, &ad, currentTally, fetchedVotes).Result

	// write updated tally
	git.ToFile(ctx, communityTree, openTallyNS.Path(), updatedTally)

	return git.Change[TallyForm]{
		Result: updatedTally,
		Msg:    fmt.Sprintf("Tally votes on ballot %v", ballotName),
	}
}

func fetchVotes[S Strategy](
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	govRepo id.OwnerRepo,
	govTree id.OwnerTree,
	ballotName ns.NS,
	user member.User,
	account member.Account,
) git.Change[[]FetchedVote] {

	fetched := []FetchedVote{}
	respond := func(ctx context.Context, req VoteEnvelope, _ id.SignedPlaintext) (resp VoteEnvelope, err error) {

		if !req.VerifyConsistency() {
			return VoteEnvelope{}, fmt.Errorf("vote envelope is not valid")
		}
		fetched = append(fetched, FetchedVote{Envelope: req, User: user})
		return req, nil
	}

	_, voterPublicTree := git.CloneBranchTree(ctx, git.Address(account.Home))
	mail.ReceiveSigned[VoteEnvelope, VoteEnvelope](
		ctx,
		govTree,
		account.Home,
		voterPublicTree,
		BallotTopic[S](ballotName),
		respond,
	)

	return git.Change[[]FetchedVote]{
		Result: fetched,
		Msg:    fmt.Sprintf("Fetched votes from user %v on ballot %v", user, ballotName),
	}
}
