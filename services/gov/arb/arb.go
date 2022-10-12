package arb

import (
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/services/gov"
	"github.com/gov4git/gov4git/services/gov/group"
)

type GovArbService struct {
	GovConfig      proto.GovConfig
	IdentityConfig proto.IdentityConfig
}

func (x GovArbService) GovService() gov.GovService {
	return gov.GovService{GovConfig: x.GovConfig}
}

func (x GovArbService) GroupService() group.GovGroupService {
	return group.GovGroupService{
		GovConfig: x.GovConfig,
	}
}
