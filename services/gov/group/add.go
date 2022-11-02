package group

import (
	"context"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/govproto"
)

type AddIn struct {
	Name            string `json:"name"`             // community unique handle for this group
	CommunityBranch string `json:"community_branch"` // branch in community repo where group will be added
}

type AddOut struct{}

func (x GovGroupService) Add(ctx context.Context, in *AddIn) (*AddOut, error) {
	// clone community repo locally
	community, err := git.MakeLocalInCtx(ctx, "community")
	if err != nil {
		return nil, err
	}
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.CommunityBranch); err != nil {
		return nil, err
	}
	// make changes to repo
	if err := x.AddLocal(ctx, community, in.Name); err != nil {
		return nil, err
	}
	// push to origin
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &AddOut{}, nil
}

func (x GovGroupService) AddLocal(ctx context.Context, community git.Local, name string) error {
	if err := x.AddLocalStageOnly(ctx, community, name); err != nil {
		return err
	}
	if err := community.Commitf(ctx, "Add group %v", name); err != nil {
		return err
	}
	return nil
}

func (x GovGroupService) AddLocalStageOnly(ctx context.Context, community git.Local, name string) error {
	// write group file
	stage := files.FormFiles{
		files.FormFile{Path: govproto.GroupInfoFilepath(name), Form: govproto.GovGroupInfo{}},
	}
	if err := community.Dir().WriteFormFiles(ctx, stage); err != nil {
		return err
	}
	// stage changes
	if err := community.Add(ctx, stage.Paths()); err != nil {
		return err
	}
	return nil
}

func GetInfo(ctx context.Context, community git.Local, name string) (*govproto.GovGroupInfo, error) {
	groupInfoPath := govproto.GroupInfoFilepath(name)
	var info govproto.GovGroupInfo
	if _, err := community.Dir().ReadFormFile(ctx, groupInfoPath, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

type UserInfo struct {
	UserName string               `json:"user_name"`
	UserInfo govproto.GovUserInfo `json:"user_info"`
}
