package mail

import (
	"context"

	"github.com/gov4git/gov4git/v2/proto/id"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
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

	_, seqnoToSentMsg := ListSent_Local[Msg](ctx, sender, receiver, topic)
	_, seqnoToReceivedMsgEffect := ListReceived_Local[Msg, Effect](ctx, sender, receiver, topic)

	// compute confirmed and not confirmed transmissions
	for seqno, sentMsg := range seqnoToSentMsg {
		if receivedMsgEffect, ok := seqnoToReceivedMsgEffect[seqno]; ok {
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
) (confirmed MsgEffects[Msg, Effect], notConfirmed MsgEffects[Msg, form.None]) {

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
) (confirmed MsgEffects[Msg, Effect], notConfirmed MsgEffects[Msg, form.None]) {

	signedConfirmed, signedNotConfirmed := Confirm_Local[id.Signed[Msg], id.Signed[Effect]](ctx, sender, receiver, topic)

	confirmed = make(MsgEffects[Msg, Effect], len(signedConfirmed))
	for i, s := range signedConfirmed {
		confirmed[i] = MsgEffect[Msg, Effect]{SeqNo: s.SeqNo, Msg: s.Msg.Value, Effect: s.Effect.Value}
	}

	notConfirmed = make(MsgEffects[Msg, form.None], len(signedNotConfirmed))
	for i, s := range signedNotConfirmed {
		notConfirmed[i] = MsgEffect[Msg, form.None]{SeqNo: s.SeqNo, Msg: s.Msg.Value}
	}

	return
}

func ConfirmCall[Req form.Form, Resp form.Form](
	ctx context.Context,
	senderAddr id.PublicAddress,
	receiverAddr id.PublicAddress,
	topic string,
) (confirmed MsgEffects[Req, Resp], notConfirmed MsgEffects[Req, form.None]) {

	return ConfirmCall_Local[Req, Resp](
		ctx,
		git.CloneOne(ctx, git.Address(senderAddr)).Tree(),
		git.CloneOne(ctx, git.Address(receiverAddr)).Tree(),
		topic,
	)
}

func ConfirmCall_Local[Req form.Form, Resp form.Form](
	ctx context.Context,
	sender *git.Tree,
	receiver *git.Tree,
	topic string,
) (confirmed MsgEffects[Req, Resp], notConfirmed MsgEffects[Req, form.None]) {

	envConfirmed, envNotConfirmed := ConfirmSigned_Local[RequestEnvelope[Req], ResponseEnvelope[Resp]](ctx, sender, receiver, topic)

	confirmed = make(MsgEffects[Req, Resp], len(envConfirmed))
	for i, s := range envConfirmed {
		confirmed[i] = MsgEffect[Req, Resp]{SeqNo: s.SeqNo, Msg: s.Msg.Request, Effect: s.Effect.Response}
	}

	notConfirmed = make(MsgEffects[Req, form.None], len(envNotConfirmed))
	for i, s := range envNotConfirmed {
		notConfirmed[i] = MsgEffect[Req, form.None]{SeqNo: s.SeqNo, Msg: s.Msg.Request}
	}

	return
}
