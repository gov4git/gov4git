package sv

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/history/metric"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func (qv SV) Tally(
	ctx context.Context,
	cloned gov.Cloned,
	ad *ballotproto.Ad,
	prior *ballotproto.Tally,
	fetched map[member.User]ballotproto.Elections, // newly fetched votes from participating users

) git.Change[form.Map, ballotproto.Tally] {

	return qv.tally(ctx, cloned, ad, prior, fetched, false)
}

func (qv SV) tally(
	ctx context.Context,
	cloned gov.Cloned, // clone of the public gov repo
	ad *ballotproto.Ad,
	prior *ballotproto.Tally,
	fetched map[member.User]ballotproto.Elections, // newly fetched votes from participating users
	strict bool, // fail if any voter has insufficient funds

) git.Change[form.Map, ballotproto.Tally] {

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
	acceptedVotes := map[member.User]ballotproto.AcceptedElections{}
	rejectedVotes := map[member.User]ballotproto.RejectedElections{}
	charges := map[member.User]float64{}
	votesByUser := map[member.User]map[string]ballotproto.StrengthAndScore{}

	for u := range users {
		oldVotes, newVotes := oldVotesByUser[u], newVotesByUser[u]
		oldScore, augmentedScore := augmentAndScoreUserVotes(ctx, cloned, ad, qv.GetScorer(ctx), oldVotes, newVotes)
		costDiff := augmentedScore.Cost - oldScore.Cost

		// try charging the user for the new votes
		err := chargeUser(ctx, cloned, ad.ID, u, costDiff, fmt.Sprintf("vote charge for ballot %v", ad.ID))
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
			metric.Log_StageOnly(
				ctx,
				cloned,
				&metric.Event{
					Vote: &metric.VoteEvent{
						By:           u.MetricUser(),
						Purpose:      ad.Purpose.MetricVotePurpose(),
						MotionPolicy: metric.MotionPolicy(ad.MotionPolicy),
						BallotPolicy: metric.BallotPolicy(ad.Policy),
						Receipts: metric.OneReceipt(
							u.MetricAccountID(),
							metric.ReceiptTypeCharge,
							account.H(account.PluralAsset, costDiff).MetricHolding(),
						),
					},
				},
			)
		}
	}

	tally := ballotproto.Tally{
		Ad:            *ad,
		Scores:        totalScore(ad.Choices, votesByUser),
		ScoresByUser:  votesByUser,
		AcceptedVotes: acceptedVotes,
		RejectedVotes: rejectedVotes,
		Charges:       charges,
	}
	return git.NewChange(
		fmt.Sprintf("Tallied QV scores for ballot %v", ad.ID),
		"ballot_qv_tally",
		form.Map{"id": ad.ID},
		tally,
		nil,
	)
}

func chargeUser(
	ctx context.Context,
	govCloned gov.Cloned,
	ballotName ballotproto.BallotID,
	user member.User,
	charge float64,
	note string,
) error {

	return account.TryTransfer_StageOnly(
		ctx,
		govCloned,
		member.UserAccountID(user),
		ballotproto.BallotEscrowAccountID(ballotName),
		account.H(account.PluralAsset, charge),
		note,
	)
}
