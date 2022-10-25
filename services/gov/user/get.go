package user

import (
	"context"
	"path/filepath"
	"strconv"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/govproto"
)

type GetIn struct {
	Name   string `json:"name"`   // community unique handle for this user
	Key    string `json:"key"`    // user property key
	Branch string `json:"branch"` // branch in community repo where user will be added
}

type GetOut struct {
	Value string `json:"value"` // user property value
}

func (x GovUserService) Get(ctx context.Context, in *GetIn) (*GetOut, error) {
	community, err := git.CloneBranch(ctx, x.GovConfig.CommunityURL, in.Branch)
	if err != nil {
		return nil, err
	}
	value, err := x.GetLocal(ctx, community, in.Name, in.Key)
	if err != nil {
		return nil, err
	}
	return &GetOut{Value: value}, nil
}

func (x GovUserService) GetLocal(ctx context.Context, community git.Local, name string, key string) (string, error) {
	propFile := filepath.Join(govproto.GovUsersDir, name, govproto.GovUserMetaDirbase, key)
	// read user property file
	data, err := community.Dir().ReadByteFile(propFile)
	if err != nil {
		return "", err
	}
	return string(data.Bytes), nil
}

func (x GovUserService) GetFloat64Local(ctx context.Context, community git.Local, name string, key string) (float64, error) {
	v, err := x.GetLocal(ctx, community, name, key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(v, 64)
}
