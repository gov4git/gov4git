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
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/ns"
)

func Tally(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	ballotName ns.NS,
) git.Change[form.Map, common.Tally] {

	govOwner := id.CloneOwner(ctx, id.OwnerAddress(govAddr))
	chg, changed := TallyStageOnly(ctx, govAddr, govOwner, ballotName)
	if !changed {
		return chg
	}
	proto.Commit(ctx, govOwner.Public.Tree(), chg)
	govOwner.Public.Push(ctx)
	return chg
}

func TallyStageOnly(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	govOwner id.OwnerCloned,
	ballotName ns.NS,
) (git.Change[form.Map, common.Tally], bool) {

	communityTree := govOwner.Public.Tree()

	ad, strat := load.LoadStrategy(ctx, communityTree, ballotName, false)

	// if the ballot is frozen, don't tally
	if ad.Frozen {
		return git.NewChange(
			"Ballot is frozen",
			"ballot_tally",
			form.Map{"ballot_name": ballotName},
			common.Tally{},
			nil,
		), false
	}

	// list participating users
	users := member.ListGroupUsersLocal(ctx, communityTree, ad.Participants)

	// get user accounts
	accounts := make([]member.Account, len(users))
	for i, user := range users {
		accounts[i] = member.GetUserLocal(ctx, communityTree, user)
	}

	// fetch votes from users
	var fetchedVotes common.FetchedVotes
	var fetchVoteChanges []git.Change[form.Map, common.FetchedVotes]
	for i, account := range accounts {
		chg := fetchVotes(ctx, govAddr, govOwner, ballotName, users[i], account)
		fetchVoteChanges = append(fetchVoteChanges, chg)
		fetchedVotes = append(fetchedVotes, chg.Result...)
	}

	// read current tally
	var currentTally *common.Tally
	if tryCurrentTally, err := must.Try1(func() common.Tally { return LoadTally(ctx, communityTree, ballotName, false) }); err == nil {
		currentTally = &tryCurrentTally
	}

	// if no votes are received, no change in tally occurs
	if len(fetchedVotes) == 0 {
		if currentTally == nil {
			currentTally = &common.Tally{}
		}
		return git.NewChange(
			"No new votes",
			"ballot_tally",
			form.Map{"ballot_name": ballotName},
			*currentTally,
			nil,
		), false
	}

	updatedTally := strat.Tally(ctx, govOwner, &ad, currentTally, fetchedVotes).Result

	// write updated tally
	openTallyNS := common.OpenBallotNS(ballotName).Sub(common.TallyFilebase)
	git.ToFileStage(ctx, communityTree, openTallyNS.Path(), updatedTally)

	return git.NewChange(
		fmt.Sprintf("Tally votes on ballot %v", ballotName),
		"ballot_tally",
		form.Map{"ballot_name": ballotName},
		updatedTally,
		nil,
	), true
}

func fetchVotes(
	ctx context.Context,
	govAddr gov.OrganizerAddress,
	govOwner id.OwnerCloned,
	ballotName ns.NS,
	user member.User,
	account member.Account,
) git.Change[form.Map, common.FetchedVotes] {

	fetched := common.FetchedVotes{}
	respond := func(ctx context.Context, req common.VoteEnvelope, _ id.SignedPlaintext) (resp common.VoteEnvelope, err error) {

		if !req.VerifyConsistency() {
			return common.VoteEnvelope{}, fmt.Errorf("vote envelope is not valid")
		}
		fetched = append(fetched,
			common.FetchedVote{
				Voter:     user,
				Address:   account.PublicAddress,
				Elections: req.Elections,
			})
		return req, nil
	}

	voterPublicTree := git.CloneOne(ctx, git.Address(account.PublicAddress)).Tree()
	mail.ReceiveSignedStageOnly(
		ctx,
		govOwner,
		account.PublicAddress,
		voterPublicTree,
		common.BallotTopic(ballotName),
		respond,
	)

	return git.NewChange(
		fmt.Sprintf("Fetched votes from user %v on ballot %v", user, ballotName),
		"ballot_fetch_votes",
		form.Map{"ballot_name": ballotName, "user": user, "account": account},
		fetched,
		nil,
	)
}

func LoadTally(
	ctx context.Context,
	communityTree *git.Tree,
	ballotName ns.NS,
	closed bool,
) common.Tally {
	var tallyNS ns.NS
	if closed {
		tallyNS = common.ClosedBallotNS(ballotName).Sub(common.TallyFilebase)
	} else {
		tallyNS = common.OpenBallotNS(ballotName).Sub(common.TallyFilebase)
	}
	return git.FromFile[common.Tally](ctx, communityTree, tallyNS.Path())
}
