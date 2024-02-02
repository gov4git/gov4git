package proposal

import (
	"bytes"
	"context"
	"fmt"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_1"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/gov4git/v2/proto/notice"
)

func cancelNotice(
	ctx context.Context,
	motion motionproto.Motion,
	outcome ballotproto.Outcome,

) notice.Notices {

	var w bytes.Buffer

	fmt.Fprintf(&w, "This unmerged PR, managed as Gov4Git proposal `%v`, has been cancelled 🌂\n\n", motion.ID)

	fmt.Fprintf(&w, "The PR approval tally was `%0.6f`.\n\n", outcome.Scores[pmp_1.ProposalBallotChoice])

	// refunded
	fmt.Fprintf(&w, "Refunds issued:\n")
	for _, refund := range ballotproto.FlattenRefunds(outcome.Refunded) {
		fmt.Fprintf(&w, "- Reviewer @%v was refunded `%0.6f` credits\n", refund.User, refund.Amount.Quantity)
	}
	fmt.Fprintln(&w, "")

	// tally by user
	fmt.Fprintf(&w, "Tally breakdown by reviewer was:\n")
	for user, ss := range outcome.ScoresByUser {
		fmt.Fprintf(&w, "- Reviewer @%v contributed `%0.6f` votes\n", user, ss[pmp_1.ProposalBallotChoice].Vote())
	}

	return notice.NewNotice(ctx, w.String())
}

func closeNotice(
	ctx context.Context,
	prop motionproto.Motion,
	accepted bool,
	againstPopular bool,
	outcome ballotproto.Outcome,
	resolved motionproto.Motions,
	// reviewers
	rewards Rewards,
	rewardDonation float64,
	// author
	bounty float64,
	escrow float64,
	award float64,
	bountyDonation float64,

) notice.Notices {

	var w bytes.Buffer

	if accepted {
		fmt.Fprintf(&w, "This PR, managed as Gov4Git proposal `%v`, has been accepted 🎉\n\n", prop.ID)
	} else {
		fmt.Fprintf(&w, "This PR, managed as Gov4Git proposal `%v`, has been rejected 🌂\n\n", prop.ID)
	}

	if againstPopular {
		if accepted {
			fmt.Fprintf(&w, "⚠️ Note that the PR was accepted against the popular vote.\n\n")
		} else {
			fmt.Fprintf(&w, "⚠️ Note that the PR was rejected against the popular vote.\n\n")
		}
	}

	fmt.Fprintf(&w, "The PR approval tally was `%0.6f`.\n\n", outcome.Scores[pmp_1.ProposalBallotChoice])

	// bounty
	XXX
	if bountyDonation > 0.0 {
		fmt.Fprintf(&w, "Bounty of `%0.6f` credits was donated to the matching fund.\n\n", bounty.Quantity)
	} else {
		fmt.Fprintf(&w, "Bounty of `%0.6f` credits was awarded to @%v.\n\n", bounty.Quantity, prop.Author)
	}

	// resolved issues
	if accepted {
		if len(resolved) > 0 {
			fmt.Fprintf(&w, "Resolved issues:\n")
			for _, con := range resolved {
				fmt.Fprintf(&w, "- [Issue #%v](%v)\n", con.ID, con.TrackerURL)
			}
			fmt.Fprintln(&w, "")
		} else {
			fmt.Fprintf(&w, "No issues were claimed.\n\n")
		}
	}

	// rewarded reviewers
	if len(rewards) > 0 {
		fmt.Fprintf(&w, "Rewarded PR reviewers:\n")
		for _, reward := range rewards {
			fmt.Fprintf(&w, "- Reviewer @%v was awarded `%0.6f` credits\n", reward.To, reward.Amount.Quantity)
		}
		fmt.Fprintln(&w, "")
	} else {
		fmt.Fprintf(&w, "No reviewers were rewarded.\n\n")
	}

	if rewardDonation > 0 {
		fmt.Fprintf(&w, "Reviewers' donation of `%0.6f` credits was made to the matching fund.\n\n", rewardDonation)
	}

	// tally by user
	fmt.Fprintf(&w, "Tally breakdown by user:\n")
	for user, ss := range outcome.ScoresByUser {
		fmt.Fprintf(&w, "- Reviewer @%v contributed `%0.6f` votes\n", user, ss[pmp_1.ProposalBallotChoice].Vote())
	}

	return notice.NewNotice(ctx, w.String())
}
