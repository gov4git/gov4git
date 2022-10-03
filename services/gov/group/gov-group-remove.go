package group

import (
	"context"
	"path/filepath"

	"github.com/petar/gov4git/lib/files"
	"github.com/petar/gov4git/lib/git"
	"github.com/petar/gov4git/proto"
)

type GovGroupRemoveIn struct {
	Name            string `json:"name"`             // community unique handle for this group
	CommunityBranch string `json:"community_branch"` // branch in community repo where group will be added
}

type GovGroupRemoveOut struct{}

func (x GovGroupRemoveOut) Human(context.Context) string {
	return ""
}

func (x GovGroupService) GroupRemove(ctx context.Context, in *GovGroupRemoveIn) (*GovGroupRemoveOut, error) {
	// clone community repo locally
	community := git.LocalFromDir(files.WorkDir(ctx).Subdir("community"))
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.CommunityBranch); err != nil {
		return nil, err
	}
	// make changes to repo
	if err := GovRemoveGroup(ctx, community, in.Name); err != nil {
		return nil, err
	}
	// push to origin
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &GovGroupRemoveOut{}, nil
}

func GovRemoveGroup(ctx context.Context, community git.Local, name string) error {
	groupFile := filepath.Join(proto.GovGroupsDir, name, proto.GovGroupInfoFilebase)
	// remove group file
	if err := community.Dir().Remove(groupFile); err != nil {
		return err
	}
	// stage changes
	if err := community.Remove(ctx, []string{groupFile}); err != nil {
		return err
	}
	// commit changes
	if err := community.Commitf(ctx, "gov: remove group %v", name); err != nil {
		return err
	}
	return nil
}
