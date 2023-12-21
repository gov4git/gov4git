package ballot

import (
	"testing"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/ballot/ballot"
	"github.com/gov4git/gov4git/v2/proto/ballot/common"
	"github.com/gov4git/gov4git/v2/proto/ballot/load"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/gov4git/v2/proto/purpose"
	"github.com/gov4git/gov4git/v2/runtime"
	"github.com/gov4git/gov4git/v2/test"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/testutil"
)

func TestTrack(t *testing.T) {

	ctx := testutil.NewCtx(t, runtime.TestWithCache)
	cty := test.NewTestCommunity(t, ctx, 2)

	ballotName := common.BallotName{"a", "b", "c"}
	choices := []string{"x", "y", "z"}

	// give voter credits
	account.Issue(ctx, cty.Gov(), cty.MemberAccountID(0), account.H(account.PluralAsset, 6.0), "test")

	// open ballot
	ballot.Open(ctx, load.QVStrategyName, cty.Organizer(), ballotName, account.NobodyAccountID, purpose.Unspecified, "ballot title", "ballot description", choices, member.Everybody)

	// vote 1: accepted vote
	elections1 := common.Elections{common.NewElection(choices[0], 1.0)}
	ballot.Vote(ctx, cty.MemberOwner(0), cty.Gov(), ballotName, elections1)

	// tally: collect and accept first vote
	ballot.Tally(ctx, cty.Organizer(), ballotName, testMaxPar)

	// vote 2: rejected vote
	elections2 := common.Elections{common.NewElection(choices[0], 2.0)}
	ballot.Vote(ctx, cty.MemberOwner(0), cty.Gov(), ballotName, elections2)

	// freeze ballot
	ballot.Freeze(ctx, cty.Organizer(), ballotName)

	// tally: collect and reject second vote (because ballot is frozen)
	ballot.Tally(ctx, cty.Organizer(), ballotName, testMaxPar)

	// unfreeze ballot
	ballot.Unfreeze(ctx, cty.Organizer(), ballotName)

	// vote 3: pending vote (never tallied)
	elections3 := common.Elections{common.NewElection(choices[0], 3.0)}
	ballot.Vote(ctx, cty.MemberOwner(0), cty.Gov(), ballotName, elections3)

	// track
	status := ballot.Track(ctx, cty.MemberOwner(0), cty.Gov(), ballotName)
	if len(status.AcceptedVotes) != 1 || status.AcceptedVotes[0].Vote.VoteStrengthChange != 1.0 {
		t.Errorf("expecting one accepted vote with strength 1.0, got %v", form.SprintJSON(status))
	}
	if len(status.RejectedVotes) != 1 || status.RejectedVotes[0].Vote.VoteStrengthChange != 2.0 {
		t.Errorf("expecting one rejected vote with strength 2.0, got %v", form.SprintJSON(status))
	}
	if len(status.PendingVotes) != 1 || status.PendingVotes[0].VoteStrengthChange != 3.0 {
		t.Errorf("expecting one pending vote with strength 3.0, got %v", form.SprintJSON(status))
	}

	// testutil.Hang()
}
