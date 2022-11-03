package bureau

import (
	"context"

	"github.com/gov4git/gov4git/proto/bureauproto"
	"github.com/gov4git/gov4git/proto/govproto"
	"github.com/gov4git/gov4git/proto/idproto"
	"github.com/gov4git/gov4git/services/id"
)

type BureauClient struct {
	GovConfig      govproto.GovConfig
	IdentityConfig idproto.IdentityConfig
}

func (x BureauClient) IdentityService() id.IdentityService {
	return id.IdentityService{IdentityConfig: x.IdentityConfig}
}

func (x BureauClient) SendRequest(ctx context.Context, topic string, reqData []byte) (seqNo int64, err error) {
	sendOut, err := x.IdentityService().SendSignedMail(ctx,
		&id.SendMailIn{
			ReceiverRepo: x.GovConfig.CommunityURL,
			Topic:        bureauproto.Topic(topic),
			Message:      string(reqData),
		})
	if err != nil {
		return -1, err
	}
	return sendOut.SeqNo, nil
}
