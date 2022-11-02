package user

import (
	"context"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/govproto"
)

type AddIn struct {
	Name            string `json:"name"`             // community unique handle for this user
	URL             string `json:"url"`              // user's public soul url
	CommunityBranch string `json:"community_branch"` // branch in community repo where user will be added
}

type AddOut struct{}

func (x GovUserService) Add(ctx context.Context, in *AddIn) (*AddOut, error) {
	// clone community repo locally
	community, err := git.MakeLocalInCtx(ctx, "community")
	if err != nil {
		return nil, err
	}
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.CommunityBranch); err != nil {
		return nil, err
	}
	// make changes to repo
	if err := x.AddLocal(ctx, community, in.Name, in.URL); err != nil {
		return nil, err
	}
	// push to origin
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &AddOut{}, nil
}

func (x GovUserService) AddLocal(ctx context.Context, community git.Local, name string, url string) error {
	if err := x.AddLocalStageOnly(ctx, community, name, url); err != nil {
		return err
	}
	if err := community.Commitf(ctx, "Add user %v", name); err != nil {
		return err
	}
	return nil
}

func (x GovUserService) AddLocalStageOnly(ctx context.Context, community git.Local, name string, url string) error {
	userFile := govproto.UserInfoFilepath(name)
	// write user file
	stage := files.FormFiles{
		files.FormFile{Path: userFile, Form: govproto.GovUserInfo{PublicURL: url}},
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

func GetInfo(ctx context.Context, community git.Local, name string) (*govproto.GovUserInfo, error) {
	userInfoPath := govproto.UserInfoFilepath(name)
	var info govproto.GovUserInfo
	if _, err := community.Dir().ReadFormFile(ctx, userInfoPath, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

type UserInfo struct {
	UserName string               `json:"user_name"`
	UserInfo govproto.GovUserInfo `json:"user_info"`
}

type UserInfos []UserInfo

func GetInfos(ctx context.Context, community git.Local, usernames []string) (UserInfos, error) {
	userInfo := make(UserInfos, len(usernames))
	for i, n := range usernames {
		u, err := GetInfo(ctx, community, n)
		if err != nil {
			return nil, err
		}
		userInfo[i] = UserInfo{UserName: n, UserInfo: *u}
	}
	return userInfo, nil
}
