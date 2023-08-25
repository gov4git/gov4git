package mail

import (
	"context"
	"strconv"

	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func Confirm[Msg form.Form, Effect form.Form](
	ctx context.Context,
	senderAddr id.PublicAddress,
	receiverAddr id.PublicAddress,
	topic string,
) (confirmed MsgEffects[Msg, Effect], notConfirmed MsgEffects[Msg, form.None]) {

	return Confirm_Local[Msg, Effect](
		ctx,
		git.CloneOne(ctx, git.Address(senderAddr)).Tree(),
		git.CloneOne(ctx, git.Address(receiverAddr)).Tree(),
		topic,
	)
}

func Confirm_Local[Msg form.Form, Effect form.Form](
	ctx context.Context,
	sender *git.Tree,
	receiver *git.Tree,
	topic string,
) (confirmed MsgEffects[Msg, Effect], notConfirmed MsgEffects[Msg, form.None]) {

	// prep
	senderCred := id.GetPublicCredentials(ctx, sender)
	receiverCred := id.GetPublicCredentials(ctx, receiver)
	senderTopicNS := SendTopicNS(receiverCred.ID, topic)
	receiverTopicNS := ReceiveTopicNS(senderCred.ID, topic)

	// read all sent messages (at the sender)
	sentMsgs := map[SeqNo]Msg{}
	senderInfos, err := sender.Filesystem.ReadDir(senderTopicNS.Path())
	must.NoError(ctx, err)
	for _, info := range senderInfos {
		if info.IsDir() {
			continue
		}
		seqno, err := strconv.Atoi(info.Name())
		if err != nil {
			continue
		}
		sentMsgs[SeqNo(seqno)] = git.FromFile[Msg](ctx, sender, senderTopicNS.Sub(info.Name()).Path())
	}

	// read all received messages and the resulting effects (at the receiver)
	receivedMsgEffects := map[SeqNo]MsgEffect[Msg, Effect]{}
	receiverInfos, err := receiver.Filesystem.ReadDir(receiverTopicNS.Path())
	must.NoError(ctx, err)
	for _, info := range receiverInfos {
		if info.IsDir() {
			continue
		}
		seqno, err := strconv.Atoi(info.Name())
		if err != nil {
			continue
		}
		msgEffect := git.FromFile[MsgEffect[Msg, Effect]](ctx, receiver, receiverTopicNS.Sub(info.Name()).Path())
		must.Assertf(ctx, msgEffect.SeqNo == SeqNo(seqno), "receiver mailbox inconsistent")
		receivedMsgEffects[SeqNo(seqno)] = msgEffect
	}

	// compute confirmed and not confirmed transmissions
	for seqno, sentMsg := range sentMsgs {
		if receivedMsgEffect, ok := receivedMsgEffects[seqno]; ok {
			confirmed = append(confirmed,
				MsgEffect[Msg, Effect]{SeqNo: seqno, Msg: sentMsg, Effect: receivedMsgEffect.Effect},
			)
		} else {
			notConfirmed = append(notConfirmed,
				MsgEffect[Msg, form.None]{SeqNo: seqno, Msg: sentMsg},
			)
		}
	}
	confirmed.Sort()
	notConfirmed.Sort()

	return confirmed, notConfirmed
}

func ConfirmSigned[Msg form.Form, Effect form.Form](
	ctx context.Context,
	senderAddr id.PublicAddress,
	receiverAddr id.PublicAddress,
	topic string,
) (confirmed MsgEffects[id.Signed[Msg], id.Signed[Effect]], notConfirmed MsgEffects[id.Signed[Msg], form.None]) {

	return ConfirmSigned_Local[Msg, Effect](
		ctx,
		git.CloneOne(ctx, git.Address(senderAddr)).Tree(),
		git.CloneOne(ctx, git.Address(receiverAddr)).Tree(),
		topic,
	)
}

func ConfirmSigned_Local[Msg form.Form, Effect form.Form](
	ctx context.Context,
	sender *git.Tree,
	receiver *git.Tree,
	topic string,
) (confirmed MsgEffects[id.Signed[Msg], id.Signed[Effect]], notConfirmed MsgEffects[id.Signed[Msg], form.None]) {

	return Confirm_Local[id.Signed[Msg], id.Signed[Effect]](ctx, sender, receiver, topic)
}
