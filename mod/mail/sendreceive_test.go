package mail

import (
	"context"
	"testing"

	"github.com/gov4git/gov4git/lib/testutil"
	"github.com/gov4git/gov4git/mod/id"
)

func TestSendReceive(t *testing.T) {
	ctx := testutil.NewCtx()
	testSenderID := id.InitTestID(ctx, t, false)
	testReceiverID := id.InitTestID(ctx, t, false)
	id.InitLocal(ctx,
		testSenderID.Public.Address, testSenderID.Private.Address,
		testSenderID.Public.Tree, testSenderID.Private.Tree,
	)
	id.InitLocal(ctx,
		testReceiverID.Public.Address, testReceiverID.Private.Address,
		testReceiverID.Public.Tree, testReceiverID.Private.Tree,
	)

	const testTopic = "topic"
	var testMsg []string = []string{"a", "b", "c"}

	s0 := Send(ctx, testSenderID.Public.Tree, testReceiverID.Public.Tree, testTopic, testMsg[0])
	if s0.Result != 0 {
		t.Fatalf("unexpected seq no")
	}

	respond := func(ctx context.Context, req string) (resp string, err error) {
		return req, nil
	}

	r0 := Receive(
		ctx,
		testReceiverID.Public.Tree,
		testSenderID.Public.Address,
		testSenderID.Public.Tree,
		testTopic,
		respond,
	)
	if len(r0.Result) != 1 {
		t.Fatalf("unexpecte length")
	}
	if r0.Result[0].Response != testMsg[0] {
		t.Fatalf("expecting %v, got %v", testMsg[0], r0.Result[0])
	}

	s1 := Send(ctx, testSenderID.Public.Tree, testReceiverID.Public.Tree, testTopic, testMsg[1])
	if s1.Result != 1 {
		t.Fatalf("expecting %v, got %v", 1, s1.Result)
	}

	s2 := Send(ctx, testSenderID.Public.Tree, testReceiverID.Public.Tree, testTopic, testMsg[2])
	if s2.Result != 2 {
		t.Fatalf("expecting %v, got %v", 2, s2.Result)
	}

	r12 := Receive(
		ctx,
		testReceiverID.Public.Tree,
		testSenderID.Public.Address,
		testSenderID.Public.Tree,
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
	testSenderID := id.InitTestID(ctx, t, false)
	testReceiverID := id.InitTestID(ctx, t, false)
	id.InitLocal(ctx,
		testSenderID.Public.Address, testSenderID.Private.Address,
		testSenderID.Public.Tree, testSenderID.Private.Tree,
	)
	id.InitLocal(ctx,
		testReceiverID.Public.Address, testReceiverID.Private.Address,
		testReceiverID.Public.Tree, testReceiverID.Private.Tree,
	)

	const testTopic = "topic"
	var testMsg []string = []string{"a", "b", "c"}

	s0 := SendSigned(ctx, testSenderID.Public.Tree, testSenderID.Private.Tree, testReceiverID.Public.Tree, testTopic, testMsg[0])
	if s0.Result != 0 {
		t.Fatalf("unexpected seq no")
	}

	respond := func(ctx context.Context, req string, _ id.SignedPlaintext) (resp string, err error) {
		return req, nil
	}

	r0 := ReceiveSigned(
		ctx,
		testReceiverID.Public.Tree,
		testReceiverID.Private.Tree,
		testSenderID.Public.Address,
		testSenderID.Public.Tree,
		testTopic,
		respond,
	)
	if len(r0.Result) != 1 {
		t.Fatalf("unexpecte length")
	}
	if r0.Result[0].Response != testMsg[0] {
		t.Fatalf("expecting %v, got %v", testMsg[0], r0.Result[0])
	}

	s1 := SendSigned(ctx, testSenderID.Public.Tree, testSenderID.Private.Tree, testReceiverID.Public.Tree, testTopic, testMsg[1])
	if s1.Result != 1 {
		t.Fatalf("expecting %v, got %v", 1, s1.Result)
	}

	s2 := SendSigned(ctx, testSenderID.Public.Tree, testSenderID.Private.Tree, testReceiverID.Public.Tree, testTopic, testMsg[2])
	if s2.Result != 2 {
		t.Fatalf("expecting %v, got %v", 2, s2.Result)
	}

	r12 := ReceiveSigned(
		ctx,
		testReceiverID.Public.Tree,
		testReceiverID.Private.Tree,
		testSenderID.Public.Address,
		testSenderID.Public.Tree,
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
