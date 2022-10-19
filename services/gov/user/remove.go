package user

import (
	"context"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/govproto"
)

type RemoveIn struct {
	Name            string `json:"name"`             // community unique handle for this user
	CommunityBranch string `json:"community_branch"` // branch in community repo where user will be added
}

type RemoveOut struct{}

func (x GovUserService) Remove(ctx context.Context, in *RemoveIn) (*RemoveOut, error) {
	// clone community repo locally
	community, err := git.MakeLocalInCtx(ctx, "community")
	if err != nil {
		return nil, err
	}
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.CommunityBranch); err != nil {
		return nil, err
	}
	// make changes to repo
	if err := Remove(ctx, community, in.Name); err != nil {
		return nil, err
	}
	// push to origin
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &RemoveOut{}, nil
}

func Remove(ctx context.Context, community git.Local, name string) error {
	userFile := filepath.Join(govproto.GovUsersDir, name, govproto.GovUserInfoFilebase)
	// remove user file
	if err := community.Dir().Remove(userFile); err != nil {
		return err
	}
	// stage changes
	if err := community.Remove(ctx, []string{userFile}); err != nil {
		return err
	}
	// commit changes
	if err := community.Commitf(ctx, "Remove user %v", name); err != nil {
		return err
	}
	return nil
}
