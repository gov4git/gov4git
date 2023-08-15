package ballot

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/load"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
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
	var fetchedVotes FetchedVotes
	var fetchVoteChanges []git.Change[form.Map, FetchedVotes]
	for i, account := range accounts {
		chg := fetchVotes(ctx, govAddr, govOwner, ballotName, users[i], account)
		fetchVoteChanges = append(fetchVoteChanges, chg)
		fetchedVotes = append(fetchedVotes, chg.Result...)
	}

	// read current tally
	currentTally := LoadTally(ctx, communityTree, ballotName, false)

	// if no votes are received, no change in tally occurs
	if len(fetchedVotes) == 0 {
		return git.NewChange(
			"No new votes",
			"ballot_tally",
			form.Map{"ballot_name": ballotName},
			currentTally,
			nil,
		), false
	}

	updatedTally := strat.Tally(ctx, govOwner, &ad, &currentTally, fetchedVotesToElections(fetchedVotes)).Result

	// write updated tally
	openTallyNS := common.OpenBallotNS(ballotName).Sub(common.TallyFilebase)
	git.ToFileStage(ctx, communityTree, openTallyNS.Path(), updatedTally)

	return git.NewChange(
		fmt.Sprintf("Tally votes on ballot %v", ballotName),
		"ballot_tally",
		form.Map{"ballot_name": ballotName},
		updatedTally,
		form.ToForms(fetchVoteChanges),
	), true
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
