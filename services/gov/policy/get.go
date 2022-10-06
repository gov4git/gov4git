package policy

import (
	"context"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto"
)

type GovPolicyGetIn struct {
	Dir             string `json:"dir"`
	CommunityBranch string `json:"community_branch"` // branch in community repo where policy will be added
}

type GovPolicyGetOut struct {
	Policy proto.GovDirPolicy `json:"policy"`
}

func (x GovPolicyGetOut) Human(ctx context.Context) string {
	data, _ := form.EncodeForm(ctx, x.Policy)
	return string(data)
}

func (x GovPolicyService) PolicyGet(ctx context.Context, in *GovPolicyGetIn) (*GovPolicyGetOut, error) {
	// clone community repo locally
	community := git.LocalFromDir(files.WorkDir(ctx).Subdir("community"))
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.CommunityBranch); err != nil {
		return nil, err
	}
	// read from repo
	policy, err := GovPolicyGet(ctx, community, in.Dir)
	if err != nil {
		return nil, err
	}
	return &GovPolicyGetOut{Policy: *policy}, nil
}

func GovPolicyGet(ctx context.Context, community git.Local, dir string) (*proto.GovDirPolicy, error) {
	policyFile := filepath.Join(dir, proto.GovRoot, proto.GovDirPolicyFilebase)
	// read policy file
	var policy proto.GovDirPolicy
	if _, err := community.Dir().ReadFormFile(ctx, policyFile, &policy); err != nil {
		return nil, err
	}
	return &policy, nil
}
