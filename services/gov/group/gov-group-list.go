package group

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"

	"github.com/petar/gitty/lib/files"
	"github.com/petar/gitty/lib/git"
	"github.com/petar/gitty/proto"
)

type GovGroupListIn struct {
	CommunityBranch string `json:"community_branch"` // branch in community repo where group will be added
}

type GovGroupListOut struct {
	Groups []string
}

func (x GovGroupListOut) Human(context.Context) string {
	var w bytes.Buffer
	for _, u := range x.Groups {
		fmt.Fprintln(&w, u)
	}
	return w.String()
}

func (x GovGroupService) GroupList(ctx context.Context, in *GovGroupListIn) (*GovGroupListOut, error) {
	// clone community repo locally
	community := git.LocalFromDir(files.WorkDir(ctx).Subdir("community"))
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.CommunityBranch); err != nil {
		return nil, err
	}
	// make changes to repo
	groups, err := GovGroupList(ctx, community)
	if err != nil {
		return nil, err
	}
	// push to origin
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &GovGroupListOut{Groups: groups}, nil
}

func GovGroupList(ctx context.Context, community git.Local) ([]string, error) {
	groupFileGlob := filepath.Join(proto.GovGroupsDir, "*", proto.GovGroupInfoFilebase)
	// glob for group files
	m, err := community.Dir().Glob(groupFileGlob)
	if err != nil {
		return nil, err
	}
	// extract group names
	groups := make([]string, len(m))
	for i := range m {
		groupDir, _ := filepath.Split(m[i])
		groups[i] = filepath.Base(groupDir)
	}
	return groups, nil
}
