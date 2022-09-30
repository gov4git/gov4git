package policy

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/petar/gitty/lib/files"
	"github.com/petar/gitty/lib/git"
	"github.com/petar/gitty/proto"
)

type GovPolicySetIn struct {
	Dir             string  `json:"dir"`
	Arb             string  `json:"arb"`
	Group           string  `json:"group"`
	Threshold       float64 `json:"threshold"`
	CommunityBranch string  `json:"community_branch"` // branch in community repo where policy will be added
}

type GovPolicySetOut struct{}

func (x GovPolicySetOut) Human() string {
	return ""
}

func (x GovPolicyService) PolicySet(ctx context.Context, in *GovPolicySetIn) (*GovPolicySetOut, error) {
	// clone community repo locally
	community := git.LocalFromDir(files.WorkDir(ctx).Subdir("community"))
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.CommunityBranch); err != nil {
		return nil, err
	}
	// make changes to repo
	if err := GovPolicySet(ctx, community, in.Dir, in.Arb, in.Group, in.Threshold); err != nil {
		return nil, err
	}
	// push to origin
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &GovPolicySetOut{}, nil
}

func GovPolicySet(ctx context.Context, community git.Local, dir string, arb string, group string, threshold float64) error {
	// TODO: verify group and threshold

	// build policy
	var policy proto.GovDirPolicy
	switch arb {
	case "quorum":
		policy = proto.GovDirPolicy{
			Change: proto.GovArbitration{
				Quorum: &proto.GovQuorum{
					Group:     group,
					Threshold: uint32(threshold),
				},
			},
		}
	default:
		return fmt.Errorf("unknown directory policy")
	}

	policyFile := filepath.Join(dir, proto.GovRoot, proto.GovDirPolicyFilebase)
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
	if err := community.Commitf(ctx, "gov: change directory %v policy to %v in group %v", dir, arb, group); err != nil {
		return err
	}
	return nil
}
