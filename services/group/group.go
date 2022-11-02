package group

import (
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/groupproto"
	"github.com/gov4git/gov4git/services/entity"
)

func Service(addr proto.Address) entity.EntityService[groupproto.Group] {
	return entity.EntityService[groupproto.Group]{
		Address:   addr,
		Namespace: groupproto.EntityNamespace,
	}
}
