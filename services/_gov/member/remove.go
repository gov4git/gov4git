package member

import (
	"context"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/govproto"
)

type RemoveIn struct {
	User            string `json:"user"`
	Group           string `json:"group"`
	CommunityBranch string `json:"community_branch"` // branch in community repo where group will be added
}

type RemoveOut struct{}

func (x GovMemberService) Remove(ctx context.Context, in *RemoveIn) (*RemoveOut, error) {
	// clone community repo locally
	community, err := git.MakeLocalInCtx(ctx, "community")
	if err != nil {
		return nil, err
	}
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.CommunityBranch); err != nil {
		return nil, err
	}
	// make changes to repo
	if err := Remove(ctx, community, in.User, in.Group); err != nil {
		return nil, err
	}
	// push to origin
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &RemoveOut{}, nil
}

func Remove(ctx context.Context, community git.Local, user string, group string) error {
	file := govproto.GroupMemberFilepath(group, user)
	// remove group file
	if err := community.Dir().Remove(file); err != nil {
		return err
	}
	// stage changes
	if err := community.Remove(ctx, []string{file}); err != nil {
		return err
	}
	// commit changes
	if err := community.Commitf(ctx, "Remove member user %v from group %v", user, group); err != nil {
		return err
	}
	return nil
}
