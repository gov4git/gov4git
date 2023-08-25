package mail

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Request_StageOnly[Req form.Form](
	ctx context.Context,
	senderCloned id.OwnerCloned,
	receiver *git.Tree,
	topic string,
	req Req,
) git.Change[form.Map, RequestEnvelope[Req]] {

	var called bool
	var msg RequestEnvelope[Req]
	mkEnv := func(_ context.Context, seqNo SeqNo) RequestEnvelope[Req] {

		must.Assertf(ctx, !called, "msg maker must be called once")
		msg = RequestEnvelope[Req]{
			SeqNo:   seqNo,
			Request: req,
		}
		return msg
	}

	chg := SendSignedMakeMsg_StageOnly(ctx, senderCloned, receiver, topic, mkEnv)

	return git.NewChange(
		fmt.Sprintf("Requested #%d", chg.Result),
		"request",
		form.Map{"topic": topic, "msg": msg},
		msg,
		form.Forms{chg},
	)
}
