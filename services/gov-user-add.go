package services

import (
	"context"
	"path/filepath"

	"github.com/petar/gitty/lib/files"
	"github.com/petar/gitty/lib/git"
	"github.com/petar/gitty/proto"
)

type GovAddUserIn struct {
	Name            string `json:"name"`             // community unique handle for this user
	URL             string `json:"url"`              // user's public soul url
	CommunityBranch string `json:"community_branch"` // branch in community repo where user will be added
}

type GovAddUserOut struct{}

func (x GovAddUserOut) Human() string {
	return ""
}

func (x GovService) AddUser(ctx context.Context, in *GovAddUserIn) (*GovAddUserOut, error) {
	// clone community repo locally
	community := git.LocalFromDir(files.WorkDir(ctx).Subdir("community"))
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.CommunityBranch); err != nil {
		return nil, err
	}
	// make changes to repo
	if err := GovAddUser(ctx, community, in.Name, in.URL); err != nil {
		return nil, err
	}
	// push to origin
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &GovAddUserOut{}, nil
}

func GovAddUser(ctx context.Context, community git.Local, name string, url string) error {
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
	if err := community.Commitf(ctx, "add user %v", name); err != nil {
		return err
	}
	return nil
}
