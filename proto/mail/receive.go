package mail

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

type MsgEffect[Msg form.Form, Effect form.Form] struct {
	Msg    Msg
	Effect Effect
}

type Receiver[Msg form.Form, Effect form.Form] func(
	ctx context.Context,
	msg Msg,
) (effect Effect, err error)

func ReceiveStageOnly[Msg form.Form, Effect form.Form](
	ctx context.Context,
	receiver *git.Tree,
	senderAddr id.PublicAddress,
	sender *git.Tree,
	topic string,
	receive Receiver[Msg, Effect],
) git.Change[form.Map, []MsgEffect[Msg, Effect]] {

	// prep
	receiverCred := id.GetPublicCredentials(ctx, receiver)
	senderCred := id.GetPublicCredentials(ctx, sender)
	senderTopicNS := SendTopicNS(receiverCred.ID, topic)
	receiverTopicNS := ReceiveTopicNS(senderCred.ID, topic)
	receiverNextNS := receiverTopicNS.Sub(NextFilebase)
	senderNextNS := senderTopicNS.Sub(NextFilebase)
	receiverInfoNS := receiverTopicNS.Sub(BoxInfoFilebase)

	// read receiver and sender next seq no
	receiverNextSeqNo, _ := git.TryFromFile[SeqNo](ctx, receiver, receiverNextNS.Path())
	senderNextSeqNo, _ := git.TryFromFile[SeqNo](ctx, sender, senderNextNS.Path())

	// write receive box info
	info := ReceiveBoxInfo{SenderID: senderCred.ID, Topic: topic}
	git.ToFileStage(ctx, receiver, receiverInfoNS.Path(), info)

	// read unread messages
	receiverLatestNextSeqNo := receiverNextSeqNo
	base.Infof("receiving receiverSeqNo=%v senderSeqNo=%v", receiverNextSeqNo, senderNextSeqNo)
	rr := []MsgEffect[Msg, Effect]{}
	for i := receiverNextSeqNo; i < senderNextSeqNo; i++ {
		msgFilebase := strconv.Itoa(int(i))
		msg := git.FromFile[Msg](ctx, sender, senderTopicNS.Sub(msgFilebase).Path())
		effect, err := receive(ctx, msg)
		if err != nil {
			base.Infof("responding to message %d in sender repo (%v)", i, err)
			continue
		}
		git.ToFileStage(ctx, receiver, receiverTopicNS.Sub(msgFilebase).Path(), effect)
		rr = append(rr, MsgEffect[Msg, Effect]{Msg: msg, Effect: effect})
		receiverLatestNextSeqNo = i + 1
	}

	// write receiver-side next seq no
	git.ToFileStage(ctx, receiver, receiverNextNS.Path(), receiverLatestNextSeqNo)

	return git.NewChange(
		fmt.Sprintf("Received mail"),
		"mail_receive",
		form.Map{"topic": topic},
		rr,
		nil,
	)
}

type SignedReceiver[Msg form.Form, Effect form.Form] func(
	ctx context.Context,
	msg Msg,
	signedReq id.SignedPlaintext,
) (effect Effect, err error)

func ReceiveSignedStageOnly[Msg form.Form, Effect form.Form](
	ctx context.Context,
	receiverCloned id.OwnerCloned,
	senderAddr id.PublicAddress,
	senderPublic *git.Tree,
	topic string,
	receive SignedReceiver[Msg, Effect],
) git.Change[form.Map, []MsgEffect[Msg, Effect]] {

	receiverPrivCred := id.GetOwnerCredentials(ctx, receiverCloned)
	rr := []MsgEffect[Msg, Effect]{}
	signRespond := func(ctx context.Context, signedReq id.SignedPlaintext) (signedResp id.SignedPlaintext, err error) {
		if !signedReq.Verify() {
			return signedResp, fmt.Errorf("signature not valid")
		}
		msg, err := form.DecodeBytes[Msg](ctx, signedReq.Plaintext)
		if err != nil {
			return signedResp, err
		}
		effect, err := receive(ctx, msg, signedReq)
		if err != nil {
			return signedResp, err
		}
		rr = append(rr, MsgEffect[Msg, Effect]{Msg: msg, Effect: effect})
		return id.Sign(ctx, receiverPrivCred, effect), nil
	}
	recvOnly := ReceiveStageOnly(ctx, receiverCloned.Public.Tree(), senderAddr, senderPublic, topic, signRespond)
	return git.NewChange(
		fmt.Sprintf("Received signed mail."),
		"mail_receive_signed",
		form.Map{"topic": topic},
		rr,
		form.Forms{recvOnly},
	)
}
