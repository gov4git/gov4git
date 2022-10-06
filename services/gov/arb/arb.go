package arb

import "github.com/gov4git/gov4git/proto"

type GovArbService struct {
	GovConfig      proto.GovConfig
	IdentityConfig proto.IdentityConfig
}
