package user

import (
	"context"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto"
	"github.com/gov4git/gov4git/proto/govproto"
)

func (x UserService) Add(ctx context.Context, name string, addr proto.Address) error {
	community, err := git.MakeLocal(ctx)
	if err != nil {
		return err
	}
	if err := community.CloneOrigin(ctx, git.Origin(x)); err != nil {
		return err
	}
	if err := AddLocal(ctx, community, name, addr); err != nil {
		return err
	}
	if err := community.PushUpstream(ctx); err != nil {
		return err
	}
	return nil
}

func AddLocal(ctx context.Context, local git.Local, name string, addr proto.Address) error {
	if err := AddLocalStageOnly(ctx, local, name, addr); err != nil {
		return err
	}
	if err := local.Commitf(ctx, "Add user %v", name); err != nil {
		return err
	}
	return nil
}

func AddLocalStageOnly(ctx context.Context, local git.Local, name string, addr proto.Address) error {
	userFile := govproto.UserInfoFilepath(name)
	stage := files.FormFiles{
		files.FormFile{
			Path: userFile,
			Form: govproto.GovUserInfo{
				Address: addr,
			}},
	}
	if err := local.Dir().WriteFormFiles(ctx, stage); err != nil {
		return err
	}
	if err := local.Add(ctx, stage.Paths()); err != nil {
		return err
	}
	return nil
}

func GetInfo(ctx context.Context, local git.Local, name string) (*govproto.GovUserInfo, error) {
	userInfoPath := govproto.UserInfoFilepath(name)
	var info govproto.GovUserInfo
	if _, err := local.Dir().ReadFormFile(ctx, userInfoPath, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

type UserInfo struct {
	UserName string               `json:"user_name"`
	UserInfo govproto.GovUserInfo `json:"user_info"`
}

type UserInfos []UserInfo

func GetInfos(ctx context.Context, local git.Local, usernames []string) (UserInfos, error) {
	userInfo := make(UserInfos, len(usernames))
	for i, n := range usernames {
		u, err := GetInfo(ctx, local, n)
		if err != nil {
			return nil, err
		}
		userInfo[i] = UserInfo{UserName: n, UserInfo: *u}
	}
	return userInfo, nil
}
