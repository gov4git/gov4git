package bureau

import (
	"context"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/proto/govproto/bureauproto"
	"github.com/gov4git/gov4git/proto/idproto"
	"github.com/gov4git/gov4git/services/id"
)

func (x GovBureauService) CallIdentify(ctx context.Context, in *bureauproto.IdentifyRequest) (seqNo int64, err error) {
	return x.SendRequest(ctx, &bureauproto.Request{Identify: in})
}

func (x GovBureauService) SendRequest(ctx context.Context, req *bureauproto.Request) (seqNo int64, err error) {
	// get private credentials
	cred, err := x.IdentityService().GetPrivateCredentials(ctx, &id.GetPrivateCredentialsIn{})
	if err != nil {
		return -1, err
	}

	// prepare a signed request
	reqData, err := form.EncodeForm(ctx, req)
	if err != nil {
		return -1, err
	}
	signed, err := idproto.SignPlaintext(ctx, &cred.PrivateCredentials, reqData)
	if err != nil {
		return -1, err
	}
	signedData, err := form.EncodeForm(ctx, signed)
	if err != nil {
		return -1, err
	}

	// send the request out
	sendOut, err := x.IdentityService().SendSignedMailWithCredentials(ctx,
		&cred.PrivateCredentials,
		&id.SendMailIn{Topic: XXX, Message: string(signedData)})
	if err != nil {
		return -1, err
	}

	return sendOut.SeqNo, nil
}
