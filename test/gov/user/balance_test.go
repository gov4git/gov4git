package user

import (
	"testing"

	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/services/gov/user"
	"github.com/gov4git/gov4git/testutil"
)

func TestBallot(t *testing.T) {
	// base.LogVerbosely()

	// create test community
	// dir := testutil.MakeStickyTestDir()
	dir := t.TempDir()
	testCommunity, err := testutil.CreateTestCommunity(dir, 1)
	if err != nil {
		t.Fatal(err)
	}
	ctx := testCommunity.Background()

	svc := user.GovUserService{GovConfig: testCommunity.CommunityGovConfig()}

	// set user balance
	const testInitBalance = 1.0
	_, err = svc.BalanceSet(ctx,
		&user.BalanceSetIn{
			User:    testCommunity.User(0),
			Balance: "credits",
			Value:   testInitBalance,
			Branch:  proto.MainBranch,
		})
	if err != nil {
		t.Fatal(err)
	}

	// get user balance
	getOut, err := svc.BalanceGet(ctx,
		&user.BalanceGetIn{
			User:    testCommunity.User(0),
			Balance: "credits",
			Branch:  proto.MainBranch,
		})
	if err != nil {
		t.Fatal(err)
	}

	if getOut.Value != testInitBalance {
		t.Errorf("expecting %v, got %v", testInitBalance, getOut.Value)
	}

	// add user balance
	const testAddBalance = 1.0
	addOut, err := svc.BalanceAdd(ctx,
		&user.BalanceAddIn{
			User:    testCommunity.User(0),
			Balance: "credits",
			Value:   testAddBalance,
			Branch:  proto.MainBranch,
		})
	if err != nil {
		t.Fatal(err)
	}

	if addOut.ValueBefore != testInitBalance {
		t.Errorf("expecting %v, got %v", testAddBalance, addOut.ValueBefore)
	}
	afterAdd := testInitBalance + testAddBalance
	if addOut.ValueAfter != afterAdd {
		t.Errorf("expecting %v, got %v", afterAdd, addOut.ValueAfter)
	}

	// mul user balance
	const testMulBalance = 2.0
	mulOut, err := svc.BalanceMul(ctx,
		&user.BalanceMulIn{
			User:    testCommunity.User(0),
			Balance: "credits",
			Value:   testMulBalance,
			Branch:  proto.MainBranch,
		})
	if err != nil {
		t.Fatal(err)
	}

	if mulOut.ValueBefore != afterAdd {
		t.Errorf("expecting %v, got %v", afterAdd, mulOut.ValueBefore)
	}
	if mulOut.ValueAfter != afterAdd*testMulBalance {
		t.Errorf("expecting %v, got %v", afterAdd*testMulBalance, mulOut.ValueAfter)
	}
}
