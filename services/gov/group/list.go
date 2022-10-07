package group

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto"
)

type ListIn struct {
	CommunityBranch string `json:"community_branch"` // branch in community repo where group will be added
}

type ListOut struct {
	Groups []string
}

func (x ListOut) Human(context.Context) string {
	var w bytes.Buffer
	for _, u := range x.Groups {
		fmt.Fprintln(&w, u)
	}
	return w.String()
}

func (x GovGroupService) List(ctx context.Context, in *ListIn) (*ListOut, error) {
	// clone community repo locally
	community := git.LocalFromDir(files.WorkDir(ctx).Subdir("community"))
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.CommunityBranch); err != nil {
		return nil, err
	}
	// make changes to repo
	groups, err := List(ctx, community)
	if err != nil {
		return nil, err
	}
	// push to origin
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &ListOut{Groups: groups}, nil
}

func List(ctx context.Context, community git.Local) ([]string, error) {
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
