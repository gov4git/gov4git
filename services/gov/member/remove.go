package member

import (
	"context"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto"
)

type GovMemberRemoveIn struct {
	User            string `json:"user"`
	Group           string `json:"group"`
	CommunityBranch string `json:"community_branch"` // branch in community repo where group will be added
}

type GovMemberRemoveOut struct{}

func (x GovMemberRemoveOut) Human(context.Context) string {
	return ""
}

func (x GovMemberService) MemberRemove(ctx context.Context, in *GovMemberRemoveIn) (*GovMemberRemoveOut, error) {
	// clone community repo locally
	community := git.LocalFromDir(files.WorkDir(ctx).Subdir("community"))
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.CommunityBranch); err != nil {
		return nil, err
	}
	// make changes to repo
	if err := GovRemoveMember(ctx, community, in.User, in.Group); err != nil {
		return nil, err
	}
	// push to origin
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &GovMemberRemoveOut{}, nil
}

func GovRemoveMember(ctx context.Context, community git.Local, user string, group string) error {
	mFile := filepath.Join(proto.GovGroupsDir, group, proto.GovMembersDirbase, user)
	// remove group file
	if err := community.Dir().Remove(mFile); err != nil {
		return err
	}
	// stage changes
	if err := community.Remove(ctx, []string{mFile}); err != nil {
		return err
	}
	// commit changes
	if err := community.Commitf(ctx, "gov: remove member user %v from group %v", user, group); err != nil {
		return err
	}
	return nil
}
