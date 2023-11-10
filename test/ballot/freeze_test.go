package ballot

import (
	"fmt"
	"testing"

	"github.com/gov4git/gov4git/proto/balance"
	"github.com/gov4git/gov4git/proto/ballot/ballot"
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/ballot/qv"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/gov4git/runtime"
	"github.com/gov4git/gov4git/test"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/testutil"
)

func TestVoteFreezeVote(t *testing.T) {
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	ballotName := common.BallotName{"a", "b", "c"}
	choices := []string{"x", "y", "z"}

	// open
	strat := qv.QV{}
	openChg := ballot.Open(
		ctx,
		strat,
		cty.Gov(),
		ballotName,
		"ballot title",
		"ballot description",
		choices,
		member.Everybody,
	)
	fmt.Println("open: ", form.SprintJSON(openChg))

	// give voter credits
	balance.Set(ctx, cty.Gov(), cty.MemberUser(0), qv.VotingCredits, 1.0)

	// vote
	elections := common.Elections{
		common.NewElection(choices[0], 1.0),
	}
	voteChg := ballot.Vote(ctx, cty.MemberOwner(0), cty.Gov(), ballotName, elections)
	fmt.Println("vote: ", form.SprintJSON(voteChg))

	// freeze ballot
	freezeChg := ballot.Freeze(ctx, cty.Organizer(), ballotName)
	fmt.Println("freeze: ", form.SprintJSON(freezeChg))

	// try voting while frozen
	if must.Try(
		func() { ballot.Vote(ctx, cty.MemberOwner(0), cty.Gov(), ballotName, elections) },
	) == nil {
		t.Fatalf("voting on a frozen ballot should have failed")
	}

	// unfreeze ballot
	unfreezeChg := ballot.Unfreeze(ctx, cty.Organizer(), ballotName)
	fmt.Println("unfreeze: ", form.SprintJSON(unfreezeChg))

	// tally
	tallyChg := ballot.Tally(ctx, cty.Organizer(), ballotName, testMaxPar)
	fmt.Println("tally: ", form.SprintJSON(tallyChg))
	if tallyChg.Result.Scores[choices[0]] != 1.0 {
		t.Errorf("expecting %v, got %v", 1.0, tallyChg.Result.Scores[choices[0]])
	}

	// close
	closeChg := ballot.Close(ctx, cty.Organizer(), ballotName, false)
	fmt.Println("close: ", form.SprintJSON(closeChg))

	// testutil.Hang()
}

// TestVoteFreezeTally tests that votes made during a freeze are consumed and discarded by tallying.
func TestVoteFreezeTally(t *testing.T) {
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	ballotName := common.BallotName{"a", "b", "c"}
	choices := []string{"x", "y", "z"}

	// open
	strat := qv.QV{}
	openChg := ballot.Open(ctx, strat, cty.Gov(), ballotName, "ballot title", "ballot description", choices, member.Everybody)
	fmt.Println("open: ", form.SprintJSON(openChg))

	// give voter credits
	balance.Set(ctx, cty.Gov(), cty.MemberUser(0), qv.VotingCredits, 1.0)

	// vote
	elections := common.Elections{
		common.NewElection(choices[0], 1.0),
	}
	voteChg := ballot.Vote(ctx, cty.MemberOwner(0), cty.Gov(), ballotName, elections)
	fmt.Println("vote: ", form.SprintJSON(voteChg))

	// freeze
	freezeChg := ballot.Freeze(ctx, cty.Organizer(), ballotName)
	fmt.Println("freeze: ", form.SprintJSON(freezeChg))

	// verify state changed
	ast := ballot.Show(ctx, gov.GovAddress(cty.Organizer().Public), ballotName)
	if !ast.Ad.Frozen {
		t.Errorf("expecting frozen")
	}

	// tally
	tallyChg := ballot.Tally(ctx, cty.Organizer(), ballotName, testMaxPar)
	fmt.Println("tally: ", form.SprintJSON(tallyChg))
	if tallyChg.Result.Scores[choices[0]] != 0.0 {
		t.Errorf("expecting %v, got %v", 0.0, tallyChg.Result.Scores[choices[0]])
	}

	// testutil.Hang()
}
