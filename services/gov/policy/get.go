package policy

import (
	"context"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto"
)

type GetIn struct {
	Dir             string `json:"dir"`
	CommunityBranch string `json:"community_branch"` // branch in community repo where policy will be added
}

type GetOut struct {
	Policy proto.GovDirPolicy `json:"policy"`
}

func (x GetOut) Human(ctx context.Context) string {
	data, _ := form.EncodeForm(ctx, x.Policy)
	return string(data)
}

func (x GovPolicyService) Get(ctx context.Context, in *GetIn) (*GetOut, error) {
	// clone community repo locally
	community := git.LocalInDir(files.WorkDir(ctx).Subdir("community"))
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

func Get(ctx context.Context, community git.Local, dir string) (*proto.GovDirPolicy, error) {
	policyFile := filepath.Join(dir, proto.GovRoot, proto.GovDirPolicyFilebase)
	// read policy file
	var policy proto.GovDirPolicy
	if _, err := community.Dir().ReadFormFile(ctx, policyFile, &policy); err != nil {
		return nil, err
	}
	return &policy, nil
}
