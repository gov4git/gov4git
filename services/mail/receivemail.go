package mail

import (
	"context"
	"path/filepath"
	"strconv"

	"github.com/gov4git/gov4git/lib/base"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/idproto"
)

type ReceiveMailIn struct {
	SenderRepo string `json:"sender_repo"`
	Topic      string `json:"topic"` // mail is sent to a topic and later picked up from the topic by the receiver
}

type ReceiveMailOut struct {
	Messages                []string                  `json:"messages"`
	SenderPublicCredentials idproto.PublicCredentials `json:"sender_public_credentials"`
}

func (x IdentityService) ReceiveMail(ctx context.Context, in *ReceiveMailIn) (*ReceiveMailOut, error) {
	receiver, err := git.CloneBranch(ctx, x.IdentityConfig.PublicURL, idproto.IdentityBranch)
	if err != nil {
		return nil, err
	}
	sender, err := git.CloneBranch(ctx, in.SenderRepo, idproto.IdentityBranch)
	if err != nil {
		return nil, err
	}
	out, err := ReceiveMailLocal(ctx, receiver, sender, in)
	if err != nil {
		return nil, err
	}
	if err := receiver.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return out, nil
}

func ReceiveMailLocal(ctx context.Context, receiver git.Local, sender git.Local, in *ReceiveMailIn) (*ReceiveMailOut, error) {
	out, err := ReceiveMailLocalStageOnly(ctx, receiver, sender, in)
	if err != nil {
		return nil, err
	}
	if err := receiver.Commitf(ctx, "Received mail on topic %v from %v", in.Topic, in.SenderRepo); err != nil {
		return nil, err
	}
	return out, nil
}

func ReceiveMailLocalStageOnly(ctx context.Context, receiver git.Local, sender git.Local, in *ReceiveMailIn) (*ReceiveMailOut, error) {
	var msgs []string
	respondFunc := func(ctx context.Context, req []byte) ([]byte, error) {
		msgs = append(msgs, string(req))
		return req, nil
	}
	senderCred, err := ReceiveRespondMailLocalStageOnly(ctx, receiver, sender, in.Topic, respondFunc)
	if err != nil {
		return nil, err
	}
	return &ReceiveMailOut{Messages: msgs, SenderPublicCredentials: *senderCred}, nil
}

type RespondFunc func(context.Context, []byte) ([]byte, error)

func ReceiveRespondMailLocalStageOnly(
	ctx context.Context,
	receiver git.Local,
	sender git.Local,
	topic string,
	respond RespondFunc,
) (*idproto.PublicCredentials, error) {

	// get receiver credentials
	receiverCred, err := GetPublicCredentialsLocal(ctx, receiver)
	if err != nil {
		return nil, err
	}

	// get sender credentials
	senderCred, err := GetPublicCredentialsLocal(ctx, sender)
	if err != nil {
		return nil, err
	}

	// sender-side
	senderTopicDirpath := idproto.SendMailTopicDirpath(receiverCred.ID, topic)
	senderTopicDir := sender.Dir().Subdir(senderTopicDirpath)

	// receiver-side
	receiverTopicDirpath := idproto.ReceiveMailTopicDirpath(senderCred.ID, topic)
	receiverTopicDir := receiver.Dir().Subdir(receiverTopicDirpath)

	// read receiver-side 'next'
	var receiverNextSeqNo int64
	receiverTopicDir.ReadFormFile(ctx, idproto.NextFilebase, &receiverNextSeqNo)

	// read sender-side 'next'
	var senderNextSeqNo int64
	senderTopicDir.ReadFormFile(ctx, idproto.NextFilebase, &senderNextSeqNo)

	// make dir for receiver-side topic
	if err := receiverTopicDir.Mk(); err != nil {
		return nil, err
	}
	// write receiver id + topic in plaintext file
	info := idproto.ReceiveBoxInfo{SenderID: senderCred.ID, Topic: topic}
	if err := receiverTopicDir.WriteFormFile(ctx, idproto.BoxInfoFilebase, info); err != nil {
		return nil, err
	}
	if err := receiver.Add(ctx, []string{filepath.Join(receiverTopicDirpath, idproto.BoxInfoFilebase)}); err != nil {
		return nil, err
	}

	// read unread messages
	receiverLatestNextSeqNo := receiverNextSeqNo
	base.Infof("r=%v s=%v", receiverNextSeqNo, senderNextSeqNo)
	for i := receiverNextSeqNo; i < senderNextSeqNo; i++ {
		msgFilebase := strconv.Itoa(int(i))
		byteFile, err := senderTopicDir.ReadByteFile(msgFilebase)
		if err != nil {
			base.Infof("reading message %d in sender repo (%v)", i, err)
			continue
		}
		resp, err := respond(ctx, byteFile.Bytes)
		if err != nil {
			base.Infof("responding to message %d in sender repo (%v)", i, err)
			continue
		}
		if err := receiverTopicDir.WriteByteFile(msgFilebase, resp); err != nil {
			return nil, err
		}
		if err := receiver.Add(ctx, []string{filepath.Join(receiverTopicDirpath, msgFilebase)}); err != nil {
			return nil, err
		}
		receiverLatestNextSeqNo = i + 1
	}

	// write receiver-side 'next'
	if err := receiverTopicDir.WriteFormFile(ctx, idproto.NextFilebase, receiverLatestNextSeqNo); err != nil {
		return nil, err
	}

	// stage receiver-side changes
	if err := receiver.Add(ctx, []string{filepath.Join(receiverTopicDirpath, idproto.NextFilebase)}); err != nil {
		return nil, err
	}

	return senderCred, nil
}

func (x IdentityService) ReceiveSignedMail(ctx context.Context, in *ReceiveMailIn) (*ReceiveMailOut, error) {
	receiveOut, err := x.ReceiveMail(ctx, in)
	if err != nil {
		return nil, err
	}
	verifiedMsgs := []string{}
	for _, m := range receiveOut.Messages {
		signed, err := idproto.ParseSignedPlaintext(ctx, []byte(m))
		if err != nil {
			base.Infof("decoding signed message from sender (%v)", err)
			continue
		}
		if !signed.Verify() {
			base.Infof("message signature verification (%v)", err)
			continue
		}
		verifiedMsgs = append(verifiedMsgs, string(signed.Plaintext))
	}
	return &ReceiveMailOut{
		Messages:                verifiedMsgs,
		SenderPublicCredentials: receiveOut.SenderPublicCredentials,
	}, nil
}
