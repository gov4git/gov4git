package group

import (
	"context"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto"
)

type SetIn struct {
	Name            string `json:"name"`             // community unique handle for this group
	Key             string `json:"key"`              // group property key
	Value           string `json:"value"`            // group property value
	CommunityBranch string `json:"community_branch"` // branch in community repo where group will be added
}

type SetOut struct{}

func (x SetOut) Human(context.Context) string {
	return ""
}

func (x GovGroupService) Set(ctx context.Context, in *SetIn) (*SetOut, error) {
	// clone community repo locally
	community := git.LocalInDir(files.WorkDir(ctx).Subdir("community"))
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.CommunityBranch); err != nil {
		return nil, err
	}
	// make changes to repo
	if err := Set(ctx, community, in.Name, in.Key, in.Value); err != nil {
		return nil, err
	}
	// push to origin
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &SetOut{}, nil
}

// XXX: sanitize key
// XXX: prevent overwrite
func Set(ctx context.Context, community git.Local, name string, key string, value string) error {
	propFile := filepath.Join(proto.GovGroupsDir, name, proto.GovGroupMetaDirbase, key)
	// write group file
	stage := files.ByteFiles{
		files.ByteFile{Path: propFile, Bytes: []byte(value)},
	}
	if err := community.Dir().WriteByteFiles(stage); err != nil {
		return err
	}
	// stage changes
	if err := community.Add(ctx, stage.Paths()); err != nil {
		return err
	}
	// commit changes
	if err := community.Commitf(ctx, "gov: change property %v of group %v", key, name); err != nil {
		return err
	}
	return nil
}