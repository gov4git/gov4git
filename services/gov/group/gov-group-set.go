package group

import (
	"context"
	"path/filepath"

	"github.com/petar/gitty/lib/files"
	"github.com/petar/gitty/lib/git"
	"github.com/petar/gitty/proto"
)

type GovGroupSetIn struct {
	Name            string `json:"name"`             // community unique handle for this group
	Key             string `json:"key"`              // group property key
	Value           string `json:"value"`            // group property value
	CommunityBranch string `json:"community_branch"` // branch in community repo where group will be added
}

type GovGroupSetOut struct{}

func (x GovGroupSetOut) Human(context.Context) string {
	return ""
}

func (x GovGroupService) GroupSet(ctx context.Context, in *GovGroupSetIn) (*GovGroupSetOut, error) {
	// clone community repo locally
	community := git.LocalFromDir(files.WorkDir(ctx).Subdir("community"))
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.CommunityBranch); err != nil {
		return nil, err
	}
	// make changes to repo
	if err := GovGroupSet(ctx, community, in.Name, in.Key, in.Value); err != nil {
		return nil, err
	}
	// push to origin
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &GovGroupSetOut{}, nil
}

// XXX: sanitize key
// XXX: prevent overwrite
func GovGroupSet(ctx context.Context, community git.Local, name string, key string, value string) error {
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
