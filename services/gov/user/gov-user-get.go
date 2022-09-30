package user

import (
	"context"
	"path/filepath"

	"github.com/petar/gitty/lib/files"
	"github.com/petar/gitty/lib/git"
	"github.com/petar/gitty/proto"
)

type GovUserGetIn struct {
	Name            string `json:"name"`             // community unique handle for this user
	Key             string `json:"key"`              // user property key
	CommunityBranch string `json:"community_branch"` // branch in community repo where user will be added
}

type GovUserGetOut struct {
	Value string `json:"value"` // user property value
}

func (x GovUserGetOut) Human(context.Context) string {
	return x.Value
}

func (x GovUserService) UserGet(ctx context.Context, in *GovUserGetIn) (*GovUserGetOut, error) {
	// clone community repo locally
	community := git.LocalFromDir(files.WorkDir(ctx).Subdir("community"))
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.CommunityBranch); err != nil {
		return nil, err
	}
	// read from repo
	value, err := GovUserGet(ctx, community, in.Name, in.Key)
	if err != nil {
		return nil, err
	}
	return &GovUserGetOut{Value: value}, nil
}

func GovUserGet(ctx context.Context, community git.Local, name string, key string) (string, error) {
	propFile := filepath.Join(proto.GovUsersDir, name, proto.GovUserMetaDirbase, key)
	// read user property file
	data, err := community.Dir().ReadByteFile(propFile)
	if err != nil {
		return "", err
	}
	return string(data.Bytes), nil
}
