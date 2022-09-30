package member

import (
	"context"
	"path/filepath"

	"github.com/petar/gitty/lib/files"
	"github.com/petar/gitty/lib/git"
	"github.com/petar/gitty/proto"
)

type GovMemberAddIn struct {
	User            string `json:"user"`
	Group           string `json:"group"`
	CommunityBranch string `json:"community_branch"` // branch in community repo where group will be added
}

type GovMemberAddOut struct{}

func (x GovMemberAddOut) Human(context.Context) string {
	return ""
}

func (x GovMemberService) MemberAdd(ctx context.Context, in *GovMemberAddIn) (*GovMemberAddOut, error) {
	// clone community repo locally
	community := git.LocalFromDir(files.WorkDir(ctx).Subdir("community"))
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.CommunityBranch); err != nil {
		return nil, err
	}
	// make changes to repo
	if err := GovMemberAdd(ctx, community, in.User, in.Group); err != nil {
		return nil, err
	}
	// push to origin
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &GovMemberAddOut{}, nil
}

func GovMemberAdd(ctx context.Context, community git.Local, user string, group string) error {
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
