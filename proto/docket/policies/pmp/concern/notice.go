package concern

import (
	"bytes"
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/docket/policies/pmp"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/notice"
)

func cancelNotice(ctx context.Context, motion schema.Motion, outcome common.Outcome) notice.Notices {

	var w bytes.Buffer

	fmt.Fprintf(&w, "This issue, managed as Gov4Git proposal `%v`, has been cancelled 🌂\n\n", motion.ID)

	fmt.Fprintf(&w, "The issue priority tally was %v.\n\n", outcome.Scores[pmp.ProposalBallotChoice])

	// refunded
	fmt.Fprintf(&w, "Refunds issued:\n")
	for _, refund := range common.FlattenRefunds(outcome.Refunded) {
		fmt.Fprintf(&w, "- User %v was refunded %v credits.\n", refund.User, refund.Amount.Quantity)
	}
	fmt.Fprintln(&w, "")

	// tally by user
	fmt.Fprintf(&w, "Tally breakdown by user:\n")
	for user, ss := range outcome.ScoresByUser {
		fmt.Fprintf(&w, "- User %v contributed %v votes.\n", user, ss[pmp.ProposalBallotChoice].Vote())
	}

	return notice.NewNotice(w.String())
}

func closeNotice(
	ctx context.Context,
	con schema.Motion,
	outcome common.Outcome,
	prop schema.Motion,

) notice.Notices {

	var w bytes.Buffer

	fmt.Fprintf(&w, "This issue, managed as Gov4Git proposal `%v`, has been closed 🎉\n\n", con.ID)

	fmt.Fprintf(&w, "The issue priority tally was %v.\n\n", outcome.Scores[pmp.ProposalBallotChoice])

	// resolved by PR
	fmt.Fprintf(&w, "Issued was resolved by [PR #%v](%v):\n\n", prop.ID, prop.TrackerURL)

	// tally by user
	fmt.Fprintf(&w, "Tally breakdown by user:\n")
	for user, ss := range outcome.ScoresByUser {
		fmt.Fprintf(&w, "- User %v contributed %v votes.\n", user, ss[pmp.ProposalBallotChoice].Vote())
	}

	return notice.NewNotice(w.String())
}
