package user

import (
	"context"
	"path/filepath"
	"strconv"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/govproto"
)

func (x UserService) Get(ctx context.Context, name string, key string) ([]byte, error) {
	local, err := git.CloneOrigin(ctx, git.Origin(x))
	if err != nil {
		return nil, err
	}
	value, err := x.GetLocal(ctx, local, name, key)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (x UserService) GetLocal(ctx context.Context, local git.Local, name string, key string) ([]byte, error) {
	propFile := filepath.Join(govproto.GovUsersDir, name, govproto.GovUserMetaDirbase, key)
	data, err := local.Dir().ReadByteFile(propFile)
	if err != nil {
		return nil, err
	}
	return data.Bytes, nil
}

func (x UserService) GetFloat64Local(ctx context.Context, local git.Local, name string, key string) (float64, error) {
	v, err := x.GetLocal(ctx, local, name, key)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(string(v), 64)
}
