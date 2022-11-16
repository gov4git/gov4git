package mail

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gov4git/gov4git/mod/id"
	"github.com/gov4git/lib4git/base"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/git"
)

type RequestResponse[Request form.Form, Response form.Form] struct {
	Request  Request
	Response Response
}

type Responder[Request form.Form, Response form.Form] func(
	ctx context.Context,
	req Request,
) (resp Response, err error)

func ReceiveStageOnly[Request form.Form, Response form.Form](
	ctx context.Context,
	receiver *git.Tree,
	senderAddr id.PublicAddress,
	sender *git.Tree,
	topic string,
	respond Responder[Request, Response],
) git.Change[[]RequestResponse[Request, Response]] {

	// prep
	receiverCred := id.GetPublicCredentials(ctx, receiver)
	senderCred := id.GetPublicCredentials(ctx, sender)
	senderTopicNS := SendTopicNS(receiverCred.ID, topic)
	receiverTopicNS := ReceiveTopicNS(senderCred.ID, topic)
	receiverNextNS := receiverTopicNS.Sub(NextFilebase)
	senderNextNS := senderTopicNS.Sub(NextFilebase)
	receiverInfoNS := receiverTopicNS.Sub(BoxInfoFilebase)

	// read receiver and sender next seq no
	receiverNextSeqNo, _ := git.TryFromFile[SeqNo](ctx, receiver, receiverNextNS.Path())
	senderNextSeqNo, _ := git.TryFromFile[SeqNo](ctx, sender, senderNextNS.Path())

	// write receive box info
	info := ReceiveBoxInfo{SenderID: senderCred.ID, Topic: topic}
	git.ToFileStage(ctx, receiver, receiverInfoNS.Path(), info)

	// read unread messages
	receiverLatestNextSeqNo := receiverNextSeqNo
	base.Infof("receiving receiverSeqNo=%v senderSeqNo=%v", receiverNextSeqNo, senderNextSeqNo)
	rr := []RequestResponse[Request, Response]{}
	for i := receiverNextSeqNo; i < senderNextSeqNo; i++ {
		msgFilebase := strconv.Itoa(int(i))
		req := git.FromFile[Request](ctx, sender, senderTopicNS.Sub(msgFilebase).Path())
		resp, err := respond(ctx, req)
		if err != nil {
			base.Infof("responding to message %d in sender repo (%v)", i, err)
			continue
		}
		git.ToFileStage(ctx, receiver, receiverTopicNS.Sub(msgFilebase).Path(), resp)
		rr = append(rr, RequestResponse[Request, Response]{Request: req, Response: resp})
		receiverLatestNextSeqNo = i + 1
	}

	// write receiver-side next seq no
	git.ToFileStage(ctx, receiver, receiverNextNS.Path(), receiverLatestNextSeqNo)

	return git.Change[[]RequestResponse[Request, Response]]{
		Result: rr,
		Msg:    fmt.Sprintf("Received mail"),
	}
}

type SignedResponder[Request form.Form, Response form.Form] func(
	ctx context.Context,
	req Request,
	signedReq id.SignedPlaintext,
) (resp Response, err error)

func ReceiveSignedStageOnly[Request form.Form, Response form.Form](
	ctx context.Context,
	receiverTree id.OwnerTree,
	senderAddr id.PublicAddress,
	senderPublic *git.Tree,
	topic string,
	respond SignedResponder[Request, Response],
) git.Change[[]RequestResponse[Request, Response]] {

	receiverPrivCred := id.GetOwnerCredentials(ctx, receiverTree)
	rr := []RequestResponse[Request, Response]{}
	signRespond := func(ctx context.Context, signedReq id.SignedPlaintext) (signedResp id.SignedPlaintext, err error) {
		if !signedReq.Verify() {
			return signedResp, fmt.Errorf("signature not valid")
		}
		req, err := form.DecodeBytes[Request](ctx, signedReq.Plaintext)
		if err != nil {
			return signedResp, err
		}
		resp, err := respond(ctx, req, signedReq)
		if err != nil {
			return signedResp, err
		}
		rr = append(rr, RequestResponse[Request, Response]{Request: req, Response: resp})
		return id.Sign(ctx, receiverPrivCred, resp), nil
	}
	ReceiveStageOnly(ctx, receiverTree.Public, senderAddr, senderPublic, topic, signRespond)
	return git.Change[[]RequestResponse[Request, Response]]{
		Result: rr,
		Msg:    fmt.Sprintf("Received signed mail"),
	}
}
