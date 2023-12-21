package sv

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/common"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func (qv SV) Tally(
	ctx context.Context,
	owner gov.OwnerCloned,
	ad *common.Advertisement,
	prior *common.Tally,
	fetched map[member.User]common.Elections, // newly fetched votes from participating users
) git.Change[form.Map, common.Tally] {

	return qv.tally(ctx, owner.PublicClone(), ad, prior, fetched, false)
}

func (qv SV) tally(
	ctx context.Context,
	govCloned gov.Cloned, // clone of the public gov repo
	ad *common.Advertisement,
	prior *common.Tally,
	fetched map[member.User]common.Elections, // newly fetched votes from participating users
	strict bool, // fail if any voter has insufficient funds

) git.Change[form.Map, common.Tally] {

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
	acceptedVotes := map[member.User]common.AcceptedElections{}
	rejectedVotes := map[member.User]common.RejectedElections{}
	charges := map[member.User]float64{}
	votesByUser := map[member.User]map[string]common.StrengthAndScore{}

	for u := range users {
		oldVotes, newVotes := oldVotesByUser[u], newVotesByUser[u]
		oldScore, augmentedScore := augmentAndScoreUserVotes(ctx, qv.GetScorer(), oldVotes, newVotes)
		costDiff := augmentedScore.Cost - oldScore.Cost

		// try charging the user for the new votes
		err := chargeUser(ctx, govCloned, ad.Name, u, costDiff, fmt.Sprintf("vote charge for ballot %v", ad.Name))
		if strict {
			must.NoError(ctx, err)
		}
		if err != nil {
			acceptedVotes[u] = oldVotes
			rejectedVotes[u] = append(prior.RejectedVotes[u], rejectVotes(newVotes, err)...)
			charges[u] = prior.Charges[u]
			votesByUser[u] = oldScore.Score
		} else {
			acceptedVotes[u] = augmentedScore.Votes
			rejectedVotes[u] = prior.RejectedVotes[u]
			charges[u] = prior.Charges[u] + costDiff
			votesByUser[u] = augmentedScore.Score

			// metrics
			history.Log_StageOnly(
				ctx,
				govCloned,
				&history.Event{
					Vote: &history.VoteEvent{
						By:      u.HistoryUser(),
						Context: history.VoteContext(""), // TODO: plumb concern/proposal context here
						Receipts: history.OneReceipt(
							u.HistoryAccountID(),
							history.ReceiptTypeCharge,
							account.H(account.PluralAsset, costDiff).HistoryHolding(),
						),
					},
				},
			)
		}
	}

	tally := common.Tally{
		Ad:            *ad,
		Scores:        totalScore(ad.Choices, votesByUser),
		ScoresByUser:  votesByUser,
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

func chargeUser(
	ctx context.Context,
	govCloned gov.Cloned,
	ballotName common.BallotName,
	user member.User,
	charge float64,
	note string,
) error {

	return account.TryTransfer_StageOnly(
		ctx,
		govCloned,
		member.UserAccountID(user),
		common.BallotEscrowAccountID(ballotName),
		account.H(account.PluralAsset, charge),
		note,
	)
}
