package group

import (
	"context"
	"path/filepath"

	"github.com/petar/gov4git/lib/files"
	"github.com/petar/gov4git/lib/git"
	"github.com/petar/gov4git/proto"
)

type GovGroupGetIn struct {
	Name            string `json:"name"`             // community unique handle for this group
	Key             string `json:"key"`              // group property key
	CommunityBranch string `json:"community_branch"` // branch in community repo where group will be added
}

type GovGroupGetOut struct {
	Value string `json:"value"` // group property value
}

func (x GovGroupGetOut) Human(context.Context) string {
	return x.Value
}

func (x GovGroupService) GroupGet(ctx context.Context, in *GovGroupGetIn) (*GovGroupGetOut, error) {
	// clone community repo locally
	community := git.LocalFromDir(files.WorkDir(ctx).Subdir("community"))
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.CommunityBranch); err != nil {
		return nil, err
	}
	// read from repo
	value, err := GovGroupGet(ctx, community, in.Name, in.Key)
	if err != nil {
		return nil, err
	}
	return &GovGroupGetOut{Value: value}, nil
}

func GovGroupGet(ctx context.Context, community git.Local, name string, key string) (string, error) {
	propFile := filepath.Join(proto.GovGroupsDir, name, proto.GovGroupMetaDirbase, key)
	// read group property file
	data, err := community.Dir().ReadByteFile(propFile)
	if err != nil {
		return "", err
	}
	return string(data.Bytes), nil
}
