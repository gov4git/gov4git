package mail

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/gov4git/gov4git/v2/proto/id"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

type MsgEffect[Msg form.Form, Effect form.Form] struct {
	SeqNo  SeqNo  `json:"seqno"`
	Msg    Msg    `json:"msg"`
	Effect Effect `json:"effect"`
}

type MsgEffects[Msg form.Form, Effect form.Form] []MsgEffect[Msg, Effect]

func (x MsgEffects[Msg, Effect]) Sort()              { sort.Sort(x) }
func (x MsgEffects[Msg, Effect]) Len() int           { return len(x) }
func (x MsgEffects[Msg, Effect]) Less(i, j int) bool { return x[i].SeqNo < x[j].SeqNo }
func (x MsgEffects[Msg, Effect]) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type Receiver[Msg form.Form, Effect form.Form] func(
	ctx context.Context,
	seqNo SeqNo,
	msg Msg,
) (effect Effect, err error)

func Receive_StageOnly[Msg form.Form, Effect form.Form](
	ctx context.Context,
	receiver *git.Tree,
	senderAddr id.PublicAddress,
	sender *git.Tree,
	topic string,
	receive Receiver[Msg, Effect], // called multiple times, once for each incoming message
) git.Change[form.Map, []MsgEffect[Msg, Effect]] {

	// prep
	receiverCred := id.GetPublicCredentials(ctx, receiver)
	senderCred := id.GetPublicCredentials(ctx, sender)
	senderTopicNS := SendTopicNS(receiverCred.ID, topic)
	receiverTopicNS := ReceiveTopicNS(senderCred.ID, topic)
	receiverNextNS := receiverTopicNS.Append(NextFilebase)
	senderNextNS := senderTopicNS.Append(NextFilebase)
	receiverInfoNS := receiverTopicNS.Append(BoxInfoFilebase)

	// read receiver and sender next seq no
	receiverNextSeqNo, _ := git.TryFromFile[SeqNo](ctx, receiver, receiverNextNS)
	senderNextSeqNo, _ := git.TryFromFile[SeqNo](ctx, sender, senderNextNS)

	// write receive box info
	if _, err := receiver.Filesystem.Stat(receiverInfoNS.GitPath()); err != nil {
		info := ReceiveBoxInfo{SenderCred: senderCred, Topic: topic}
		git.ToFileStage(ctx, receiver, receiverInfoNS, info)
	}

	// read unread messages
	receiverLatestNextSeqNo := receiverNextSeqNo
	changed := false
	base.Infof("receiving receiverSeqNo=%v senderSeqNo=%v", receiverNextSeqNo, senderNextSeqNo)
	msgEffects := []MsgEffect[Msg, Effect]{}
	for i := receiverNextSeqNo; i < senderNextSeqNo; i++ {
		msgFilebase := strconv.Itoa(int(i)) + ".json"
		msg := git.FromFile[Msg](ctx, sender, senderTopicNS.Append(msgFilebase))
		effect, err := receive(ctx, i, msg)
		if err != nil {
			base.Infof("responding to message %d in sender repo (%v)", i, err)
			continue
		}
		msgEffect := MsgEffect[Msg, Effect]{SeqNo: i, Msg: msg, Effect: effect}
		git.ToFileStage(ctx, receiver, receiverTopicNS.Append(msgFilebase), msgEffect)
		msgEffects = append(msgEffects, msgEffect)
		receiverLatestNextSeqNo = i + 1
		changed = true
	}

	// write receiver-side next seq no
	if changed {
		git.ToFileStage(ctx, receiver, receiverNextNS, receiverLatestNextSeqNo)
	}

	return git.NewChange(
		fmt.Sprintf("Received %d messages", len(msgEffects)),
		"receive",
		form.Map{"topic": topic},
		msgEffects,
		nil,
	)
}

type SignedReceiver[Msg form.Form, Effect form.Form] func(
	ctx context.Context,
	seqNo SeqNo,
	signedMsg id.Signed[Msg],
) (effect Effect, err error)

func ReceiveSigned_StageOnly[Msg form.Form, Effect form.Form](
	ctx context.Context,
	receiverCloned id.OwnerCloned,
	senderAddr id.PublicAddress,
	senderPublic *git.Tree,
	topic string,
	receive SignedReceiver[Msg, Effect],
) git.Change[form.Map, []MsgEffect[Msg, Effect]] {

	receiverPrivCred := id.GetOwnerCredentials(ctx, receiverCloned)
	var receiver Receiver[id.Signed[Msg], id.Signed[Effect]] = func(
		ctx context.Context,
		seqNo SeqNo,
		signedReq id.Signed[Msg],
	) (signedResp id.Signed[Effect], err error) {
		if !signedReq.Verify(ctx) {
			return signedResp, fmt.Errorf("signature not valid")
		}
		effect, err := receive(ctx, seqNo, signedReq)
		if err != nil {
			return signedResp, err
		}
		return id.Sign(ctx, receiverPrivCred, effect), nil
	}
	recvOnly := Receive_StageOnly[id.Signed[Msg], id.Signed[Effect]](ctx, receiverCloned.Public.Tree(), senderAddr, senderPublic, topic, receiver)
	msgEffects := make([]MsgEffect[Msg, Effect], len(recvOnly.Result))
	for i, signedMsgEffect := range recvOnly.Result {
		msgEffects[i] = MsgEffect[Msg, Effect]{
			SeqNo:  signedMsgEffect.SeqNo,
			Msg:    signedMsgEffect.Msg.Value,
			Effect: signedMsgEffect.Effect.Value,
		}
	}
	return git.NewChange(
		fmt.Sprintf("Received %d signed messages", len(msgEffects)),
		"receive_signed",
		form.Map{"topic": topic},
		msgEffects,
		form.Forms{recvOnly},
	)
}
