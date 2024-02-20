package proposal

import (
	"bytes"
	"context"
	"fmt"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/pmp_1"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/waimea"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/gov4git/v2/proto/notice"
)

func cancelNotice(
	ctx context.Context,
	prop motionproto.Motion,
	propState *waimea.ProposalState,
	outcome ballotproto.Outcome,

) notice.Notices {

	var w bytes.Buffer

	fmt.Fprintf(&w, "This unmerged PR, managed as Gov4Git proposal `%v`, has been cancelled 🌂\n\n", prop.ID)

	fmt.Fprintf(&w, "The PR __approval score__ was `%0.6f`.\n\n", propState.ApprovalScore)

	fmt.Fprintf(&w, "##\n\n") // refunds

	refunds := ballotproto.FlattenRefunds(outcome.Refunded)
	if len(refunds) > 0 {
		fmt.Fprintf(&w, "Refunds were issued:\n")
		for _, refund := range refunds {
			fmt.Fprintf(&w, "- Reviewer @%v was refunded `%0.6f` credits\n", refund.User, refund.Amount.Quantity)
		}
		fmt.Fprintln(&w, "")
	}

	fmt.Fprintf(&w, "##\n\n") // tally

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

	fmt.Fprintf(&w, "The PR __approval score__ was `%0.6f`.\n\n", r.ApprovalScore)

	if r.Accepted {
		if len(r.Resolved) > 0 {
			fmt.Fprintf(&w, "Issues claimed and __resolved__ by this PR were:\n")
			for _, con := range r.Resolved {
				fmt.Fprintf(&w, "- [Issue #%v](%v)\n", con.ID, con.TrackerURL)
			}
			fmt.Fprintln(&w, "")
		} else {
			fmt.Fprintf(&w, "No issues were claimed by this PR.\n\n")
		}
	}

	// fmt.Fprintf(&w, "##\n\n") // author awards

	if r.CostOfPriority > 0 {
		fmt.Fprintf(&w, "The __cost of priority__ of issues claimed by this PR was `%0.6f`.\n\n", r.CostOfPriority)
	}

	if r.RealizedBounty > 0 {
		fmt.Fprintf(&w, "The __realized bounty__ for the author of this PR was `%0.6f`, comprising:"+
			"- A __priority bounty__ of `%0.6f`, and\n"+
			"- A __review bounty__ of `%0.6f`.\n\n",
			r.RealizedBounty, r.ProjectedPriorityBounty, r.ProjectedReviewBounty)
	}

	if r.BountyDonation > 0 {
		fmt.Fprintf(&w, "A __donation__ of `%0.6f` credits from the bounty funds was made to the penny jar.\n\n", r.BountyDonation)
	}

	// fmt.Fprintf(&w, "##\n\n") // reviewer awards

	if len(r.Rewarded) > 0 {
		fmt.Fprintf(&w, "PR __reviewers__ were rewarded as follows:\n")
		for _, reward := range r.Rewarded {
			fmt.Fprintf(&w, "- Reviewer @%v was rewarded `%0.6f` credits\n", reward.To, reward.Amount.Quantity)
		}
		fmt.Fprintln(&w, "")
	}

	if r.CostOfReview > 0 {
		fmt.Fprintf(&w, "The __cost of review__ of this PR was `%0.6f`.\n\n", r.CostOfReview)
	}

	if r.RewardDonation > 0 {
		fmt.Fprintf(&w, "A donation of `%0.6f` credits from the __cost of review__ was made to the penny jar.\n\n", r.RewardDonation)
	}

	// fmt.Fprintf(&w, "##\n\n") // tally

	scoresByReviewer := r.ApprovalPollOutcome.ScoresByUser
	if len(scoresByReviewer) > 0 {
		fmt.Fprintf(&w, "PR approval tally by reviewer was:\n")
		for user, ss := range scoresByReviewer {
			fmt.Fprintf(&w, "- Reviewer @%v contributed `%0.6f` votes\n", user, ss[pmp_1.ProposalBallotChoice].Vote())
		}
	}

	return notice.NewNotice(ctx, w.String())
}
