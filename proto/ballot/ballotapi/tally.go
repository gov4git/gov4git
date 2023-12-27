package ballotapi

import (
	"context"
	"fmt"
	"time"

	"github.com/gov4git/gov4git/v2/proto"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotio"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history/trace"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Tally(
	ctx context.Context,
	govAddr gov.OwnerAddress,
	ballotName ballotproto.BallotName,
	maxPar int,

) git.Change[form.Map, ballotproto.Tally] {

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
	govAddr gov.OwnerAddress,
	govOwner gov.OwnerCloned,
	ballotName ballotproto.BallotName,
	maxPar int,

) (git.Change[form.Map, ballotproto.Tally], bool) {

	communityTree := govOwner.Public.Tree()
	ad, _ := ballotio.LoadStrategy(ctx, communityTree, ballotName)

	pv := loadParticipatingVoters(ctx, govOwner.PublicClone(), ad)
	votersCloned := clonePar(ctx, pv.VoterAccounts, maxPar)
	pv.attachVoterClones(ctx, votersCloned)

	return tallyVotersCloned_StageOnly(ctx, govOwner, ballotName, pv.VoterAccounts, pv.VoterClones)
}

func tallyVotersCloned_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	ballotName ballotproto.BallotName,
	voterAccounts map[member.User]member.UserProfile,
	votersCloned map[member.User]git.Cloned,

) (git.Change[form.Map, ballotproto.Tally], bool) {

	communityTree := cloned.Public.Tree()
	ad, strat := ballotio.LoadStrategy(ctx, communityTree, ballotName)
	must.Assertf(ctx, !ad.Closed, "ballot is closed")

	// read current tally
	currentTally := loadTally_Local(ctx, communityTree, ballotName)

	var fetchedVotes FetchedVotes
	var fetchVoteChanges []git.Change[form.Map, FetchedVotes]
	for user, account := range voterAccounts {
		chg := fetchVotesCloned(ctx, cloned, ballotName, user, account, votersCloned[user])
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
		openTallyNS := ballotproto.BallotPath(ballotName).Append(ballotproto.TallyFilebase)
		git.ToFileStage(ctx, communityTree, openTallyNS, currentTally)

		return git.NewChange(
			"Ballot is frozen, discarding pending votes",
			"ballot_tally",
			form.Map{"ballot_name": ballotName},
			currentTally,
			form.ToForms(fetchVoteChanges),
		), true
	}

	updatedTally := strat.Tally(ctx, cloned, &ad, &currentTally, fetchedVotesToElections(fetchedVotes)).Result

	// write updated tally
	openTallyNS := ballotproto.BallotPath(ballotName).Append(ballotproto.TallyFilebase)
	git.ToFileStage(ctx, communityTree, openTallyNS, updatedTally)

	// log
	trace.Log_StageOnly(ctx, cloned.PublicClone(), &trace.Event{
		Op:     "ballot_tally",
		Args:   trace.M{"ballot": ballotName},
		Result: trace.M{"ad": ad, "tally": updatedTally},
	})

	return git.NewChange(
		fmt.Sprintf("Tally votes on ballot %v", ballotName),
		"ballot_tally",
		form.Map{"ballot_name": ballotName},
		updatedTally,
		form.ToForms(fetchVoteChanges),
	), true
}

func rejectFetchedVotes(fv FetchedVotes, rej map[member.User]ballotproto.RejectedElections) {
	for _, fv := range fv {
		for _, el := range fv.Elections {
			rej[fv.Voter] = append(
				rej[fv.Voter],
				ballotproto.RejectedElection{Time: time.Now(), Vote: el, Reason: "ballot is frozen"},
			)
		}
	}
}

func loadTally_Local(
	ctx context.Context,
	communityTree *git.Tree,
	ballotName ballotproto.BallotName,

) ballotproto.Tally {

	tallyNS := ballotproto.BallotPath(ballotName).Append(ballotproto.TallyFilebase)
	return git.FromFile[ballotproto.Tally](ctx, communityTree, tallyNS)
}
