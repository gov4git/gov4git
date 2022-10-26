package arb

import (
	"github.com/gov4git/gov4git/proto/govproto"
	"github.com/gov4git/gov4git/proto/idproto"
	"github.com/gov4git/gov4git/services/gov"
	"github.com/gov4git/gov4git/services/gov/group"
)

type GovArbService struct {
	GovConfig      govproto.GovConfig
	IdentityConfig idproto.IdentityConfig
}

func (x GovArbService) GovService() gov.GovService {
	return gov.GovService{GovConfig: x.GovConfig}
}

func (x GovArbService) GroupService() group.GovGroupService {
	return group.GovGroupService{
		GovConfig: x.GovConfig,
	}
}
