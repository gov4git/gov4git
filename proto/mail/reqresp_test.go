package mail

import (
	"context"
	"testing"

	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/testutil"
)

func TestReqResp(t *testing.T) {
	ctx := testutil.NewCtx()
	testSenderID := id.NewTestID(ctx, t, git.MainBranch, false)
	testReceiverID := id.NewTestID(ctx, t, git.MainBranch, false)
	id.InitLocal(ctx, testSenderID.OwnerAddress(), testSenderID.OwnerCloned())
	id.InitLocal(ctx, testReceiverID.OwnerAddress(), testReceiverID.OwnerCloned())

	const testTopic = "topic"
	var testMsg []string = []string{"a", "b", "c"}

	s0 := Request_StageOnly(ctx, testSenderID.OwnerCloned(), testReceiverID.Public.Tree(), testTopic, testMsg[0])
	if s0.Result.Request != testMsg[0] {
		t.Errorf("expecting %v, got %v", testMsg[0], s0.Result.Request)
	}

	respond := func(ctx context.Context, _ SeqNo, req string) (resp string, err error) {
		return req, nil
	}

	r0 := Respond_StageOnly[string, string](
		ctx,
		testReceiverID.OwnerCloned(),
		testSenderID.PublicAddress(),
		testSenderID.Public.Tree(),
		testTopic,
		respond,
	)
	if len(r0.Result) != 1 {
		t.Fatalf("unexpected length")
	}
	if r0.Result[0].Response != testMsg[0] {
		t.Fatalf("expecting %v, got %v", testMsg[0], r0.Result[0].Response)
	}

	Request_StageOnly(ctx, testSenderID.OwnerCloned(), testReceiverID.Public.Tree(), testTopic, testMsg[1])
	Request_StageOnly(ctx, testSenderID.OwnerCloned(), testReceiverID.Public.Tree(), testTopic, testMsg[2])

	r12 := Respond_StageOnly[string, string](
		ctx,
		testReceiverID.OwnerCloned(),
		testSenderID.PublicAddress(),
		testSenderID.Public.Tree(),
		testTopic,
		respond,
	)
	if len(r12.Result) != 2 {
		t.Fatalf("unexpected length")
	}
	if r12.Result[0].Response != testMsg[1] {
		t.Fatalf("expecting %v, got %v", testMsg[1], r12.Result[0].Response)
	}
	if r12.Result[1].Response != testMsg[2] {
		t.Fatalf("expecting %v, got %v", testMsg[2], r12.Result[1].Response)
	}
}
