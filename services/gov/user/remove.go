package user

import (
	"context"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/govproto"
)

func (x UserService) Remove(ctx context.Context, name string) error {
	local, err := git.MakeLocal(ctx)
	if err != nil {
		return err
	}
	if err := local.CloneOrigin(ctx, git.Origin(x)); err != nil {
		return err
	}
	if err := Remove(ctx, local, name); err != nil {
		return err
	}
	if err := local.PushUpstream(ctx); err != nil {
		return err
	}
	return nil
}

func Remove(ctx context.Context, local git.Local, name string) error {
	userFile := filepath.Join(govproto.GovUsersDir, name, govproto.GovUserInfoFilebase)
	if err := local.Dir().Remove(userFile); err != nil {
		return err
	}
	if err := local.Remove(ctx, []string{userFile}); err != nil {
		return err
	}
	if err := local.Commitf(ctx, "Remove user %v", name); err != nil {
		return err
	}
	return nil
}
