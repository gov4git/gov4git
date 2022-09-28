package services

import (
	"context"

	"github.com/petar/gitty/proto"
)

type GovService struct {
	GovConfig proto.GovConfig
}

// XXX: output a result object, which has a human-readable printing

type GovAddUserIn struct {
	UserURL string `json:"user_url"`
}

type GovAddUserOut struct {
	//XXX
}

func (x GovAddUserOut) Human() string {
	return "XXX"
}

func (x GovService) AddUser(ctx context.Context, in *GovAddUserIn) (*GovAddUserOut, error) {
	panic("u")
}
