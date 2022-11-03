package entity

import (
	"context"

	"github.com/gov4git/gov4git/lib/files"
	"github.com/gov4git/gov4git/lib/form"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/entityproto"
)

func (x EntityService[V]) Add(ctx context.Context, name string, value V) error {
	local, err := git.MakeLocal(ctx)
	if err != nil {
		return err
	}
	if err := local.CloneOrigin(ctx, x.Address); err != nil {
		return err
	}
	if err := x.AddLocal(ctx, local, name, value); err != nil {
		return err
	}
	if err := local.PushUpstream(ctx); err != nil {
		return err
	}
	return nil
}

func (x EntityService[V]) AddLocal(ctx context.Context, local git.Local, name string, value V) error {
	if err := x.AddLocalStageOnly(ctx, local, name, value); err != nil {
		return err
	}
	if err := local.Commitf(ctx, "Add %v entity %v", x.Namespace, name); err != nil {
		return err
	}
	return nil
}

func (x EntityService[V]) AddLocalStageOnly(ctx context.Context, local git.Local, name string, value V) error {
	valueFilepath := entityproto.ValueFilepath(x.Namespace, name)
	stage := files.FormFiles{
		files.FormFile{
			Path: valueFilepath,
			Form: value,
		},
	}
	if err := local.Dir().WriteFormFiles(ctx, stage); err != nil {
		return err
	}
	if err := local.Add(ctx, stage.Paths()); err != nil {
		return err
	}
	return nil
}

func (x EntityService[V]) Get(ctx context.Context, local git.Local, name string) (V, error) {
	valuePath := entityproto.ValueFilepath(x.Namespace, name)
	var value V
	if _, err := local.Dir().ReadFormFile(ctx, valuePath, &value); err != nil {
		return value, err
	}
	return value, nil
}

type EntityValue[V form.Form] struct {
	Name  string
	Value V
}

type EntityValues[V form.Form] []EntityValue[V]

func (x EntityService[V]) GetMany(ctx context.Context, local git.Local, names []string) (EntityValues[V], error) {
	entityValues := make(EntityValues[V], len(names))
	for i, name := range names {
		u, err := x.Get(ctx, local, name)
		if err != nil {
			return nil, err
		}
		entityValues[i] = EntityValue[V]{Name: name, Value: u}
	}
	return entityValues, nil
}
