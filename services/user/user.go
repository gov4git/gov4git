package user

import (
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/userproto"
	"github.com/gov4git/gov4git/services/entity"
)

func Service(addr proto.Address) entity.EntityService[userproto.User] {
	return entity.EntityService[userproto.User]{
		Address:   addr,
		Namespace: userproto.EntityNamespace,
	}
}
