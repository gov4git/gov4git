package mail

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/lib/must"
	"github.com/gov4git/gov4git/mod/id"
)

type SeqNo int64

func Send[M form.Form](
	ctx context.Context,
	sender *git.Tree,
	receiver *git.Tree,
	topic string,
	msg M,
) git.Change[SeqNo] {

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
	git.ToFileStage(ctx, sender, msgNS.Path(), msg)

	// write + stage next file
	var newNextSeqNo SeqNo = nextSeqNo + 1
	if newNextSeqNo < nextSeqNo {
		must.Errorf(ctx, "mailbox size exceeded")
	}
	git.ToFileStage(ctx, sender, nextNS.Path(), newNextSeqNo)

	return git.Change[SeqNo]{
		Result: nextSeqNo,
		Msg:    fmt.Sprintf("Sent mail #%d", nextSeqNo),
	}
}

func SendSigned[M form.Form](
	ctx context.Context,
	senderTree id.OwnerTree,
	receiver *git.Tree,
	topic string,
	msg M,
) git.Change[SeqNo] {
	senderPrivCred := id.GetPrivateCredentials(ctx, senderTree.Private)
	signed := id.Sign(ctx, senderPrivCred, msg)
	return Send(ctx, senderTree.Public, receiver, topic, signed)
}
