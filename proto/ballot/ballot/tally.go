package ballot

import (
	"context"
	"fmt"
	"time"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/load"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Tally(
	ctx context.Context,
	govAddr gov.GovOwnerAddress,
	ballotName common.BallotName,
	maxPar int,
) git.Change[form.Map, common.Tally] {

	govOwner := gov.CloneOwner(ctx, govAddr)
	chg, changed := Tally_StageOnly(ctx, govAddr, govOwner, ballotName, maxPar)
	if !changed {
		return chg
	}
	proto.Commit(ctx, govOwner.Public.Tree(), chg)
	govOwner.Public.Push(ctx)
	return chg
}

func Tally_StageOnly(
	ctx context.Context,
	govAddr gov.GovOwnerAddress,
	govOwner gov.GovOwnerCloned,
	ballotName common.BallotName,
	maxPar int,
) (git.Change[form.Map, common.Tally], bool) {

	communityTree := govOwner.Public.Tree()
	ad, _ := load.LoadStrategy(ctx, communityTree, ballotName)

	pv := loadParticipatingVoters(ctx, communityTree, ad)
	votersCloned := clonePar(ctx, pv.VoterAccounts, maxPar)
	pv.attachVoterClones(ctx, votersCloned)

	return tallyVotersCloned_StageOnly(ctx, govAddr, govOwner, ballotName, pv.VoterAccounts, pv.VoterClones)
}

func tallyVotersCloned_StageOnly(
	ctx context.Context,
	govAddr gov.GovOwnerAddress,
	govOwner gov.GovOwnerCloned,
	ballotName common.BallotName,
	voterAccounts map[member.User]member.Account,
	votersCloned map[member.User]git.Cloned,
) (git.Change[form.Map, common.Tally], bool) {

	communityTree := govOwner.Public.Tree()
	ad, strat := load.LoadStrategy(ctx, communityTree, ballotName)
	must.Assertf(ctx, !ad.Closed, "ballot is closed")

	// read current tally
	currentTally := LoadTally(ctx, communityTree, ballotName)

	var fetchedVotes FetchedVotes
	var fetchVoteChanges []git.Change[form.Map, FetchedVotes]
	for user, account := range voterAccounts {
		chg := fetchVotesCloned(ctx, govAddr, govOwner, ballotName, user, account, votersCloned[user])
		fetchVoteChanges = append(fetchVoteChanges, chg)
		fetchedVotes = append(fetchedVotes, chg.Result...)
	}

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

	// if the ballot is frozen, consume and reject pending votes
	if ad.Frozen {
		rejectFetchedVotes(fetchedVotes, currentTally.RejectedVotes)

		// write updated tally
		openTallyNS := common.BallotPath(ballotName).Append(common.TallyFilebase)
		git.ToFileStage(ctx, communityTree, openTallyNS, currentTally)

		return git.NewChange(
			"Ballot is frozen, discarding pending votes",
			"ballot_tally",
			form.Map{"ballot_name": ballotName},
			currentTally,
			form.ToForms(fetchVoteChanges),
		), true
	}

	updatedTally := strat.Tally(ctx, govOwner, &ad, &currentTally, fetchedVotesToElections(fetchedVotes)).Result

	// write updated tally
	openTallyNS := common.BallotPath(ballotName).Append(common.TallyFilebase)
	git.ToFileStage(ctx, communityTree, openTallyNS, updatedTally)

	return git.NewChange(
		fmt.Sprintf("Tally votes on ballot %v", ballotName),
		"ballot_tally",
		form.Map{"ballot_name": ballotName},
		updatedTally,
		form.ToForms(fetchVoteChanges),
	), true
}

func rejectFetchedVotes(fv FetchedVotes, rej map[member.User]common.RejectedElections) {
	for _, fv := range fv {
		for _, el := range fv.Elections {
			rej[fv.Voter] = append(
				rej[fv.Voter],
				common.RejectedElection{Time: time.Now(), Vote: el, Reason: "ballot is frozen"},
			)
		}
	}
}

func LoadTally(ctx context.Context, communityTree *git.Tree, ballotName common.BallotName) common.Tally {
	tallyNS := common.BallotPath(ballotName).Append(common.TallyFilebase)
	return git.FromFile[common.Tally](ctx, communityTree, tallyNS)
}
