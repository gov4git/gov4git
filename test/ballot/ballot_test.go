package core

import (
	"fmt"
	"testing"

	"github.com/gov4git/gov4git/mod/ballot/core"
	"github.com/gov4git/gov4git/mod/ballot/proto"
	"github.com/gov4git/gov4git/mod/ballot/qv"
	"github.com/gov4git/gov4git/mod/member"
	"github.com/gov4git/gov4git/test"
	"github.com/gov4git/lib4git/ns"
	"github.com/gov4git/lib4git/testutil"
)

func TestBallot(t *testing.T) {
	ctx := testutil.NewCtx()
	cty := test.NewTestCommunity(t, ctx, 2)

	ballotName := ns.NS("a/b/c")
	choices := []string{"x", "y", "z"}

	// open
	strat := qv.PriorityPoll{UseVotingCredits: false}
	openChg := core.Open(
		ctx,
		strat,
		cty.Community(),
		ballotName,
		"ballot_name",
		"ballot description",
		choices,
		member.Everybody,
	)
	fmt.Println("open: ", openChg)

	// list
	ads := core.ListOpen(ctx, cty.Community())
	if len(ads) != 1 {
		t.Errorf("expecting 1 ad, got %v", len(ads))
	}
	fmt.Println("ads: ", ads)

	// vote
	elections := []proto.Election{
		{
			VoteChoice:         choices[0],
			VoteStrengthChange: 1.0,
		},
	}
	voteChg := core.Vote(
		ctx,
		cty.MemberOwner(0),
		cty.Community(),
		ballotName,
		elections,
	)
	fmt.Println("vote: ", voteChg)

	// tally
	tallyChg := core.Tally(
		ctx,
		cty.Organizer(),
		ballotName,
	)
	fmt.Println("tally: ", tallyChg)
	if len(tallyChg.Result.FetchedVotes) != 1 {
		t.Errorf("expecting 1 vote, got %v", len(tallyChg.Result.FetchedVotes))
	}

	// close
	closeChg := core.Close(ctx, cty.Community(), ballotName)
	fmt.Println("close: ", closeChg)

	// testutil.Hang()
}
