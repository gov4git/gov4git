package group

import (
	"context"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/govproto"
)

type ListIn struct {
	CommunityBranch string `json:"community_branch"` // branch in community repo where group will be added
}

type ListOut struct {
	Groups []string
}

func (x GovGroupService) List(ctx context.Context, in *ListIn) (*ListOut, error) {
	// clone community repo locally
	community, err := git.MakeLocalInCtx(ctx, "community")
	if err != nil {
		return nil, err
	}
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
	// glob for group files
	m, err := community.Dir().Glob(govproto.GroupInfoFilepath("*"))
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
