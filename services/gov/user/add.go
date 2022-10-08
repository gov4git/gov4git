package user

import (
	"context"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto"
)

type AddIn struct {
	Name            string `json:"name"`             // community unique handle for this user
	URL             string `json:"url"`              // user's public soul url
	CommunityBranch string `json:"community_branch"` // branch in community repo where user will be added
}

type AddOut struct{}

func (x AddOut) Human(context.Context) string {
	return ""
}

func (x GovUserService) Add(ctx context.Context, in *AddIn) (*AddOut, error) {
	// clone community repo locally
	community := git.LocalInDir(files.WorkDir(ctx).Subdir("community"))
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.CommunityBranch); err != nil {
		return nil, err
	}
	// make changes to repo
	if err := Add(ctx, community, in.Name, in.URL); err != nil {
		return nil, err
	}
	// push to origin
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &AddOut{}, nil
}

func Add(ctx context.Context, community git.Local, name string, url string) error {
	userFile := filepath.Join(proto.GovUsersDir, name, proto.GovUserInfoFilebase)
	// write user file
	stage := files.FormFiles{
		files.FormFile{Path: userFile, Form: proto.GovUserInfo{URL: url}},
	}
	if err := community.Dir().WriteFormFiles(ctx, stage); err != nil {
		return err
	}
	// stage changes
	if err := community.Add(ctx, stage.Paths()); err != nil {
		return err
	}
	// commit changes
	if err := community.Commitf(ctx, "gov: add user %v", name); err != nil {
		return err
	}
	return nil
}