package group

import (
	"context"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/govproto"
)

type GetIn struct {
	Name            string `json:"name"`             // community unique handle for this group
	Key             string `json:"key"`              // group property key
	CommunityBranch string `json:"community_branch"` // branch in community repo where group will be added
}

type GetOut struct {
	Value string `json:"value"` // group property value
}

func (x GovGroupService) Get(ctx context.Context, in *GetIn) (*GetOut, error) {
	// clone community repo locally
	community, err := git.MakeLocalInCtx(ctx, "community")
	if err != nil {
		return nil, err
	}
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.CommunityBranch); err != nil {
		return nil, err
	}
	// read from repo
	value, err := x.GetLocal(ctx, community, in.Name, in.Key)
	if err != nil {
		return nil, err
	}
	return &GetOut{Value: value}, nil
}

func (x GovGroupService) GetLocal(ctx context.Context, community git.Local, name string, key string) (string, error) {
	propFile := filepath.Join(govproto.GovGroupsDir, name, govproto.GovGroupMetaDirbase, key)
	// read group property file
	data, err := community.Dir().ReadByteFile(propFile)
	if err != nil {
		return "", err
	}
	return string(data.Bytes), nil
}
