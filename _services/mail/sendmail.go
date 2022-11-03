package mail

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/idproto"
	"github.com/gov4git/gov4git/services/id"
)

type SeqNo int64

func SendMail(
	ctx context.Context,
	senderAddr proto.PairAddress,
	receiverAddr proto.Address,
	topic string,
	msg string,
) (SeqNo, error) {
	sender, err := git.CloneOrigin(ctx, senderAddr.PublicAddress)
	if err != nil {
		return 0, err
	}
	out, err := SendMailLocal(ctx, sender, receiverAddr, topic, msg)
	if err != nil {
		return 0, err
	}
	if err := sender.PushUpstream(ctx); err != nil {
		return 0, err
	}
	return out, nil
}

func SendMailLocal(
	ctx context.Context,
	sender git.Local,
	receiverAddr proto.Address,
	topic string,
	msg string,
) (SeqNo, error) {
	out, err := SendMailLocalStageOnly(ctx, sender, receiverAddr, topic, msg)
	if err != nil {
		return 0, err
	}
	if err := sender.Commitf(ctx, "Sent mail on topic %v", topic); err != nil {
		return 0, err
	}
	return out, nil
}

func SendMailLocalStageOnly(
	ctx context.Context,
	sender git.Local,
	receiverAddr proto.Address,
	topic string,
	msg string,
) (SeqNo, error) {

	// fetch receiver id
	receiverCred, err := id.IdentityPublicService(receiverAddr).GetPublicCredentials(ctx)
	if err != nil {
		return 0, err
	}

	// make outgoing mail directory in sender's repo
	topicDirpath := idproto.SendMailTopicDirpath(receiverCred.ID, topic)
	topicDir := sender.Dir().Subdir(topicDirpath)

	if err := topicDir.Mk(); err != nil {
		return 0, err
	}
	// write receiver id + topic in plaintext file
	info := idproto.SendBoxInfo{ReceiverID: receiverCred.ID, Topic: topic}
	if err := topicDir.WriteFormFile(ctx, idproto.BoxInfoFilebase, info); err != nil {
		return 0, err
	}
	if err := sender.Add(ctx, []string{filepath.Join(topicDirpath, idproto.BoxInfoFilebase)}); err != nil {
		return 0, err
	}

	// read the next message number
	var nextSeqNo int64
	topicDir.ReadFormFile(ctx, idproto.NextFilebase, &nextSeqNo) // if file is missing, nextSeqNo = 0

	// write + stage message
	msgFileBase := strconv.Itoa(int(nextSeqNo))
	if err := topicDir.WriteByteFile(msgFileBase, []byte(msg)); err != nil {
		return 0, err
	}
	if err := sender.Add(ctx, []string{filepath.Join(topicDirpath, msgFileBase)}); err != nil {
		return 0, err
	}

	// write + stage next file
	var newNextSeqNo int64 = nextSeqNo + 1
	if newNextSeqNo < nextSeqNo {
		return 0, fmt.Errorf("mailbox size exceeded")
	}
	if err := topicDir.WriteFormFile(ctx, idproto.NextFilebase, newNextSeqNo); err != nil {
		return 0, err
	}
	if err := sender.Add(ctx, []string{filepath.Join(topicDirpath, idproto.NextFilebase)}); err != nil {
		return 0, err
	}
	return SeqNo(nextSeqNo), nil
}

func SendSignedMail(
	ctx context.Context,
	senderAddr proto.PairAddress,
	receiverAddr proto.Address,
	topic string,
	msg string,
) (SeqNo, error) {
	cred, err := id.IdentityPrivateService(senderAddr).GetPrivateCredentials(ctx)
	if err != nil {
		return 0, err
	}
	return SendSignedMailWithCredentials(ctx, cred, senderAddr, receiverAddr, topic, msg)
}

func SendSignedMailWithCredentials(
	ctx context.Context,
	senderCred *idproto.PrivateCredentials,
	senderAddr proto.PairAddress,
	receiverAddr proto.Address,
	topic string,
	msg string,
) (SeqNo, error) {
	signed, err := idproto.SignPlaintext(ctx, senderCred, []byte(msg))
	if err != nil {
		return 0, err
	}
	signedData, err := form.EncodeForm(ctx, signed)
	if err != nil {
		return 0, err
	}
	return SendMail(ctx, senderAddr, receiverAddr, topic, string(signedData))
}

func SendSignedMailLocalStageOnlyWithCredentials(
	ctx context.Context,
	sender git.Local,
	senderCred *idproto.PrivateCredentials,
	receiverAddr proto.Address,
	topic string,
	msg string,
) (SeqNo, error) {

	signed, err := idproto.SignPlaintext(ctx, senderCred, []byte(msg))
	if err != nil {
		return 0, err
	}
	signedData, err := form.EncodeForm(ctx, signed)
	if err != nil {
		return 0, err
	}
	return SendMailLocalStageOnly(ctx, sender, receiverAddr, topic, string(signedData))
}
