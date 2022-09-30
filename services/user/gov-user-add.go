package user

import (
	"context"
	"path/filepath"

	"github.com/petar/gitty/lib/files"
	"github.com/petar/gitty/lib/git"
	"github.com/petar/gitty/proto"
)

type GovUserAddIn struct {
	Name            string `json:"name"`             // community unique handle for this user
	URL             string `json:"url"`              // user's public soul url
	CommunityBranch string `json:"community_branch"` // branch in community repo where user will be added
}

type GovUserAddOut struct{}

func (x GovUserAddOut) Human(context.Context) string {
	return ""
}

func (x GovUserService) UserAdd(ctx context.Context, in *GovUserAddIn) (*GovUserAddOut, error) {
	// clone community repo locally
	community := git.LocalFromDir(files.WorkDir(ctx).Subdir("community"))
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.CommunityBranch); err != nil {
		return nil, err
	}
	// make changes to repo
	if err := GovUserAdd(ctx, community, in.Name, in.URL); err != nil {
		return nil, err
	}
	// push to origin
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &GovUserAddOut{}, nil
}

func GovUserAdd(ctx context.Context, community git.Local, name string, url string) error {
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
