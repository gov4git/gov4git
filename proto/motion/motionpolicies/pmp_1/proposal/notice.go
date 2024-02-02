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

	fmt.Fprintln(&w, "<hr/>") // refunds

	fmt.Fprintf(&w, "Refunds issued:\n")
	for _, refund := range ballotproto.FlattenRefunds(outcome.Refunded) {
		fmt.Fprintf(&w, "- Reviewer @%v was refunded `%0.6f` credits\n", refund.User, refund.Amount.Quantity)
	}
	fmt.Fprintln(&w, "")

	fmt.Fprintln(&w, "<hr/>") // tally

	fmt.Fprintf(&w, "PR approval tally by reviewer was:\n")
	for user, ss := range outcome.ScoresByUser {
		fmt.Fprintf(&w, "- Reviewer @%v contributed `%0.6f` votes\n", user, ss[pmp_1.ProposalBallotChoice].Vote())
	}

	return notice.NewNotice(ctx, w.String())
}

func closeNotice(
	ctx context.Context,
	prop motionproto.Motion,
	r *CloseReport,

) notice.Notices {

	var w bytes.Buffer

	if r.Accepted {
		fmt.Fprintf(&w, "This PR, managed as Gov4Git proposal `%v`, has been accepted 🎉\n\n", prop.ID)
	} else {
		fmt.Fprintf(&w, "This PR, managed as Gov4Git proposal `%v`, has been rejected 🌂\n\n", prop.ID)
	}

	if r.AgainstPopular {
		if r.Accepted {
			fmt.Fprintf(&w, "⚠️ Note that the PR was accepted against the popular vote.\n\n")
		} else {
			fmt.Fprintf(&w, "⚠️ Note that the PR was rejected against the popular vote.\n\n")
		}
	}

	fmt.Fprintf(&w, "The PR approval tally was `%0.6f`.\n\n", r.ApprovalPollOutcome.Scores[pmp_1.ProposalBallotChoice])

	if r.Accepted {
		if len(r.Resolved) > 0 {
			fmt.Fprintf(&w, "Resolved issues:\n")
			for _, con := range r.Resolved {
				fmt.Fprintf(&w, "- [Issue #%v](%v)\n", con.ID, con.TrackerURL)
			}
			fmt.Fprintln(&w, "")
		} else {
			fmt.Fprintf(&w, "No issues were claimed by this PR.\n\n")
		}
	}

	fmt.Fprintln(&w, "<hr/>") // author awards

	if r.CostOfPriority > 0 {
		fmt.Fprintf(&w, "The _cost of priority_ of issues claimed by this PR (the cost of prioritization) was `%0.6f`.\n\n", r.CostOfPriority)
	}

	if r.ProjectedBounty > 0 {
		fmt.Fprintf(&w, "The projected bounty, after matching, for the author of this PR was `%0.6f`.\n\n", r.ProjectedBounty)
	}

	if r.RealizedBounty > 0 {
		fmt.Fprintf(&w, "The realized bounty fot the author of this PR was `%0.6f`.\n\n", r.RealizedBounty)
	}

	if r.BountyDonation > 0 {
		fmt.Fprintf(&w, "A donation of `%0.6f` credits from the _cost of priority_ was made to the matching fund.\n\n", r.BountyDonation)
	}

	fmt.Fprintln(&w, "<hr/>") // reviewer awards

	if len(r.Rewarded) > 0 {
		fmt.Fprintf(&w, "PR reviewers were rewarded:\n")
		for _, reward := range r.Rewarded {
			fmt.Fprintf(&w, "- Reviewer @%v was rewarded `%0.6f` credits\n", reward.To, reward.Amount.Quantity)
		}
		fmt.Fprintln(&w, "")
	}

	if r.CostOfReview > 0 {
		fmt.Fprintf(&w, "The _cost of review_ of this PR was `%0.6f`.\n\n", r.CostOfReview)
	}

	if r.RewardDonation > 0 {
		fmt.Fprintf(&w, "A donation of `%0.6f` credits from the _cost of review_ was made to the matching fund.\n\n", r.RewardDonation)
	}

	fmt.Fprintln(&w, "<hr/>") // tally

	scoresByReviewer := r.ApprovalPollOutcome.ScoresByUser
	if len(scoresByReviewer) > 0 {
		fmt.Fprintf(&w, "PR approval tally by reviewer was:\n")
		for user, ss := range scoresByReviewer {
			fmt.Fprintf(&w, "- Reviewer @%v contributed `%0.6f` votes\n", user, ss[pmp_1.ProposalBallotChoice].Vote())
		}
	}

	return notice.NewNotice(ctx, w.String())
}
