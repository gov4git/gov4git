package concern

import (
	"bytes"
	"context"
	"fmt"

	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/motion/motionpolicies/waimea"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/gov4git/v2/proto/notice"
)

func cancelNotice(
	ctx context.Context,
	con motionproto.Motion,
	conState *waimea.ConcernState,
	outcome ballotproto.Outcome,

) notice.Notices {

	var w bytes.Buffer

	fmt.Fprintf(&w, "This issue, managed as Gov4Git proposal `%v`, has been cancelled 🌂\n\n", con.ID)

	fmt.Fprintf(&w, "The __priority score__ of the issue was `%0.6f`.\n", conState.PriorityScore)
	fmt.Fprintf(&w, "The __cost of priority__ of the issue was `%0.6f`.\n\n", conState.CostOfPriority)

	// refunded
	refunds := ballotproto.FlattenRefunds(outcome.Refunded)
	if len(refunds) > 0 {
		fmt.Fprintf(&w, "Refunds issued:\n")
		for _, refund := range refunds {
			fmt.Fprintf(&w, "- User @%v was refunded `%0.6f` credits\n", refund.User, refund.Amount.Quantity)
		}
		fmt.Fprintln(&w, "")
	}

	// tally by user
	if len(outcome.ScoresByUser) > 0 {
		fmt.Fprintf(&w, "Tally breakdown by user:\n")
		for user, ss := range outcome.ScoresByUser {
			fmt.Fprintf(&w, "- User @%v contributed `%0.6f` votes\n", user, ss[waimea.ConcernBallotChoice].Vote())
		}
	}

	return notice.NewNotice(ctx, w.String())
}

func closeNotice(
	ctx context.Context,
	con motionproto.Motion,
	conState *waimea.ConcernState,
	outcome ballotproto.Outcome,
	prop motionproto.Motion,

) notice.Notices {

	var w bytes.Buffer

	fmt.Fprintf(&w, "This issue, managed as Gov4Git concern `%v`, has been closed 🎉\n\n", con.ID)

	fmt.Fprintf(&w, "The __priority score__ of the issue was `%0.6f`.\n\n", conState.PriorityScore)
	fmt.Fprintf(&w, "The __cost of priority__ of the issue was `%0.6f`.\n\n", conState.CostOfPriority)

	// resolved by PR
	fmt.Fprintf(&w, "Ths issue was resolved by [PR #%v](%v):\n\n", prop.ID, prop.TrackerURL)

	// refunded
	refunds := ballotproto.FlattenRefunds(outcome.Refunded)
	if len(refunds) > 0 {
		fmt.Fprintf(&w, "Refunds issued:\n")
		for _, refund := range refunds {
			fmt.Fprintf(&w, "- User @%v was refunded `%0.6f` credits\n", refund.User, refund.Amount.Quantity)
		}
		fmt.Fprintln(&w, "")
	}

	// tally by user
	if len(outcome.ScoresByUser) > 0 {
		fmt.Fprintf(&w, "Tally breakdown by user:\n")
		for user, ss := range outcome.ScoresByUser {
			fmt.Fprintf(&w, "- User @%v contributed `%0.6f` votes\n", user, ss[waimea.ConcernBallotChoice].Vote())
		}
	}

	return notice.NewNotice(ctx, w.String())
}
