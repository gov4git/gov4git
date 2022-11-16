package ballot

import (
	"fmt"
	"testing"

	"github.com/gov4git/gov4git/mod/ballot"
	"github.com/gov4git/gov4git/mod/member"
	"github.com/gov4git/gov4git/mod/qv"
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
	openChg := ballot.Open[qv.PriorityPoll](
		ctx,
		cty.Community(),
		ballotName,
		"ballot_name",
		"ballot description",
		choices,
		member.Everybody,
	)
	fmt.Println("open: ", openChg)

	// list
	ads := ballot.ListOpen[qv.PriorityPoll](ctx, cty.Community())
	if len(ads) != 1 {
		t.Errorf("expecting 1 ad, got %v", len(ads))
	}
	fmt.Println("ads: ", ads)

	// vote
	elections := []ballot.Election{
		{
			VoteChoice:         choices[0],
			VoteStrengthChange: 1.0,
		},
	}
	voteChg := ballot.Vote[qv.PriorityPoll](
		ctx,
		cty.MemberOwner(0),
		cty.Community(),
		ballotName,
		elections,
	)
	fmt.Println("vote: ", voteChg)

	// tally
	tallyChg := ballot.Tally[qv.PriorityPoll](
		ctx,
		cty.Organizer(),
		ballotName,
	)
	fmt.Println("tally: ", tallyChg)
	if len(tallyChg.Result.FetchedVotes) != 1 {
		t.Errorf("expecting 1 vote, got %v", len(tallyChg.Result.FetchedVotes))
	}

	// close
	closeChg := ballot.Close[qv.PriorityPoll](ctx, cty.Community(), ballotName)
	fmt.Println("close: ", closeChg)

	// testutil.Hang()
}
