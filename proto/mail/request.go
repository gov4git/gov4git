package mail

import (
	"context"
	"fmt"

	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

func RequestStageOnly[Req form.Form](
	ctx context.Context,
	sender *git.Tree,
	receiver *git.Tree,
	topic string,
	req Req,
) git.Change[form.Map, RequestEnvelope[Req]] {

	var msg RequestEnvelope[Req]
	mkEnv := func(_ context.Context, seqNo SeqNo) RequestEnvelope[Req] {
		msg = RequestEnvelope[Req]{
			SeqNo:   seqNo,
			Request: req,
		}
		return msg
	}
	chg := SendMakeMsgStageOnly(ctx, sender, receiver, topic, mkEnv)
	return git.NewChange(
		fmt.Sprintf("Requested #%d", chg.Result),
		"request",
		form.Map{"topic": topic, "msg": msg},
		msg,
		form.Forms{chg},
	)
}
