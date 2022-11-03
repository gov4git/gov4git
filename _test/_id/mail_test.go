package id

import (
	"testing"

	"github.com/gov4git/gov4git/lib/base"
	"github.com/gov4git/gov4git/services/id"
	"github.com/gov4git/gov4git/testutil"
)

func TestSendReceivePlaintext(t *testing.T) {
	base.LogVerbosely()

	// create test community
	// dir := testutil.MakeStickyTestDir()
	dir := t.TempDir()
	testCommunity, err := testutil.CreateTestCommunity(dir, 2)
	if err != nil {
		t.Fatal(err)
	}
	ctx := testCommunity.Background()

	svc0 := id.IdentityService{IdentityConfig: testCommunity.UserIdentityConfig(0)}
	svc1 := id.IdentityService{IdentityConfig: testCommunity.UserIdentityConfig(1)}

	testMsg := []string{"a", "b", "c"}
	testTopic := "topic"

	// send msg 1
	_, err = svc0.SendMail(ctx, &id.SendMailIn{
		ReceiverRepo: testCommunity.UserPublicRepoURL(1),
		Topic:        testTopic,
		Message:      testMsg[0],
	})
	if err != nil {
		t.Fatal(err)
	}

	// receive msg 1
	receiveOut, err := svc1.ReceiveMail(ctx,
		&id.ReceiveMailIn{
			SenderRepo: testCommunity.UserPublicRepoURL(0),
			Topic:      testTopic,
		})
	if err != nil {
		t.Fatal(err)
	}

	if len(receiveOut.Messages) != 1 {
		t.Fatalf("expecting 1 message, got %v", len(receiveOut.Messages))
	}

	if receiveOut.Messages[0] != testMsg[0] {
		t.Fatalf("expecting %v, got %v", testMsg[0], receiveOut.Messages[0])
	}

	// send msg 2, 3
	_, err = svc0.SendMail(ctx, &id.SendMailIn{
		ReceiverRepo: testCommunity.UserPublicRepoURL(1),
		Topic:        testTopic,
		Message:      testMsg[1],
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc0.SendMail(ctx, &id.SendMailIn{
		ReceiverRepo: testCommunity.UserPublicRepoURL(1),
		Topic:        testTopic,
		Message:      testMsg[2],
	})
	if err != nil {
		t.Fatal(err)
	}

	// receive msg 2, 3
	receive23Out, err := svc1.ReceiveMail(ctx,
		&id.ReceiveMailIn{
			SenderRepo: testCommunity.UserPublicRepoURL(0),
			Topic:      testTopic,
		})
	if err != nil {
		t.Fatal(err)
	}

	if len(receive23Out.Messages) != 2 {
		t.Fatalf("expecting 2 messages, got %v", len(receive23Out.Messages))
	}

	if receive23Out.Messages[0] != testMsg[1] {
		t.Fatalf("expecting %v, got %v", testMsg[1], receive23Out.Messages[0])
	}

	if receive23Out.Messages[1] != testMsg[2] {
		t.Fatalf("expecting %v, got %v", testMsg[2], receive23Out.Messages[1])
	}

}
