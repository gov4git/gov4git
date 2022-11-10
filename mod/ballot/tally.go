package ballot

import (
	"context"

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

	// process votes from users
	for i, account := range accounts {
		processUserVotes[S](ctx, govAddr, govRepo, govTree, ballotName, users[i], account)
	}

	XXX
}

func processUserVotes[S Strategy](
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	govRepo id.OwnerRepo,
	govTree id.OwnerTree,
	ballotName ns.NS,
	user member.User,
	account member.Account,
) {

	respond := func(ctx context.Context, req VoteEnvelope, _ id.SignedPlaintext) (resp VoteEnvelope, err error) {
		XXX
	}

	_, voterPublicTree := git.CloneBranchTree(ctx, git.Address(account.Home))

	chg := mail.ReceiveSigned[VoteEnvelope, VoteEnvelope](
		ctx,
		govTree,
		account.Home,
		voterPublicTree,
		BallotTopic[S](ballotName),
		respond,
	)

}
