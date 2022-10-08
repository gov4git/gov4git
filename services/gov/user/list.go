package user

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto"
)

type ListIn struct {
	CommunityBranch string `json:"community_branch"` // branch in community repo where user will be added
}

type ListOut struct {
	Users []string
}

func (x ListOut) Human(context.Context) string {
	var w bytes.Buffer
	for _, u := range x.Users {
		fmt.Fprintln(&w, u)
	}
	return w.String()
}

func (x GovUserService) List(ctx context.Context, in *ListIn) (*ListOut, error) {
	// clone community repo locally
	community, err := git.MakeLocalInCtx(ctx, "community")
	if err != nil {
		return nil, err
	}
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.CommunityBranch); err != nil {
		return nil, err
	}
	// make changes to repo
	users, err := List(ctx, community)
	if err != nil {
		return nil, err
	}
	// push to origin
	if err := community.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return &ListOut{Users: users}, nil
}

func List(ctx context.Context, community git.Local) ([]string, error) {
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
