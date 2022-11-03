package userproto

import "github.com/gov4git/gov4git/proto"

const EntityNamespace = "user"

type User struct {
	Address proto.Address `json:"address"`
}
