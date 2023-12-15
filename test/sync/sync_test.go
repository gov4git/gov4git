package sync

import (
	"fmt"
	"math"
	"testing"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballot"
	"github.com/gov4git/gov4git/v2/proto/ballot/common"
	"github.com/gov4git/gov4git/v2/proto/ballot/load"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/gov4git/v2/proto/sync"
	"github.com/gov4git/gov4git/v2/runtime"
	"github.com/gov4git/gov4git/v2/test"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/testutil"
)

func TestSync(t *testing.T) {
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	ballotName0 := common.BallotName{"a", "b", "c"}
	ballotName1 := common.BallotName{"d", "e", "f"}
	choices := []string{"x", "y", "z"}

	// open two ballots
	strat := load.QVStrategyName
	openChg0 := ballot.Open(ctx, strat, cty.Organizer(), ballotName0, "ballot_0", "ballot 0", choices, member.Everybody)
	fmt.Println("open 0: ", form.SprintJSON(openChg0))
	openChg1 := ballot.Open(ctx, strat, cty.Organizer(), ballotName1, "ballot_1", "ballot 1", choices, member.Everybody)
	fmt.Println("open 1: ", form.SprintJSON(openChg1))

	// give credits to users
	account.Issue(ctx, cty.Gov(), cty.MemberAccountID(0), account.H(account.PluralAsset, 5.0), "test")
	account.Issue(ctx, cty.Gov(), cty.MemberAccountID(1), account.H(account.PluralAsset, 5.0), "test")

	// vote
	elections0 := common.Elections{common.NewElection(choices[0], 5.0)}
	elections1 := common.Elections{common.NewElection(choices[0], -5.0)}
	voteChg0 := ballot.Vote(ctx, cty.MemberOwner(0), cty.Gov(), ballotName0, elections0)
	fmt.Println("vote 0: ", form.SprintJSON(voteChg0))
	voteChg1 := ballot.Vote(ctx, cty.MemberOwner(1), cty.Gov(), ballotName1, elections1)
	fmt.Println("vote 1: ", form.SprintJSON(voteChg1))

	// tally
	syncChg := sync.Sync(ctx, cty.Organizer(), 2)
	fmt.Println("sync: ", form.SprintJSON(syncChg))

	// verify tallies are correct
	ast0 := ballot.Show(ctx, cty.Gov(), ballotName0)
	if ast0.Tally.Scores[choices[0]] != math.Sqrt(5.0) {
		t.Errorf("expecting %v, got %v", math.Sqrt(5.0), ast0.Tally.Scores[choices[0]])
	}
	ast1 := ballot.Show(ctx, cty.Gov(), ballotName1)
	if ast1.Tally.Scores[choices[0]] != -math.Sqrt(5.0) {
		t.Errorf("expecting %v, got %v", -math.Sqrt(5.0), ast1.Tally.Scores[choices[0]])
	}

	// testutil.Hang()
}
