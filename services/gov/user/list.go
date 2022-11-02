package user

import (
	"context"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/govproto"
)

func (x UserService) List(ctx context.Context) ([]string, error) {
	local, err := git.MakeLocal(ctx)
	if err != nil {
		return nil, err
	}
	if err := local.CloneOrigin(ctx, git.Origin(x)); err != nil {
		return nil, err
	}
	users, err := List(ctx, local)
	if err != nil {
		return nil, err
	}
	if err := local.PushUpstream(ctx); err != nil {
		return nil, err
	}
	return users, nil
}

func List(ctx context.Context, community git.Local) ([]string, error) {
	userFileGlob := filepath.Join(govproto.GovUsersDir, "*", govproto.GovUserInfoFilebase)
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
