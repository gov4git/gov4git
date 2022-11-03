package entity

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/entityproto"
)

func (x EntityService[V]) SetProp(
	ctx context.Context,
	name string,
	key string,
	value []byte,
) error {
	local, err := git.CloneOrigin(ctx, x.Address)
	if err != nil {
		return err
	}
	if err := x.SetPropLocal(ctx, local, name, key, value); err != nil {
		return err
	}
	if err := local.PushUpstream(ctx); err != nil {
		return err
	}
	return nil
}

func (x EntityService[V]) SetPropLocal(
	ctx context.Context,
	local git.Local,
	name string,
	key string,
	value []byte,
) error {
	if err := x.SetPropLocalStageOnly(ctx, local, name, key, value); err != nil {
		return err
	}
	if err := local.Commitf(ctx, "Change property %v of %v entity %v", key, x.Namespace, name); err != nil {
		return err
	}
	return nil
}

// XXX: sanitize key
// XXX: prevent overwrite
func (x EntityService[V]) SetPropLocalStageOnly(
	ctx context.Context,
	local git.Local,
	name string,
	key string,
	value []byte,
) error {
	propFilepath := entityproto.PropValueFilepath(x.Namespace, name, key)
	stage := files.ByteFiles{
		files.ByteFile{
			Path:  propFilepath,
			Bytes: value,
		},
	}
	if err := local.Dir().WriteByteFiles(stage); err != nil {
		return err
	}
	if err := local.Add(ctx, stage.Paths()); err != nil {
		return err
	}
	return nil
}

func (x EntityService[V]) SetPropFloat64Local(
	ctx context.Context,
	local git.Local,
	name string,
	key string,
	value float64,
) error {
	return x.SetPropLocal(ctx, local, name, key, []byte(fmt.Sprintf("%v", value)))
}

func (x EntityService[V]) SetPropFloat64LocalStageOnly(
	ctx context.Context,
	local git.Local,
	name string,
	key string,
	value float64,
) error {
	return x.SetPropLocalStageOnly(ctx, local, name, key, []byte(fmt.Sprintf("%v", value)))
}
