package entity

import (
	"context"
	"path/filepath"

	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/entityproto"
)

func (x EntityService[V]) List(ctx context.Context) ([]string, error) {
	local, err := git.MakeLocal(ctx)
	if err != nil {
		return nil, err
	}
	if err := local.CloneOrigin(ctx, x.Address); err != nil {
		return nil, err
	}
	users, err := x.ListLocal(ctx, local)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (x EntityService[V]) ListLocal(ctx context.Context, local git.Local) ([]string, error) {
	valueFileGlob := entityproto.ValueFilepath(x.Namespace, "*")
	m, err := local.Dir().Glob(valueFileGlob)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(m))
	for i := range m {
		entityDir, _ := filepath.Split(m[i])
		names[i] = filepath.Base(entityDir)
	}
	return names, nil
}
