package bureau

import (
	"github.com/gov4git/gov4git/proto/govproto"
	"github.com/gov4git/gov4git/proto/idproto"
	"github.com/gov4git/gov4git/services/gov"
	"github.com/gov4git/gov4git/services/id"
)

type GovBureauService struct {
	GovConfig      govproto.GovConfig
	IdentityConfig idproto.IdentityConfig
}

func (x GovBureauService) GovService() gov.GovService {
	return gov.GovService{GovConfig: x.GovConfig}
}

func (x GovBureauService) IdentityService() id.IdentityService {
	return id.IdentityService{IdentityConfig: x.IdentityConfig}
}
