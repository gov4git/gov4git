package mail

import (
	"context"
	"testing"

	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/testutil"
)

func TestSendReceive(t *testing.T) {
	ctx := testutil.NewCtx()
	testSenderID := id.NewTestID(ctx, t, git.MainBranch, false)
	testReceiverID := id.NewTestID(ctx, t, git.MainBranch, false)
	id.InitLocal(ctx, testSenderID.OwnerAddress(), testSenderID.OwnerCloned())
	id.InitLocal(ctx, testReceiverID.OwnerAddress(), testReceiverID.OwnerCloned())

	const testTopic = "topic"
	var testMsg []string = []string{"a", "b", "c"}

	s0 := SendStageOnly(ctx, testSenderID.Public.Tree(), testReceiverID.Public.Tree(), testTopic, testMsg[0])
	if s0.Result != 0 {
		t.Fatalf("unexpected seq no")
	}

	respond := func(ctx context.Context, req string) (resp string, err error) {
		return req, nil
	}

	r0 := ReceiveStageOnly(
		ctx,
		testReceiverID.Public.Tree(),
		testSenderID.PublicAddress(),
		testSenderID.Public.Tree(),
		testTopic,
		respond,
	)
	if len(r0.Result) != 1 {
		t.Fatalf("unexpecte length")
	}
	if r0.Result[0].Response != testMsg[0] {
		t.Fatalf("expecting %v, got %v", testMsg[0], r0.Result[0])
	}

	s1 := SendStageOnly(ctx, testSenderID.Public.Tree(), testReceiverID.Public.Tree(), testTopic, testMsg[1])
	if s1.Result != 1 {
		t.Fatalf("expecting %v, got %v", 1, s1.Result)
	}

	s2 := SendStageOnly(ctx, testSenderID.Public.Tree(), testReceiverID.Public.Tree(), testTopic, testMsg[2])
	if s2.Result != 2 {
		t.Fatalf("expecting %v, got %v", 2, s2.Result)
	}

	r12 := ReceiveStageOnly(
		ctx,
		testReceiverID.Public.Tree(),
		testSenderID.PublicAddress(),
		testSenderID.Public.Tree(),
		testTopic,
		respond,
	)
	if len(r12.Result) != 2 {
		t.Fatalf("unexpecte length")
	}
	if r12.Result[0].Response != testMsg[1] {
		t.Fatalf("expecting %v, got %v", testMsg[1], r12.Result[0])
	}
	if r12.Result[1].Response != testMsg[2] {
		t.Fatalf("expecting %v, got %v", testMsg[2], r12.Result[1])
	}
}

func TestSendReceiveSigned(t *testing.T) {
	ctx := testutil.NewCtx()
	testSenderID := id.NewTestID(ctx, t, git.MainBranch, false)
	testReceiverID := id.NewTestID(ctx, t, git.MainBranch, false)
	id.InitLocal(ctx, testSenderID.OwnerAddress(), testSenderID.OwnerCloned())
	id.InitLocal(ctx, testReceiverID.OwnerAddress(), testReceiverID.OwnerCloned())

	const testTopic = "topic"
	var testMsg []string = []string{"a", "b", "c"}

	s0 := SendSignedStageOnly(ctx, testSenderID.OwnerCloned(), testReceiverID.Public.Tree(), testTopic, testMsg[0])
	if s0.Result != 0 {
		t.Fatalf("unexpected seq no")
	}

	respond := func(ctx context.Context, req string, _ id.SignedPlaintext) (resp string, err error) {
		return req, nil
	}

	r0 := ReceiveSignedStageOnly(
		ctx,
		testReceiverID.OwnerCloned(),
		testSenderID.PublicAddress(),
		testSenderID.Public.Tree(),
		testTopic,
		respond,
	)
	if len(r0.Result) != 1 {
		t.Fatalf("unexpecte length")
	}
	if r0.Result[0].Response != testMsg[0] {
		t.Fatalf("expecting %v, got %v", testMsg[0], r0.Result[0])
	}

	s1 := SendSignedStageOnly(ctx, testSenderID.OwnerCloned(), testReceiverID.Public.Tree(), testTopic, testMsg[1])
	if s1.Result != 1 {
		t.Fatalf("expecting %v, got %v", 1, s1.Result)
	}

	s2 := SendSignedStageOnly(ctx, testSenderID.OwnerCloned(), testReceiverID.Public.Tree(), testTopic, testMsg[2])
	if s2.Result != 2 {
		t.Fatalf("expecting %v, got %v", 2, s2.Result)
	}

	r12 := ReceiveSignedStageOnly(
		ctx,
		testReceiverID.OwnerCloned(),
		testSenderID.PublicAddress(),
		testSenderID.Public.Tree(),
		testTopic,
		respond,
	)
	if len(r12.Result) != 2 {
		t.Fatalf("unexpecte length")
	}
	if r12.Result[0].Response != testMsg[1] {
		t.Fatalf("expecting %v, got %v", testMsg[1], r12.Result[0])
	}
	if r12.Result[1].Response != testMsg[2] {
		t.Fatalf("expecting %v, got %v", testMsg[2], r12.Result[1])
	}
}
