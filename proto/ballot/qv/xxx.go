package qv

import (
	"context"

	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Tally_XXX(
	ctx context.Context,
	govOwner id.OwnerCloned,
	ad *common.Advertisement,
	prior *common.Tally,
	fetched []common.FetchedVote, // newly fetched votes from participating users
) git.Change[form.Map, common.Tally] {

	// group old and new votes by user
	oldVotesByUser := map[member.User]common.FetchedVote{}
	if prior != nil { // load prior participant votes
		oldVotesByUser = GroupVotesByUser(ctx, prior.Votes)
	}
	newVotesByUser := GroupVotesByUser(ctx, fetched)

	// compute set of all users
	users := map[member.User]bool{}
	for u := range oldVotesByUser {
		users[u] = true
	}
	for u := range newVotesByUser {
		users[u] = true
	}

	// for every user, try adding the new votes to the old ones and charge the user
	scoredVotesByUser := map[member.User]ScoredVotes{}
	discardedVotes := common.DiscardedVotes{}
	XXX // keep track of charges in tally

	for u := range users {
		oldVotes, newVotes := oldVotesByUser[u], newVotesByUser[u]
		oldScore, newScore := AugmentUserVotes(ctx, oldVotes, newVotes)
		costDiff := newScore.Cost - oldScore.Cost
		// try charging the user for the new votes
		if err := ChargeUser(ctx, govOwner.Public.Tree(), u, costDiff); err != nil {
			scoredVotesByUser[u] = oldScore
			discardedVotes = append(discardedVotes, common.DiscardedVote{Vote: newVotes, Reason: err.Error()})
			XXX // keep track of charges in tally
		} else {
			scoredVotesByUser[u] = newScore
			XXX // keep track of charges in tally
		}
	}

	panic("XXX")
}

func AugmentUserVotes(
	ctx context.Context,
	oldVotes common.FetchedVote,
	newVotes common.FetchedVote,
) (oldScore, newScore ScoredVotes) {
	oldScore = ScoreAndCostOfUserVotes(ctx, oldVotes)
	augmentedVotes := common.FetchedVote{
		Voter:     mergeUser(ctx, oldVotes.Voter, newVotes.Voter),
		Address:   mergeAddress(oldVotes.Address, newVotes.Address),
		Elections: append(append(common.Elections{}, oldVotes.Elections...), newVotes.Elections...),
	}
	newScore = ScoreAndCostOfUserVotes(ctx, augmentedVotes)
	return
}

func mergeUser(ctx context.Context, u, v member.User) member.User {
	if u != "" && v != "" {
		must.Assertf(ctx, u == v, "users must be same")
	}
	if v == "" {
		return u
	}
	return v
}

func mergeAddress(a id.PublicAddress, b id.PublicAddress) id.PublicAddress {
	if b.IsEmpty() {
		return a
	}
	return b
}
