package member

import (
	"context"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/govproto"
)

type AddIn struct {
	User            string `json:"user"`
	Group           string `json:"group"`
	CommunityBranch string `json:"community_branch"` // branch in community repo where group will be added
}

type AddOut struct{}

func (x GovMemberService) Add(ctx context.Context, in *AddIn) (*AddOut, error) {
	// clone community repo locally
	community, err := git.MakeLocalInCtx(ctx, "community")
	if err != nil {
		return nil, err
	}
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.CommunityBranch); err != nil {
		return nil, err
	}
	// make changes to repo
	if err := x.AddLocal(ctx, community, in.User, in.Group); err != nil {
		return nil, err
	}
	// push to origin
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &AddOut{}, nil
}

func (x GovMemberService) AddLocal(ctx context.Context, community git.Local, user string, group string) error {
	if err := x.AddLocalStageOnly(ctx, community, user, group); err != nil {
		return err
	}
	if err := community.Commitf(ctx, "Add member user %v to group %v", user, group); err != nil {
		return err
	}
	return nil
}

func (x GovMemberService) AddLocalStageOnly(ctx context.Context, community git.Local, user string, group string) error {
	file := govproto.GroupMemberFilepath(group, user)
	// write group file
	stage := files.ByteFiles{
		files.ByteFile{Path: file},
	}
	if err := community.Dir().WriteByteFiles(stage); err != nil {
		return err
	}
	// stage changes
	if err := community.Add(ctx, stage.Paths()); err != nil {
		return err
	}
	return nil
}
