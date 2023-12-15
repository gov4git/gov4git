package mail

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/v2/proto/id"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func Request_StageOnly[Req form.Form](
	ctx context.Context,
	senderCloned id.OwnerCloned,
	receiver *git.Tree,
	topic string,
	req Req,
) git.Change[form.Map, RequestEnvelope[Req]] {

	mkReqEnv := func(_ context.Context, seqNo SeqNo) RequestEnvelope[Req] {
		return RequestEnvelope[Req]{
			SeqNo:   seqNo,
			Request: req,
		}
	}

	chg := SendSignedMakeMsg_StageOnly(ctx, senderCloned, receiver, topic, mkReqEnv)
	msg := chg.Result.Msg

	return git.NewChange(
		fmt.Sprintf("Requested #%d", chg.Result.SeqNo),
		"request",
		form.Map{"topic": topic, "msg": msg},
		msg,
		form.Forms{chg},
	)
}
