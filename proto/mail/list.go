package mail

import (
	"context"
	"path/filepath"
	"strconv"

	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func ListSent_Local[Msg form.Form](
	ctx context.Context,
	sender *git.Tree,
	receiver *git.Tree,
	topic string,
) (SentMsgs[Msg], map[SeqNo]Msg) {

	// prep
	receiverCred := id.GetPublicCredentials(ctx, receiver)
	senderTopicNS := SendTopicNS(receiverCred.ID, topic)

	// read all sent messages (at the sender)
	sentMsgs := SentMsgs[Msg]{}
	seqnoToMsg := map[SeqNo]Msg{}
	senderInfos, err := git.TreeReadDir(ctx, sender, senderTopicNS)
	must.NoError(ctx, err)
	for _, info := range senderInfos {
		if info.IsDir() {
			continue
		}
		n := info.Name()
		if filepath.Ext(n) != ".json" {
			continue
		}
		seqno, err := strconv.Atoi(n[:len(n)-len(".json")])
		if err != nil {
			continue
		}
		msg := git.FromFile[Msg](ctx, sender, senderTopicNS.Append(n))
		seqnoToMsg[SeqNo(seqno)] = msg
		sentMsgs = append(sentMsgs, SentMsg[Msg]{SeqNo: SeqNo(seqno), Msg: msg})
	}
	sentMsgs.Sort()

	return sentMsgs, seqnoToMsg
}

func ListReceived_Local[Msg form.Form, Effect form.Form](
	ctx context.Context,
	sender *git.Tree,
	receiver *git.Tree,
	topic string,
) (MsgEffects[Msg, Effect], map[SeqNo]MsgEffect[Msg, Effect]) {

	// prep
	senderCred := id.GetPublicCredentials(ctx, sender)
	receiverTopicNS := ReceiveTopicNS(senderCred.ID, topic)

	// read all received messages and the resulting effects (at the receiver)
	seqnoToMsgEffect := map[SeqNo]MsgEffect[Msg, Effect]{}
	msgEffects := MsgEffects[Msg, Effect]{}

	receiverInfos, err := git.TreeReadDir(ctx, receiver, receiverTopicNS)
	must.NoError(ctx, err)
	for _, info := range receiverInfos {
		if info.IsDir() {
			continue
		}
		n := info.Name()
		if filepath.Ext(n) != ".json" {
			continue
		}
		seqno, err := strconv.Atoi(n[:len(n)-len(".json")])
		if err != nil {
			continue
		}
		msgEffect := git.FromFile[MsgEffect[Msg, Effect]](ctx, receiver, receiverTopicNS.Append(n))
		must.Assertf(ctx, msgEffect.SeqNo == SeqNo(seqno), "receiver mailbox inconsistent")
		seqnoToMsgEffect[SeqNo(seqno)] = msgEffect
		msgEffects = append(msgEffects, msgEffect)
	}
	msgEffects.Sort()

	return msgEffects, seqnoToMsgEffect
}
