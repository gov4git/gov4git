package entity

import (
	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/proto"
)

type EntityService[V form.Form] struct {
	Address   proto.Address
	Namespace string
}
