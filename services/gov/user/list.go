package user

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto"
)

type GovUserListIn struct {
	CommunityBranch string `json:"community_branch"` // branch in community repo where user will be added
}

type GovUserListOut struct {
	Users []string
}

func (x GovUserListOut) Human(context.Context) string {
	var w bytes.Buffer
	for _, u := range x.Users {
		fmt.Fprintln(&w, u)
	}
	return w.String()
}

func (x GovUserService) UserList(ctx context.Context, in *GovUserListIn) (*GovUserListOut, error) {
	// clone community repo locally
	community := git.LocalFromDir(files.WorkDir(ctx).Subdir("community"))
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.CommunityBranch); err != nil {
		return nil, err
	}
	// make changes to repo
	users, err := GovUserList(ctx, community)
	if err != nil {
		return nil, err
	}
	// push to origin
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &GovUserListOut{Users: users}, nil
}

func GovUserList(ctx context.Context, community git.Local) ([]string, error) {
	userFileGlob := filepath.Join(proto.GovUsersDir, "*", proto.GovUserInfoFilebase)
	// glob for user files
	m, err := community.Dir().Glob(userFileGlob)
	if err != nil {
		return nil, err
	}
	// extract user names
	users := make([]string, len(m))
	for i := range m {
		userDir, _ := filepath.Split(m[i])
		users[i] = filepath.Base(userDir)
	}
	return users, nil
}
