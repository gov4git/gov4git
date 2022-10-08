package member

import (
	"context"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto"
)

type AddIn struct {
	User            string `json:"user"`
	Group           string `json:"group"`
	CommunityBranch string `json:"community_branch"` // branch in community repo where group will be added
}

type AddOut struct{}

func (x AddOut) Human(context.Context) string {
	return ""
}

func (x GovMemberService) Add(ctx context.Context, in *AddIn) (*AddOut, error) {
	// clone community repo locally
	community, err := git.MakeLocalCtx(ctx, "community")
	if err != nil {
		return nil, err
	}
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.CommunityBranch); err != nil {
		return nil, err
	}
	// make changes to repo
	if err := Add(ctx, community, in.User, in.Group); err != nil {
		return nil, err
	}
	// push to origin
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &AddOut{}, nil
}

func Add(ctx context.Context, community git.Local, user string, group string) error {
	mFile := filepath.Join(proto.GovGroupsDir, group, proto.GovMembersDirbase, user)
	// write group file
	stage := files.ByteFiles{
		files.ByteFile{Path: mFile},
	}
	if err := community.Dir().WriteByteFiles(stage); err != nil {
		return err
	}
	// stage changes
	if err := community.Add(ctx, stage.Paths()); err != nil {
		return err
	}
	// commit changes
	if err := community.Commitf(ctx, "gov: add member user %v to group %v", user, group); err != nil {
		return err
	}
	return nil
}
