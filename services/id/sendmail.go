package id

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/idproto"
)

type SendMailIn struct {
	ReceiverRepo string `json:"receiver_repo"`
	Topic        string `json:"topic"` // mail is sent to a topic and later picked up from the topic by the receiver
	Message      string `json:"message"`
}

type SendMailOut struct {
	SeqNo int64 `json:"seqno"`
}

func (x IdentityService) SendMail(ctx context.Context, in *SendMailIn) (*SendMailOut, error) {
	sender, err := git.CloneBranch(ctx, x.IdentityConfig.PublicURL, idproto.IdentityBranch)
	if err != nil {
		return nil, err
	}
	out, err := x.SendMailLocal(ctx, sender, in)
	if err != nil {
		return nil, err
	}
	if err := sender.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return out, nil
}

func (x IdentityService) SendMailLocal(ctx context.Context, sender git.Local, in *SendMailIn) (*SendMailOut, error) {
	out, err := SendMailLocalStageOnly(ctx, sender, in)
	if err != nil {
		return nil, err
	}
	if err := sender.Commitf(ctx, "Sent mail on topic %v", in.Topic); err != nil {
		return nil, err
	}
	return out, nil
}

func SendMailLocalStageOnly(ctx context.Context, sender git.Local, in *SendMailIn) (*SendMailOut, error) {

	// fetch receiver id
	receiverCred, err := GetPublicCredentials(ctx, in.ReceiverRepo)
	if err != nil {
		return nil, err
	}

	// make outgoing mail directory in sender's repo
	topicDirpath := idproto.SendMailTopicDirpath(receiverCred.ID, in.Topic)
	topicDir := sender.Dir().Subdir(topicDirpath)

	if err := topicDir.Mk(); err != nil {
		return nil, err
	}
	// write receiver id + topic in plaintext file
	info := idproto.SendBoxInfo{ReceiverID: receiverCred.ID, Topic: in.Topic}
	if err := topicDir.WriteFormFile(ctx, idproto.BoxInfoFilebase, info); err != nil {
		return nil, err
	}
	if err := sender.Add(ctx, []string{filepath.Join(topicDirpath, idproto.BoxInfoFilebase)}); err != nil {
		return nil, err
	}

	// read the next message number
	var nextSeqNo int64
	topicDir.ReadFormFile(ctx, idproto.NextFilebase, &nextSeqNo) // if file is missing, nextSeqNo = 0

	// write + stage message
	msgFileBase := strconv.Itoa(int(nextSeqNo))
	if err := topicDir.WriteByteFile(msgFileBase, []byte(in.Message)); err != nil {
		return nil, err
	}
	if err := sender.Add(ctx, []string{filepath.Join(topicDirpath, msgFileBase)}); err != nil {
		return nil, err
	}

	// write + stage next file
	var newNextSeqNo int64 = nextSeqNo + 1
	if newNextSeqNo < nextSeqNo {
		return nil, fmt.Errorf("mailbox size exceeded")
	}
	if err := topicDir.WriteFormFile(ctx, idproto.NextFilebase, newNextSeqNo); err != nil {
		return nil, err
	}
	if err := sender.Add(ctx, []string{filepath.Join(topicDirpath, idproto.NextFilebase)}); err != nil {
		return nil, err
	}
	return &SendMailOut{SeqNo: nextSeqNo}, nil
}

func (x IdentityService) SendSignedMail(ctx context.Context, in *SendMailIn) (*SendMailOut, error) {
	cred, err := x.GetPrivateCredentials(ctx, &GetPrivateCredentialsIn{})
	if err != nil {
		return nil, err
	}
	return x.SendSignedMailWithCredentials(ctx, &cred.PrivateCredentials, in)
}

func (x IdentityService) SendSignedMailWithCredentials(ctx context.Context, priv *idproto.PrivateCredentials, in *SendMailIn) (*SendMailOut, error) {
	signed, err := idproto.SignPlaintext(ctx, priv, []byte(in.Message))
	if err != nil {
		return nil, err
	}
	signedData, err := form.EncodeForm(ctx, signed)
	if err != nil {
		return nil, err
	}
	return x.SendMail(ctx, &SendMailIn{ReceiverRepo: in.ReceiverRepo, Topic: in.Topic, Message: string(signedData)})
}

func (x IdentityService) SendSignedMailLocalStageOnlyWithCredentials(
	ctx context.Context,
	public git.Local,
	priv *idproto.PrivateCredentials,
	in *SendMailIn,
) (*SendMailOut, error) {
	signed, err := idproto.SignPlaintext(ctx, priv, []byte(in.Message))
	if err != nil {
		return nil, err
	}
	signedData, err := form.EncodeForm(ctx, signed)
	if err != nil {
		return nil, err
	}
	return SendMailLocalStageOnly(ctx, public, &SendMailIn{Topic: in.Topic, Message: string(signedData)})
}
