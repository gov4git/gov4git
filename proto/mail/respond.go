package mail

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

type Responder[Req form.Form, Resp form.Form] func(
	ctx context.Context,
	seqNo SeqNo,
	req Req,
) (resp Resp, err error)

func RespondStageOnly[Req form.Form, Resp form.Form](
	ctx context.Context,
	receiverCloned id.OwnerCloned,
	senderAddr id.PublicAddress,
	senderPublic *git.Tree,
	topic string,
	respond Responder[Req, Resp],
) git.Change[form.Map, []ResponseEnvelope[Req, Resp]] {

	receive := func(
		ctx context.Context,
		seqNo SeqNo,
		reqEnv RequestEnvelope[Req],
		signedReqEnv id.SignedPlaintext,
	) (ResponseEnvelope[Req, Resp], error) {

		must.Assertf(ctx, reqEnv.SeqNo == seqNo, "request seqno %d does not match response seqno %d", reqEnv.SeqNo, seqNo)

		resp, err := respond(ctx, seqNo, reqEnv.Request)
		if err != nil {
			return ResponseEnvelope[Req, Resp]{}, err
		}
		return ResponseEnvelope[Req, Resp]{
			SeqNo:    seqNo,
			Request:  reqEnv.Request,
			Response: resp,
		}, nil
	}

	chg := ReceiveSignedStageOnly(ctx, receiverCloned, senderAddr, senderPublic, topic, receive)
	respEnvs := make([]ResponseEnvelope[Req, Resp], len(chg.Result))
	for i, msgEffect := range chg.Result {
		respEnvs[i] = msgEffect.Effect
	}

	return git.NewChange(
		fmt.Sprintf("Responded to %d requests", len(respEnvs)),
		"respond",
		form.Map{"topic": topic},
		respEnvs,
		form.Forms{chg},
	)
}
