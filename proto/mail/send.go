package mail

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
	"github.com/gov4git/lib4git/must"
)

func SendMakeMsg_StageOnly[Msg form.Form](
	ctx context.Context,
	sender *git.Tree,
	receiver *git.Tree,
	topic string,
	mkMsg func(context.Context, SeqNo) Msg,
) git.Change[form.Map, SentMsg[Msg]] {

	// fetch receiver id
	receiverCred := id.GetPublicCredentials(ctx, receiver)
	topicNS := SendTopicNS(receiverCred.ID, topic)

	// write receiver id + topic in send box file
	infoValue := SendBoxInfo{ReceiverCred: receiverCred, Topic: topic}
	infoNS := topicNS.Append(BoxInfoFilebase)
	git.ToFileStage(ctx, sender, infoNS, infoValue)

	// read the next message number
	// if file is missing, nextSeqNo = 0
	nextNS := topicNS.Append(NextFilebase)
	nextSeqNo, _ := git.TryFromFile[SeqNo](ctx, sender, nextNS)

	// write message
	msgNS := topicNS.Append(strconv.Itoa(int(nextSeqNo)))
	msg := mkMsg(ctx, nextSeqNo)
	git.ToFileStage(ctx, sender, msgNS, msg)

	// write + stage next file
	var newNextSeqNo SeqNo = nextSeqNo + 1
	if newNextSeqNo < nextSeqNo {
		must.Errorf(ctx, "mailbox size exceeded")
	}
	git.ToFileStage(ctx, sender, nextNS, newNextSeqNo)

	return git.NewChange(
		fmt.Sprintf("Sent #%d", nextSeqNo),
		"send",
		form.Map{"topic": topic, "msg": msg},
		SentMsg[Msg]{SeqNo: nextSeqNo, Msg: msg},
		nil,
	)
}

func Send_StageOnly[Msg form.Form](
	ctx context.Context,
	sender *git.Tree,
	receiver *git.Tree,
	topic string,
	msg Msg,
) git.Change[form.Map, SentMsg[Msg]] {

	return SendMakeMsg_StageOnly[Msg](ctx, sender, receiver, topic, func(context.Context, SeqNo) Msg { return msg })
}

func SendSignedMakeMsg_StageOnly[Msg form.Form](
	ctx context.Context,
	senderCloned id.OwnerCloned,
	receiver *git.Tree,
	topic string,
	mkMsg func(context.Context, SeqNo) Msg,
) git.Change[form.Map, SentMsg[Msg]] {

	senderPrivCred := id.GetPrivateCredentials(ctx, senderCloned.Private.Tree())
	mkSignedMsg := func(ctx context.Context, seqNo SeqNo) id.Signed[Msg] {
		return id.Sign(ctx, senderPrivCred, mkMsg(ctx, seqNo))
	}
	sendOnly := SendMakeMsg_StageOnly[id.Signed[Msg]](ctx, senderCloned.Public.Tree(), receiver, topic, mkSignedMsg)
	return git.NewChange(
		fmt.Sprintf("Sent signed #%d", sendOnly.Result.SeqNo),
		"send_signed",
		form.Map{"topic": topic},
		SentMsg[Msg]{SeqNo: sendOnly.Result.SeqNo, Msg: sendOnly.Result.Msg.Value},
		form.Forms{sendOnly},
	)
}

func SendSigned_StageOnly[Msg form.Form](
	ctx context.Context,
	senderCloned id.OwnerCloned,
	receiver *git.Tree,
	topic string,
	msg Msg,
) git.Change[form.Map, SentMsg[Msg]] {

	return SendSignedMakeMsg_StageOnly[Msg](ctx, senderCloned, receiver, topic, func(context.Context, SeqNo) Msg { return msg })
}
