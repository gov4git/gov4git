package ballot

import (
	"fmt"
	"testing"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotapi"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotio"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballotproto"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/gov4git/v2/proto/purpose"
	"github.com/gov4git/gov4git/v2/runtime"
	"github.com/gov4git/gov4git/v2/test"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/must"
	"github.com/gov4git/lib4git/testutil"
)

func TestInsufficientCredits(t *testing.T) {
	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	ballotName := ballotproto.BallotName{"a", "b", "c"}
	choices := []string{"x", "y", "z"}

	// open
	strat := ballotio.QVStrategyName
	openChg := ballotapi.Open(
		ctx,
		strat,
		cty.Organizer(),
		ballotName,
		account.NobodyAccountID,
		purpose.Unspecified,
		"ballot title",
		"ballot description",
		choices,
		member.Everybody,
	)
	fmt.Println("open: ", form.SprintJSON(openChg))

	// give voter credits
	account.Issue(ctx, cty.Gov(), cty.MemberAccountID(0), account.H(account.PluralAsset, 1.0), "test")

	// vote
	elections := ballotproto.Elections{
		ballotproto.NewElection(choices[0], 2.0),
	}
	if err := must.Try(
		func() { ballotapi.Vote(ctx, cty.MemberOwner(0), cty.Gov(), ballotName, elections) },
	); err != nil {
		fmt.Println("vote rejected: ", err.Error())
	} else {
		t.Fatalf("vote must fail")
	}

	// tally
	tallyChg := ballotapi.Tally(ctx, cty.Organizer(), ballotName, testMaxPar)
	fmt.Println("tally: ", form.SprintJSON(tallyChg))
	if tallyChg.Result.Scores[choices[0]] != 0.0 {
		t.Errorf("expecting %v, got %v", 0.0, tallyChg.Result.Scores[choices[0]])
	}

	// close
	closeChg := ballotapi.Close(ctx, cty.Organizer(), ballotName, account.BurnAccountID)
	fmt.Println("close: ", form.SprintJSON(closeChg))

	// testutil.Hang()
}
