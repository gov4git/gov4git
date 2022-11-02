package user

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/govproto"
)

func (x UserService) Set(
	ctx context.Context,
	name string,
	key string,
	value []byte,
) error {
	local, err := git.CloneOrigin(ctx, git.Origin(x))
	if err != nil {
		return err
	}
	if err := x.SetLocal(ctx, local, name, key, value); err != nil {
		return err
	}
	if err := local.PushUpstream(ctx); err != nil {
		return err
	}
	return nil
}

func (x UserService) SetLocal(
	ctx context.Context,
	local git.Local,
	name string,
	key string,
	value []byte,
) error {
	if err := x.SetLocalStageOnly(ctx, local, name, key, value); err != nil {
		return err
	}
	if err := local.Commitf(ctx, "Change property %v of user %v", key, name); err != nil {
		return err
	}
	return nil
}

// XXX: sanitize key
// XXX: prevent overwrite
func (x UserService) SetLocalStageOnly(
	ctx context.Context,
	community git.Local,
	name string,
	key string,
	value []byte,
) error {
	propFile := filepath.Join(govproto.GovUsersDir, name, govproto.GovUserMetaDirbase, key)
	stage := files.ByteFiles{
		files.ByteFile{Path: propFile, Bytes: value},
	}
	if err := community.Dir().WriteByteFiles(stage); err != nil {
		return err
	}
	if err := community.Add(ctx, stage.Paths()); err != nil {
		return err
	}
	return nil
}

func (x UserService) SetFloat64Local(
	ctx context.Context,
	local git.Local,
	name string,
	key string,
	value float64,
) error {
	return x.SetLocal(ctx, local, name, key, []byte(fmt.Sprintf("%v", value)))
}

func (x UserService) SetFloat64LocalStageOnly(
	ctx context.Context,
	local git.Local,
	name string,
	key string,
	value float64,
) error {
	return x.SetLocalStageOnly(ctx, local, name, key, []byte(fmt.Sprintf("%v", value)))
}
