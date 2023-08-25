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

type SeqNo int64

func SendMakeMsgStageOnly[M form.Form](
	ctx context.Context,
	sender *git.Tree,
	receiver *git.Tree,
	topic string,
	mkMsg func(context.Context, SeqNo) M,
) git.Change[form.Map, SeqNo] {

	// fetch receiver id
	receiverCred := id.GetPublicCredentials(ctx, receiver)
	topicNS := SendTopicNS(receiverCred.ID, topic)

	// write receiver id + topic in send box file
	infoValue := SendBoxInfo{ReceiverID: receiverCred.ID, Topic: topic}
	infoNS := topicNS.Sub(BoxInfoFilebase)
	git.ToFileStage(ctx, sender, infoNS.Path(), infoValue)

	// read the next message number
	// if file is missing, nextSeqNo = 0
	nextNS := topicNS.Sub(NextFilebase)
	nextSeqNo, _ := git.TryFromFile[SeqNo](ctx, sender, nextNS.Path())

	// write message
	msgNS := topicNS.Sub(strconv.Itoa(int(nextSeqNo)))
	msg := mkMsg(ctx, nextSeqNo)
	git.ToFileStage(ctx, sender, msgNS.Path(), msg)

	// write + stage next file
	var newNextSeqNo SeqNo = nextSeqNo + 1
	if newNextSeqNo < nextSeqNo {
		must.Errorf(ctx, "mailbox size exceeded")
	}
	git.ToFileStage(ctx, sender, nextNS.Path(), newNextSeqNo)

	return git.NewChange(
		fmt.Sprintf("Sent #%d", nextSeqNo),
		"send",
		form.Map{"topic": topic, "msg": msg},
		nextSeqNo,
		nil,
	)
}

func SendStageOnly[M form.Form](
	ctx context.Context,
	sender *git.Tree,
	receiver *git.Tree,
	topic string,
	msg M,
) git.Change[form.Map, SeqNo] {

	return SendMakeMsgStageOnly(ctx, sender, receiver, topic, func(context.Context, SeqNo) M { return msg })
}

func SendSignedMakeMsgStageOnly[M form.Form](
	ctx context.Context,
	senderCloned id.OwnerCloned,
	receiver *git.Tree,
	topic string,
	mkMsg func(context.Context, SeqNo) M,
) git.Change[form.Map, SeqNo] {

	senderPrivCred := id.GetPrivateCredentials(ctx, senderCloned.Private.Tree())
	mkSignedMsg := func(ctx context.Context, seqNo SeqNo) id.SignedPlaintext {
		return id.Sign(ctx, senderPrivCred, mkMsg(ctx, seqNo))
	}
	sendOnly := SendMakeMsgStageOnly(ctx, senderCloned.Public.Tree(), receiver, topic, mkSignedMsg)
	return git.NewChange(
		fmt.Sprintf("Sent signed #%d", sendOnly.Result),
		"send_signed",
		form.Map{"topic": topic},
		sendOnly.Result,
		form.Forms{sendOnly},
	)
}

func SendSignedStageOnly[M form.Form](
	ctx context.Context,
	senderCloned id.OwnerCloned,
	receiver *git.Tree,
	topic string,
	msg M,
) git.Change[form.Map, SeqNo] {

	return SendSignedMakeMsgStageOnly(ctx, senderCloned, receiver, topic, func(context.Context, SeqNo) M { return msg })
}
