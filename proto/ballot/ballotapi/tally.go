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
	addr gov.OwnerAddress,
	id ballotproto.BallotID,
	maxPar int,

) git.Change[form.Map, ballotproto.Tally] {

	cloned := gov.CloneOwner(ctx, addr)
	chg, changed := Tally_StageOnly(ctx, addr, cloned, id, maxPar)
	if !changed {
		return chg
	}
	proto.Commit(ctx, cloned.Public.Tree(), chg)
	cloned.Public.Push(ctx)
	return chg
}

func Tally_StageOnly(
	ctx context.Context,
	addr gov.OwnerAddress,
	cloned gov.OwnerCloned,
	id ballotproto.BallotID,
	maxPar int,

) (git.Change[form.Map, ballotproto.Tally], bool) {

	t := cloned.Public.Tree()
	ad, _ := ballotio.LoadPolicy(ctx, t, id)

	pv := loadParticipatingVoters(ctx, cloned.PublicClone(), ad)
	votersCloned := clonePar(ctx, pv.VoterAccounts, maxPar)
	pv.attachVoterClones(ctx, votersCloned)

	return tallyVotersCloned_StageOnly(ctx, cloned, id, pv.VoterAccounts, pv.VoterClones)
}

func tallyVotersCloned_StageOnly(
	ctx context.Context,
	cloned gov.OwnerCloned,
	id ballotproto.BallotID,
	voterAccounts map[member.User]member.UserProfile,
	votersCloned map[member.User]git.Cloned,

) (git.Change[form.Map, ballotproto.Tally], bool) {

	t := cloned.Public.Tree()
	ad, strat := ballotio.LoadPolicy(ctx, t, id)
	must.Assertf(ctx, !ad.Closed, "ballot is closed")

	// read current tally
	currentTally := loadTally_Local(ctx, t, id)

	var fetchedVotes FetchedVotes
	var fetchVoteChanges []git.Change[form.Map, FetchedVotes]
	for user, account := range voterAccounts {
		chg := fetchVotesCloned(ctx, cloned, id, user, account, votersCloned[user])
		fetchVoteChanges = append(fetchVoteChanges, chg)
		fetchedVotes = append(fetchedVotes, chg.Result...)
	}

	// if no votes are received, no change in tally occurs
	if len(fetchedVotes) == 0 {
		return git.NewChange(
			"No new votes",
			"ballot_tally",
			form.Map{"id": id},
			currentTally,
			nil,
		), false
	}

	// if the ballot is frozen, consume and reject pending votes
	if ad.Frozen {
		rejectFetchedVotes(fetchedVotes, currentTally.RejectedVotes)

		// write updated tally
		git.ToFileStage(ctx, t, id.TallyNS(), currentTally)

		return git.NewChange(
			"Ballot is frozen, discarding pending votes",
			"ballot_tally",
			form.Map{"id": id},
			currentTally,
			form.ToForms(fetchVoteChanges),
		), true
	}

	updatedTally := strat.Tally(ctx, cloned, &ad, &currentTally, fetchedVotesToElections(fetchedVotes)).Result

	// write updated tally
	git.ToFileStage(ctx, t, id.TallyNS(), updatedTally)

	// log
	trace.Log_StageOnly(ctx, cloned.PublicClone(), &trace.Event{
		Op:     "ballot_tally",
		Args:   trace.M{"ballot": id},
		Result: trace.M{"ad": ad, "tally": updatedTally},
	})

	return git.NewChange(
		fmt.Sprintf("Tally votes on ballot %v", id),
		"ballot_tally",
		form.Map{"id": id},
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
	t *git.Tree,
	id ballotproto.BallotID,

) ballotproto.Tally {

	return git.FromFile[ballotproto.Tally](ctx, t, id.TallyNS())
}
