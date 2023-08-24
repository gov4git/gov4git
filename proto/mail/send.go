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

func SendPkgStageOnly(
	ctx context.Context,
	sender *git.Tree,
	receiver *git.Tree,
	topic string,
	msg Pkg,
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
	msg.stagePkg(ctx, sender, msgNS)

	// write + stage next file
	var newNextSeqNo SeqNo = nextSeqNo + 1
	if newNextSeqNo < nextSeqNo {
		must.Errorf(ctx, "mailbox size exceeded")
	}
	git.ToFileStage(ctx, sender, nextNS.Path(), newNextSeqNo)

	return git.NewChange(
		fmt.Sprintf("Sent mail #%d", nextSeqNo),
		"mail_send",
		form.Map{"topic": topic},
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

	return SendPkgStageOnly(ctx, sender, receiver, topic, PkgFile(msg))
}

func SendSignedPkgStageOnly(
	ctx context.Context,
	senderCloned id.OwnerCloned,
	receiver *git.Tree,
	topic string,
	msg Pkg,
) git.Change[form.Map, SeqNo] {

	senderPrivCred := id.GetPrivateCredentials(ctx, senderCloned.Private.Tree())
	signed := id.Sign(ctx, senderPrivCred, msg) // sign package encoded as a single form
	envelope := PkgDir{"signature": PkgFile(signed), "msg": msg}
	sendOnly := SendPkgStageOnly(ctx, senderCloned.Public.Tree(), receiver, topic, envelope)
	return git.NewChange(
		fmt.Sprintf("Sent signed mail #%d", sendOnly.Result),
		"mail_receive_signed",
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

	return SendSignedPkgStageOnly(ctx, senderCloned, receiver, topic, PkgFile(msg))
}
