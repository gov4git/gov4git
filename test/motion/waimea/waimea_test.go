package waimea

import (
	"math"
	"testing"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/motion/motionapi"
	"github.com/gov4git/gov4git/v2/proto/motion/motionproto"
	"github.com/gov4git/lib4git/must"
)

func TestCancelConcernCancelProposal(t *testing.T) {
	c := testCancelConcernCancelProposal
	ctx, cty := SetupTest(t, c)

	motionapi.CancelMotion(ctx, cty.Organizer(), testConcernID)  // issue
	motionapi.CancelMotion(ctx, cty.Organizer(), testProposalID) // pr

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

	if math.Abs(u0.Quantity-c.Voter0EndBalance) > 0.01 {
		t.Errorf("expecting %v, got %v", c.Voter0EndBalance, u0.Quantity)
	}

	if math.Abs(u1.Quantity-c.Voter1EndBalance) > 0.01 {
		t.Errorf("expecting %v, got %v", c.Voter1EndBalance, u1.Quantity)
	}

	if math.Abs(u2.Quantity-c.AuthorEndBalance) > 0.01 {
		t.Errorf("expecting %v, got %v", c.AuthorEndBalance, u2.Quantity)
	}
}

func TestCancelProposalCancelConcern(t *testing.T) {
	c := testCancelProposalCancelConcern
	ctx, cty := SetupTest(t, c)

	motionapi.CancelMotion(ctx, cty.Organizer(), testProposalID) // pr
	motionapi.CancelMotion(ctx, cty.Organizer(), testConcernID)  // issue

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

	if math.Abs(u0.Quantity-c.Voter0EndBalance) > 0.01 {
		t.Errorf("expecting %v, got %v", c.Voter0EndBalance, u0.Quantity)
	}

	if math.Abs(u1.Quantity-c.Voter1EndBalance) > 0.01 {
		t.Errorf("expecting %v, got %v", c.Voter1EndBalance, u1.Quantity)
	}

	if math.Abs(u2.Quantity-c.AuthorEndBalance) > 0.01 {
		t.Errorf("expecting %v, got %v", c.AuthorEndBalance, u2.Quantity)
	}
}

func TestCancelConcernAcceptProposal(t *testing.T) {
	c := testCancelConcernAcceptProposal
	ctx, cty := SetupTest(t, c)

	motionapi.CancelMotion(ctx, cty.Organizer(), testConcernID)                     // issue
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

	if math.Abs(u0.Quantity-c.Voter0EndBalance) > 0.01 {
		t.Errorf("expecting %v, got %v", c.Voter0EndBalance, u0.Quantity)
	}

	if math.Abs(u1.Quantity-c.Voter1EndBalance) > 0.01 {
		t.Errorf("expecting %v, got %v", c.Voter1EndBalance, u1.Quantity)
	}

	if math.Abs(u2.Quantity-c.AuthorEndBalance) > 0.01 {
		t.Errorf("expecting %v, got %v", c.AuthorEndBalance, u2.Quantity)
	}
}

func TestAcceptProposal(t *testing.T) {
	c := testAcceptProposal
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

	if math.Abs(u0.Quantity-c.Voter0EndBalance) > 0.01 {
		t.Errorf("expecting %v, got %v", c.Voter0EndBalance, u0.Quantity)
	}

	if math.Abs(u1.Quantity-c.Voter1EndBalance) > 0.01 {
		t.Errorf("expecting %v, got %v", c.Voter1EndBalance, u1.Quantity)
	}

	if math.Abs(u2.Quantity-c.AuthorEndBalance) > 0.01 {
		t.Errorf("expecting %v, got %v", c.AuthorEndBalance, u2.Quantity)
	}
}

func TestRejectProposal(t *testing.T) {
	c := testRejectProposal
	ctx, cty := SetupTest(t, c)

	motionapi.CloseMotion(ctx, cty.Organizer(), testProposalID, motionproto.Reject) // pr

	// user accounts
	u0 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(0)).Balance(account.PluralAsset)
	u1 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(1)).Balance(account.PluralAsset)
	u2 := account.Get(ctx, cty.Gov(), cty.MemberAccountID(2)).Balance(account.PluralAsset)

	if math.Abs(u0.Quantity-c.Voter0EndBalance) > 0.01 {
		t.Errorf("expecting %v, got %v", c.Voter0EndBalance, u0.Quantity)
	}

	if math.Abs(u1.Quantity-c.Voter1EndBalance) > 0.01 {
		t.Errorf("expecting %v, got %v", c.Voter1EndBalance, u1.Quantity)
	}

	if math.Abs(u2.Quantity-c.AuthorEndBalance) > 0.01 {
		t.Errorf("expecting %v, got %v", c.AuthorEndBalance, u2.Quantity)
	}
}
