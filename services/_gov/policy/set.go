package policy

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/govproto"
)

type SetIn struct {
	Dir             string  `json:"dir"`
	Arb             string  `json:"arb"`
	Group           string  `json:"group"`
	Threshold       float64 `json:"threshold"`
	CommunityBranch string  `json:"community_branch"` // branch in community repo where policy will be added
}

type SetOut struct{}

func (x GovPolicyService) Set(ctx context.Context, in *SetIn) (*SetOut, error) {
	// clone community repo locally
	community, err := git.MakeLocalInCtx(ctx, "community")
	if err != nil {
		return nil, err
	}
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.CommunityBranch); err != nil {
		return nil, err
	}
	// make changes to repo
	if err := Set(ctx, community, in.Dir, in.Arb, in.Group, in.Threshold); err != nil {
		return nil, err
	}
	// push to origin
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &SetOut{}, nil
}

func Set(ctx context.Context, community git.Local, dir string, arb string, group string, threshold float64) error {
	// TODO: verify group and threshold

	// build policy
	var policy govproto.GovDirPolicy
	switch arb {
	case "quorum":
		policy = govproto.GovDirPolicy{
			Change: govproto.GovArbitration{
				Quorum: &govproto.GovQuorum{
					Group:     group,
					Threshold: uint32(threshold),
				},
			},
		}
	default:
		return fmt.Errorf("unknown directory policy")
	}

	policyFile := filepath.Join(dir, govproto.GovRoot, govproto.GovDirPolicyFilebase)
	// write policy file
	stage := files.FormFiles{
		files.FormFile{Path: policyFile, Form: policy},
	}
	if err := community.Dir().WriteFormFiles(ctx, stage); err != nil {
		return err
	}
	// stage changes
	if err := community.Add(ctx, stage.Paths()); err != nil {
		return err
	}
	// commit changes
	if err := community.Commitf(ctx, "Change directory %v policy to %v in group %v", dir, arb, group); err != nil {
		return err
	}
	return nil
}
