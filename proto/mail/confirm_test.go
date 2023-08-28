package mail

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/testutil"
)

func TestConfirm(t *testing.T) {
	ctx := testutil.NewCtx()
	testSenderID := id.NewTestID(ctx, t, git.MainBranch, false)
	testReceiverID := id.NewTestID(ctx, t, git.MainBranch, false)
	id.InitLocal(ctx, testSenderID.OwnerAddress(), testSenderID.OwnerCloned())
	id.InitLocal(ctx, testReceiverID.OwnerAddress(), testReceiverID.OwnerCloned())

	const testTopic = "topic"
	var testMsg []string = []string{"a", "b", "c"}

	Send_StageOnly(ctx, testSenderID.Public.Tree(), testReceiverID.Public.Tree(), testTopic, testMsg[0])

	respond := func(ctx context.Context, seqNo SeqNo, req string) (resp string, err error) {
		return req, nil
	}

	Receive_StageOnly(ctx, testReceiverID.Public.Tree(), testSenderID.PublicAddress(), testSenderID.Public.Tree(), testTopic, respond)

	Send_StageOnly(ctx, testSenderID.Public.Tree(), testReceiverID.Public.Tree(), testTopic, testMsg[1])
	Send_StageOnly(ctx, testSenderID.Public.Tree(), testReceiverID.Public.Tree(), testTopic, testMsg[2])

	confirmed, notConfirmed := Confirm_Local[string, string](ctx, testSenderID.Public.Tree(), testReceiverID.Public.Tree(), testTopic)
	expConfirmed := MsgEffects[string, string]{
		{SeqNo: SeqNo(0), Msg: testMsg[0], Effect: testMsg[0]},
	}
	if !reflect.DeepEqual(confirmed, expConfirmed) {
		fmt.Errorf("expecting %v, got %v", expConfirmed, confirmed)
	}
	expNotConfirmed := MsgEffects[string, form.None]{
		{SeqNo: SeqNo(1), Msg: testMsg[1]},
		{SeqNo: SeqNo(2), Msg: testMsg[2]},
	}
	if !reflect.DeepEqual(notConfirmed, expNotConfirmed) {
		fmt.Errorf("expecting %v, got %v", expNotConfirmed, notConfirmed)
	}
}
