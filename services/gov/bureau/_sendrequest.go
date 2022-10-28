package bureau

import (
	"context"

	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/proto/govproto/bureauproto"
	"github.com/gov4git/gov4git/proto/idproto"
	"github.com/gov4git/gov4git/services/id"
)

func (x GovBureauService) Identify(ctx context.Context, in *bureauproto.IdentifyRequest) (*bureauproto.IdentifyResponse, error) {
	resp, err := x.SendRequest(ctx, &bureauproto.Request{Identify: in})
	if err != nil {
		return nil, err
	}
	return resp.Identify, nil
}

func (x GovBureauService) SendRequest(ctx context.Context, req *bureauproto.Request) (*bureauproto.Response, error) {
	// get private credentials
	cred, err := x.IdentityService().GetPrivateCredentials(ctx, &id.GetPrivateCredentialsIn{})
	if err != nil {
		return nil, err
	}

	// prepare a signed request
	reqData, err := form.EncodeForm(ctx, req)
	if err != nil {
		return nil, err
	}

	signed, err := idproto.SignPlaintext(ctx, &cred.PrivateCredentials, reqData)
	if err != nil {
		return nil, err
	}
	signedData, err := form.EncodeForm(ctx, signed)
	if err != nil {
		return nil, err
	}

	//
	_, err = x.IdentityService().SendSignedMailWithCredentials(ctx,
		&cred.PrivateCredentials,
		&id.SendMailIn{Topic: XXX, Message: string(signedData)})
	if err != nil {
		return nil, err
	}

	panic("")
}
