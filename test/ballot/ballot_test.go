package ballot

import (
	"fmt"
	"testing"

	"github.com/gov4git/gov4git/proto/ballot/ballot"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/qv"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/gov4git/test"
	"github.com/gov4git/lib4git/ns"
	"github.com/gov4git/lib4git/testutil"
)

func TestBallot(t *testing.T) {
	ctx := testutil.NewCtx()
	cty := test.NewTestCommunity(t, ctx, 2)

	ballotName := ns.NS{"a", "b", "c"}
	choices := []string{"x", "y", "z"}

	// open
	strat := qv.PriorityPoll{UseVotingCredits: false}
	openChg := ballot.Open(
		ctx,
		strat,
		cty.Gov(),
		ballotName,
		"ballot_name",
		"ballot description",
		choices,
		member.Everybody,
	)
	fmt.Println("open: ", openChg)

	// list
	ads := ballot.List(ctx, cty.Gov(), false)
	if len(ads) != 1 {
		t.Errorf("expecting 1 ad, got %v", len(ads))
	}
	fmt.Println("ads: ", ads)

	// vote
	elections := common.Elections{
		{
			VoteChoice:         choices[0],
			VoteStrengthChange: 1.0,
		},
	}
	voteChg := ballot.Vote(
		ctx,
		cty.MemberOwner(0),
		cty.Gov(),
		ballotName,
		elections,
	)
	fmt.Println("vote: ", voteChg)

	// tally
	tallyChg := ballot.Tally(
		ctx,
		cty.Organizer(),
		ballotName,
	)
	fmt.Println("tally: ", tallyChg)
	if len(tallyChg.Result.Votes) != 1 {
		t.Errorf("expecting 1 vote, got %v", len(tallyChg.Result.Votes))
	}

	// close
	closeChg := ballot.Close(ctx, cty.Organizer(), ballotName, common.Summary("ok"))
	fmt.Println("close: ", closeChg)

	// testutil.Hang()
}
