package policy

import (
	"context"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/govproto"
)

type GetIn struct {
	Dir             string `json:"dir"`
	CommunityBranch string `json:"community_branch"` // branch in community repo where policy will be added
}

type GetOut struct {
	Policy govproto.GovDirPolicy `json:"policy"`
}

func (x GovPolicyService) Get(ctx context.Context, in *GetIn) (*GetOut, error) {
	// clone community repo locally
	community, err := git.MakeLocalInCtx(ctx, "community")
	if err != nil {
		return nil, err
	}
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.CommunityBranch); err != nil {
		return nil, err
	}
	// read from repo
	policy, err := Get(ctx, community, in.Dir)
	if err != nil {
		return nil, err
	}
	return &GetOut{Policy: *policy}, nil
}

func Get(ctx context.Context, community git.Local, dir string) (*govproto.GovDirPolicy, error) {
	policyFile := filepath.Join(dir, govproto.GovRoot, govproto.GovDirPolicyFilebase)
	// read policy file
	var policy govproto.GovDirPolicy
	if _, err := community.Dir().ReadFormFile(ctx, policyFile, &policy); err != nil {
		return nil, err
	}
	return &policy, nil
}
