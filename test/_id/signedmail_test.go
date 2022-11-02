package id

import (
	"testing"

	"github.com/gov4git/gov4git/lib/base"
	"github.com/gov4git/gov4git/services/id"
	"github.com/gov4git/gov4git/testutil"
)

func TestSignedSendReceive(t *testing.T) {
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

	testMsg := "hello world"
	testTopic := "topic"

	_, err = svc0.SendSignedMail(ctx, &id.SendMailIn{
		ReceiverRepo: testCommunity.UserPublicRepoURL(1),
		Topic:        testTopic,
		Message:      testMsg,
	})
	if err != nil {
		t.Fatal(err)
	}

	receiveOut, err := svc1.ReceiveSignedMail(ctx,
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

	if receiveOut.Messages[0] != testMsg {
		t.Fatalf("expecting %v, got %v", testMsg, receiveOut.Messages[0])
	}
}
