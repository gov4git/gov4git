package pmp

import (
	"math"
	"testing"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/motion/motionapi"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/lib4git/must"
)

func TestOpenCancelConcernCloseProposal(t *testing.T) {
	c := testCaseWithoutMatch
	ctx, cty := SetupTest(t, c)

	motionapi.CancelMotion(ctx, cty.Organizer(), testConcernID)                     // issue
	motionapi.CloseMotion(ctx, cty.Organizer(), testProposalID, motionproto.Accept) // pr

	// user accounts
	u0 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(0)).Balance(account.PluralAsset)
	u1 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(1)).Balance(account.PluralAsset)

	exp0 := c.User0Credits - math.Abs(c.User0ProposalStrength) + 16.733
	if math.Abs(u0.Quantity-exp0) > 0.1 {
		t.Errorf("expecting %v, got %v", exp0, u0.Quantity)
	}

	exp1 := c.User1Credits + c.User1ProposalStrength
	if math.Abs(u1.Quantity-exp1) > 0.1 {
		t.Errorf("expecting %v, got %v", exp1, u1.Quantity)
	}
}

func TestOpenCancelConcernCancelProposal(t *testing.T) {
	c := testCaseWithoutMatch
	ctx, cty := SetupTest(t, c)

	motionapi.CancelMotion(ctx, cty.Organizer(), testConcernID)        // issue
	motionapi.CancelMotion(ctx, cty.Organizer(), testProposalID, true) // pr

	// user accounts
	u0 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(0)).Balance(account.PluralAsset)
	u1 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(1)).Balance(account.PluralAsset)

	exp0 := c.User0Credits
	if u0.Quantity != exp0 {
		t.Errorf("expecting %v, got %v", exp0, u0.Quantity)
	}

	exp1 := c.User1Credits
	if u1.Quantity != exp1 {
		t.Errorf("expecting %v, got %v", exp1, u1.Quantity)
	}
}

func TestOpenCancelProposalCancelConcern(t *testing.T) {
	c := testCaseWithoutMatch
	ctx, cty := SetupTest(t, c)

	motionapi.CancelMotion(ctx, cty.Organizer(), testProposalID, true) // pr
	motionapi.CancelMotion(ctx, cty.Organizer(), testConcernID)        // issue

	// user accounts
	u0 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(0)).Balance(account.PluralAsset)
	u1 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(1)).Balance(account.PluralAsset)

	exp0 := c.User0Credits
	if u0.Quantity != exp0 {
		t.Errorf("expecting %v, got %v", exp0, u0.Quantity)
	}

	exp1 := c.User1Credits
	if u1.Quantity != exp1 {
		t.Errorf("expecting %v, got %v", exp1, u1.Quantity)
	}
}

func TestAcceptLinkedProposalWithoutMatch(t *testing.T) {
	c := testCaseWithoutMatch
	ctx, cty := SetupTest(t, c)

	motionapi.CloseMotion(ctx, cty.Organizer(), testProposalID, motionproto.Accept) // pr

	err := must.Try(
		func() {
			motionapi.CancelMotion(ctx, cty.Organizer(), testConcernID)
		},
	) // issue
	if err == nil {
		t.Errorf("expecting error")
	}

	// user accounts
	u0 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(0)).Balance(account.PluralAsset)
	u1 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(1)).Balance(account.PluralAsset)
	u2 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(2)).Balance(account.PluralAsset)

	if math.Abs(u0.Quantity-c.User0EndBalance) > 0.01 {
		t.Errorf("expecting %v, got %v", c.User0EndBalance, u0.Quantity)
	}

	if math.Abs(u1.Quantity-c.User1EndBalance) > 0.01 {
		t.Errorf("expecting %v, got %v", c.User1EndBalance, u1.Quantity)
	}

	if math.Abs(u2.Quantity-c.User2EndBalance) > 0.01 {
		t.Errorf("expecting %v, got %v", c.User2EndBalance, u2.Quantity)
	}
}

func TestAcceptLinkedProposalWithMatch(t *testing.T) {
	c := testCaseWithMatch
	ctx, cty := SetupTest(t, c)

	motionapi.CloseMotion(ctx, cty.Organizer(), testProposalID, motionproto.Accept) // pr

	err := must.Try(
		func() {
			motionapi.CancelMotion(ctx, cty.Organizer(), testConcernID)
		},
	) // issue
	if err == nil {
		t.Errorf("expecting error")
	}

	// user accounts
	u0 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(0)).Balance(account.PluralAsset)
	u1 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(1)).Balance(account.PluralAsset)
	u2 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(2)).Balance(account.PluralAsset)

	if math.Abs(u0.Quantity-c.User0EndBalance) > 0.01 {
		t.Errorf("expecting %v, got %v", c.User0EndBalance, u0.Quantity)
	}

	if math.Abs(u1.Quantity-c.User1EndBalance) > 0.01 {
		t.Errorf("expecting %v, got %v", c.User1EndBalance, u1.Quantity)
	}

	if math.Abs(u2.Quantity-c.User2EndBalance) > 0.01 {
		t.Errorf("expecting %v, got %v", c.User2EndBalance, u2.Quantity)
	}
}
