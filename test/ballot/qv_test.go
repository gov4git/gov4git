package ballot

import (
	"fmt"
	"testing"

	"github.com/gov4git/gov4git/mod/balance"
	"github.com/gov4git/gov4git/mod/ballot/core"
	"github.com/gov4git/gov4git/mod/ballot/proto"
	"github.com/gov4git/gov4git/mod/ballot/qv"
	"github.com/gov4git/gov4git/mod/member"
	"github.com/gov4git/gov4git/test"
	"github.com/gov4git/lib4git/ns"
	"github.com/gov4git/lib4git/testutil"
)

func TestQV(t *testing.T) {
	ctx := testutil.NewCtx()
	cty := test.NewTestCommunity(t, ctx, 2)

	ballotName := ns.NS("a/b/c")
	choices := []string{"x", "y", "z"}

	// init voter credits
	balance.Set(ctx, cty.Community(), cty.MemberUser(0), qv.VotingCredits, 3.0)

	// open
	strat := qv.PriorityPoll{UseVotingCredits: true}
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

	// vote
	elections := proto.Elections{
		{
			VoteChoice:         choices[0],
			VoteStrengthChange: 2.0,
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
	if len(tallyChg.Result.Votes) != 1 {
		t.Errorf("expecting 1 vote, got %v", len(tallyChg.Result.Votes))
	}

	// close
	closeChg := core.Close(ctx, cty.Organizer(), ballotName, qv.SummaryAbandoned)
	fmt.Println("close: ", closeChg)

	// verify voter credits
	u0 := balance.Get(ctx, cty.Community(), cty.MemberUser(0), qv.VotingCredits)
	if u0 != 1.0 {
		t.Errorf("expecting 1, got %v", u0)
	}

	// testutil.Hang()
}
