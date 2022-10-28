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
	Topic   string `json:"topic"` // mail is sent to a topic and later picked up from the topic by the receiver
	Message string `json:"message"`
}

type SendMailOut struct {
	SeqNo int64 `json:"seqno"`
}

func (x IdentityService) SendMail(ctx context.Context, in *SendMailIn) (*SendMailOut, error) {
	public, err := git.CloneBranch(ctx, x.IdentityConfig.PublicURL, idproto.IdentityBranch)
	if err != nil {
		return nil, err
	}
	out, err := x.SendMailLocal(ctx, public, in)
	if err != nil {
		return nil, err
	}
	if err := public.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return out, nil
}

func (x IdentityService) SendMailLocal(ctx context.Context, public git.Local, in *SendMailIn) (*SendMailOut, error) {
	out, err := x.SendMailLocalStageOnly(ctx, public, in)
	if err != nil {
		return nil, err
	}
	if err := public.Commitf(ctx, "Sent mail on topic %v", in.Topic); err != nil {
		return nil, err
	}
	return out, nil
}

const nextFilebase = "next"

func (x IdentityService) SendMailLocalStageOnly(ctx context.Context, public git.Local, in *SendMailIn) (*SendMailOut, error) {
	topicDirpath := idproto.SendMailTopicDirpath(in.Topic)
	topicDir := public.Dir().Subdir(topicDirpath)

	if err := topicDir.Mk(); err != nil {
		return nil, err
	}

	var nextSeqNo int64
	topicDir.ReadFormFile(ctx, nextFilebase, &nextSeqNo) // if file is missing, nextSeqNo = 0

	// write + stage message
	msgFileBase := strconv.Itoa(int(nextSeqNo))
	if err := topicDir.WriteByteFile(msgFileBase, []byte(in.Message)); err != nil {
		return nil, err
	}
	if err := public.Add(ctx, []string{filepath.Join(topicDirpath, msgFileBase)}); err != nil {
		return nil, err
	}

	// write + stage next file
	var newNextSeqNo int64 = nextSeqNo + 1
	if newNextSeqNo < nextSeqNo {
		return nil, fmt.Errorf("mailbox size exceeded")
	}
	if err := topicDir.WriteFormFile(ctx, nextFilebase, newNextSeqNo); err != nil {
		return nil, err
	}
	if err := public.Add(ctx, []string{filepath.Join(topicDirpath, nextFilebase)}); err != nil {
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
	return x.SendMail(ctx, &SendMailIn{Topic: in.Topic, Message: string(signedData)})
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
	return x.SendMailLocalStageOnly(ctx, public, &SendMailIn{Topic: in.Topic, Message: string(signedData)})
}
