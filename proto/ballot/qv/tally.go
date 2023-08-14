package qv

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto/balance"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func Tally_XXX(
	ctx context.Context,
	owner id.OwnerCloned,
	ad *common.Advertisement,
	prior *common.Tally2,
	fetched map[member.User]common.Elections, // newly fetched votes from participating users
) git.Change[form.Map, common.Tally2] {

	oldVotesByUser, newVotesByUser := prior.AcceptedVotes, fetched

	// compute set of all users
	users := map[member.User]bool{}
	for u := range oldVotesByUser {
		users[u] = true
	}
	for u := range newVotesByUser {
		users[u] = true
	}

	// for every user, try augmenting the old votes with the new ones and charging the user
	acceptedVotes := map[member.User]common.Elections{}
	rejectedVotes := map[member.User]common.RejectedElections{}
	charges := map[member.User]float64{}
	votesByUser := map[member.User]map[string]common.StrengthAndScore{}

	for u := range users {
		oldVotes, newVotes := oldVotesByUser[u], newVotesByUser[u]
		oldScore, augmentedScore := augmentAndScoreUserVotes(ctx, oldVotes, newVotes)
		costDiff := augmentedScore.Cost - oldScore.Cost

		// try charging the user for the new votes
		if err := ChargeUser(ctx, owner, u, costDiff); err != nil {
			acceptedVotes[u] = oldVotes
			rejectedVotes[u] = append(prior.RejectedVotes[u], rejectVotes(newVotes, err)...)
			charges[u] = prior.Charges[u]
			votesByUser[u] = oldScore.Score
		} else {
			acceptedVotes[u] = augmentedScore.Votes
			rejectedVotes[u] = prior.RejectedVotes[u]
			charges[u] = prior.Charges[u] + costDiff
			votesByUser[u] = augmentedScore.Score
		}
	}

	tally := common.Tally2{
		Ad:            *ad,
		Scores:        totalScore(ad.Choices, votesByUser),
		VotesByUser:   votesByUser,
		AcceptedVotes: acceptedVotes,
		RejectedVotes: rejectedVotes,
		Charges:       charges,
	}
	return git.NewChange(
		fmt.Sprintf("Tallied QV scores for ballot %v", ad.Name),
		"ballot_qv_tally",
		form.Map{"ballot_name": ad.Name},
		tally,
		nil,
	)
}

func ChargeUser(ctx context.Context, owner id.OwnerCloned, user member.User, charge float64) error {
	return balance.TryChargeStageOnly(ctx, owner.Public.Tree(), user, VotingCredits, charge)
}
